package boot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evan-forbes/ethlab/cmd"
	"github.com/evan-forbes/ethlab/contracts/ens"
	"github.com/evan-forbes/ethlab/server"
	"github.com/evan-forbes/ethlab/thereum"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

/* TODO:
- impl the ability to export the chain / use different databases
*/
func Boot(c *cli.Context) error {
	// load config
	config := thereum.DefaultConfig()
	if cPath := c.String("config"); cPath != "" {
		custConfig, err := thereum.ConfigFromFile(cPath)
		if err != nil {
			return errors.Wrapf(err, "failure to load config from path %s:", cPath)
		}
		config = custConfig
	}

	// listen for ctrl + c cancels and start a global context/waitgroup for the app
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()

	// start the thereum backend
	eth, err := thereum.New(config, nil)
	if err != nil {
		log.Fatal(err)
	}
	mngr.WG.Add(1)
	go eth.Run(mngr.Ctx, mngr.WG)

	// start the http server
	srvr := server.NewServer(mngr.Ctx, fmt.Sprintf("%s:%d", config.Host, config.Port), eth)
	// start the
	go func() {
		log.Fatal(srvr.ListenAndServe())
	}()

	// start the websocket server
	go func() {
		log.Fatal(srvr.ServeWS(fmt.Sprintf("%s:%d", config.WSHost, config.WSPort)))
	}()

	// wait for a hot second to make sure the servers are up and running // TODO: use something more definitive
	time.Sleep(time.Millisecond * 88)

	// dial for a client to
	client, err := ethclient.Dial(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
	if err != nil {
		log.Fatal("failed to connect to ethlab", err)
	}

	// prepare root to deploy the ethlab contracts
	root := eth.Accounts["root"]
	root.TxOpts.GasLimit = 1000000000

	err = deployBaseContracts(client, root.TxOpts)
	if err != nil {
		return err
	}

	<-mngr.Done()
	return nil
}

func deployBaseContracts(client *ethclient.Client, opts *bind.TransactOpts) error {
	// depoly the ethlab version of the ens
	ensAddr, _, _, err := ens.DeployENS(opts, client)
	if err != nil {
		log.Fatal("failed to deploy ENS ", err)
	}
	fmt.Println("ENS deployed: ", ensAddr.Hex())
	return nil
}
