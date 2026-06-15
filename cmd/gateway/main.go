package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/naren-m/panchangam/cache"
	"github.com/naren-m/panchangam/gateway"
	"github.com/naren-m/panchangam/log"
)

var logger = log.Logger()

func main() {
	// Command line flags
	var (
		grpcEndpoint = flag.String("grpc-endpoint", "localhost:50051", "gRPC server endpoint")
		httpPort     = flag.String("http-port", "8080", "HTTP server port")
		logLevel     = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
		enableCache  = flag.Bool("enable-cache", true, "Enable Redis caching")
		redisAddr    = flag.String("redis-addr", "localhost:6379", "Redis server address")
		redisDB      = flag.Int("redis-db", 0, "Redis database number")
		cacheTTL     = flag.Duration("cache-ttl", 30*time.Minute, "Cache TTL duration")
		healthCheck  = flag.Bool("health-check", false, "Run a gateway health check and exit")
	)
	flag.Parse()

	if *healthCheck {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := runGatewayHealthCheck(ctx, gatewayHealthAddress(*httpPort)); err != nil {
			fmt.Fprintf(os.Stderr, "gateway health check failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	cacheAddr, cacheDB, ttl, cacheEnabled := applyCacheEnvOverrides(*redisAddr, *redisDB, *cacheTTL, *enableCache)

	// Set log level
	// Note: This would typically be implemented in the log package
	logger.Info("Starting Panchangam HTTP Gateway",
		"grpc_endpoint", *grpcEndpoint,
		"http_port", *httpPort,
		"log_level", *logLevel,
	)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create wait group for goroutines
	var wg sync.WaitGroup

	// Initialize Redis cache if enabled
	var redisCache *cache.RedisCache
	if cacheEnabled {
		logger.Info("Initializing Redis cache", "addr", cacheAddr, "db", cacheDB, "ttl", ttl)

		redisPassword := os.Getenv("REDIS_PASSWORD")
		var err error
		redisCache, err = cache.NewRedisCache(cacheAddr, redisPassword, cacheDB, ttl)
		if err != nil {
			logger.Error("Failed to initialize Redis cache, continuing without cache", "error", err)
			redisCache = nil
		} else {
			logger.Info("Redis cache initialized successfully")
		}
	} else {
		logger.Info("Cache disabled")
	}

	// Create gateway server
	var gatewayServer *gateway.GatewayServer
	if redisCache != nil {
		gatewayServer = gateway.NewGatewayServerWithCache(*grpcEndpoint, *httpPort, redisCache)
	} else {
		gatewayServer = gateway.NewGatewayServer(*grpcEndpoint, *httpPort)
	}

	// Start gateway server in a goroutine
	gatewayErr := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := gatewayServer.Start(ctx); err != nil {
			gatewayErr <- err
		}
	}()

	logger.Info("HTTP Gateway server started successfully",
		"address", fmt.Sprintf("http://localhost:%s", *httpPort),
		"cache_enabled", redisCache != nil,
		"endpoints", startupEndpoints(*httpPort, redisCache != nil),
	)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var sig os.Signal
	select {
	case sig = <-sigChan:
	case err := <-gatewayErr:
		logger.Error("Gateway server error", "error", err)
		cancel()
		if redisCache != nil {
			if closeErr := redisCache.Close(); closeErr != nil {
				logger.Error("Error closing Redis cache", "error", closeErr)
			}
		}
		os.Exit(1)
	}
	logger.Info("Received shutdown signal", "signal", sig)

	// Cancel context to signal shutdown
	cancel()

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown gateway server
	logger.Info("Shutting down gateway server")
	if err := gatewayServer.Stop(shutdownCtx); err != nil {
		logger.Error("Error during gateway shutdown", "error", err)
	}

	// Close Redis cache if initialized
	if redisCache != nil {
		logger.Info("Closing Redis cache connection")
		if err := redisCache.Close(); err != nil {
			logger.Error("Error closing Redis cache", "error", err)
		}
	}

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for shutdown to complete or timeout
	select {
	case <-done:
		logger.Info("Gateway shutdown completed successfully")
	case <-shutdownCtx.Done():
		logger.Warn("Gateway shutdown timed out")
	}
}

func startupEndpoints(httpPort string, cacheEnabled bool) []string {
	endpoints := []string{
		fmt.Sprintf("http://localhost:%s/api/v1/health", httpPort),
		fmt.Sprintf("http://localhost:%s/api/v1/panchangam", httpPort),
	}

	if cacheEnabled {
		endpoints = append(endpoints,
			fmt.Sprintf("http://localhost:%s/api/v1/cache/health", httpPort),
			fmt.Sprintf("http://localhost:%s/api/v1/cache/stats", httpPort),
		)
	}

	return endpoints
}

func applyCacheEnvOverrides(redisAddr string, redisDB int, cacheTTL time.Duration, enableCache bool) (string, int, time.Duration, bool) {
	if env := strings.TrimSpace(os.Getenv("REDIS_ADDR")); env != "" {
		redisAddr = env
	}
	if env := strings.TrimSpace(os.Getenv("REDIS_DB")); env != "" {
		if db, err := strconv.Atoi(env); err == nil {
			redisDB = db
		} else {
			logger.Warn("Ignoring invalid REDIS_DB", "value", env, "error", err)
		}
	}
	if env := strings.TrimSpace(os.Getenv("CACHE_TTL")); env != "" {
		if ttl, err := time.ParseDuration(env); err == nil {
			cacheTTL = ttl
		} else {
			logger.Warn("Ignoring invalid CACHE_TTL", "value", env, "error", err)
		}
	}
	if env := strings.TrimSpace(os.Getenv("ENABLE_CACHE")); env != "" {
		switch strings.ToLower(env) {
		case "true", "1":
			enableCache = true
		case "false", "0":
			enableCache = false
		default:
			logger.Warn("Ignoring invalid ENABLE_CACHE", "value", env)
		}
	}
	return redisAddr, redisDB, cacheTTL, enableCache
}

func gatewayHealthAddress(httpPort string) string {
	if strings.Contains(httpPort, ":") {
		return httpPort
	}
	return "localhost:" + httpPort
}

func runGatewayHealthCheck(ctx context.Context, address string) error {
	url := fmt.Sprintf("http://%s/api/v1/health", address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build gateway health request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("call gateway health endpoint %s: %w", url, err)
	}
	defer func() {
		_ = resp.Body.Close() // best-effort cleanup after a small health response
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("gateway health endpoint %s returned %s", url, resp.Status)
	}

	return nil
}
