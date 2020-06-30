package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/evan-forbes/ethlab/contracts/ens"
	"github.com/evan-forbes/ethlab/module"
	"github.com/evan-forbes/ethlab/thereum"
	"github.com/gorilla/mux"
)

/* TODO:
- handle errors during streaming
- implement the ability to unsubscribe cleanly
*/

////////////////////////////////
// 		RPC Server
//////////////////////////////

// Server connects traditional ethereum clients to the thereum backend via
// the standardized ethereum json rpc.
type Server struct {
	http.Server
	router *mux.Router      // handles api endpoints
	muxer  *muxer           // connects msg to procedure
	back   *thereum.Thereum // backend to serve
	ctx    context.Context
}

// NewServer issues a new server with the rpc handler already registered
func NewServer(ctx context.Context, addr string, back *thereum.Thereum) *Server {
	rtr := mux.NewRouter()
	srv := &Server{
		Server: http.Server{
			Addr:         addr,
			WriteTimeout: time.Second * 20, // TODO: read from config
			ReadTimeout:  time.Second * 10,
			IdleTimeout:  time.Second * 100,
			Handler:      rtr,
			// TLSConfig: ,
		},
		router: rtr,
		muxer:  newMuxer(),
		ctx:    ctx,
	}
	// install the universal rpc handler to the router
	srv.router.HandleFunc("/", srv.rpcHandler())
	srv.router.HandleFunc("/requestETH", srv.faucetHandler())
	var ens ensHandler
	err := ens.Install("http://" + addr)
	if err != nil {
		log.Fatalf("failure to deploy ens contract %s", err)
	}
	srv.router.HandleFunc("/ens", ens.Handle)
	srv.back = back
	return srv
	// set write timeouts
}

// rpcHandler returns the main http handler function that processes *all* rpc requests
func (s *Server) rpcHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// add the header to the response
		// w.Header().Set("content-type", "application/json")

		// read the body of the request
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			resp := rpcError(500, fmt.Sprintf("could not read request: %s", string(body)))
			w.Write(resp)
			return
		}

		// unmarshal the rpc request
		var req rpcMessage
		err = json.Unmarshal(body, &req)
		if err != nil {
			// send an error back if rpcMessage cannot be unmarshaled
			w.Write(rpcError(500, fmt.Sprintf("could not unmarshal rpc message: %s", err)))
			return
		}

		// use the method's procdure to perform the remote procedure call
		pro, has := s.muxer.Route(req.Method)
		if !has {
			log.Println("no procedure for method: ", req.Method)
			w.Write(rpcError(0, fmt.Sprintf("method %s not supported", req.Method)))
			return
		}
		resp, err := pro(s.back, &req)
		if err != nil {
			w.Write(rpcError(500, fmt.Sprintf("error calling %s: %s", req.Method, err)))
			return
		}

		// marshal the processed response
		out, err := json.Marshal(resp)
		if err != nil {
			w.Write(rpcError(500, fmt.Sprintf("interanl marshaling error calling %s: %s", req.Method, err)))
			log.Println(err)
			log.Printf("failed to marshal response from %s procedure: %+v\n", req.Method, resp)
		}

		// write the response to the client
		_, err = w.Write(out)
		if err != nil {
			w.Write(rpcError(500, fmt.Sprintf("interanl writing error calling %s: %s", req.Method, err)))
			log.Println(err)
			return
		}
	}
}

////////////////////////////////
// 	Routing Msg to Procedure
//////////////////////////////

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

////////////////////////////////
// 		RPC Messaging
//////////////////////////////

// A value of this type can be a JSON-RPC request, notification, successful response or
// error response. Which one it is depends on the fields.
type rpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      int             `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (s *Server) baseRPCMessage() rpcMessage {
	return rpcMessage{
		Version: "2.0",
		ID:      1,
	}
}

func rpcError(code int, msg string) []byte {
	out, _ := json.Marshal(
		rpcMessage{
			Version: "2.0",
			ID:      1,
			Error: &jsonError{
				Code:    code,
				Message: msg,
			},
		},
	)
	return out
}

////////////////////////////////
//	Ethereum Naming Server
//////////////////////////////

type ensHandler struct {
	address  common.Address
	deployed bool
}

func (e *ensHandler) Install(host string) error {
	usr, err := module.StarterKit(host)
	if err != nil {
		return err
	}
	addr, err := ens.Deploy(usr)
	if err != nil {
		return err
	}
	e.address = addr
	e.deployed = true
	return nil
}

// Handle responds to an http request with the address of the ENS contract
func (e *ensHandler) Handle(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(e.address.Hex()))
	if err != nil {
		log.Println("failure to respond via ens handler:", err)
		return
	}
	return
}

////////////////////////////////
//	Faucet
//////////////////////////////

func (s *Server) faucetHandler() http.HandlerFunc {

	type faucetPay struct {
		Address string   `json:"address"`
		Amount  *big.Int `json:"amount"`
	}

	type faucetResp struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// read the body of the request
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			resp, _ := json.Marshal(faucetResp{Message: err.Error()})
			w.Write(resp)
			return
		}

		// unmarshal the rpc request
		var pay faucetPay
		err = json.Unmarshal(body, &pay)
		if err != nil {
			resp, _ := json.Marshal(faucetResp{Message: err.Error()})
			w.Write(resp)
			return
		}
		// check if the recipient qualifies for free eth and the address is valid
		root := s.back.Accounts["root"]
		oneETH, _ := new(big.Int).SetString("1000000000000000000", 10)
		tx, err := root.CreateSend(common.HexToAddress(pay.Address), oneETH)
		if err != nil {
			log.Printf("could not send one ETH upon request to %s: %s\n", pay.Address, err.Error())
			return
		}
		err = s.back.AddTx(tx)
		if err != nil {
			log.Printf("could not send one ETH upon request to %s: %s failure to add transaction\n", pay.Address, err.Error())
			return
		}
		resp, _ := json.Marshal(faucetResp{Message: "success"})
		_, err = w.Write(resp)
		if err != nil {
			log.Println("failure to respond with successful faucet response:", err)
			return
		}
		return
	}
}
