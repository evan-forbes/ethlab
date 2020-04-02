package txpool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Pooler descibes the methods expected by Thereum to interact with a pool of transactions.
type Pooler interface {
	Next() (*types.Transaction, bool)
	Insert(author common.Address, tx *types.Transaction) error
}

// TxPool is the most basic
type TxPool struct {
	Pool map[common.Address]map[uint64]*types.Transaction
	// SortedTransactions map[]
	Order []*txID
}

func New() *TxPool {
	return &TxPool{
		Pool: make(map[common.Address]map[uint64]*types.Transaction),
	}
}

// Next pops the highest gas price transaction
// TODO: add comments and test
func (pool *TxPool) Next() (*types.Transaction, bool) {
	if len(pool.Order) == 0 {
		return nil, false
	}
	nextID := pool.Order[len(pool.Order)-1]
	pool.Order[len(pool.Order)-1] = nil
	pool.Order = pool.Order[:len(pool.Order)-1
	tx, has := pool.Pool[nextID.address][nextID.nonce]
	return tx, has
}

// Insert adds a transaction to the pool, replaceing the old transaction if the nonce is the same.
// Transactions should be verified before insertion.
func (pool *TxPool) Insert(author common.Address, tx *types.Transaction) error {
	nonce := tx.Nonce()
	id := &txID{address: author, nonce: nonce, gasPrice: tx.GasPrice()}
	// check to see if this transaction already exists
	oldtx, has := pool.Pool[author][nonce]
	if has {
		// if the gas price is not larger, don't do anything
		if oldtx.GasPrice().Cmp(tx.GasPrice()) != 1 {
			return nil
		}
		// remove the old transaction from the ordered set before placing
		remove(pool.Order, id, 0, len(pool.Order)-1)
	}
	//// place the transaction in the pool ////
	pool.Pool[author][nonce] = tx
	id := &txID{address: author, nonce: nonce, gasPrice: tx.GasPrice()}

	// don't attempt to search and insert the txID if there none to search
	if len(pool.Order) == 0 {
		pool.Order = append(pool.Order, id)
		return nil
	}
	place(pool.Order, id, 0, len(pool.Order)-1)
	return nil
}

type txID struct {
	address  common.Address
	nonce    uint64
	gasPrice *big.Int
}

// place inserts a transaction id in order (sorted by gas price, using recursion)
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

// remove deletes a transaction id from the ordered slice
func remove(s []*txID, n *txID, head, tail int) {
	diff := tail - head
	// if our search has finished
	if diff == 1 {
		// check to see if it is the same txID
		if n.address == s[head].address && n.nonce == s[head].nonce {
			// remove head
			sliceDelete(s, head)
			return
		}
		if n.address == s[tail].address && n.nonce == s[tail].nonce {
			// remove tail
			sliceDelete(s, tail)
			return
		}
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

func sliceDelete(a []interface, i int) {
	if i < len(a)-1 {
		copy(a[i:], a[i+1:])
	  }
	  a[len(a)-1] = nil
	  a = a[:len(a)-1]
}