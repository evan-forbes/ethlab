package thereum

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

type Config struct {
	InMemory      bool              `json:"in_memory"`
	GenesisConfig *core.Genesis     `json:"genesis"`
	Allocation    map[string]string `json:"allocation"` // "Name": "100000000000000000"
}

// DB returns the proper database specified by the config
// currently only supports in memory databases
func (c Config) DB() ethdb.Database {
	if c.InMemory {
		return rawdb.NewMemoryDatabase()
	}
	return rawdb.NewMemoryDatabase()
}

// Genesis issues a new genesis configuration specified in the config
func (c Config) Genesis() (core.Genesis, Accounts, error) {
	var out core.Genesis
	if c.GenesisConfig == nil {
		out = defaultGenesis()
	}
	accnts := make(Accounts)
	var err error
	for name, sbal := range c.Allocation {
		bal, ok := new(big.Int).SetString(sbal, 10)
		if !ok {
			err = errors.New("could set string balance during genesis allocations")
		}
		acc, aerr := NewAccount(name, bal)
		if aerr != nil {
			err = aerr
		}
		accnts[name] = acc
	}
	out.Alloc = accnts.Genesis()
	return out, accnts, err
}

func defaultGenesis() core.Genesis {
	alloc := core.GenesisAlloc(
		make(map[common.Address]core.GenesisAccount),
	)
	genesis := core.Genesis{
		Config:     params.AllEthashProtocolChanges,
		GasLimit:   10485760,
		Alloc:      alloc,
		Difficulty: new(big.Int).SetInt64(1),
	}
	return genesis
}

// func readAllocation(in map[string]string) (map[string]*big.Int, error) {
// 	out := make(map[string]*big.Int)
// 	var err error
// 	for name, bal := range in {
// 		nbal, nerr := new(big.Int).String(bal)
// 		if nerr != nil {
// 			err = nerr
// 			continue
// 		}
// 		out[name] = nbal
// 	}
// 	return out, err
// }
