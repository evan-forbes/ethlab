package thereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/evan-forbes/ethlab/txpool"
)

// Thereum contains and controls the processes needed to run a single node
// PoA ethereum blockchain.
type Thereum struct {
	ctx      context.Context
	wg       *sync.WaitGroup
	root     *Account
	txPool   *txpool.LinkedPool
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

	Events   *filters.EventSystem // Event system for filtering logs and events
	Accounts Accounts             // access to initial accounts specified in config.Allocations

	chainConfig *params.ChainConfig
}

// New using a config and root signing address to make a new Thereum blockchain
func New(config Config, root *Account) (*Thereum, error) {
	// init the configured db
	db := config.DB()
	// delay := config.Delayer()

	// init the genesis block + any accounts designated in config.Allocaiton
	genesis, accounts, err := config.Genesis()
	if err != nil {
		return nil, err
	}

	genBlock := genesis.MustCommit(db)

	if root == nil {
		root, _ = NewAccount("defaultRoot", big.NewInt(100))
	}
	chainConfig := params.AllEthashProtocolChanges
	chainConfig.ChainID = big.NewInt(1)
	bc, _ := core.NewBlockChain(db, nil, chainConfig, ethash.NewFaker(), vm.Config{}, nil)
	for _, acc := range accounts {
		fmt.Printf("%s\t\t%s\t%s\n", acc.Name, acc.Address.Hex(), acc.Balance.String())
	}
	t := &Thereum{
		txPool:     txpool.NewLinkedPool(),
		database:   db,
		blockchain: bc,
		signer:     types.NewEIP155Signer(big.NewInt(1)),
		root:       root,
		gasLimit:   config.GenesisConfig.GasLimit, // TODO: config and make more flexible
		delay:      config.Delay,
		Events:     filters.NewEventSystem(&filterBackend{db: db, bc: bc}, false),
		Accounts:   accounts,
	}
	t.latestBlock = genBlock
	t.chainConfig = chainConfig
	return t, nil
}

////////////////////////////////////
// 		Growing the Chain
//////////////////////////////////

// Run starts issuing new blocks using transactions in the transaction pool
func (t *Thereum) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer t.Shutdown(wg)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			t.Commit()
			// fmt.Println("block hash:", t.latestBlock.Hash().Hex())
			time.Sleep(time.Millisecond * time.Duration(t.delay))
		}
	}
}

// Commit creates a new block using existing transaction from the txpool
func (t *Thereum) Commit() {
	// TODO: 1)this is fugly 2) add custom delay 3) add ability to pause
	// create a new block using existing transaction in the pool
	block, state := t.nextBlock()
	t.mu.Lock()
	t.latestBlock = block
	t.latestState = state
	t.mu.Unlock()
	// add optional delay before adding block to simulate pending state
	t.appendBlock(block)

}

// nextBlock mints a new block, filling it with transactions from the transaction pool
func (t *Thereum) nextBlock() (*types.Block, *state.StateDB) {
	// make new blocks using the transaction pool
	t.mu.Lock()
	defer t.mu.Unlock()
	blocks, _ := core.GenerateChain(
		t.chainConfig,
		t.blockchain.CurrentBlock(),
		ethash.NewFaker(),
		t.database,
		1,
		func(i int, b *core.BlockGen) {
			b.SetCoinbase(t.root.Address)
			// get the next set of highest paying transactions
			txs := txpool.Batch(t.gasLimit, t.txPool)
			// add them to the new block.
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
	t.mu.Lock()
	defer t.mu.Unlock()
	_, err := t.blockchain.InsertChain([]*types.Block{block})
	// TODO: get rid of panic and handle the errors
	if err != nil {
		panic(err)
	}
	return
}

////////////////////////////////////
// 	Transaction Reading and Pooling
//////////////////////////////////

// AddTx validates and inserts the transaction into the txpool
func (t *Thereum) AddTx(tx *types.Transaction) error {
	// validate tx
	from, err := t.validateTx(tx)
	if err != nil {
		return fmt.Errorf("could not validate transaction: %s", err)
	}
	t.txPool.Insert(from, tx)
	fmt.Println("added tx", tx.Hash().Hex())
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

// TxReceipt returns the receipt, if any, from a mined transaction's hash
func (t *Thereum) TxReceipt(hash common.Hash) (*types.Receipt, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	receipt, _, _, _ := rawdb.ReadReceipt(t.database, hash, t.chainConfig)
	return receipt, nil
}

// LatestState returns the latest state
func (t *Thereum) LatestState() *state.StateDB {
	t.mu.Lock()
	state := t.latestState
	t.mu.Unlock()
	return state
}

// LatestBlock returns the latest block. Not guarenteed to be final
func (t *Thereum) LatestBlock() *types.Block {
	t.mu.Lock()
	block := t.latestBlock
	t.mu.Unlock()
	return block
}

// Shutdown begins the procedure to stop the Thereum blockchain
func (t *Thereum) Shutdown(wg *sync.WaitGroup) {
	defer wg.Done()
	t.blockchain.Stop()
}

// TransactionCount returns the number of transactions in a given block
func (t *Thereum) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if blockHash == t.latestBlock.Hash() {
		t.latestBlock.Transactions()
		return uint(t.latestBlock.Transactions().Len()), nil
	}

	block := t.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return uint(0), errors.New("block does not exist")
	}

	return uint(block.Transactions().Len()), nil
}

// TransactionCountByAddress returns the number of transactions sent by an address at a given block
func (t *Thereum) TransactionCountByAddress(ctx context.Context, addr common.Address, blockHash common.Hash) (*hexutil.Uint64, error) {
	state, err := t.blockchain.StateAt(blockHash)
	if err != nil {
		return nil, err
	}
	count := state.GetNonce(addr)
	return (*hexutil.Uint64)(&count), state.Error()
}

// stateByBlockNumber retrieves a state by a given blocknumber.
func (t *Thereum) stateByBlockNumber(ctx context.Context, blockNumber *big.Int) (*state.StateDB, error) {
	if blockNumber == nil || blockNumber.Cmp(t.blockchain.CurrentBlock().Number()) == 0 {
		return t.blockchain.State()
	}
	block, err := t.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	return t.blockchain.StateAt(block.Hash())
}

// BlockByNumber retrieves a block from the database by number, caching it
// (associated with its hash) if found.
func (t *Thereum) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if number == nil || number.Cmp(t.latestBlock.Number()) == 0 {
		return t.blockchain.CurrentBlock(), nil
	}

	block := t.blockchain.GetBlockByNumber(uint64(number.Int64()))
	if block == nil {
		return nil, errors.New("block does not exist")
	}

	return block, nil
}

// BalanceAt returns the wei balance of a certain account in the blockchain.
func (t *Thereum) BalanceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (*big.Int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	statedb, err := t.stateByBlockNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	return statedb.GetBalance(contract), nil
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

// This will only be needed if we want to support using as a native client
// // callContract implements common code between normal and pending contract calls.
// // state is modified during execution, make sure to copy it if necessary.
// func (t *Thereum) callContract(ctx context.Context, call ethereum.CallMsg, block *types.Block, statedb *state.StateDB) ([]byte, uint64, bool, error) {
// 	// Ensure message is initialized properly.
// 	if call.GasPrice == nil {
// 		call.GasPrice = big.NewInt(1)
// 	}
// 	if call.Gas == 0 {
// 		call.Gas = 50000000
// 	}
// 	if call.Value == nil {
// 		call.Value = new(big.Int)
// 	}
// 	// Set infinite balance to the fake caller account.
// 	// from := statedb.GetOrNewStateObject(call.From)
// 	// from.SetBalance(math.MaxBig256)
// 	// Execute the call.
// 	msg := callmsg{call}

// 	evmContext := core.NewEVMContext(msg, block.Header(), t.blockchain, nil)
// 	// Create a new environment which holds all relevant information
// 	// about the transaction and calling mechanisms.
// 	vmenv := vm.NewEVM(evmContext, statedb, params.AllEthashProtocolChanges, vm.Config{})
// 	gaspool := new(core.GasPool).AddGas(math.MaxUint64)

// 	return core.NewStateTransition(vmenv, msg, gaspool).TransitionDb()
// }

// // callmsg implements core.Message to allow passing it as a transaction simulator.
// type callmsg struct {
// 	ethereum.CallMsg
// }

// func (m callmsg) From() common.Address { return m.CallMsg.From }
// func (m callmsg) Nonce() uint64        { return 0 }
// func (m callmsg) CheckNonce() bool     { return false }
// func (m callmsg) To() *common.Address  { return m.CallMsg.To }
// func (m callmsg) GasPrice() *big.Int   { return m.CallMsg.GasPrice }
// func (m callmsg) Gas() uint64          { return m.CallMsg.Gas }
// func (m callmsg) Value() *big.Int      { return m.CallMsg.Value }
// func (m callmsg) Data() []byte         { return m.CallMsg.Data }
