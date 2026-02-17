package algorithm

import (
	"errors"
	"loadbalancer/pkg/server"
)

type RoundRobin struct {
	curServer int
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		curServer: -1,
	}
}

func (r *RoundRobin) Next(servers []server.Server) (int, error) {
	if len(servers) == 0 {
		r.curServer = -1
		return -1, errors.New("List of servers is empty. Cannot select next server")
	}
	r.curServer = (r.curServer + 1) % len(servers)
	return r.curServer, nil
}
