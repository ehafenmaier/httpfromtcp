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
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.True(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("    Host: localhost:12345    \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:12345", headers["host"])
	assert.Equal(t, 31, n)
	assert.True(t, done)

	// Test: Two valid headers
	headers = NewHeaders()
	data = []byte("    Host: localhost:12345    \r\n  User-Agent:  test-agent   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:12345", headers["host"])
	assert.Equal(t, "test-agent", headers["user-agent"])
	assert.Equal(t, 61, n)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host  : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid characters in header field name
	headers = NewHeaders()
	data = []byte("    H@st: localhost:42069    \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
