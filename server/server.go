package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/evan-forbes/ethlab/thereum"
	"github.com/gorilla/mux"
)

type Server struct {
	http.Server
	router *mux.Router
	muxer  *muxer
	back   *thereum.Thereum
}

func NewServer() *Server {
	return &Server{
		router: mux.NewRouter(),
		muxer:  newMuxer(),
	}
	// set write timeouts
}

func (s *Server) rpcHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// add the header
		w.Header().Set("content-type", "application/json")
		// read the body
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			resp := rpcError(500, fmt.Sprintf("could not read request: %s", string(body)))
			w.Write(resp)
			return
		}
		// unmarshal the body
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
			w.Write(rpcError(400, "method not supported"))
			return
		}
		resp, err := pro(s.back, req)
		if err != nil {
			w.Write(rpcError(500, fmt.Sprintf("error calling %s: %s", req.Method, err)))
		}
		// marshal the response
		out, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
		}
		// write the response to the body
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

func Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called handler: ", r.Method)
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println("could not read body", err)
		return
	}
	fmt.Println(string(body))
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(`{"jsonrpc": "2.0", "id": 1, "result": {"message": "this is a reply"}}`))
}

func run() {
	srv := NewServer()
	srv.router.HandleFunc("/", http.TimeoutHandleFunc(s.rpcHandler()))
	http.ListenAndServe("127.0.0.1:8000", srv.router)
}

// A value of this type can be a JSON-RPC request, notification, successful response or
// error response. Which one it is depends on the fields.
type rpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      int             `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
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
