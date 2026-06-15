package main

import (
	"errors"
	"testing"
)

type closeClientConnection struct {
	calls int
	err   error
}

func (c *closeClientConnection) Close() error {
	c.calls++
	return c.err
}

func TestCloseClientConnectionIgnoresNilConnection(t *testing.T) {
	closeClientConnectionSafely(nil)
}

func TestCloseClientConnectionCallsConnectionOnce(t *testing.T) {
	conn := &closeClientConnection{}

	closeClientConnectionSafely(conn)

	if conn.calls != 1 {
		t.Fatalf("expected one close call, got %d", conn.calls)
	}
}

func TestCloseClientConnectionHandlesCloseError(t *testing.T) {
	conn := &closeClientConnection{err: errors.New("close failed")}

	closeClientConnectionSafely(conn)

	if conn.calls != 1 {
		t.Fatalf("expected one close call, got %d", conn.calls)
	}
}
