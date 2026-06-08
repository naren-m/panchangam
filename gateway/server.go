package gateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/naren-m/panchangam/cache"
	"github.com/naren-m/panchangam/log"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logger = log.Logger()

// GatewayServer represents the HTTP gateway server.
type GatewayServer struct {
	grpcEndpoint      string
	httpPort          string
	server            *http.Server
	cache             *cache.RedisCache
	ephemerisProvider ephemeris.EphemerisProvider
}

// NewGatewayServer creates a new HTTP gateway server.
func NewGatewayServer(grpcEndpoint, httpPort string) *GatewayServer {
	return &GatewayServer{
		grpcEndpoint: grpcEndpoint,
		httpPort:     httpPort,
	}
}

// NewGatewayServerWithCache creates a new HTTP gateway server with Redis cache.
func NewGatewayServerWithCache(grpcEndpoint, httpPort string, redisCache *cache.RedisCache) *GatewayServer {
	return &GatewayServer{
		grpcEndpoint: grpcEndpoint,
		httpPort:     httpPort,
		cache:        redisCache,
	}
}

// SetEphemerisProvider sets the ephemeris provider for sky view functionality.
func (g *GatewayServer) SetEphemerisProvider(provider ephemeris.EphemerisProvider) {
	g.ephemerisProvider = provider
}

// Start starts the HTTP gateway server.
func (g *GatewayServer) Start(ctx context.Context) error {
	conn, err := grpc.NewClient(
		g.grpcEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	defer closeGRPCConnectionSafely(conn)

	client := ppb.NewPanchangamClient(conn)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/panchangam", g.handlePanchangam(client))
	mux.HandleFunc("/api/v1/tithi/current", g.handleCurrentTithi(client))
	mux.HandleFunc("/api/v1/sky-view", g.handleSkyView())

	if g.cache != nil {
		mux.HandleFunc("/api/v1/cache/health", g.handleCacheHealth())
		mux.HandleFunc("/api/v1/cache/stats", g.handleCacheStats())
	}

	handler := loggingMiddleware(mux)
	handler = addHealthCheck(handler)

	allowedOrigins := getCORSOrigins()
	logger.Info("CORS configuration", "allowed_origins", allowedOrigins)

	c := cors.New(newCORSOptions(allowedOrigins))
	handler = c.Handler(handler)

	g.server = &http.Server{
		Addr:              ":" + g.httpPort,
		Handler:           handler,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("HTTP Gateway server starting", "port", g.httpPort, "grpc_endpoint", g.grpcEndpoint)
	if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

// Stop gracefully stops the HTTP gateway server.
func (g *GatewayServer) Stop(ctx context.Context) error {
	if g.server == nil {
		return nil
	}

	logger.Info("Shutting down HTTP Gateway server")
	return g.server.Shutdown(ctx)
}

func closeGRPCConnectionSafely(conn interface{ Close() error }) {
	if conn == nil {
		return
	}
	if err := conn.Close(); err != nil {
		logger.Error("Failed to close gRPC connection", "error", err)
	}
}

func newCORSOptions(allowedOrigins []string) cors.Options {
	return cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"*",
		},
		ExposedHeaders: []string{
			"X-Request-Id",
			"X-Response-Time",
		},
		AllowCredentials: false,
		MaxAge:           300,
		Debug:            false,
	}
}
