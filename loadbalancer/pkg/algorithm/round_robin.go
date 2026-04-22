package algorithm

import (
	"errors"
	"loadbalancer/pkg/server"
	"sync"
)

type RoundRobin struct {
	curServer int
	mu        sync.Mutex
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		curServer: -1,
	}
}

func (r *RoundRobin) Next(servers []server.Server) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(servers) == 0 {
		r.curServer = -1
		return -1, errors.New("List of servers is empty. Cannot select next server")
	}

	if r.curServer >= len(servers) {
		r.curServer = -1
	}

	tries := len(servers)
	r.curServer = (r.curServer + 1) % len(servers)
	tries -= 1
	for servers[r.curServer].GetHealth() == false {
		if tries <= 0 {
			return -1, errors.New("No healthy servers available")
		}
		r.curServer = (r.curServer + 1) % len(servers)
		tries -= 1
	}
	return r.curServer, nil
}
