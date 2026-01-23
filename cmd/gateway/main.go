package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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
	)
	flag.Parse()

	// Override flags with environment variables if present
	if env := os.Getenv("REDIS_ADDR"); env != "" {
		*redisAddr = env
	}
	if env := os.Getenv("REDIS_DB"); env != "" {
		if db, err := strconv.Atoi(env); err == nil {
			*redisDB = db
		}
	}
	if env := os.Getenv("REDIS_PASSWORD"); env != "" {
		// Redis password will be used in cache initialization
	}
	if env := os.Getenv("CACHE_TTL"); env != "" {
		if ttl, err := time.ParseDuration(env); err == nil {
			*cacheTTL = ttl
		}
	}
	if env := os.Getenv("ENABLE_CACHE"); env != "" {
		*enableCache = env == "true" || env == "1"
	}

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
	if *enableCache {
		logger.Info("Initializing Redis cache", "addr", *redisAddr, "db", *redisDB, "ttl", *cacheTTL)
		
		redisPassword := os.Getenv("REDIS_PASSWORD")
		var err error
		redisCache, err = cache.NewRedisCache(*redisAddr, redisPassword, *redisDB, *cacheTTL)
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
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := gatewayServer.Start(ctx); err != nil {
			logger.Error("Gateway server error", "error", err)
		}
	}()

	// Wait a moment for server to start
	time.Sleep(100 * time.Millisecond)
	endpoints := []string{
		fmt.Sprintf("http://localhost:%s/api/v1/health", *httpPort),
		fmt.Sprintf("http://localhost:%s/api/v1/panchangam", *httpPort),
	}
	
	if redisCache != nil {
		endpoints = append(endpoints,
			fmt.Sprintf("http://localhost:%s/api/v1/cache/health", *httpPort),
			fmt.Sprintf("http://localhost:%s/api/v1/cache/stats", *httpPort),
		)
	}
	
	logger.Info("HTTP Gateway server started successfully",
		"address", fmt.Sprintf("http://localhost:%s", *httpPort),
		"cache_enabled", redisCache != nil,
		"endpoints", endpoints,
	)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	sig := <-sigChan
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