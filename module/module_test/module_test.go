package module_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evan-forbes/ethlab/cmd"
	"github.com/evan-forbes/ethlab/module"
	"github.com/evan-forbes/ethlab/server"
)

func TestRequestETH(t *testing.T) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()
	_, err := server.LaunchServer(mngr.Ctx, mngr.WG)
	if err != nil {
		t.Error(err)
	}
	usr, err := module.StarterKit("127.0.0.1:8000")
	if err != nil {
		t.Error(err)
		return
	}
	err = module.RequestETH("127.0.0.1:8000", usr.NewTxOpts().From.Hex(), big.NewInt(2000000000000000000))
	if err != nil {
		t.Error(err)
		return
	}
	err = module.RequestETH("127.0.0.1:8000", usr.NewTxOpts().From.Hex(), big.NewInt(2000000000000000000))
	if err != nil {
		t.Error(err)
		return
	}
	client, err := ethclient.Dial("http://127.0.0.1:8000")
	if err != nil {
		t.Error(err)
		return
	}
	// time.Sleep(time.Millisecond * 500)
	bal, err := client.BalanceAt(mngr.Ctx, usr.NewTxOpts().From, nil)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(bal.String())
	mngr.Cancel()
	<-mngr.Done()
}

// This is where I'm getting a funky error.
// Without sleeping, there are no errors, but the balance is still zero
// - an error should be thrown or some process should block until the server is ready
// With sleeping, we get an nonce too low error, so the transaction cannot be added to the pool
// - is the nonce not increasing? which tx cannot be added to the pool? root's or the new users?
// If we don't wait, then none of the balances are recorded, if we do wait, then only the second balance is set...
func TestStarterKit(t *testing.T) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()
	_, err := server.LaunchServer(mngr.Ctx, mngr.WG)
	if err != nil {
		t.Error(err)
	}
	// time.Sleep(time.Second)
	usr, err := module.StarterKit("127.0.0.1:8000")
	if err != nil {
		t.Error(err)
		return
	}
	usr2, err := module.StarterKit("127.0.0.1:8000")
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second * 1)
	bal2, err := usr2.Balance()
	if err != nil {
		t.Error(err)
		return
	}

	bal, err := usr.Balance()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("balances:", bal2.String(), bal.String())
	if bal2.Cmp(big.NewInt(0)) != 1 {
		fmt.Println(bal2.String())
		t.Error("starter kit balance is incorrect")
		return
	}
	if bal.Cmp(big.NewInt(0)) != 1 {
		fmt.Println(bal.String())
		t.Error("starter kit balance is incorrect")
		return
	}
	mngr.Cancel()
	<-mngr.Done()
}

// func TestDeploy(t *testing.T) {
// 	deploy := ens.Deploy
// 	mngr := cmd.NewManager(context.Background(), nil)
// 	go mngr.Listen()
// 	_, err := server.LaunchServer(mngr.Ctx, mngr.WG)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	usr, err := module.StarterKit("127.0.0.1:8000")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// }
