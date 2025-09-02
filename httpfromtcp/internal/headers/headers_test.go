package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	v, ok := headers.Get("HOSt")
	assert.Equal(t, "localhost:42069", v)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single headers with extra whitespace
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	v, ok = headers.Get("hoST")
	assert.Equal(t, "localhost:42069", v)
	assert.Equal(t, 37, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	data = []byte("Host: localhost\r\nPort: :42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	v, ok = headers.Get("hoST")
	assert.Equal(t, "localhost", v)
	v, ok = headers.Get("Port")
	assert.False(t, ok)
	assert.Equal(t, 17, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	v, ok = headers.Get("hosT")
	assert.Equal(t, "localhost", v)
	v, ok = headers.Get("pOrT")
	assert.Equal(t, ":42069", v)
	assert.Equal(t, 14, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go\r\nSet-Person: prime-loves-zig\r\nSet-Person: tj-loves-ocaml\r\n\r\n")
	cur := 0
	n, done, err = headers.Parse(data)
	cur += n
	n, done, err = headers.Parse(data[cur:])
	cur += n
	n, done, err = headers.Parse(data[cur:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	v, ok = headers.Get("set-PeRson")
	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", v)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
