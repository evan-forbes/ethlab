package txpool

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Pooler descibes the methods expected by Thereum to interact with a pool of transactions.
// Note: Pooler is currently unused, but could replace TxPool
type Pooler interface {
	// Batch returns some number of txs whose total gas cost is less than or equal to limit
	Batch(limit uint64) []*types.Transaction
	// Insert adds the transaction to the pool
	Insert(author common.Address, tx *types.Transaction)
}

// TxPool maintains a queue of transactions sorted by gas price
type TxPool struct {
	Pool         map[common.Address]map[uint64]setTx
	Order        []*txID // Order maintains
	mu           sync.RWMutex
	invalidCount int // invalidCount keeps track of the number of replaced transactions
	maxSize      int // maxSize specs the max number of txs in the pool
}

type setTx struct {
	tx   *types.Transaction
	txID *txID
}

type txID struct {
	address  common.Address
	nonce    uint64
	gasPrice *big.Int
	gasUsed  uint64
	valid    bool
}

// New inits a new TxPool
func New() *TxPool {
	return &TxPool{
		Pool:    make(map[common.Address]map[uint64]setTx),
		maxSize: 100000000000,
	}
}

// Next pops and returns the highest gas price transaction, along with
// a bool determining if there are any transactions left. nil txs are
// returned for non-existant and invalid transaction that have yet to
// be cleaned.
func (pool *TxPool) Next() (*types.Transaction, bool) {
	if len(pool.Order) == 0 {
		return nil, false
	}
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// pop the highest gas price transaction off
	nextID := pool.Order[len(pool.Order)-1]
	pool.Order[len(pool.Order)-1] = nil
	pool.Order = pool.Order[:len(pool.Order)-1]

	if !nextID.valid {
		return nil, true
	}

	// get the tx from the pool
	tx, has := pool.Pool[nextID.address][nextID.nonce]
	if !has {
		// if a tx has somehow been removed from the pool but not from the order
		return nil, true
	}

	// remove the transaction from the pool
	delete(pool.Pool[nextID.address], nextID.nonce)

	return tx.tx, true
}

// Insert adds a transaction to the pool, replaceing the old transaction if the nonce is the same.
// Transactions should be verified before insertion.
func (pool *TxPool) Insert(author common.Address, tx *types.Transaction) {
	nonce := tx.Nonce()
	id := &txID{address: author, nonce: nonce, gasPrice: tx.GasPrice(), gasUsed: tx.Gas(), valid: true}

	pool.mu.Lock()
	defer pool.mu.Unlock()
	// check to see if this transaction already exists
	_, has := pool.Pool[author]
	if !has {
		pool.Pool[author] = make(map[uint64]setTx)
	}
	oldtx, has := pool.Pool[author][nonce]
	if has {
		// if the gas price is not larger, don't do anything
		if oldtx.tx.GasPrice().Cmp(tx.GasPrice()) != 1 {
			return
		}
		// mark the old transaction as invalid
		oldtx.txID.valid = false

		// count invalid txs to make sure we don't call clean too often
		pool.invalidCount++

		if pool.invalidCount > 100 {
			clean(pool.Order)
			pool.invalidCount = 0
		}
	}

	//// add the transaction in the pool ////
	pool.Pool[author][nonce] = setTx{tx: tx, txID: id}

	// don't attempt to search and insert the txID if there're none to search
	if len(pool.Order) == 0 {
		pool.Order = append(pool.Order, id)
		return
	}

	// put the transaction into the ordered set
	place(pool.Order, id, 0, len(pool.Order)-1)
	return
}

// Batch returns the next set of transactions based on the provided limit
func (pool *TxPool) Batch(limit uint64) []*types.Transaction {
	var out []*types.Transaction
	var used uint64
	for {
		if len(pool.Order) == 0 {
			break
		}
		id := pool.Order[len(pool.Order)-1]
		used = used + id.gasUsed
		// ensure that the limit is not reached
		if used > limit {
			break
		}
		// get the next transactions
		tx, ok := pool.Next()
		if !ok {
			break
		}
		if tx == nil {
			continue
		}
		out = append(out, tx)
	}
	return out
}

// cull removes the lowest price txs from the pool
func (pool *TxPool) cull() {

}

// place inserts a transaction id in order, sorted by gas price
func place(s []*txID, n *txID, head, tail int) {
	diff := tail - head
	// if our search has finished
	if diff == 1 {
		// insert in front of head
		s = append(s, nil)
		copy(s[head+2:], s[head+1:])
		s[head+1] = n
		return
	}
	delta := head + (diff / 2)
	node := s[delta]
	switch node.gasPrice.Cmp(n.gasPrice) {
	case 0:
		s = append(s, nil)
		copy(s[delta+1:], s[delta:])
		s[delta] = n
	case -1:
		// try againg above delta
		place(s, n, delta, tail)
	case 1:
		// try again below delta
		place(s, n, head, delta)
	}
	return
}

func clean(s []*txID) {
	for i, id := range s {
		if !id.valid {
			sliceDelete(s, i)
		}
	}
}

func sameTxID(a, b *txID) bool {
	return a.address == b.address && a.nonce == b.nonce
}

func sliceDelete(a []*txID, i int) {
	if i < len(a)-1 {
		copy(a[i:], a[i+1:])
	}
	a[len(a)-1] = nil
	a = a[:len(a)-1]
}

// RECYCLE

// // remove deletes a transaction id from the ordered slice
// // its kind of slow, so try not to use too much.
// func remove(s []*txID, n *txID, head, tail int) {
// 	diff := tail - head
// 	// if our search has finished
// 	if diff == 1 {
// 		// check to see if it is the same txID
// 		if n.address == s[head].address && n.nonce == s[head].nonce {
// 			// remove head
// 			sliceDelete(s, head)
// 			return
// 		}
// 		if n.address == s[tail].address && n.nonce == s[tail].nonce {
// 			// remove tail
// 			sliceDelete(s, tail)
// 			return
// 		}
// 	}
// 	delta := head + (diff / 2)
// 	node := s[delta]
// 	switch node.gasPrice.Cmp(n.gasPrice) {
// 	case 0:
// 		// cross our fingers and hope this doesn't happen when the txpool filled with many txs of the same price...
// 		// check all txs with the potentially same gas price
// 		for i := head; i < tail; i++ {
// 			if sameTxID(s[i], n) {
// 				sliceDelete(s, i)
// 			}
// 		}
// 	case -1:
// 		// try again above delta
// 		remove(s, n, delta, tail)
// 	case 1:
// 		// try again below delta
// 		remove(s, n, head, delta)
// 	}
// 	return
// }
