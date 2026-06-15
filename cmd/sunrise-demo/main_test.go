package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type closeConnection struct {
	calls int
	err   error
}

func (c *closeConnection) Close() error {
	c.calls++
	return c.err
}

func TestCloseConnectionIgnoresNilConnection(t *testing.T) {
	closeConnectionSafely(nil)
}

func TestCloseConnectionCallsConnectionOnce(t *testing.T) {
	conn := &closeConnection{}

	closeConnectionSafely(conn)

	if conn.calls != 1 {
		t.Fatalf("expected one close call, got %d", conn.calls)
	}
}

func TestCloseConnectionHandlesCloseError(t *testing.T) {
	conn := &closeConnection{err: errors.New("close failed")}

	closeConnectionSafely(conn)

	if conn.calls != 1 {
		t.Fatalf("expected one close call, got %d", conn.calls)
	}
}

func TestWriteSunriseDemoError(t *testing.T) {
	var out bytes.Buffer

	if err := writeSunriseDemoError(&out, errors.New("server unavailable")); err != nil {
		t.Fatalf("writeSunriseDemoError returned error: %v", err)
	}

	want := "server unavailable\n"
	if got := out.String(); got != want {
		t.Fatalf("writeSunriseDemoError() output = %q, want %q", got, want)
	}
}

func TestResolveSunriseDemoLocationUsesPreset(t *testing.T) {
	name, lat, lon, timezone, err := resolveSunriseDemoLocation(" mumbai ", 1, 2, "UTC")
	if err != nil {
		t.Fatalf("resolveSunriseDemoLocation returned error: %v", err)
	}
	if name != "Mumbai" {
		t.Fatalf("expected Mumbai preset name, got %q", name)
	}
	if lat != 19.0760 || lon != 72.8777 || timezone != "Asia/Kolkata" {
		t.Fatalf("expected Mumbai preset, got lat=%f lon=%f timezone=%s", lat, lon, timezone)
	}
}

func TestResolveSunriseDemoLocationPreservesCustomValues(t *testing.T) {
	name, lat, lon, timezone, err := resolveSunriseDemoLocation("", 12.3, 45.6, "UTC")
	if err != nil {
		t.Fatalf("resolveSunriseDemoLocation returned error: %v", err)
	}
	if name != "" {
		t.Fatalf("expected no preset name, got %q", name)
	}
	if lat != 12.3 || lon != 45.6 || timezone != "UTC" {
		t.Fatalf("expected custom values, got lat=%f lon=%f timezone=%s", lat, lon, timezone)
	}
}

func TestResolveSunriseDemoLocationReportsAvailablePresets(t *testing.T) {
	_, _, _, _, err := resolveSunriseDemoLocation("unknown", 1, 2, "UTC")
	if err == nil {
		t.Fatal("expected unknown location to return an error")
	}
	if !strings.Contains(err.Error(), "Available: nyc, london, tokyo, sydney, mumbai, capetown") {
		t.Fatalf("expected available preset list in error, got %v", err)
	}
}
