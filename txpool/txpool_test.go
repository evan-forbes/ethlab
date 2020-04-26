package txpool

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/matryer/is"
)

func placeSetup() []*txID {
	var s []*txID
	for i := 0; i < 20; i = i + 2 {
		s = append(s, &txID{gasPrice: new(big.Int).SetInt64(int64(i))})
	}
	return s
}

func TestInsert(t *testing.T) {
	pool := New()
	txBytes, err := hex.DecodeString("f85d01808094a52a2b202f7fc08fff1cb7d7bfab7f2248780b17018025a00208b2f26400f6147f6707ebd1af94b6b234cfa7bbbab34d12edcbe933a1cca5a0747c67a0d736882a02d9ecae83a0c10f3acc2ef5136fe4f4e2e4a4f8d1c184d5")
	if err != nil {
		t.Error(err)
	}
	var tx types.Transaction
	err = rlp.DecodeBytes(txBytes, &tx)
	if err != nil {
		t.Error(err)
	}
	pool.Insert(common.Address{}, &tx)
	regurge, has := pool.Next()
	if !has {
		t.Error("tx not inserted")
	}
	if regurge == nil {
		t.Error("tx not returned")
	}
}

func TestPlace(t *testing.T) {
	is := is.New(t)
	// test case 1: insert 11
	s := placeSetup()
	n := &txID{gasPrice: new(big.Int).SetInt64(int64(11))}
	place(s, n, 0, len(s)-1)
	is.Equal(n.gasPrice.String(), s[6].gasPrice.String())

	// test case 1: insert 3
	s = placeSetup()
	n = &txID{gasPrice: new(big.Int).SetInt64(int64(3))}
	place(s, n, 0, len(s)-1)
	is.Equal(n.gasPrice.String(), s[2].gasPrice.String())

	// test case 2: insert 0
	s = placeSetup()
	n = &txID{gasPrice: new(big.Int).SetInt64(int64(0))}
	place(s, n, 0, len(s)-1)
	is.Equal(n.gasPrice.String(), s[1].gasPrice.String())
	// is.Equal(n.gasPrice.String(), s[6].gasPrice.String())
}

func genTxs(author common.Address, nonce uint64) []*types.Transaction {
	var out []*types.Transaction
	return out
}

func TestRemove(t *testing.T) {
	is := is.New(t)
	// test case 1: insert 11
	s := placeSetup()
	n := &txID{gasPrice: new(big.Int).SetInt64(int64(2))}
	place(s, n, 0, len(s)-1)
	is.Equal(n.gasPrice.String(), s[6].gasPrice.String())

	// test case 2: insert 0
	s = placeSetup()
	n = &txID{gasPrice: new(big.Int).SetInt64(int64(0))}
	place(s, n, 0, len(s)-1)
	is.Equal(n.gasPrice.String(), s[1].gasPrice.String())
	// is.Equal(n.gasPrice.String(), s[6].gasPrice.String())
}
