package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/evan-forbes/ethlab/thereum"
	"github.com/gorilla/websocket"
)

//
// func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {

// 	// unmarshal the rpc request
// 	var req rpcMessage
// 	err = json.Unmarshal(body, &req)
// 	if err != nil {
// 		// send an error back if rpcMessage cannot be unmarshaled
// 		w.Write(rpcError(500, fmt.Sprintf("could not unmarshal rpc message: %s", err)))
// 		return
// 	}
// 	fmt.Println("method", req.Method)
// 	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 	}
// 	rawPrms := string(req.Params)
// 	switch {
// 	case strings.Contains(rawPrms, "logs"):
// 		ctx, _ := context.WithTimeout(s.ctx, time.Hour)
// 		subLogs(ctx, s.back, conn, json.RawMessage(rawPrms))

// 	case strings.Contains(rawPrms, "newHeads"):
// 		ctx, _ := context.WithTimeout(s.ctx, time.Hour)
// 		subHeads(ctx, s.back, conn, json.RawMessage(rawPrms))
// 	default:
// 		w.Write(rpcError(500, fmt.Sprintf("no subscription for %s", rawPrms)))
// 		log.Println(err)
// 	}
// }

////////////////////////////////
//	Streaming Heads
//////////////////////////////

func subHeads(ctx context.Context, eth *thereum.Thereum, conn *websocket.Conn) {
	sink := make(chan *types.Header)
	sub := eth.Events.SubscribeNewHeads(sink)
	feedHeads(ctx, conn, sub, sink)
}

func feedHeads(ctx context.Context, conn *websocket.Conn, sub *filters.Subscription, heads <-chan *types.Header) {
	defer conn.Close()
	defer sub.Unsubscribe()
	for {
		select {
		case <-sub.Err():
			return
		case <-ctx.Done():
			return
		case head := <-heads:
			result := headPacket{
				Subscription: string(sub.ID),
				Result:       head,
			}
			// marshal the log in the expected formate
			resultBytes, err := json.Marshal(result)
			if err != nil {
				log.Println("failed to marshal head during streaming")
				return
			}
			// create the rpc msg using the marshaled log
			msg := *&rpcMessage{
				Version: "2.0",
				ID:      1,
				Method:  "eth_subscribe",
				Params:  resultBytes,
			}
			// write over the websocket
			err = conn.WriteJSON(msg)
			if err != nil {
				log.Println("failed to marshal head during streaming")
				return
			}
		}
	}
}

type headPacket struct {
	Subscription string        `json:"subscription"`
	Result       *types.Header `json:"result"`
}

////////////////////////////////
//	Streaming Logs
//////////////////////////////

// subLogs is the procedure to stream logs via websocket
func subLogs(ctx context.Context, eth *thereum.Thereum, conn *websocket.Conn, rawPrms json.RawMessage) error {
	// attempt to unmarshal in
	params := []interface{}{}
	var method string
	var query query
	params = append(params, method, query)
	err := json.Unmarshal(rawPrms, &params)
	if err != nil {
		return err
	}
	// subscribe via the backend's EventSystem
	sink := make(chan []*types.Log)
	sub, err := eth.Events.SubscribeLogs(query.FilterQuery(), sink)
	if err != nil {
		return err
	}
	// Write the logs to the connection
	go feedLogs(ctx, conn, sub, sink)
	return nil
}

func feedLogs(ctx context.Context, conn *websocket.Conn, sub *filters.Subscription, logs <-chan []*types.Log) {
	defer conn.Close()
	defer sub.Unsubscribe()
	for {
		select {
		case <-sub.Err():
			return
		case <-ctx.Done():
			return
		case ls := <-logs:
			for _, l := range ls {
				result := logPacket{
					Subscription: string(sub.ID),
					Result:       l,
				}
				resultBytes, err := json.Marshal(result)
				if err != nil {
					log.Println("failed to marshal log during streaming")
					return
				}
				msg := *&rpcMessage{
					Version: "2.0",
					ID:      1,
					Method:  "eth_subscribe",
					Params:  resultBytes,
				}
				err = conn.WriteJSON(msg)
				if err != nil {
					log.Println("failed to marshal log during streaming")
					return
				}
			}
		}
	}
}

type logPacket struct {
	Subscription string     `json:"subscription"`
	Result       *types.Log `json:"result"`
}

// query helps unmarshall a json filter query
// TODO: ensure that any format of query sent via rpc works
type query struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	From    string   `json:"from_block"`
	To      string   `json:"to_block"`
	Block   string   `json:"blockHash"`
}

// FilterQuery converts query into a standardized ethereum.FilterQuery
func (q *query) FilterQuery() ethereum.FilterQuery {
	var topics []common.Hash
	for _, t := range q.Topics {
		topics = append(topics, common.HexToHash(t))
	}
	return ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(q.Address)},
		Topics:    [][]common.Hash{topics},
	}
}
