package module_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/go-ethereum/ethclient"
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
	err = module.RequestETH("127.0.0.1:8000", usr.NewTxOpts().From.Hex(), big.NewInt(1))
	if err != nil {
		t.Error(err)
		return
	}
	client, err := ethclient.Dial("http://127.0.0.1:8000")
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Millisecond * 100)
	bal, err := client.BalanceAt(mngr.Ctx, usr.NewTxOpts().From, nil)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(bal.String())
	mngr.Cancel()
	<-mngr.Done()
}
