package main

import (
	"context"
	"net"
	"time"

	"github.com/naren-m/panchangam/aaa"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/naren-m/panchangam/log"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	ps "github.com/naren-m/panchangam/services/panchangam"
	"google.golang.org/grpc"
)

var logger = log.Logger()

func main() {
	// Step 1: Initialize OpenTelemetry
	// Set up OpenTelemetry.
	o, err := observability.NewObserver("localhost:4317")
	defer o.Shutdown(context.Background())

	// Create ephemeris providers and manager
	jplProvider := ephemeris.NewJPLProvider()
	swissProvider := ephemeris.NewSwissProvider()
	cache := ephemeris.NewMemoryCache(1000, 1*time.Hour)
	manager := ephemeris.NewManager(jplProvider, swissProvider, cache)
	defer manager.Close()

	// Create panchangam config
	config := ps.DefaultConfig()

	// Create a listener on TCP port 50052 (avoid conflicts)
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		logger.With("error", err).Error("Failed to listen:")
		return
	}
	a := aaa.NewAuth()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			observability.UnaryServerInterceptor(),
			a.AuthInterceptor(),
			a.AccountingInterceptor(),
		),
	)

	pService := ps.NewPanchangamServer(manager, config)
	ppb.RegisterPanchangamServer(grpcServer, pService)

	logger.Info("Server started on", "port", "50052")
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
