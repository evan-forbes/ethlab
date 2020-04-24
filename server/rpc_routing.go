package server

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/evan-forbes/ethlab/thereum"
)

// muxer maps supported methods to their appropriate procedures.
// thread safe.
type muxer struct {
	routes map[string]procedure
	mut    sync.RWMutex
}

type procedure func(eth *thereum.Thereum, msg rpcMessage) (*rpcMessage, error)

func newMuxer() *muxer {
	return &muxer{
		routes: map[string]procedure{
			// add rpc methods here!
			"":                       nullProcedure,
			"eth_chainId":            nullProcedure,
			"eth_protocolVersion":    nullProcedure,
			"eth_gasPrice":           nullProcedure,
			"eth_blockNumber":        nullProcedure,
			"eth_getBalance":         nullProcedure,
			"eth_getStorageAt":       nullProcedure,
			"eth_sendTransaction":    nullProcedure,
			"eth_call":               nullProcedure,
			"eth_getLogs":            nullProcedure,
			"eth_getFilterLogs":      nullProcedure,
			"eth_sendRawTransaction": sendETH,
		},
	}
}

func (m *muxer) Route(method string) (procedure, bool) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	pro, has := m.routes[method]
	return pro, has
}

func nullProcedure(eth *thereum.Thereum, msg rpcMessage) (*rpcMessage, error) {
	nullMessage := &rpcMessage{
		Version: "2.0",
		ID:      60,
		Error: &jsonError{
			Code:    999,
			Message: "specified method is not registered or supported",
		},
	}
	return nullMessage, nil
}

type sendTxParams struct {
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Gas      uint64         `json:"gas"`
	GasPrice *big.Int       `json:"gasPrice"`
	Value    *big.Int       `json:"value"`
	Data     string         `json:"data"`
}

type sendTxParamsJSON struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
	Data     string `json:"data"`
}

func (p *sendTxParams) UnmarshalJSON(in []byte) error {
	var data sendTxParamsJSON
	err := json.Unmarshal(in, &data)
	if err != nil {
		return err
	}
	p.From = common.HexToAddress(data.From)
	p.To = common.HexToAddress(data.To)
	p.Gas, err = strconv.ParseUint(data.To, 10, 64)
	if err != nil {
		return err
	}
	// REMOVE '0x'
	_, ok := p.Value.SetString(data.Value, 16)
	if !ok {
		return fmt.Errorf("could not parse big integer")
	}

	return nil
}

func (p *sendTxParams) MarshalJSON() ([]byte, error) {
	var out []byte
	return out, nil
}

func sendTx(eth *thereum.Thereum, msg rpcMessage) (*rpcMessage, error) {
	out := &rpcMessage{
		Version: "2.0",
		ID:      60,
	}
	return out, nil
}

func chainID(eth *thereum.Thereum, msg rpcMessage) (*rpcMessage, error) {
	out := &rpcMessage{
		Version: "2.0",
		ID:      60,
	}

	// parse rpc msg into tx
	// maybe do some extra tx validation
	// add tx to the pool

	// sender, err := types.Sender(types.NewEIP155Signer(b.config.ChainID), tx)
	// if err != nil {
	// 	return out, fmt.Errorf("invalid transaction: %v", err)
	// }
	return out, nil
}

/*
I need to figure out how I'm going to parse incoming parameters
things I have:
 	- the types required for that rpc to be fullfilled
	- examples of the result

options
	- make a parser for each type and attempt to call it for that specific unmarshaller/parser
	-
*/