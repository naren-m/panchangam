package main

import (
	"fmt"
	"net"

	ppb "github.com/naren-m/panchangam/proto/panchangam"
	ps "github.com/naren-m/panchangam/services/panchangam"


	"google.golang.org/grpc"
	logging "github.com/naren-m/panchangam/logging"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

)


func main() {

	// Step 1: Initialize OpenTelemetry

	// Step 2: Use the logging package to create spans and log messages
	_, span := logging.CreateSpan("main")
	defer span.End()

	// Perform some operation here...

	// Log a message without trace context
	logging.Logger.Info("This is a regular log message.")

	// Log a message with trace context
	logging.LogWithTrace(span, logging.Logger.Level, "This is a log message with trace context.", nil)
	// Initialize the logger
	logging.Logger.Info("Starting server...")
	logging.Logger.Debug("Starting server...")
	// Create a listener on TCP port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logging.Logger.Fatalf("Failed to listen: %v", err)
		return
	}
	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

	// Create a new gRPC server
	// Register the Panchangam service with the server
	ppb.RegisterPanchangamServer(grpcServer, &ps.PanchangamServer{})

	logging.Logger.Info("Server started on port :50051")

	// Start serving requests
	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Println("Failed to serve:", err)
		return
	}
}
