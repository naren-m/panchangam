package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/naren-m/panchangam/log"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"github.com/naren-m/panchangam/services/panchangam"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var logger = log.Logger()

func main() {
	// Command line flags
	var (
		grpcPort = flag.String("grpc-port", "50051", "gRPC server port")
		logLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	)
	flag.Parse()

	// Initialize observability
	ctx := context.Background()
	observer := observability.Observer()
	defer func() {
		if err := observer.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown observability", "error", err)
		}
	}()

	logger.Info("Starting Panchangam gRPC Server",
		"grpc_port", *grpcPort,
		"log_level", *logLevel,
	)

	// Create gRPC server with observability interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(observer.UnaryServerInterceptor()),
		grpc.StreamInterceptor(observer.StreamServerInterceptor()),
	)

	// Register Panchangam service
	panchangamService := panchangam.NewPanchangamServer()
	ppb.RegisterPanchangamServer(grpcServer, panchangamService)

	// Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("panchangam.Panchangam", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection service (for grpcurl and other tools)
	reflection.Register(grpcServer)

	// Create TCP listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *grpcPort))
	if err != nil {
		logger.Error("Failed to create TCP listener", "error", err, "port", *grpcPort)
		os.Exit(1)
	}

	// Start server in a goroutine
	go func() {
		logger.Info("gRPC server started successfully",
			"address", fmt.Sprintf("localhost:%s", *grpcPort),
			"health_check", "grpc://localhost:"+*grpcPort+"/grpc.health.v1.Health/Check",
		)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("gRPC server error", "error", err)
			os.Exit(1)
		}
	}()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	sig := <-sigChan
	logger.Info("Received shutdown signal", "signal", sig)

	// Graceful shutdown with timeout
	logger.Info("Shutting down gRPC server gracefully")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop accepting new requests
	healthServer.Shutdown()

	// Wait for existing requests to complete or timeout
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	// Wait for graceful shutdown or force stop after timeout
	select {
	case <-stopped:
		logger.Info("gRPC server shutdown completed successfully")
	case <-shutdownCtx.Done():
		logger.Warn("gRPC server shutdown timed out, forcing stop")
		grpcServer.Stop()
	}

	logger.Info("Panchangam gRPC Server stopped")
}
