package thereum

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/evan-forbes/ethlab/cmd"
)

func TestBoot(t *testing.T) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()
	heads := make(chan *types.Header)

	eth, err := New(defaultConfig(), common.Address{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("made thereum")
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
