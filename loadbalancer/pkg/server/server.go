package server

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Stats struct {
	cpuUsage    float64
	connections atomic.Int32
}

type Server interface {
	SetHealth(healthy bool)
	GetHealth() bool
	SetWeight(w int) error
	GetWeight() int
	AddConnection()
	RemoveConnection()
	Connections() int32
	SetCPUUsage(usage float64) error
	GetCPUUsage() float64
	SetURL(url string)
	GetURL() string
}

type server struct {
	mu      sync.Mutex
	stats   Stats
	healthy bool
	url     string
	weight  int
}

func NewServer(url string) *server {
	return &server{healthy: true, url: url, weight: 1}
}

func (s *server) GetURL() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.url
}

func (s *server) SetURL(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.url = url
}

func (s *server) SetHealth(healthy bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.healthy = healthy
}

func (s *server) GetHealth() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.healthy
}

func (s *server) SetWeight(w int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if w <= 0 {
		return errors.New("Weight has to be >= 1")
	}
	s.weight = w
	return nil
}

func (s *server) GetWeight() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.weight
}

func (s *server) AddConnection() {
	s.stats.connections.Add(1)
}

func (s *server) RemoveConnection() {
	s.stats.connections.Add(-1)
}

func (s *server) Connections() int32 {
	return s.stats.connections.Load()
}

func (s *server) SetCPUUsage(usage float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if usage < 0 {
		return errors.New("CPUUsage cannot be a negative number")
	}
	s.stats.cpuUsage = usage
	return nil
}

func (s *server) GetCPUUsage() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stats.cpuUsage
}
