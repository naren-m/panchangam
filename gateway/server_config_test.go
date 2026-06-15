package gateway

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORSConfiguration(t *testing.T) {
	server := NewGatewayServer("localhost:50052", "8080")
	assert.NotNil(t, server)
	assert.Equal(t, "localhost:50052", server.grpcEndpoint)
	assert.Equal(t, "8080", server.httpPort)
}

func TestCORSOptionsDisableDebugLogging(t *testing.T) {
	allowedOrigins := []string{"http://localhost:5173"}

	options := newCORSOptions(allowedOrigins)

	assert.Equal(t, allowedOrigins, options.AllowedOrigins)
	assert.Equal(t, []string{http.MethodGet, http.MethodPost, http.MethodOptions}, options.AllowedMethods)
	assert.Equal(t, []string{"*"}, options.AllowedHeaders)
	assert.False(t, options.Debug)
}

type closeConnection struct {
	calls int
	err   error
}

func (c *closeConnection) Close() error {
	c.calls++
	return c.err
}

func TestCloseGRPCConnectionIgnoresNilConnection(t *testing.T) {
	closeGRPCConnectionSafely(nil)
}

func TestCloseGRPCConnectionCallsConnectionOnce(t *testing.T) {
	conn := &closeConnection{}

	closeGRPCConnectionSafely(conn)

	assert.Equal(t, 1, conn.calls)
}

func TestCloseGRPCConnectionHandlesCloseError(t *testing.T) {
	conn := &closeConnection{err: errors.New("close failed")}

	closeGRPCConnectionSafely(conn)

	assert.Equal(t, 1, conn.calls)
}
