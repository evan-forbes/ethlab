package thereum

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/evan-forbes/ethlab/cmd"
)

func setupThereum(t *testing.T) (*Thereum, *cmd.Manager) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()
	heads := make(chan *types.Header)

	eth, err := New(defaultConfig(), common.Address{})
	if err != nil {
		t.Error(err)
	}
	mngr.WG.Add(1)
	go eth.Run(mngr.Ctx, mngr.WG)
	return eth, mngr, nil
}

func TestBoot(t *testing.T) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()
	heads := make(chan *types.Header)

	eth, err := New(defaultConfig(), common.Address{})
	if err != nil {
		t.Error(err)
	}
	mngr.WG.Add(1)
	go eth.Run(mngr.Ctx, mngr.WG)
	sub := eth.Events.SubscribeNewHeads(heads)
	mngr.WG.Add(1)
	go func() {
		defer mngr.WG.Done()
		for {
			select {
			case head := <-heads:
				fmt.Println(head.Hash().Hex())
			case err := <-sub.Err():
				t.Error(err)
				return
			case <-mngr.Ctx.Done():
				return
			}
		}
	}()
	mngr.WG.Wait()
}

// ContractCaller defines the methods needed to allow operating with contract on a read
// only basis.
type ContractCaller interface {
	// CodeAt returns the code of the given account. This is needed to differentiate
	// between contract internal errors and the local chain being out of sync.
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	// ContractCall executes an Ethereum contract call with the specified data as the
	// input.
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

func TestSendEth(t *testing.T) {
	eth, _ := setupThereum(t)
	acc := eth.Accounts["root"]
	acc.SendETH()
}

// func TestDeploy(t *testing.T) {

// }

// func TestContract(t *testing.T) {

// }

// func TestLog(t *testing.T) {

// }
