package main

import (
	"context"
	"net"
	"os"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestRunGRPCHealthCheckAcceptsServingStatus(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on local port: %v", err)
	}
	defer func() {
		_ = listener.Close() // best-effort test cleanup
	}()

	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(panchangamHealthServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
	defer grpcServer.Stop()

	go func() {
		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("serve gRPC health endpoint: %v", err)
		}
	}()

	err = runGRPCHealthCheck(context.Background(), listener.Addr().String())
	if err != nil {
		t.Fatalf("expected serving gRPC health check, got %v", err)
	}
}

func TestRunGRPCHealthCheckRejectsNotServingStatus(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on local port: %v", err)
	}
	defer func() {
		_ = listener.Close() // best-effort test cleanup
	}()

	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(panchangamHealthServiceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	defer grpcServer.Stop()

	go func() {
		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("serve gRPC health endpoint: %v", err)
		}
	}()

	err = runGRPCHealthCheck(context.Background(), listener.Addr().String())
	if err == nil {
		t.Fatal("expected not-serving gRPC health check to return an error")
	}
}

func TestMainHandlesServeFailureOutsideServeGoroutine(t *testing.T) {
	source, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}

	for _, want := range []string{
		"serverErr := make(chan error, 1)",
		"case err := <-serverErr:",
		"err != grpc.ErrServerStopped",
	} {
		if !strings.Contains(string(source), want) {
			t.Fatalf("main should handle gRPC serve failures through %q", want)
		}
	}

	if strings.Contains(string(source), "logger.Error(\"gRPC server error\", \"error\", err)\n\t\t\tos.Exit(1)") {
		t.Fatal("gRPC Serve goroutine must not call os.Exit directly")
	}
}
