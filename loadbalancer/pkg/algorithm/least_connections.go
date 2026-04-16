package algorithm

import (
	"errors"
	"loadbalancer/pkg/server"
	"log"
	"slices"
)

type LeastConnections struct {
}

type ConnectionsServerInfo struct {
	Connections int
	Healthy     bool
	Index       int
}

func NewLeastConnections() *LeastConnections {
	return &LeastConnections{}
}

func (r *LeastConnections) Next(servers []*server.Server) (int, error) {
	if len(servers) == 0 {
		return -1, errors.New("List of servers is empty. Cannot select next server")
	}
	server_stats := []ConnectionsServerInfo{}

	for i, server := range servers {
		server_stats = append(server_stats, ConnectionsServerInfo{Connections: int(server.Connections()), Index: i, Healthy: server.Healthy})
	}

	slices.SortFunc(server_stats, func(s1, s2 ConnectionsServerInfo) int {
		if s1.Connections == s2.Connections {
			return 0
		}
		if s1.Connections > s2.Connections {
			return 1
		}
		return -1
	})

	log.Println(server_stats)

	picked_server := 0
	for server_stats[picked_server].Healthy == false {
		picked_server += 1
	}

	return server_stats[picked_server].Index, nil
}
