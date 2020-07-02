package ens_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/evan-forbes/ethlab/cmd"
	"github.com/evan-forbes/ethlab/contracts/ens"
	"github.com/evan-forbes/ethlab/module"
	"github.com/evan-forbes/ethlab/server"
)

func TestENS(t *testing.T) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()
	_, err := server.LaunchServer(mngr.Ctx, mngr.WG)
	if err != nil {
		t.Error(err)
	}
	// get the ens address
	addr, err := module.ENSAddress("127.0.0.1:8000")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 122)
	// create a new user that has some eth
	usr, err := module.StarterKit("127.0.0.1:8000")
	if err != nil {
		t.Error(err)
		fmt.Println(err)
		return
	}
	time.Sleep(time.Millisecond * 122)
	bal, err := usr.Client.BalanceAt(mngr.Ctx, usr.NewTxOpts().From, nil)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(bal.String())
	// bind to the ens contract
	ens, err := ens.NewENS(common.HexToAddress(addr), usr.Client)
	if err != nil {
		t.Error(err)
	}
	// create a query and log channel
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(addr)},
	}
	logs := make(chan types.Log)
	// start streaming logs
	_, err = usr.Client.SubscribeFilterLogs(mngr.Ctx, query, logs)
	if err != nil {
		t.Error(err)
	}
	// print streamed logs
	go func() {
		defer mngr.Cancel()
		for {
			select {
			// case <-sub.Err():
			// 	return
			case l, ok := <-logs:
				if !ok {
					return
				}
				fmt.Println(l)
				return
			}
		}
	}()
	name := common.LeftPadBytes([]byte("uniswap"), 32)
	var name32 [32]byte
	copy(name32[:], name[:32])
	tx, err := ens.Add(usr.NewTxOpts(), name32, common.HexToAddress("0x514910771af9ca656af840dff83e8264ecf986ca"))
	fmt.Println("error", err)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("transaction", tx)
	// wait for cancel
	<-mngr.Done()
}
