package server

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
