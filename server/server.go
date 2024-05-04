package main

import (
	"context"
	"log/slog"
	"github.com/naren-m/panchangam/observability"
	"github.com/naren-m/panchangam/aaa"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	ps "github.com/naren-m/panchangam/services/panchangam"
	"google.golang.org/grpc"
	"net"
)

func main() {
	// Step 1: Initialize OpenTelemetry
	// Set up OpenTelemetry.
	o := observability.NewObserver("")
	o = observability.NewObserver("")
	defer o.Shutdown(context.Background())

	// Create a listener on TCP port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.With("error", err).Error("Failed to listen:")
		return
	}

    grpcServer := grpc.NewServer(
        grpc.StatsHandler(observability.NewServerHandler()),
        grpc.ChainUnaryInterceptor(
            observability.UnaryServerInterceptor(),
            aaa.AuthInterceptor(),
        ),
    )

	pService := ps.NewPanchangamServer()
	ppb.RegisterPanchangamServer(grpcServer, pService)

	slog.Info("Server started on", "port", "50051")
	// Start serving requests
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- grpcServer.Serve(listener)
	}()
	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		grpcServer.Stop()
		return
	}
}
