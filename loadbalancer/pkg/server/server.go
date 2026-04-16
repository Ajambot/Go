package server

import (
	"sync"
	"sync/atomic"
)

type Stats struct {
	CPUUsage    float64 `json:"CPUUsage"`
	Connections atomic.Int32
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

func (s *Server) AddConnection() {
	s.Stats.Connections.Add(1)
}

func (s *Server) RemoveConnection() {
	s.Stats.Connections.Add(-1)
}

func (s *Server) Connections() int32 {
	return s.Stats.Connections.Load()
}
