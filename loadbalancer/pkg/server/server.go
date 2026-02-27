package server

import "sync"

type Stats struct {
	CPUUsage float64 `json:"CPUUsage"`
}

type Server struct {
	mu      sync.Mutex
	Stats   Stats
	Healthy bool
	Url     string
	Weight  int
}

func (s *Server) SetHealthy(healthy bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Healthy = healthy
}
