package algorithm

import (
	"errors"
	"loadbalancer/pkg/server"
	"sync"
)

type WeightedRoundRobin struct {
	curServer int
	repeat    int
	mu        sync.Mutex
}

func NewWeightedRoundRobin() *WeightedRoundRobin {
	return &WeightedRoundRobin{
		curServer: -1,
	}
}

func (r *WeightedRoundRobin) Next(servers []server.Server) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(servers) == 0 {
		r.curServer = -1
		r.repeat = 0
		return -1, errors.New("List of servers is empty. Cannot select next server")
	}

	if r.curServer >= len(servers) {
		r.curServer = -1
		r.repeat = 0
	}

	tries := len(servers)
	for r.repeat <= 0 || !servers[r.curServer].GetHealth() {
		if tries <= 0 {
			return -1, errors.New("Couldn't find a healthy server")
		}
		r.curServer = (r.curServer + 1) % len(servers)
		r.repeat = servers[r.curServer].GetWeight()
		tries -= 1
	}
	r.repeat -= 1

	return r.curServer, nil
}
