package thereum

import (
	"sync"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
)

// API maps *Thereum methods to go-ethereum's bind.ContractBackend
type API struct {
	*Thereum
}

func NewBackendAPI(config *Config, root *opts.Trans) (out API{}, err error) {
	if config == nil {	
		config = &defaultConfig()
	}
	New(*config, )
	return API{}
}