package thereum

import "github.com/ethereum/go-ethereum/core"

// TODO: find which functions need to be called to commit to the chain.

/*
so far the appears to be
- make a new block?
- sign header of new block and insert into header.Extra
- seal block
-
*/

type Thereum struct {
	txPool *core.TxPool
}

func New() (*Thereum, error) {
	return &TherThereum{}, nil
}

func (t *Thereum) Start() {

}

func (t *Thereum) Shutdown() {

}
