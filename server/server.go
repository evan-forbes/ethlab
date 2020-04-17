package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO: figure out what the rpc client expects as far as headers/whatever
// using json encoding/decoding for msgs

type HelloArgs struct {
	Who string `json:"who"`
}

type HelloReply struct {
	Message string `json:"message"`
}

type HelloService struct {
	router *mux.Router
}

func Handl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called handler: ", r.Method)

}

func run() {
	srv := HelloService{router: mux.NewRouter()}
	srv.router.HandleFunc("/", Handl)
	http.ListenAndServe("127.0.0.1:8000", srv.router)
}
