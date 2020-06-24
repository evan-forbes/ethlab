package thereum

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

// Config contains the standard variables for creating a new Thereum chain/node
type Config struct {
	InMemory      bool              `json:"in_memory"`
	GenesisConfig core.Genesis      `json:"genesis"`
	Allocation    map[string]string `json:"allocation"` // "Name": "100000000000000000"
	GasLimit      uint64            `json:"gas_limit"`
	Delay         uint
	Host          string `json:"host"`
	Port          uint   `json:"port"`
	WSHost        string `json:"ws_host"`
	WSPort        uint   `json:"ws_port"`
}

// ConfigFromFile opens and decodes a config.json file
func ConfigFromFile(path string) (Config, error) {
	var out Config
	file, err := os.Open(path)
	if err != nil {
		return out, err
	}
	err = json.NewDecoder(file).Decode(&out)
	return out, err
}

// DB returns the proper database specified by the config
// currently only supports in memory databases
func (c Config) DB() ethdb.Database {
	if c.InMemory {
		return rawdb.NewMemoryDatabase()
	}
	return rawdb.NewMemoryDatabase()
}

// // GasLimiter uses the config to init a new GasLimiter
// // TODO: alter to include more ways to limit gas
// func (c Config) GasLimiter() GasLimiter {
// 	out, ok := new(big.Int).SetString(c.GasLimit, 10)
// 	if !ok {
// 		return &ConstantGasLimit{big.NewInt(10485760)}
// 	}
// 	return &ConstantGasLimit{limit: out}
// }

// Genesis issues a new genesis configuration specified in the config
func (c Config) Genesis() (core.Genesis, Accounts, error) {
	var out core.Genesis
	// if c.GenesisConfig == nil {
	// 	out = defaultGenesis()
	// }
	accnts := make(Accounts)
	var err error
	for name, sbal := range c.Allocation {
		bal, ok := new(big.Int).SetString(sbal, 10)
		if !ok {
			err = errors.New("could set string balance during genesis allocations")
		}
		acc, aerr := NewAccount(name, bal)
		if aerr != nil {
			fmt.Println("problem making new account for", name, bal.String(), aerr)
			err = aerr
		}
		accnts[name] = acc
	}
	accnts.SetGasPrice(big.NewInt(10000))
	out.Alloc = accnts.Genesis()
	return out, accnts, err
}

func defaultGenesis() core.Genesis {
	alloc := core.GenesisAlloc(
		make(map[common.Address]core.GenesisAccount),
	)
	genesis := core.Genesis{
		Config:     params.AllEthashProtocolChanges,
		GasLimit:   104857600,
		Alloc:      alloc,
		Difficulty: new(big.Int).SetInt64(1),
	}
	genesis.Config.ChainID = big.NewInt(1)
	return genesis
}

func DefaultConfig() Config {
	return Config{
		InMemory:      true,
		GenesisConfig: defaultGenesis(),
		Allocation: map[string]string{
			"root":   "999999999999999999999999999999",
			"Alice":  "10000000000000000000",
			"Bob":    "10000000000000000000",
			"Celine": "10000000000000000000",
			"Dobby":  "10000000000000000000",
			"Erin":   "10000000000000000000",
			"Frank":  "10000000000000000000",
		},
		GasLimit: 10485760,
		Delay:    100,
		Host:     "127.0.0.1",
		Port:     8438,
		WSHost:   "127.0.0.1",
		WSPort:   8439,
	}
}
