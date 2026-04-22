package server

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialization(t *testing.T) {
	// Test: expected default values
	server := NewServer("https://localhost:12345")
	assert.Equal(t, true, server.healthy)
	assert.Equal(t, "https://localhost:12345", server.url)
	assert.Equal(t, 1, server.weight)
	assert.Equal(t, int32(0), server.Connections())
}

func TestSetters(t *testing.T) {
	// Test: SetCPUUsage valid value
	server := NewServer("https://localhost:12345")
	err := server.SetCPUUsage(0.69)
	require.NoError(t, err)
	assert.Equal(t, 0.69, server.stats.cpuUsage)

	// Test: SetCPUUsage negative cpuusage
	err = server.SetCPUUsage(-0.420)
	require.Error(t, err)
	assert.Equal(t, 0.69, server.stats.cpuUsage)

	// Test: SetHealthy
	server.SetHealth(false)
	assert.False(t, server.healthy)

	// Test: SetUrl
	server.SetURL("https://instagram.com")
	assert.Equal(t, "https://instagram.com", server.url)

	// Test: SetWeight valid value
	err = server.SetWeight(69)
	require.NoError(t, err)
	assert.Equal(t, 69, server.weight)

	// Test: SetWeight value 0
	err = server.SetWeight(0)
	require.Error(t, err)
	assert.Equal(t, 69, server.weight)

	// Test: SetWeight value negative
	err = server.SetWeight(-420)
	require.Error(t, err)
	assert.Equal(t, 69, server.weight)

	// Test: AddConections
	server.AddConnection()
	assert.Equal(t, int32(1), server.Connections())
	server.AddConnection()
	assert.Equal(t, int32(2), server.Connections())

	// Test: RemoveConnections
	server.RemoveConnection()
	assert.Equal(t, int32(1), server.Connections())
	server.RemoveConnection()
	assert.Equal(t, int32(0), server.Connections())
}

func TestConcurrency(t *testing.T) {
	server := NewServer("https://localhost:12345")
	var wg sync.WaitGroup
	iterations := 1000

	// Concurrent connections management
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			server.AddConnection()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			server.RemoveConnection()
		}
	}()

	// Concurrent attribute updates
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			server.SetHealth(i%2 == 0)
			_ = server.SetWeight((i % 10) + 1)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = server.SetCPUUsage(float64(i) / 1000.0)
			server.SetURL("https://updated-url.com")
		}
	}()

	// Concurrent reads
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = server.GetHealth()
			_ = server.GetWeight()
			_ = server.GetCPUUsage()
			_ = server.GetURL()
			_ = server.Connections()
		}
	}()

	wg.Wait()

	// Final consistency checks
	assert.Equal(t, int32(0), server.Connections(), "Connections should be balanced back to 0")
}
