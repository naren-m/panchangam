package main

import (
	"context"
	logging "github.com/naren-m/panchangam/logging"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	ps "github.com/naren-m/panchangam/services/panchangam"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

var tracer observability.Tracer

func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Step 1: Initialize OpenTelemetry
	// Set up OpenTelemetry.
	tp := observability.InitTracerProvider()
	defer tp.Shutdown(ctx)

	mp := observability.InitMeterProvider()
	defer mp.Shutdown(context.Background())

	tracer := tp.Tracer("panchangam")

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
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	grpcServer.Stop()

	return
}
