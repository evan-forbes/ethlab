package txpool

import (
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// slight redesign of pool, using a channel instead of next
//  - return a nil transaction when there are no longer any transactions
//  - handle any invalid transactions here, not in thereum
//  - let thereum handle batching, cause that makes more sense.
//  - thereum should close the transaction feed? yes.
//  	so pass the feed to

// txSet wraps around multiple types.Trasnactions to help simulate running multiple
// transactions at once without having to make a ds proxy contract.
type txSet struct {
	Transactions []*types.Transaction
	ID           *txID
}

// LinkedPool is an ordered pool of transactions sorted by gas price. It also allows for
// 'linked' transactions
type LinkedPool struct {
	pool         map[common.Address]map[uint64]txSet
	order        []*txID // maintain gas price order
	mu           sync.RWMutex
	invalidCount int // invalidCount keeps track of the number of replaced transactions
}

func NewLinkedPool() *LinkedPool {
	return &LinkedPool{
		pool: make(map[common.Address]map[uint64]txSet),
	}
}

// Next retrieves the highest priced transaction/set of transactions
func (pool *LinkedPool) Next() (txSet, bool) {
	if len(pool.order) == 0 {
		return txSet{}, false
	}
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// pop the highest gas price transaction off
	nextID := pool.order[len(pool.order)-1]
	pool.order[len(pool.order)-1] = nil
	pool.order = pool.order[:len(pool.order)-1]

	if !nextID.valid {
		// try again if the tx set has been marked
		return pool.Next()
	}

	// get the tx from the pool
	set, has := pool.pool[nextID.address][nextID.nonce]
	if !has {
		// if a tx has somehow been removed from the pool but not from the order
		return pool.Next()
	}

	// remove the transaction from the pool
	delete(pool.pool[nextID.address], nextID.nonce)

	return set, true
}

// Insert adds a set of transactions to the ordered pool. If multiple transactions are provided
// they are treated as 'linked'. (Linked transactions run individually one after another and will be sorted
// using the lowest gas price of all txs provided).
func (pool *LinkedPool) Insert(author common.Address, txs ...*types.Transaction) {
	// don't insert nothing
	if len(txs) == 0 {
		return
	}
	// combine gas prices and limits for multiple txs
	gsprc := big.NewInt(1)
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
	fmt.Println("gas price of:", gsprc, gslmt)
	id := &txID{
		address:  author,
		nonce:    nonce,
		gasPrice: gsprc,
		gasUsed:  gslmt,
		valid:    true,
	}
	// form set
	set := txSet{
		Transactions: txs,
		ID:           id,
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()
	// check to see if this transaction already exists
	_, has := pool.pool[author]
	if !has {
		pool.pool[author] = make(map[uint64]txSet)
	}
	oldtx, has := pool.pool[author][nonce]
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
			clean(pool.order)
			pool.invalidCount = 0
		}
	}

	//// add the transaction in the pool ////
	pool.pool[author][nonce] = set

	// don't attempt to search and insert the txID if there're none to search
	if len(pool.order) == 0 {
		pool.order = append(pool.order, id)
		return
	}

	if len(pool.order) == 1 {

	}

	// put the transaction into the ordered set
	search(pool.order, gsprc)
	return
}

// Batch will get the maximum transactions from a linked pool for the provided gas limit
func Batch(gasLimit uint64, pool *LinkedPool) []*types.Transaction {
	var gasCount uint64
	var out []*types.Transaction
	for {
		if gasCount >= gasLimit {
			break
		}
		set, has := pool.Next()
		if !has {
			break
		}
		gasCount = gasCount + set.ID.gasUsed
		for _, tx := range set.Transactions {
			out = append(out, tx)
		}
	}
	return out
}

func search(order []*txID, price *big.Int) (n int) {
	sfunc := func(i int) bool {
		return order[i].gasPrice.Cmp(price) >= 0
	}
	return sort.Search(len(order), sfunc)
}

// RECYCLE BIN

//
// forcing the user to use goroutine to fill a channel just didn't feel right in this context.
//
// func (pool *LinkedPool) Feed(ctx context.Context, feed chan<- *types.Transaction) {
// 	defer close(feed)
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			txs, has := pool.next()
// 			if !has {
// 				// send a nil to signal the pool is empty
// 				feed <- nil
// 			}
// 			for _, tx := range txs {
// 				if tx == nil {
// 					continue
// 				}
// 				feed <- tx
// 			}
// 		}
// 	}
// }
