package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var wsPool = new(sync.Pool)

// RunWS starts a seperate websocket server for ethereum json rpc pub/sub
func (s *Server) RunWS(endpoint string) error {
	lstnr, err := net.Listen("tcp", endpoint)
	if err != nil {
		return err
	}
	srv := http.Server{Handler: s.wsHandler()}
	return srv.Serve(lstnr)
}

// wsHandler is the main handler for issuing subscriptions via websocket
func (s *Server) wsHandler() http.Handler {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		WriteBufferPool: wsPool,
		CheckOrigin:     originValidator,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// upgrade the connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("websocket upgrade failed", err)
			return
		}
		var msg rpcMessage
		err = conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("error reading json with ws:", err)
		}
		rawPrms := string(msg.Params)
		switch {
		case strings.Contains(rawPrms, "logs"):
			ctx, _ := context.WithTimeout(s.ctx, time.Hour)
			subLogs(ctx, s.back, conn, json.RawMessage(rawPrms))
		case strings.Contains(rawPrms, "newHeads"):
			ctx, _ := context.WithTimeout(s.ctx, time.Hour)
			subHeads(ctx, s.back, conn)
		default:
			w.Write(rpcError(500, fmt.Sprintf("no subscription for %s", rawPrms)))
			conn.Close()
		}
	})
}

func originValidator(*http.Request) bool {
	return true
}
