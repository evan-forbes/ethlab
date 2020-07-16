package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/evan-forbes/ethlab/thereum"
	"github.com/pkg/errors"
)

/////////////////////////////
// 		Procedures
///////////////////////////

// faucet facilitates handing out eth
func faucet(eth *thereum.Thereum, msg *rpcMessage) (*rpcMessage, error) {
	return nil, nil
}
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
		ID:      1,
		Result:  tx.Hash().Hex(),
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
	// ensure that some data was passed throught the rpc msg
	if len(hexTx) == 0 {
		return nil, errors.New("no parameters provided for raw transaction")
	}
	// unmarshal the hex bytes into a transaction
	hash := common.HexToHash(hexTx[0])
	fmt.Println("hash", hash)
	// fetch the receipt
	receipt, err := eth.TxReceipt(hash)
	if err != nil {
		return nil, err
	}
	// marshal the result

	if err != nil {
		return nil, errors.Wrap(err, "failed getTxReceipt")
	}
	out := &rpcMessage{
		Version: "2.0",
		ID:      1,
		Result:  receipt,
	}
	return out, nil
}

// getTxCount returns the number of transaction sent from an address at a given
// block
func getTxCount(eth *thereum.Thereum, msg *rpcMessage) (*rpcMessage, error) {
	// "params":["0x407d73d8a49eeb85d32cf465507dd71d507100c1","latest"]
	var params []interface{}
	err := json.Unmarshal(msg.Params, &params)
	if err != nil {
		return nil, err
	}
	if len(params) != 2 {
		return nil, errors.New("2 arguments needed in parameters")
	}
	hexAddr, ok := params[0].(string)
	if !ok {
		return nil, errors.New("first arg in params must be a string representing an address")
	}
	addr := common.HexToAddress(hexAddr)

	hexHash, ok := params[1].(string)
	if !ok {
		return nil, errors.New("block numbers not yet supported, try a block hash")
	}
	// type switch, to add numbers later?
	var hsh common.Hash
	if hexHash == "latest" {
		hsh = eth.LatestBlock().Hash()
	} else {
		hsh = common.HexToHash(hexHash)
	}
	count, err := eth.TransactionCountByAddress(context.Background(), addr, hsh)
	if err != nil {
		return nil, err
	}
	out := &rpcMessage{
		Version: "2.0",
		ID:      1,
		Result:  count,
	}
	return out, nil
}

func getBalanceAt(eth *thereum.Thereum, msg *rpcMessage) (*rpcMessage, error) {
	// "params":["0x407d73d8a49eeb85d32cf465507dd71d507100c1", "latest"]
	var params []interface{}
	err := json.Unmarshal(msg.Params, &params)
	if err != nil {
		return nil, err
	}
	if len(params) != 2 {
		return nil, errors.New("2 arguments needed in parameters")
	}
	hexAddr, ok := params[0].(string)
	if !ok {
		return nil, errors.New("first arg in params must be a string representing an address")
	}
	addr := common.HexToAddress(hexAddr)

	// hexHash, ok := params[1].(string)
	// if !ok {
	// 	return nil, errors.New("block numbers not yet supported, try a block hash")
	// }
	bal, err := eth.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		return nil, err
	}
	out := &rpcMessage{
		Version: "2.0",
		ID:      1,
		Result:  fmt.Sprintf("0x%x", bal),
	}
	return out, nil
}

func getNonce(eth *thereum.Thereum, msg *rpcMessage) (*rpcMessage, error) {
	return nil, nil
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
