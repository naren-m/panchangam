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
	"github.com/sirupsen/logrus"
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
	// Step 2: Use the logging package to create spans and log messages
	mainContext := context.Background()
	mainSpan := logging.NewSpan(mainContext, "main", logrus.DebugLevel)
    defer mainSpan.End()

	mainSpan.Trace("Starting server...", nil)

	// Create a listener on TCP port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logging.Logger.Fatalf("Failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

	pService := ps.NewPanchangamServer(mainSpan, tracer)
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
