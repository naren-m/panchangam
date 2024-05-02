package main

import (
	"net"
	"context"
	"errors"
	"os"
	"os/signal"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	ps "github.com/naren-m/panchangam/services/panchangam"
	"github.com/naren-m/panchangam/observability"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	logging "github.com/naren-m/panchangam/logging"

)


func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Step 1: Initialize OpenTelemetry
	// Set up OpenTelemetry.
	otelShutdown, err := observability.SetupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()
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

	grpcServer := grpc.NewServer()

	pService := ps.NewPanchangamServer(mainSpan)
	ppb.RegisterPanchangamServer(grpcServer, pService)

	mainSpan.Trace("Server started on port :50051", nil)
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
