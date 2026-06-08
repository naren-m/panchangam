package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestStartupEndpointsWithoutCache(t *testing.T) {
	endpoints := startupEndpoints("8080", false)

	expected := []string{
		"http://localhost:8080/api/v1/health",
		"http://localhost:8080/api/v1/panchangam",
	}
	if !reflect.DeepEqual(endpoints, expected) {
		t.Fatalf("expected endpoints %v, got %v", expected, endpoints)
	}
}

func TestStartupEndpointsWithCache(t *testing.T) {
	endpoints := startupEndpoints("8080", true)

	expected := []string{
		"http://localhost:8080/api/v1/health",
		"http://localhost:8080/api/v1/panchangam",
		"http://localhost:8080/api/v1/cache/health",
		"http://localhost:8080/api/v1/cache/stats",
	}
	if !reflect.DeepEqual(endpoints, expected) {
		t.Fatalf("expected endpoints %v, got %v", expected, endpoints)
	}
}

func TestRunGatewayHealthCheckAcceptsHealthyGateway(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on local port: %v", err)
	}
	defer func() {
		_ = listener.Close() // best-effort test cleanup
	}()

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/health" {
				t.Fatalf("expected /api/v1/health, got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}),
	}
	defer func() {
		_ = server.Close() // best-effort test cleanup
	}()

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			t.Errorf("serve health check endpoint: %v", err)
		}
	}()

	err = runGatewayHealthCheck(context.Background(), listener.Addr().String())
	if err != nil {
		t.Fatalf("expected healthy gateway, got %v", err)
	}
}

func TestRunGatewayHealthCheckRejectsUnhealthyGateway(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on local port: %v", err)
	}
	defer func() {
		_ = listener.Close() // best-effort test cleanup
	}()

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}),
	}
	defer func() {
		_ = server.Close() // best-effort test cleanup
	}()

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			t.Errorf("serve health check endpoint: %v", err)
		}
	}()

	err = runGatewayHealthCheck(context.Background(), listener.Addr().String())
	if err == nil {
		t.Fatal("expected unhealthy gateway to return an error")
	}
}

func TestApplyCacheEnvOverridesUsesValidValues(t *testing.T) {
	t.Setenv("REDIS_ADDR", " redis.example:6380 ")
	t.Setenv("REDIS_DB", "3")
	t.Setenv("CACHE_TTL", "45s")
	t.Setenv("ENABLE_CACHE", "1")

	addr, db, ttl, enabled := applyCacheEnvOverrides("localhost:6379", 0, 30*time.Minute, false)

	if addr != "redis.example:6380" {
		t.Fatalf("expected redis address override, got %q", addr)
	}
	if db != 3 {
		t.Fatalf("expected redis db override, got %d", db)
	}
	if ttl != 45*time.Second {
		t.Fatalf("expected cache ttl override, got %s", ttl)
	}
	if !enabled {
		t.Fatal("expected cache to be enabled")
	}
}

func TestApplyCacheEnvOverridesDisablesCache(t *testing.T) {
	t.Setenv("ENABLE_CACHE", "false")

	_, _, _, enabled := applyCacheEnvOverrides("localhost:6379", 0, 30*time.Minute, true)

	if enabled {
		t.Fatal("expected ENABLE_CACHE=false to disable cache")
	}
}

func TestApplyCacheEnvOverridesKeepsDefaultsForInvalidValues(t *testing.T) {
	t.Setenv("REDIS_DB", "bad")
	t.Setenv("CACHE_TTL", "bad")
	t.Setenv("ENABLE_CACHE", "bad")

	addr, db, ttl, enabled := applyCacheEnvOverrides("localhost:6379", 0, 30*time.Minute, true)

	if addr != "localhost:6379" {
		t.Fatalf("expected default redis address, got %q", addr)
	}
	if db != 0 {
		t.Fatalf("expected default redis db, got %d", db)
	}
	if ttl != 30*time.Minute {
		t.Fatalf("expected default cache ttl, got %s", ttl)
	}
	if !enabled {
		t.Fatal("expected invalid ENABLE_CACHE to keep the default")
	}
}

func TestGatewayStartupDoesNotUseNoopRedisPasswordOverride(t *testing.T) {
	source := readGatewayMainSource(t)

	if strings.Contains(source, `if env := os.Getenv("REDIS_PASSWORD"); env != ""`) {
		t.Fatal("gateway startup should not keep a REDIS_PASSWORD env block that does no work")
	}
}

func TestGatewayStartupReportsInvalidCacheEnv(t *testing.T) {
	source := readGatewayMainSource(t)

	for _, want := range []string{
		`logger.Warn("Ignoring invalid REDIS_DB"`,
		`logger.Warn("Ignoring invalid CACHE_TTL"`,
		`logger.Warn("Ignoring invalid ENABLE_CACHE"`,
	} {
		if !strings.Contains(source, want) {
			t.Fatalf("gateway startup should log %s", want)
		}
	}
}

func TestGatewayStartupExitsOnServerStartError(t *testing.T) {
	source := readGatewayMainSource(t)

	for _, want := range []string{
		"gatewayErr := make(chan error, 1)",
		"gatewayErr <- err",
		"case err := <-gatewayErr:",
		`logger.Error("Gateway server error", "error", err)`,
	} {
		if !strings.Contains(source, want) {
			t.Fatalf("gateway startup should handle server errors through %q", want)
		}
	}
}

func readGatewayMainSource(t *testing.T) string {
	t.Helper()

	source, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}

	return string(source)
}
