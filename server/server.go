package main

import (
	"net"
	"context"
	"os"
	"os/signal"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	ps "github.com/naren-m/panchangam/services/panchangam"
	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	logging "github.com/naren-m/panchangam/logging"
	"go.opentelemetry.io/otel/trace"
)
var tracer trace.Tracer

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

	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

	pService := ps.NewPanchangamServer(tracer)
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
