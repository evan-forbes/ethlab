package txpool

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// slight redesign of pool, using a channel instead of next
//  - return a nil transaction when there are no longer any transactions
//  - handle any invalid transactions here, not in thereum
//  - let thereum handle batching, cause that makes more sense.
//  - thereum should close the transaction feed? yes.
//  	so pass the feed to

// TxSet wraps around multiple types.Trasnactions to help simulate running multiple
// transactions at once without having to make a ds proxy contract.
type TxSet struct {
	Transactions []*types.Transaction
	ID           *txID
}

// LinkedPool is an ordered pool of transactions sorted by gas price. It also allows for
// 'linked' transactions
type LinkedPool struct {
	ctx          context.Context
	Pool         map[common.Address]map[uint64]TxSet
	Order        []*txID // Order maintains
	mu           sync.RWMutex
	invalidCount int // invalidCount keeps track of the number of replaced transactions
}

func NewLinkedPool(ctx context.Context) *LinkedPool {
	return &LinkedPool{
		ctx:  ctx,
		Pool: make(map[common.Address]map[uint64]TxSet),
	}
}

// next retrieves the highest priced transaction/set of transactions
func (pool *LinkedPool) next() ([]*types.Transaction, bool) {
	if len(pool.Order) == 0 {
		return []*types.Transaction{}, false
	}
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// pop the highest gas price transaction off
	nextID := pool.Order[len(pool.Order)-1]
	pool.Order[len(pool.Order)-1] = nil
	pool.Order = pool.Order[:len(pool.Order)-1]

	if !nextID.valid {
		// try again if the tx set has been marked
		return pool.next()
	}

	// get the tx from the pool
	tx, has := pool.Pool[nextID.address][nextID.nonce]
	if !has {
		// if a tx has somehow been removed from the pool but not from the order
		return nil, true
	}

	// remove the transaction from the pool
	delete(pool.Pool[nextID.address], nextID.nonce)

	return tx.Transactions, true
}

// Insert adds a set of transactions to the ordered pool. If multiple transactions are provided
// they are treated as 'linked'. (linked txs simulate using a DSProxy contract, running multiple
// transaction in a single block) Linked txs will use the lowest gas price of all txs provided.
func (pool *LinkedPool) Insert(author common.Address, txs ...*types.Transaction) {
	// don't insert nothing
	if len(txs) == 0 {
		return
	}
	// combine gas prices and limits for multiple txs
	var gsprc *big.Int
	var gslmt uint64
	for _, tx := range txs {
		// use the lowest gas price of all the transactions
		if gsprc.Cmp(tx.GasPrice()) < 0 {
			gsprc = tx.GasPrice()
		}
		// add up the total gas limit of all transactions
		gslmt = gslmt + tx.Gas()
	}
	// use nonce of the first tx.
	nonce := txs[0].Nonce()
	// form id
	id := &txID{
		address:  author,
		nonce:    nonce,
		gasPrice: gsprc,
		gasUsed:  gslmt,
		valid:    true,
	}
	// form set
	set := TxSet{
		Transactions: txs,
		ID:           id,
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()
	// check to see if this transaction already exists
	_, has := pool.Pool[author]
	if !has {
		pool.Pool[author] = make(map[uint64]TxSet)
	}
	oldtx, has := pool.Pool[author][nonce]
	if has {
		// if the gas price is not larger, don't do anything
		if oldtx.ID.gasPrice.Cmp(gsprc) != 1 {
			return
		}
		// mark the old transaction as invalid
		oldtx.ID.valid = false

		// count invalid txs to make sure we don't call clean too often
		pool.invalidCount++

		if pool.invalidCount > 100 {
			clean(pool.Order)
			pool.invalidCount = 0
		}
	}

	//// add the transaction in the pool ////
	pool.Pool[author][nonce] = set

	// don't attempt to search and insert the txID if there're none to search
	if len(pool.Order) == 0 {
		pool.Order = append(pool.Order, id)
		return
	}

	// put the transaction into the ordered set
	place(pool.Order, id, 0, len(pool.Order)-1)
	return
}

func (pool *LinkedPool) Feed(ctx context.Context, feed chan<- *types.Transaction) {
	defer close(feed)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			txs, has := pool.next()
			if !has {
				// send a nil to signal the pool is empty
				feed <- nil
			}
			for _, tx := range txs {
				if tx == nil {
					continue
				}
				feed <- tx
			}
		}
	}
}

// I need to signal to batch that there are no more transactions
// rather than set a time limit
// use nil to signal that the pool is empty

func Batch(gasLimit uint64, feed <-chan *types.Transaction) []*types.Transaction {
	var gasCount uint64
	var out []*types.Transaction
	for {
		if gasCount >= gasLimit {
			break
		}
		select {
		case tx := <-feed:
			out = append(out, tx)
		case <-time.After(wait):
			break
		}
	}

}

// setup feed
// custom txSet with multiple txs
// accept multiple txs
