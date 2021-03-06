package server

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evan-forbes/ethlab/cmd"
	"github.com/evan-forbes/ethlab/contracts/ens"
	"github.com/evan-forbes/ethlab/module"
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

// func TestServer(t *testing.T) {
// 	mngr := cmd.NewManager(context.Background(), nil)
// 	go mngr.Listen()

// 	srvr := NewServer("127.0.0.1:8000", nil)
// 	go func() {
// 		t.Log(srvr.ListenAndServe())
// 	}()
// 	time.Sleep(time.Second * 3)
// 	client, err := ethclient.Dial("http://127.0.0.1:8000")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// create signed txs
// 	txs, _, err := genTxs(2, thereum.NewAccounts("alice", "bob"))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	for _, tx := range txs {
// 		err := client.SendTransaction(mngr.Ctx, tx)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	fmt.Println("done")
// 	<-mngr.Done()
// }

// // generates transactions, all of which send 1 ETH to a freshly created 'sink' account
// func genTxs(count int, accs thereum.Accounts) ([]*types.Transaction, common.Address, error) {
// 	sinkAccout, _ := thereum.NewAccount("sink", big.NewInt(0))
// 	var i int
// 	var out []*types.Transaction
// 	for _, acc := range accs {
// 		if i >= count {
// 			break
// 		}
// 		tx, err := acc.CreateSend(sinkAccout.Address, big.NewInt(1))
// 		if err != nil {
// 			return out, sinkAccout.Address, err
// 		}
// 		out = append(out, tx)
// 		i++
// 	}
// 	return out, sinkAccout.Address, nil
// }

// func TestTxMarshal(t *testing.T) {
// 	txs, _, err := genTxs(2, thereum.NewAccounts("alice", "bob"))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	for _, tx := range txs {
// 		j, err := json.Marshal(txToRPC(tx))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		fmt.Println(string(j))
// 	}
// }

// func txToRPC(tx *types.Transaction) *rpcMessage {
// 	rpcParam := sendETHrpc{
// 		To:       *tx.To(),
// 		Gas:      tx.Gas(),
// 		GasPrice: tx.GasPrice(),
// 		Value:    tx.Value(),
// 	}
// 	return &rpcMessage{
// 		Version: "2.0",
// 		ID:      60,
// 		Params:  []interface{}{rpcParam},
// 	}
// }

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

func TestTxSend(t *testing.T) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()

	eth, err := thereum.New(thereum.DefaultConfig(), nil)
	if err != nil {
		t.Error(err)
	}
	mngr.WG.Add(1)
	go eth.Run(mngr.Ctx, mngr.WG)

	srvr := NewServer(mngr.Ctx, "127.0.0.1:8000", eth)
	go func() {
		t.Log(srvr.ListenAndServe())
	}()
	time.Sleep(time.Second * 1)
	client, err := ethclient.Dial("http://127.0.0.1:8000")
	if err != nil {
		t.Error(err)
	}
	alice := eth.Accounts["Alice"]
	bob := eth.Accounts["Bob"]

	oneETH, ok := new(big.Int).SetString("1000000000000000000", 10)
	if !ok {
		t.Error("could not set string")
	}
	tx, err := alice.CreateSend(bob.Address, oneETH)
	if err != nil {
		t.Error(err)
	}

	err = client.SendTransaction(mngr.Ctx, tx)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 3)
	receipt, err := client.TransactionReceipt(mngr.Ctx, tx.Hash())
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("receipt: %+v\n", receipt)
	<-mngr.Done()
}

func TestStream(t *testing.T) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()

	eth, err := thereum.New(thereum.DefaultConfig(), nil)
	if err != nil {
		t.Error(err)
	}
	mngr.WG.Add(1)
	go eth.Run(mngr.Ctx, mngr.WG)

	srvr := NewServer(mngr.Ctx, "127.0.0.1:8000", eth)
	go func() {
		t.Log(srvr.ListenAndServe())
	}()
	go func() {
		t.Log(srvr.ServeWS("127.0.0.1:8001"))
	}()
	time.Sleep(time.Second * 1)
	client, err := ethclient.Dial("ws://127.0.0.1:8001")
	if err != nil {
		t.Error(err)
	}
	sink := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(mngr.Ctx, sink)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}
	fmt.Println("subscribed", sub)
	mngr.WG.Add(1)
	go func() {
		defer mngr.WG.Done()
		for {
			select {
			case err := <-sub.Err():
				t.Error(err)
				return
			case <-mngr.Ctx.Done():
				return
			case head := <-sink:
				fmt.Printf("%+v\n", head.Hash().Hex())
			}
		}
	}()
	mngr.WG.Wait()
}

func TestENSHandler(t *testing.T) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()
	_, err := LaunchServer(mngr.Ctx, mngr.WG)
	if err != nil {
		t.Error(err)
	}
	_, err = module.ENSAddress("127.0.0.1:8000")
	if err != nil {
		t.Error(err)
	}
}

/*
I need to use users instead of thereum accounts for root
its simpler one, and two it increments nonce every use.
*/

func TestGetNonce(t *testing.T) {
	mngr := cmd.NewManager(context.Background(), nil)
	go mngr.Listen()
	_, err := LaunchServer(mngr.Ctx, mngr.WG)
	if err != nil {
		t.Error(err)
	}
	client, err := ethclient.Dial("http://127.0.0.1:8000")
	if err != nil {
		t.Error(err)
		return
	}
	usr, err := module.NewUser()
	client.NonceAt(mngr.Ctx, usr.NewTxOpts().From, nil)
}

func bootWithUser(t *testing.T) (*module.User, *Server, error) {
	// setup
	mngr := cmd.NewManager(context.Background(), nil)
	// listen for cancels
	go mngr.Listen()
	// start the server
	srv, err := LaunchServer(mngr.Ctx, mngr.WG)
	if err != nil {
		t.Error(err)
		return nil, srv, err
	}
	// make a user
	usr, err := module.StarterKit("127.0.0.1:8000")
	if err != nil {
		t.Error(err)
		return nil, srv, err
	}
	err = module.RequestETH("127.0.0.1:8000", usr.From.Hex(), big.NewInt(1000000000000000000))
	if err != nil {
		t.Error(err)
		return nil, srv, err
	}
	return usr, srv, nil
}

// so for some reason the transactions are being finalized but not being caught by the logger.
func TestSubscribeLogs(t *testing.T) {
	usr, srv, err := bootWithUser(t)
	if err != nil {
		t.Error(err)
		return
	}
	addr, err := ens.Deploy(usr)
	if err != nil {
		t.Error(err)
		return
	}
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			addr,
		},
	}
	logs := make(chan []*types.Log)
	sub, err := srv.back.Events.SubscribeLogs(query, logs)
	if err != nil {
		t.Error(err)
		return
	}
	go func() {
		for {
			select {
			case err := <-sub.Err():
				t.Error(err)
				return
			case log := <-logs:
				for _, l := range log {
					fmt.Println("found a log !!!!!!! ", l.TxHash.Hex())
				}
			}
		}
	}()
	time.Sleep(100 * time.Millisecond)
	ens, err := ens.NewENS(addr, usr.Client)
	if err != nil {
		t.Error(err)
	}
	tx, err := ens.LogTest(usr.NewTxOpts())
	if err != nil {
		t.Error(err)
	}
	fmt.Println("homebrewed tx", tx.Hash().Hex())
	time.Sleep(5 * time.Second)

}
