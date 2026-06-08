package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/naren-m/panchangam/log"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/naren-m/panchangam/services/panchangam"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const panchangamHealthServiceName = "panchangam.Panchangam"

var logger = log.Logger()

func main() {
	// Command line flags
	var (
		grpcPort    = flag.String("grpc-port", "50051", "gRPC server port")
		logLevel    = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
		healthCheck = flag.Bool("health-check", false, "Run a gRPC health check and exit")
	)
	flag.Parse()

	if *healthCheck {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := runGRPCHealthCheck(ctx, grpcHealthAddress(*grpcPort)); err != nil {
			fmt.Fprintf(os.Stderr, "gRPC health check failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

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
		grpc.UnaryInterceptor(observability.UnaryServerInterceptor()),
	)

	// Register Panchangam service
	panchangamService := panchangam.NewPanchangamServer()
	ppb.RegisterPanchangamServer(grpcServer, panchangamService)

	// Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(panchangamHealthServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection service (for grpcurl and other tools)
	reflection.Register(grpcServer)

	// Create TCP listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *grpcPort))
	if err != nil {
		logger.Error("Failed to create TCP listener", "error", err, "port", *grpcPort)
		os.Exit(1)
	}

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("gRPC server started successfully",
			"address", fmt.Sprintf("localhost:%s", *grpcPort),
			"health_check", "grpc://localhost:"+*grpcPort+"/grpc.health.v1.Health/Check",
		)
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			serverErr <- err
		}
	}()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var sig os.Signal
	select {
	case sig = <-sigChan:
	case err := <-serverErr:
		logger.Error("gRPC server error", "error", err)
		os.Exit(1)
	}
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

func grpcHealthAddress(grpcPort string) string {
	if strings.Contains(grpcPort, ":") {
		return grpcPort
	}
	return "localhost:" + grpcPort
}

func runGRPCHealthCheck(ctx context.Context, address string) error {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create gRPC health client for %s: %w", address, err)
	}
	defer func() {
		_ = conn.Close() // best-effort cleanup for a short-lived health check
	}()

	client := grpc_health_v1.NewHealthClient(conn)
	resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: panchangamHealthServiceName})
	if err != nil {
		return fmt.Errorf("call gRPC health endpoint %s: %w", address, err)
	}
	if resp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("gRPC health status for %s is %s", panchangamHealthServiceName, resp.GetStatus().String())
	}

	return nil
}
