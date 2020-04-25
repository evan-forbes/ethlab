package server

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evan-forbes/ethlab/cmd"
	"github.com/evan-forbes/ethlab/thereum"
	"github.com/matryer/is"
)

// func TestServer(t *testing.T) {
// 	go run()
// 	time.Sleep(100 * time.Millisecond)
// 	client, err := rpc.Dial("http://127.0.0.1:8000")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	fmt.Println("connected")
// 	// do something with the client like send a transaction or log in
// 	// fmt.Println(client.Call(&rsp, "deploy", ))
// 	// fmt.Println("response", rsp)
// 	time.Sleep(1 * time.Second)
// }

func TestServer(t *testing.T) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()

	srvr := NewServer("127.0.0.1:8000")
	go func() {
		t.Log(srvr.ListenAndServe())
	}()
	time.Sleep(time.Second * 3)
	client, err := ethclient.Dial("http://127.0.0.1:8000")
	if err != nil {
		t.Error(err)
	}
	// create signed txs
	txs, _, err := genTxs(2, thereum.NewAccounts("alice", "bob"))
	if err != nil {
		t.Error(err)
	}
	for _, tx := range txs {
		fmt.Println("attempting to send txs")
		err := client.SendTransaction(mngr.Ctx, tx)
		if err != nil {
			t.Error(err)
		}
	}
	fmt.Println("done")
	<-mngr.Done()
}

// generates transactions, all of which send 1 ETH to a freshly created 'sink' account
func genTxs(count int, accs thereum.Accounts) ([]*types.Transaction, common.Address, error) {
	sinkAccout, _ := thereum.NewAccount("sink", big.NewInt(0))
	var i int
	var out []*types.Transaction
	for _, acc := range accs {
		if i >= count {
			break
		}
		tx, err := acc.CreateSend(sinkAccout.Address, big.NewInt(1))
		if err != nil {
			return out, sinkAccout.Address, err
		}
		out = append(out, tx)
		i++
	}
	return out, sinkAccout.Address, nil
}

func TestTxMarshal(t *testing.T) {
	txs, _, err := genTxs(2, thereum.NewAccounts("alice", "bob"))
	if err != nil {
		t.Error(err)
	}
	for _, tx := range txs {
		j, err := json.Marshal(txToRPC(tx))
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(j))
	}
}

func txToRPC(tx *types.Transaction) *rpcMessage {
	rpcParam := sendETHrpc{
		To:       *tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
	}
	return &rpcMessage{
		Version: "2.0",
		ID:      60,
		Params:  []interface{}{rpcParam},
	}
}

type sendETHrpc struct {
	From     common.Address
	To       common.Address
	Gas      uint64
	GasPrice *big.Int
	Value    *big.Int
}

type sendETHrpcJSONwrap struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
}

func (msg sendETHrpc) MarshalJSON() ([]byte, error) {
	out := sendETHrpcJSONwrap{
		From:     msg.From.Hex(),
		To:       msg.To.Hex(),
		Gas:      "0x" + strconv.FormatUint(msg.Gas, 16),
		GasPrice: fmt.Sprintf("0x%x", msg.GasPrice),
		Value:    fmt.Sprintf("0x%x", msg.Value),
	}
	fmt.Println(out)
	return json.Marshal(out)
}

// checks the MarshalJSON method for sendETHrpc
func TestSendETHrpcMarshal(t *testing.T) {
	is := is.New(t)
	gsprc, ok := new(big.Int).SetString("10000000000000", 10)
	if !ok {
		t.Error("set string error:")
	}
	val, ok := new(big.Int).SetString("2441406250", 10)
	if !ok {
		t.Error("set string error")
	}
	expected := `{"from":"0xb60E8dD61C5d32be8058BB8eb970870F07233155","to":"0xd46E8dD67C5d32be8058Bb8Eb970870F07244567","gas":"0x76c0","gasPrice":"0x9184e72a000","value":"0x9184e72a"}`
	v1 := sendETHrpc{
		From:     common.HexToAddress("0xb60e8dd61c5d32be8058bb8eb970870f07233155"),
		To:       common.HexToAddress("0xd46e8dd67c5d32be8058bb8eb970870f07244567"),
		Gas:      30400,
		GasPrice: gsprc,
		Value:    val,
	}
	result, err := v1.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	is.Equal(string(result), expected)
}

type tj1 struct {
	A string `json:"a"`
	B string `json:"b"`
	C int    `json:"c"`
}

type tj2 struct {
	A string  `json:"a"`
	B string  `json:"b"`
	C float64 `json:"c"`
}

func TestSliceUnmarshal(t *testing.T) {
	out := []interface{}{}
	out = append(out, &tj1{})
	out = append(out, &tj2{})
	data := []byte(`[{"a": "cat", "b": "dog", "c": 42}, {"a": "cat", "b": "dog", "c": 42.42}]`)
	err := json.Unmarshal(data, &out)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(out[0], "woooo", out[1])
}
