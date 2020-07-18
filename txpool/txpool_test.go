package txpool

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/evan-forbes/ethlab/module"
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

func TestInsertMany(t *testing.T) {
	// make a bunch of transactions
	sender, err := module.NewUser()
	if err != nil {
		t.Error(err)
		return
	}
	recvr, err := module.NewUser()
	if err != nil {
		t.Error(err)
		return
	}
	// Make a bunch of txs (these don't need to be valid txs for this test)
	var txs []*types.Transaction
	for i := 0; i < 100; i++ {
		tx, err := sender.NewSend(recvr.NewTxOpts().From, big.NewInt(10000000000000), big.NewInt(1000000000+int64(i)), 300000)
		if err != nil {
			t.Error(err)
			return
		}
		txs = append(txs, tx)
	}
	pool := NewLinkedPool()
	for i := 0; i < len(txs); i++ {
		pool.Insert(sender.From, txs[i])
	}
	fmt.Println("pool size", pool.Len())
	for i, tx := range pool.Batch(10000000) {
		fmt.Printf("tx number %d:\n gas %d price %s \n", i, tx.Gas(), tx.GasPrice().String())
		fmt.Println(pool.Len())
	}
}

func TestPlace(t *testing.T) {
	is := is.New(t)
	// test case 1: insert 11
	s := placeSetup()
	n := &txID{gasPrice: new(big.Int).SetInt64(int64(11))}
	place(&s, n, 0, len(s)-1)
	is.Equal(n.gasPrice.String(), s[6].gasPrice.String())

	// test case 1: insert 3
	s = placeSetup()
	n = &txID{gasPrice: new(big.Int).SetInt64(int64(3))}
	place(&s, n, 0, len(s)-1)
	is.Equal(n.gasPrice.String(), s[2].gasPrice.String())

	// test case 2: insert 0
	s = placeSetup()
	n = &txID{gasPrice: new(big.Int).SetInt64(int64(0))}
	place(&s, n, 0, len(s)-1)
	is.Equal(n.gasPrice.String(), s[1].gasPrice.String())
	// is.Equal(n.gasPrice.String(), s[6].gasPrice.String())
}

func makeTransactions() []*types.Transaction {
	raw := []string{
		// "f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772",
		"f85f8064825208946df11c6a93177f0a351daee93dfaca5177345ab2018025a04b652f9ff38fc9c8c5347da03c9fea7a8baf5ffa0c2c206402a38c70b32b157fa00e5e5ebac190c751f824f5f12fa7324e6cf7aee1022b65496bb2724d5eb8d190",
		"f85f806f825213946df11c6a93177f0a351daee93dfaca5177345ab2018026a0c9807a2f71f0e41a45521ff9817e60b847d50cb4b7b694a378cc5378ba8b0b7ba0611acbfc055c0286dbfbfc27cd0811ae1d2371690a523424af419a698b2ef6c8",
		"f85f807a82521e946df11c6a93177f0a351daee93dfaca5177345ab2018025a066df31e66fc372c0f7dd761e6779fd73454c98e06fed2932f5e1560fa59d8a49a02e637f85e1c15584e2deca936c7cdf7ad0b06400fd063d53cc65da5a65871c15",
		"f860808185825229946df11c6a93177f0a351daee93dfaca5177345ab2018025a00d0eb76da0a820340bbd194ac6f544a64f9cc3ed0138dc2bd4a138b25fa6314ba03d5e8302f3dd1a8e0893f81dfa6f33ba9339168bba8c5655eb9118fa12773fd6",
		"f860808190825234946df11c6a93177f0a351daee93dfaca5177345ab2018026a0f3d2f3202df2bd2f26c16b49d718822ef28c4e1f33c50079faedbcdf5d721c7aa06b773432c51c16fa4f86fba102250cff0cfa7621d272b8cd16f42779ed5cda26",
		"f86080819b82523f946df11c6a93177f0a351daee93dfaca5177345ab2018026a0d784d15e469a15ffe57cad005151e339ad3588a8c14ca40dd8ebac94fb9aae33a04bb9b434b6a2ff4d4e786a88e98e03d8dbe1e2253dfe4117324359e3b1b28d60",
	}
	var out []*types.Transaction
	for _, s := range raw {

		txBytes, err := hex.DecodeString(s)
		if err != nil {
			fmt.Println("error during reading of transactions")
		}
		tx := new(types.Transaction)
		rlp.DecodeBytes(txBytes, &tx)
		out = append(out, tx)
	}
	return out
}

func TestLinkedInsert(t *testing.T) {
	is := is.New(t)
	txs := makeTransactions()
	signer := types.NewEIP155Signer(big.NewInt(1))
	pool := NewLinkedPool()
	for _, tx := range txs {
		from, err := signer.Sender(tx)
		if err != nil {
			t.Error(err)
		}
		pool.Insert(from, tx)
	}
	is.Equal(len(pool.order), len(txs))
	lastPrice := big.NewInt(1000000000000000000)
	// ensure that each tx is sorted properly
	for i := 0; i < len(txs); i++ {
		set, _ := pool.next()
		is.True(set.ID.gasPrice.Cmp(lastPrice) < 0)
		lastPrice = set.ID.gasPrice
	}
}

func TestBatching(t *testing.T) {
	// is := is.New(t)
	txs := makeTransactions()
	signer := types.NewEIP155Signer(big.NewInt(1))
	pool := NewLinkedPool()
	for _, tx := range txs {
		from, err := signer.Sender(tx)
		if err != nil {
			t.Error(err)
		}
		pool.Insert(from, tx)
	}
	btxs := pool.Batch(120000)
	bbtx := pool.Batch(120000)
	fmt.Println(btxs)
	fmt.Println(bbtx)

}

func TestRemove(t *testing.T) {
	is := is.New(t)
	// test case 1: insert 11
	s := placeSetup()
	n := &txID{gasPrice: new(big.Int).SetInt64(int64(2))}
	place(&s, n, 0, len(s)-1)
	is.Equal(n.gasPrice.String(), s[6].gasPrice.String())

	// test case 2: insert 0
	s = placeSetup()
	n = &txID{gasPrice: new(big.Int).SetInt64(int64(0))}
	place(&s, n, 0, len(s)-1)
	is.Equal(n.gasPrice.String(), s[1].gasPrice.String())
	// is.Equal(n.gasPrice.String(), s[6].gasPrice.String())
}
