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

func (s *Server) RunWSServer(endpoint string) error {
	var (
		listener net.Listener
		err      error
	)
	if listener, err = net.Listen("tcp", endpoint); err != nil {
		return err
	}
	srv := http.Server{Handler: s.WebsocketHandler()}
	go srv.Serve(listener)
	return nil
}

func (s *Server) WebsocketHandler() http.Handler {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		WriteBufferPool: wsPool,
		CheckOrigin:     originValidator,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// log.Debug("WebSocket upgrade failed", "err", err)
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
			fmt.Println("attempting to server new heads")
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
