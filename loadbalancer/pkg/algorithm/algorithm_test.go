package algorithm

import (
	"loadbalancer/pkg/server"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoundRobin(t *testing.T) {
	rr := NewRoundRobin()
	serverList := make([]server.Server, 0)
	// Test: No servers in server list
	i, err := rr.Next(serverList)
	require.Error(t, err)
	assert.Equal(t, -1, i)
	assert.Equal(t, -1, rr.curServer)

	// Test: all servers are healthy
	rr = NewRoundRobin()
	serverList = append(serverList, server.NewServer("localhost:1"))
	serverList = append(serverList, server.NewServer("localhost:2"))
	serverList = append(serverList, server.NewServer("localhost:3"))
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 1, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 2, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)

	// Test: One of the servers unhealthy
	rr = NewRoundRobin()
	serverList[1].SetHealth(false)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 2, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)

	// Test: No healthy servers
	rr = NewRoundRobin()
	serverList[0].SetHealth(false)
	serverList[2].SetHealth(false)
	i, err = rr.Next(serverList)
	require.Error(t, err)

	// Test: Server list shrunk
	rr = NewRoundRobin()
	serverList[0].SetHealth(true)
	serverList[1].SetHealth(true)
	serverList[2].SetHealth(true)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 1, i)
	serverList = slices.Delete(serverList, 1, 3)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)

}

func TestWeightedRoundRobin(t *testing.T) {
	rr := NewWeightedRoundRobin()
	serverList := make([]server.Server, 0)
	// Test: No servers in server list
	i, err := rr.Next(serverList)
	require.Error(t, err)
	assert.Equal(t, -1, i)
	assert.Equal(t, -1, rr.curServer)

	// Test: all servers are healthy and weight 1
	rr = NewWeightedRoundRobin()
	serverList = append(serverList, server.NewServer("localhost:1"))
	serverList = append(serverList, server.NewServer("localhost:2"))
	serverList = append(serverList, server.NewServer("localhost:3"))
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 1, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 2, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)

	// Test: One of the servers unhealthy
	rr = NewWeightedRoundRobin()
	serverList[1].SetHealth(false)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 2, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)

	// Test: No healthy servers
	rr = NewWeightedRoundRobin()
	serverList[0].SetHealth(false)
	serverList[2].SetHealth(false)
	i, err = rr.Next(serverList)
	require.Error(t, err)

	// Test: all servers are healthy and variable weights
	rr = NewWeightedRoundRobin()
	serverList[0].SetHealth(true)
	serverList[1].SetHealth(true)
	serverList[2].SetHealth(true)
	serverList[0].SetWeight(3)
	serverList[1].SetWeight(2)
	for range 3 {
		i, err = rr.Next(serverList)
		require.NoError(t, err)
		assert.Equal(t, 0, i)
	}

	for range 2 {
		i, err = rr.Next(serverList)
		require.NoError(t, err)
		assert.Equal(t, 1, i)
	}

	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 2, i)

	// Test: Server to be repeated becomes unhealthy
	rr = NewWeightedRoundRobin()
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)
	serverList[0].SetHealth(false)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 1, i)

	// Test: Server list shrunk
	rr = NewWeightedRoundRobin()
	serverList[0].SetHealth(true)
	serverList[0].SetWeight(1)
	serverList[1].SetWeight(5)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 1, i)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 1, i)
	serverList = slices.Delete(serverList, 1, 3)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)
}

func TestLeastConnections(t *testing.T) {
	rr := NewLeastConnections()
	serverList := make([]server.Server, 0)
	// Test: No servers in server list
	i, err := rr.Next(serverList)
	require.Error(t, err)
	assert.Equal(t, -1, i)

	// Test: all servers are healthy arbitrary connections
	rr = NewLeastConnections()
	serverList = append(serverList, server.NewServer("localhost:1"))
	serverList = append(serverList, server.NewServer("localhost:2"))
	serverList = append(serverList, server.NewServer("localhost:3"))
	for range 1 {
		serverList[1].AddConnection()
	}
	for range 3 {
		serverList[2].AddConnection()
	}
	for range 5 {
		serverList[0].AddConnection()
	}
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 1, i)
	for range 3 {
		serverList[1].AddConnection()
	}
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 2, i)

	for range 5 {
		serverList[1].AddConnection()
		serverList[2].AddConnection()
	}

	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)

	// Test: One of the servers unhealthy
	rr = NewLeastConnections()
	serverList[0].SetHealth(false)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 2, i)

	// Test: No healthy servers
	rr = NewLeastConnections()
	serverList[1].SetHealth(false)
	serverList[2].SetHealth(false)
	i, err = rr.Next(serverList)
	require.Error(t, err)
}

func TestResourceBased(t *testing.T) {
	rr := NewResourceBased()
	serverList := make([]server.Server, 0)
	// Test: No servers in server list
	i, err := rr.Next(serverList)
	require.Error(t, err)
	assert.Equal(t, -1, i)

	// Test: all servers are healthy arbitrary connections
	rr = NewResourceBased()
	serverList = append(serverList, server.NewServer("localhost:1"))
	serverList = append(serverList, server.NewServer("localhost:2"))
	serverList = append(serverList, server.NewServer("localhost:3"))
	serverList[1].SetCPUUsage(0.1)
	serverList[2].SetCPUUsage(0.4)
	serverList[0].SetCPUUsage(0.7)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 1, i)
	serverList[1].SetCPUUsage(0.5)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 2, i)

	serverList[1].SetCPUUsage(0.9)
	serverList[2].SetCPUUsage(0.8)

	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 0, i)

	// Test: One of the servers unhealthy
	rr = NewResourceBased()
	serverList[0].SetHealth(false)
	i, err = rr.Next(serverList)
	require.NoError(t, err)
	assert.Equal(t, 2, i)

	// Test: No healthy servers
	rr = NewResourceBased()
	serverList[1].SetHealth(false)
	serverList[2].SetHealth(false)
	i, err = rr.Next(serverList)
	require.Error(t, err)
}
