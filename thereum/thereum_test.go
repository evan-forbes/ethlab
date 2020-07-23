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
	root, err := NewAccount("root", big.NewInt(5))
	if err != nil {
		t.Error(err)
	}

	eth, err := New(DefaultConfig(), root)
	if err != nil {
		t.Error(err)
	}
	mngr.WG.Add(1)
	go eth.Run(mngr.Ctx, mngr.WG)
	return eth, mngr
}

func TestBoot(t *testing.T) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()

	heads := make(chan *types.Header)
	root, err := NewAccount("root", big.NewInt(100))

	eth, err := New(DefaultConfig(), root)
	if err != nil {
		t.Error(err)
	}
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
	<-mngr.Done()
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

func genTxs(count int, accs Accounts) ([]*types.Transaction, common.Address, error) {
	sinkAccout, _ := NewAccount("sink", big.NewInt(0))
	var i int
	var out []*types.Transaction
	for _, acc := range accs {
		if i >= count {
			break
		}
		if acc == nil {
			fmt.Println("nil account?")
			continue
		}
		tx, err := acc.CreateSend(sinkAccout.Address, big.NewInt(1))
		if err != nil {
			return out, sinkAccout.Address, err
		}
		out = append(out, tx)
		i++
	}
	return out, sinkAccout.Address, nil
}

func TestCreateSend(t *testing.T) {
	acc, err := NewAccount("test", big.NewInt(4))
	if err != nil {
		t.Error(err)
	}
	sink, err := NewAccount("sink", big.NewInt(1))
	if err != nil {
		t.Error(err)
	}

	tx, err := acc.CreateSend(sink.Address, big.NewInt(2))
	if err != nil {
		t.Error(err)
	}

	j, err := tx.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(j))

}

func TestSubscribeLogs(t *testing.T) {
	eth, mngr := setupThereum(t)
	
	// deploy a contract to deploy logs
	eth.
}
