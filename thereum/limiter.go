package thereum

import "math/big"

type GasLimiter interface {
	Limit() *big.Int
}

type ConstantGasLimit struct {
	limit *big.Int
}

func (l *ConstantGasLimit) Limit() *big.Int {
	return l.limit
}
