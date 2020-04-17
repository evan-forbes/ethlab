package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

func TestServer(t *testing.T) {
	go run()
	time.Sleep(100 * time.Millisecond)
	client, err := rpc.Dial("http://127.0.0.1:8000")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("connected")
	var rsp HelloReply
	fmt.Println(client.Call(&rsp, "", HelloArgs{Who: "evan"}))
	fmt.Println("response", rsp)
	time.Sleep(1 * time.Second)
}
