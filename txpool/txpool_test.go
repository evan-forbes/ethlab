package txpool

import (
	"math/big"
	"testing"

	"github.com/matryer/is"
)

func placeSetup() []*txID {
	var s []*txID
	for i := 0; i < 20; i = i + 2 {
		s = append(s, &txID{gasPrice: new(big.Int).SetInt64(int64(i))})
	}
	return s
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
