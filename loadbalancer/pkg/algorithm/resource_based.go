package algorithm

import (
	"errors"
	"loadbalancer/pkg/server"
	"slices"
)

type ResourceBased struct {
}

type ResourceServerInfo struct {
	CPUUsage float64
	Healthy  bool
	Index    int
}

func NewResourceBased() *ResourceBased {
	return &ResourceBased{}
}

func (r *ResourceBased) Next(servers []server.Server) (int, error) {
	if len(servers) == 0 {
		return -1, errors.New("List of servers is empty. Cannot select next server")
	}
	server_stats := []ResourceServerInfo{}
	// [float64, int]

	for i, server := range servers {
		server_stats = append(server_stats, ResourceServerInfo{CPUUsage: server.GetCPUUsage(), Index: i, Healthy: server.GetHealth()})
	}

	slices.SortFunc(server_stats, func(s1, s2 ResourceServerInfo) int {
		if s1.CPUUsage == s2.CPUUsage {
			return 0
		}
		if s1.CPUUsage > s2.CPUUsage {
			return 1
		}
		return -1
	})

	picked_server := 0
	for server_stats[picked_server].Healthy == false {
		picked_server += 1
		if picked_server >= len(server_stats) {
			return -1, errors.New("No healthy servers")
		}
	}

	return server_stats[picked_server].Index, nil
}
