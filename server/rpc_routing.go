package server

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/evan-forbes/ethlab/thereum"
)

// muxer maps supported methods to their appropriate procedures.
// thread safe.
type muxer struct {
	routes map[string]procedure
	mut    sync.RWMutex
}

type procedure func(eth *thereum.Thereum, msg *rpcMessage) (*rpcMessage, error)

func newMuxer() *muxer {
	return &muxer{
		routes: map[string]procedure{
			// add rpc methods here!
			"":                          nullProcedure,
			"eth_chainId":               nullProcedure,
			"eth_protocolVersion":       nullProcedure,
			"eth_gasPrice":              nullProcedure,
			"eth_blockNumber":           nullProcedure,
			"eth_getBalance":            nullProcedure,
			"eth_getStorageAt":          nullProcedure,
			"eth_sendTransaction":       sendRawTx, // account management shouldn't really be a feature
			"eth_sendRawTransaction":    sendRawTx,
			"eth_getTransactionReceipt": getTxReceipt,
			"eth_call":                  nullProcedure,
			"eth_getLogs":               nullProcedure,
			"eth_getFilterLogs":         nullProcedure,
		},
	}
}

func (m *muxer) Route(method string) (procedure, bool) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	pro, has := m.routes[method]
	return pro, has
}

/////////////////////////////
// 		Procedures
///////////////////////////

func nullProcedure(eth *thereum.Thereum, msg *rpcMessage) (*rpcMessage, error) {
	nullMessage := rpcMessage{
		Version: "2.0",
		ID:      1,
		Error: &jsonError{
			Code:    999,
			Message: "specified method is not registered or supported",
		},
	}
	return &nullMessage, nil
}

type sendTxParams struct {
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Gas      uint64         `json:"gas"`
	GasPrice *big.Int       `json:"gasPrice"`
	Value    *big.Int       `json:"value"`
	Data     string         `json:"data"`
}

// sendRawTx handles a singed raw transaction provided in an rpc message
func sendRawTx(eth *thereum.Thereum, msg *rpcMessage) (*rpcMessage, error) {
	// unmarshal into temp data structs (passed via json as a slice of a single hex string)
	var hexTx []string
	err := json.Unmarshal(msg.Params, &hexTx)
	if err != nil {
		return nil, err
	}
	// ensure that some data was passed throught the rpc msg
	if len(hexTx) == 0 {
		return nil, errors.New("no parameters provided for raw transaction")
	}
	fmt.Println(hexTx[0])
	// unmarshal the hex bytes into a transaction
	var tx types.Transaction
	txBytes, err := hex.DecodeString(strings.Replace(hexTx[0], "0x", "", 1))
	if err != nil {
		return nil, err
	}
	err = rlp.DecodeBytes(txBytes, &tx)
	if err != nil {
		return nil, err
	}
	fmt.Println("tx decoded", tx, err)
	if err != nil {
		return nil, err
	}

	// add the transaction to thereum
	err = eth.AddTx(&tx)
	if err != nil {
		return nil, err
	}
	out := &rpcMessage{
		Version: "2.0",
		ID:      int(eth.Config.ChainID.Int64()),
		Result:  []byte(""),
	}

	return out, nil
}

// getTxReceipt attempts to fetch receipt data from the thereum object based on the hash
// provided in the rpc message
func getTxReceipt(eth *thereum.Thereum, msg *rpcMessage) (*rpcMessage, error) {
	// unmarshal into temp data structs (passed via json as a slice of a single hex string)
	var hexTx []string
	err := json.Unmarshal(msg.Params, &hexTx)
	if err != nil {
		return nil, err
	}
	fmt.Println("hexTx", hexTx[0])
	// ensure that some data was passed throught the rpc msg
	if len(hexTx) == 0 {
		return nil, errors.New("no parameters provided for raw transaction")
	}
	// unmarshal the hex bytes into a transaction
	hash := common.HexToHash(hexTx[0])
	fmt.Println("hash", hash)
	// fetch the receipt
	receipt, err := eth.TxReceipt(hash)
	fmt.Println("txr", err)
	if err != nil {
		return nil, err
	}
	fmt.Println("internal receipt", receipt)
	// marshal the result
	res, err := json.Marshal(receipt)
	if err != nil {
		return nil, err
	}
	out := &rpcMessage{
		Version: "2.0",
		ID:      int(eth.Config.ChainID.Int64()),
		Result:  res,
	}
	return out, nil
}

// func (p *sendTxParams) UnmarshalJSON(in []byte) error {
// 	type params struct {
// 		From     string `json:"from"`
// 		To       string `json:"to"`
// 		Gas      string `json:"gas"`
// 		GasPrice string `json:"gasPrice"`
// 		Value    string `json:"value"`
// 		Data     string `json:"data"`
// 	}
// 	var data params
// 	err := json.Unmarshal(in, &data)
// 	if err != nil {
// 		return err
// 	}
// 	p.From = common.HexToAddress(data.From)
// 	p.To = common.HexToAddress(data.To)
// 	p.Gas, err = strconv.ParseUint(data.To, 10, 64)
// 	if err != nil {
// 		return err
// 	}
// 	// REMOVE '0x'
// 	_, ok := p.Value.SetString(data.Value, 16)
// 	if !ok {
// 		return fmt.Errorf("could not parse big integer")
// 	}

// 	return nil
// }

// func chainID(eth *thereum.Thereum, msg rpcMessage) (*rpcMessage, error) {
// 	out := &rpcMessage{
// 		Version: "2.0",
// 		ID:      60,
// 	}

// 	// parse rpc msg into tx
// 	// maybe do some extra tx validation
// 	// add tx to the pool

// 	// sender, err := types.Sender(types.NewEIP155Signer(b.config.ChainID), tx)
// 	// if err != nil {
// 	// 	return out, fmt.Errorf("invalid transaction: %v", err)
// 	// }
// 	return out, nil
// }

/*
I need to figure out how I'm going to parse incoming parameters
things I have:
 	- the types required for that rpc to be fullfilled
	- examples of the result

options
	- make a parser for each type and attempt to call it for that specific unmarshaller/parser
	-
*/
