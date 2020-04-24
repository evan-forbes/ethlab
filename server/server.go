package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/evan-forbes/ethlab/thereum"
	"github.com/gorilla/mux"
)

// Server connects traditional ethereum clients to the thereum backend via
// the standardized ethereum json rpc.
type Server struct {
	http.Server
	router *mux.Router
	muxer  *muxer
	back   *thereum.Thereum
}

// NewServer issues a new server with the rpc handler already registered
func NewServer(addr string) *Server {
	rtr := mux.NewRouter()
	srv := &Server{
		Server: http.Server{
			Addr:         addr,
			WriteTimeout: time.Second * 20,
			ReadTimeout:  time.Second * 10,
			IdleTimeout:  time.Second * 100,
			Handler:      rtr,
			// TLSConfig: ,
		},
		router: rtr,
		muxer:  newMuxer(),
	}
	// install the universal rpc handler to the router
	srv.router.HandleFunc("/", srv.rpcHandler())
	return srv
	// set write timeouts
}

// rpcHandler returns the http handler function that processes all rpc requests
func (s *Server) rpcHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// add the header to the response
		w.Header().Set("content-type", "application/json")

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
		resp, err := pro(s.back, req)
		if err != nil {
			w.Write(rpcError(500, fmt.Sprintf("error calling %s: %s", req.Method, err)))
			return
		}

		// marshal the processed response
		out, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
		}

		// write the response to the client
		_, err = w.Write(out)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func rpcError(code int, msg string) []byte {
	out, _ := json.Marshal(
		rpcMessage{
			Version: "2.0",
			ID:      60,
			Error: &jsonError{
				Code:    code,
				Message: msg,
			},
		},
	)
	return out
}

// func Handle(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("called handler: ", r.Method)
// 	body, err := ioutil.ReadAll(r.Body)
// 	defer r.Body.Close()
// 	if err != nil {
// 		fmt.Println("could not read body", err)
// 		return
// 	}
// 	fmt.Println(string(body))
// 	w.Header().Set("content-type", "application/json")
// 	w.Write([]byte(`{"jsonrpc": "2.0", "id": 1, "result": {"message": "this is a reply"}}`))
// }

// A value of this type can be a JSON-RPC request, notification, successful response or
// error response. Which one it is depends on the fields.
type rpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      int             `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  []interface{}   `json:"params,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (s *Server) baseRPCMessage() rpcMessage {
	return rpcMessage{
		Version: "2.0",
		ID:      int(s.back.Config.ChainID.Int64()),
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
