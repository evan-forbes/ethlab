package thereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

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
	ctx      context.Context
	wg       *sync.WaitGroup
	root     common.Address
	TxPool   *txpool.TxPool
	gasLimit uint64
	// gasLimit GasLimiter
	delay      uint
	signer     types.Signer
	database   ethdb.Database   // In memory database to store our testing data
	blockchain *core.BlockChain // Ethereum blockchain to handle the consensus

	mu sync.Mutex

	// use the locked wrapper methods to access these!
	// no one like global state, but I also don't appreciate how they're returned from the ethereum data structure, blockchain
	latestBlock *types.Block   // pending block
	latestState *state.StateDB // pending state

	events *filters.EventSystem // Event system for filtering log events live

	config *params.ChainConfig
}

// New using a config and root signing address to make a new Thereum blockchain
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
		signer:     types.NewEIP155Signer(big.NewInt(1)),
		root:       root,
		gasLimit:   config.GenesisConfig.GasLimit, // TODO: config and make more flexible
		delay:      config.Delay,
	}, nil
}

// Run starts issuing new blocks using transactions in the transaction pool
func (t *Thereum) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer t.Shutdown(wg)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			t.Commit()
			time.Sleep(time.Millisecond * time.Duration(t.delay))
		}
	}
}

// Commit creates a new block using existing transaction from the TxPool
func (t *Thereum) Commit() {
	// TODO: 1)this is fugly 2) add custom delay 3) add ability to pause
	// create a new block using existing transaction in the pool
	t.mu.Lock()
	defer t.mu.Unlock()
	block, state := t.nextBlock()
	t.latestBlock = block
	t.latestState = state
	// add optional delay before adding block to simulate pending state
	t.appendBlock(block)
}

// nextBlock mints a new block, filling it with transactions from the transaction pool
func (t *Thereum) nextBlock() (*types.Block, *state.StateDB) {
	blocks, _ := core.GenerateChain(
		t.config,
		t.blockchain.CurrentBlock(),
		ethash.NewFaker(),
		t.database,
		1,
		func(i int, b *core.BlockGen) {
			b.SetCoinbase(t.root)
			txs := t.TxPool.Batch(t.gasLimit)
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

func (t *Thereum) appendBlock(block *types.Block) {
	_, err := t.blockchain.InsertChain([]*types.Block{block})
	// TODO: get rid of panic and handle the errors
	if err != nil {
		panic(err)
	}
	return
}

// AddTx validates and inserts the transaction into the TxPool
func (t *Thereum) AddTx(tx *types.Transaction) error {
	// validate tx
	from, err := t.validateTx(tx)
	if err != nil {
		return fmt.Errorf("could not validate transaction: %s", err)
	}
	t.TxPool.Insert(from, tx)
	return nil
}

// validateTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits (price and size).
func (t *Thereum) validateTx(tx *types.Transaction) (common.Address, error) {
	// Reject transactions over defined size to prevent DOS attacks
	if uint64(tx.Size()) > txMaxSize {
		return common.Address{}, errors.New("invalid transaction: too large")
	}
	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value().Sign() < 0 {
		return common.Address{}, errors.New("invalid transaction: negative value")
	}
	// Ensure the transaction doesn't exceed the current block limit gas.
	if t.blockchain.GasLimit() < tx.Gas() {
		return common.Address{}, errors.New("invalid transaction: gas limit broken")
	}
	// Make sure the transaction is signed properly
	from, err := types.Sender(t.signer, tx)
	if err != nil {
		return from, errors.New("invalid transaction: signature could not be verified")
	}
	state := t.LatestState()
	// Ensure the transaction adheres to nonce ordering
	if state.GetNonce(from) > tx.Nonce() {
		return from, errors.New("invalid transaction: nonce too low")
	}
	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	if state.GetBalance(from).Cmp(tx.Cost()) < 0 {
		return from, errors.New("invalid transaction: not enough funds")
	}
	// Ensure the transaction has more gas than the basic tx fee.
	intrGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, true, true)
	if err != nil {
		return from, err
	}
	if tx.Gas() < intrGas {
		return from, errors.New("invalid transaction: not enough gas to cover intrinsic transaction function")
	}
	return from, nil
}

func (t *Thereum) LatestState() *state.StateDB {
	t.mu.Lock()
	state := t.latestState
	t.mu.Unlock()
	return state
}

func (t *Thereum) LatestBlock() *types.Block {
	t.mu.Lock()
	block := t.latestBlock
	t.mu.Unlock()
	return block
}

func (t *Thereum) Shutdown(wg *sync.WaitGroup) {
	defer wg.Done()
	t.blockchain.Stop()
}

const (
	// txSlotSize is used to calculate how many data slots a single transaction
	// takes up based on its size. The slots are used as DoS protection, ensuring
	// that validating a new transaction remains a constant operation (in reality
	// O(maxslots), where max slots are 4 currently).
	txSlotSize = 32 * 1024

	// txMaxSize is the maximum size a single transaction can have. This field has
	// non-trivial consequences: larger transactions are significantly harder and
	// more expensive to propagate; larger transactions also take more resources
	// to validate whether they fit into the pool or not.
	txMaxSize = 2 * txSlotSize // 64KB, don't bump without EIP-2464 support
)
