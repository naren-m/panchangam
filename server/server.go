package main

import (
	"context"
	logging "github.com/naren-m/panchangam/logging"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	ps "github.com/naren-m/panchangam/services/panchangam"
	"google.golang.org/grpc"
	"net"
)

func main() {
	// Step 1: Initialize OpenTelemetry
	// Set up OpenTelemetry.
	o := observability.Observer("")
	o = observability.Observer("")
	defer o.Shutdown(context.Background())

	// Create a listener on TCP port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logging.Logger.Fatalf("Failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer(grpc.StatsHandler(observability.NewServerHandler()))

	pService := ps.NewPanchangamServer()
	ppb.RegisterPanchangamServer(grpcServer, pService)

	logging.Logger.Info("Server started on port :50051", nil)
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
