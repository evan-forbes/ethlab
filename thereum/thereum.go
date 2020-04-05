package thereum

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/evan-forbes/ethlab/txpool"
)

//

// TODO: verify transactions before adding to the transactoin pool

/*
so far the appears to be
- incorporate txPool into something similar to the simulated backend.
- add auto run
*/

type Thereum struct {
	ctx     context.Context
	wg      *sync.WaitGroup
	root    common.Address
	TxPool  txpool.Pooler
	gasLimt *big.Int
	// gasLimit GasLimiter
	// delay    Delayer
	// Signer types.Signer
	database   ethdb.Database   // In memory database to store our testing data
	blockchain *core.BlockChain // Ethereum blockchain to handle the consensus

	mu sync.Mutex

	// try and see if I can get away without using this shit vvv
	pendingBlock *types.Block   // Currently pending block that will be imported on request
	pendingState *state.StateDB // Currently pending state that will be the active on on request

	events *filters.EventSystem // Event system for filtering log events live

	config *params.ChainConfig
}

func New(config *Config, root common.Address) (*Thereum, error) {
	// init the configured db
	db := config.DB()
	// delay := config.Delayer()

	// init the genesis block + any accounts designated in config.Allocaiton
	genesis, accounts, err := config.Genesis()
	if err != nil {
		return nil, err
	}
	genesis.MustCommit(db)
	bc, _ := core.NewBlockChain(db, nil, genesis.Config, ethash.NewFaker(), vm.Config{}, nil)
	for _, acc := range accounts {
		fmt.Printf("%s\t\t%s\t%s\n", acc.Name, acc.Address, acc.Balance)
	}
	return &Thereum{
		TxPool:     txpool.New(),
		database:   db,
		blockchain: bc,
		root:       root,
		gasLimt:    big.NewInt(10485760), // TODO: config and make more flexible
	}, nil
}

func (t *Thereum) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer t.Shutdown(wg)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// issue a new pending block
			// t.Delay()
			t.mu.Lock()
			block, state := t.NewPendingBlock()

			t.Commit(block)
			t.mu.Unlock()

		}
	}
}

// BatchTxs will add the max number of transaction to the provided block.
func (t *Thereum) BatchTxs() []*types.Transaction {
	// fetch the number a transactions, until all of the gas is used
	var txs []*types.Transaction
	for gas := new(big.Int); gas.Cmp(t.gasLimt); {
	}
	return txs
}

// NewPendingBlock mints a new block, filling it with transactions from the transaction pool
func (t *Thereum) NewPendingBlock() (*types.Block, *state.StateDB) {
	blocks, _ := core.GenerateChain(
		t.config,
		t.blockchain.CurrentBlock(),
		ethash.NewFaker(),
		t.database,
		1,
		func(i int, b *core.BlockGen) {
			b.SetCoinbase(t.root)
			txs := t.BatchTxs()
			for _, tx := range txs {
				b.AddTxWithChain(t.blockchain, tx)
			}
		},
	)
	statedb, _ := t.blockchain.State()

	freshBlock := blocks[0]
	freshState, _ := state.New(freshBlock.Root(), statedb.Database())
	return freshBlock, freshState
}

func (t *Thereum) Commit(block *types.Block) {
	_, err := t.blockchain.InsertChain([]*types.Block{block})
	if err != nil {
		panic(err)
	}
	return
}

func (t *Thereum) Start() {

}

func (t *Thereum) Shutdown(wg *sync.WaitGroup) {
	defer wg.Done()
}

// // validateTx checks whether a transaction is valid according to the consensus
// // rules and adheres to some heuristic limits of the local node (price and size).
// func (t *Thereum) validateTx(tx *types.Transaction, local bool) error {
// 	// Reject transactions over defined size to prevent DOS attacks
// 	if uint64(tx.Size()) > txMaxSize {
// 		return ErrOversizedData
// 	}
// 	// Transactions can't be negative. This may never happen using RLP decoded
// 	// transactions but may occur if you create a transaction using the RPC.
// 	if tx.Value().Sign() < 0 {
// 		return ErrNegativeValue
// 	}
// 	// Ensure the transaction doesn't exceed the current block limit gas.
// 	if t.currentMaxGas < tx.Gas() {
// 		return ErrGasLimit
// 	}
// 	// Make sure the transaction is signed properly
// 	from, err := types.Sender(pool.signer, tx)
// 	if err != nil {
// 		return ErrInvalidSender
// 	}
// 	// Drop non-local transactions under our own minimal accepted gas price
// 	local = local || pool.locals.contains(from) // account may be local even if the transaction arrived from the network
// 	if !local && pool.gasPrice.Cmp(tx.GasPrice()) > 0 {
// 		return ErrUnderpriced
// 	}
// 	// Ensure the transaction adheres to nonce ordering
// 	if pool.currentState.GetNonce(from) > tx.Nonce() {
// 		return ErrNonceTooLow
// 	}
// 	// Transactor should have enough funds to cover the costs
// 	// cost == V + GP * GL
// 	if pool.currentState.GetBalance(from).Cmp(tx.Cost()) < 0 {
// 		return ErrInsufficientFunds
// 	}
// 	// Ensure the transaction has more gas than the basic tx fee.
// 	intrGas, err := IntrinsicGas(tx.Data(), tx.To() == nil, true, pool.istanbul)
// 	if err != nil {
// 		return err
// 	}
// 	if tx.Gas() < intrGas {
// 		return ErrIntrinsicGas
// 	}
// 	return nil
// }
