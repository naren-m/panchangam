package main

import (
	"fmt"
	"net"

	ppb "github.com/naren-m/panchangam/proto/panchangam"
	ps "github.com/naren-m/panchangam/services/panchangam"


	"google.golang.org/grpc"
	logging "github.com/naren-m/panchangam/logging"

	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

)





// initTracing sets up the OpenTelemetry tracing.
func initTracing() {
    // Create a new exporter to logging to stdout
    exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
    if err != nil {
        logging.Logger.Fatalf("failed to initialize stdouttrace export pipeline: %v", err)
    }

    // Create a new tracer provider with a batch span processor and the stdout exporter
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("panchangam service"),
        )),
    )

    // Set the global TracerProvider
    otel.SetTracerProvider(tp)
}


func main() {

	// Step 1: Initialize OpenTelemetry
	initTracing()

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
