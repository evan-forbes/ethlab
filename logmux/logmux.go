package logmux

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
)

type Muxer struct {
	Handlers map[string]Handler
}

// Mux routes a queue of logs to their respective handlers, sorted by Topic
func (m *Muxer) Mux(ctx context.Context, wg *sync.WaitGroup, logs <-chan *types.Log) {
	defer wg.Done()
	for l := range logs {
		handler, has := m.Handlers[l.Topics[0].Hex()]
		if !has {
			continue // ignore log
		}
		go handler(ctx, l)
	}
}

func (m *Muxer) Merge(secondary *Muxer) {
	for topic, handlr := range m.Handlers {
		m.Handlers[topic] = handlr
	}
}

// Handler describes some function that acts upon a log
type Handler func(context.Context, *types.Log)
