package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/evan-forbes/ethlab/thereum"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

// use the websocket.Handler
// defined at https://rollout.io/blog/getting-started-with-websockets-in-go/

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
		fmt.Println("method", req.Method)

		// // forward any pub/sub requests to the websocket handler
		// if req.Method == "eth_subscribe" {
		// 	s.wsHandler(w, r)
		// 	return
		// }

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
// 	Routing Subscriptions
//////////////////////////////

// there are only two supported methods, so making an entire seperate router doesn't
// quite make sense

type wsProcedure func(ctx context.Context, eth *thereum.Thereum, conn *websocket.Conn, params []string) error

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

/*
web socket notes
pretty sure the handler func upgrades the connection to a websocket connection?

still not sure if there are two seperate servers running on two different ports or what.

the handler may start a goroutine to handle any transmission of data to the client using
the conn.WriteJSON() method.

having a seperate server might not be a big deal.
*/

////////////////////////////////
// 		Web Socket
//////////////////////////////

func (s *Server) socket(conn *websocket.Conn) {
	for {
		var recv []byte
		_, err := conn.Read(recv)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("calling socket! reviecing: ", string(recv))
	}
}
