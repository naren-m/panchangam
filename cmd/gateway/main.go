package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	)
	flag.Parse()

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

	// Create gateway server
	gatewayServer := gateway.NewGatewayServer(*grpcEndpoint, *httpPort)

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
	logger.Info("HTTP Gateway server started successfully",
		"address", fmt.Sprintf("http://localhost:%s", *httpPort),
		"health_check", fmt.Sprintf("http://localhost:%s/api/v1/health", *httpPort),
		"api_endpoint", fmt.Sprintf("http://localhost:%s/api/v1/panchangam", *httpPort),
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