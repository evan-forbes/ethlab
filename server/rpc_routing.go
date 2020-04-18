package server

import (
	"sync"

	"github.com/evan-forbes/ethlab/thereum"
)

// muxer maps supported methods to their appropriate procedures.
// thread safe.
type muxer struct {
	routes map[string]procedure
	mut    sync.RWMutex
}

type procedure func(*thereum.Thereum, rpcMessage) (rpcMessage, error)

func newMuxer(s *Server) *muxer {
	return &muxer{
		routes: map[string]procedure{
			// add rpc methods here!
			"":        nullProcedure,
			"getLogs": getLogs,
		},
	}
}

func (m *muxer) Route(method string) (procedure, bool) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	pro, has := m.routes[method]
	return pro, has
}

func nullProcedure(*rpcMessage) (*rpcMessage, error) {
	nullMessage := rpcMessage{
		Version: "2.0",
		ID:      60,
		Error: &jsonError{
			Code:    999,
			Message: "specified method is not registered or supported",
		},
	}
	return &nullMessage, nil
}
