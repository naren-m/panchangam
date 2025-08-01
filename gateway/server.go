package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/naren-m/panchangam/cache"
	"github.com/naren-m/panchangam/log"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logger = log.Logger()

// GatewayServer represents the HTTP gateway server
type GatewayServer struct {
	grpcEndpoint string
	httpPort     string
	server       *http.Server
	cache        *cache.RedisCache
}

// NewGatewayServer creates a new HTTP gateway server
func NewGatewayServer(grpcEndpoint, httpPort string) *GatewayServer {
	return &GatewayServer{
		grpcEndpoint: grpcEndpoint,
		httpPort:     httpPort,
	}
}

// NewGatewayServerWithCache creates a new HTTP gateway server with Redis cache
func NewGatewayServerWithCache(grpcEndpoint, httpPort string, redisCache *cache.RedisCache) *GatewayServer {
	return &GatewayServer{
		grpcEndpoint: grpcEndpoint,
		httpPort:     httpPort,
		cache:        redisCache,
	}
}

// Start starts the HTTP gateway server
func (g *GatewayServer) Start(ctx context.Context) error {
	// Create a client connection to the gRPC server
	conn, err := grpc.NewClient(
		g.grpcEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	defer conn.Close()

	// Create gRPC client
	client := ppb.NewPanchangamClient(conn)

	// Create HTTP router
	mux := http.NewServeMux()

	// Add panchangam endpoint
	mux.HandleFunc("/api/v1/panchangam", g.handlePanchangam(client))
	
	// Add cache health endpoint
	if g.cache != nil {
		mux.HandleFunc("/api/v1/cache/health", g.handleCacheHealth())
		mux.HandleFunc("/api/v1/cache/stats", g.handleCacheStats())
	}

	// Add custom middleware for logging and monitoring
	handler := loggingMiddleware(mux)
	
	// Add health check endpoint
	handler = addHealthCheck(handler)

	// Configure CORS with dynamic origins from environment
	allowedOrigins := getCORSOrigins()
	
	// Log configured origins for debugging
	logger.Info("CORS configuration", "allowed_origins", allowedOrigins)
	
	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"*", // Allow all headers for debugging
		},
		ExposedHeaders: []string{
			"X-Request-Id",
			"X-Response-Time",
		},
		AllowCredentials: false,
		MaxAge:           300,
		Debug: true, // Enable CORS debugging
	})

	// Apply CORS middleware
	handler = c.Handler(handler)

	// Create HTTP server
	g.server = &http.Server{
		Addr:    ":" + g.httpPort,
		Handler: handler,
		// Security and performance settings
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("HTTP Gateway server starting", "port", g.httpPort, "grpc_endpoint", g.grpcEndpoint)

	// Start server
	if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

// Stop gracefully stops the HTTP gateway server
func (g *GatewayServer) Stop(ctx context.Context) error {
	if g.server == nil {
		return nil
	}

	logger.Info("Shutting down HTTP Gateway server")
	return g.server.Shutdown(ctx)
}

// handlePanchangam handles HTTP requests to the panchangam endpoint with caching
func (g *GatewayServer) handlePanchangam(client ppb.PanchangamClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract query parameters
		query := r.URL.Query()
		
		// Required parameters
		date := query.Get("date")
		if date == "" {
			writeErrorResponse(w, r, http.StatusBadRequest, "MISSING_PARAMETER", "Missing required parameter: date", nil)
			return
		}

		latStr := query.Get("lat")
		if latStr == "" {
			writeErrorResponse(w, r, http.StatusBadRequest, "MISSING_PARAMETER", "Missing required parameter: lat", nil)
			return
		}

		lngStr := query.Get("lng")
		if lngStr == "" {
			writeErrorResponse(w, r, http.StatusBadRequest, "MISSING_PARAMETER", "Missing required parameter: lng", nil)
			return
		}

		// Parse latitude
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_PARAMETER", "Invalid latitude format", map[string]interface{}{
				"parameter": "lat",
				"value": latStr,
				"expected": "float64",
			})
			return
		}

		// Parse longitude
		lng, err := strconv.ParseFloat(lngStr, 64)
		if err != nil {
			writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_PARAMETER", "Invalid longitude format", map[string]interface{}{
				"parameter": "lng",
				"value": lngStr,
				"expected": "float64",
			})
			return
		}

		// Optional parameters
		timezone := query.Get("tz")
		if timezone == "" {
			timezone = "UTC" // Default timezone
		}

		region := query.Get("region")
		if region == "" {
			region = "global" // Default region
		}
		
		calculationMethod := query.Get("method")
		if calculationMethod == "" {
			calculationMethod = "traditional" // Default method
		}
		
		locale := query.Get("locale")
		if locale == "" {
			locale = "en" // Default locale
		}

		// Set timeout for operations
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		var responseData interface{}
		cacheHit := false

		// Try cache first if available
		if g.cache != nil {
			cacheKey := g.cache.GenerateCacheKey(date, region, calculationMethod, lat, lng)
			
			cachedData, err := g.cache.Get(ctx, cacheKey)
			if err != nil {
				logger.Error("Cache get error", "error", err, "key", cacheKey)
			} else if cachedData != nil {
				// Cache hit - convert to response format
				responseData = convertCacheToResponse(cachedData)
				cacheHit = true
				logger.Debug("Cache hit", "key", cacheKey, "date", date)
			}
		}

		// If cache miss, fetch from gRPC service
		if !cacheHit {
			// Create gRPC request
			req := &ppb.GetPanchangamRequest{
				Date:              date,
				Latitude:          lat,
				Longitude:         lng,
				Timezone:          timezone,
				Region:            region,
				CalculationMethod: calculationMethod,
				Locale:            locale,
			}

			// Make gRPC call
			resp, err := client.Get(ctx, req)
			if err != nil {
				handleGRPCError(w, r, err)
				return
			}

			responseData = resp.PanchangamData

			// Store in cache if available
			if g.cache != nil && resp.PanchangamData != nil {
				cacheData := convertResponseToCache(resp.PanchangamData)
				cacheKey := g.cache.GenerateCacheKey(date, region, calculationMethod, lat, lng)
				
				if err := g.cache.Set(ctx, cacheKey, cacheData); err != nil {
					logger.Error("Cache set error", "error", err, "key", cacheKey)
				} else {
					logger.Debug("Cache set", "key", cacheKey, "date", date)
				}
			}
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		if cacheHit {
			w.Header().Set("X-Cache", "HIT")
			w.Header().Set("Cache-Control", "public, max-age=1800") // 30 minutes for cached data
		} else {
			w.Header().Set("X-Cache", "MISS")
			w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes for fresh data
		}

		// Write successful response
		if err := json.NewEncoder(w).Encode(responseData); err != nil {
			logger.Error("Failed to encode response", "error", err)
			writeErrorResponse(w, r, http.StatusInternalServerError, "ENCODING_ERROR", "Failed to encode response", nil)
			return
		}
	}
}

// writeErrorResponse writes a standardized error response
func writeErrorResponse(w http.ResponseWriter, r *http.Request, status int, code, message string, details map[string]interface{}) {
	requestID := r.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = generateRequestID()
	}

	errorResp := APIError{
		Error: ErrorDetails{
			Code:      code,
			Message:   message,
			Details:   details,
			RequestID: requestID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Path:      r.URL.Path,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-Id", requestID)
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		logger.Error("Failed to encode error response", "error", err)
	}
}

// handleGRPCError converts gRPC errors to HTTP responses
func handleGRPCError(w http.ResponseWriter, r *http.Request, err error) {
	requestID := r.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = generateRequestID()
	}

	httpStatus, apiError := convertGRPCError(err, requestID, r.URL.Path)
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-Id", requestID)
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(apiError); err != nil {
		logger.Error("Failed to encode gRPC error response", "error", err)
	}
}

// loggingMiddleware adds request logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate request ID if not present
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Add request ID to response headers
		w.Header().Set("X-Request-Id", requestID)

		// Create a response writer wrapper to capture status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapper, r)

		// Calculate response time
		duration := time.Since(start)
		w.Header().Set("X-Response-Time", duration.String())

		// Log the request
		logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"status", wrapper.statusCode,
			"duration", duration,
			"request_id", requestID,
			"user_agent", r.Header.Get("User-Agent"),
			"remote_addr", r.RemoteAddr,
		)
	})
}

// addHealthCheck adds a health check endpoint
func addHealthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/health" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"status": "healthy",
				"timestamp": "%s",
				"service": "panchangam-gateway",
				"version": "1.0.0"
			}`, time.Now().UTC().Format(time.RFC3339))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// getCORSOrigins returns the list of allowed CORS origins from environment configuration
func getCORSOrigins() []string {
	// Default origins for development
	defaultOrigins := []string{
		"http://localhost:5173", // Vite dev server
		"http://localhost:3000", // React dev server
		"http://localhost:8086", // Docker frontend container
	}
	
	// Get origins from environment variable
	corsOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOriginsEnv == "" {
		// If not set, return default origins for backward compatibility
		return defaultOrigins
	}
	
	// Parse comma-separated origins
	envOrigins := strings.Split(corsOriginsEnv, ",")
	origins := make([]string, 0, len(envOrigins))
	
	for _, origin := range envOrigins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	
	// If no valid origins were provided, return defaults
	if len(origins) == 0 {
		logger.Warn("No valid CORS origins found in CORS_ALLOWED_ORIGINS, using defaults")
		return defaultOrigins
	}
	
	return origins
}

// convertResponseToCache converts gRPC response to cache format
func convertResponseToCache(data interface{}) *cache.PanchangamCacheData {
	// This is a simplified conversion - you may need to adjust based on your actual gRPC response structure
	responseMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}

	cacheData := &cache.PanchangamCacheData{}
	
	if date, ok := responseMap["date"].(string); ok {
		cacheData.Date = date
	}
	if tithi, ok := responseMap["tithi"].(string); ok {
		cacheData.Tithi = tithi
	}
	if nakshatra, ok := responseMap["nakshatra"].(string); ok {
		cacheData.Nakshatra = nakshatra
	}
	if yoga, ok := responseMap["yoga"].(string); ok {
		cacheData.Yoga = yoga
	}
	if karana, ok := responseMap["karana"].(string); ok {
		cacheData.Karana = karana
	}
	if sunriseTime, ok := responseMap["sunrise_time"].(string); ok {
		cacheData.SunriseTime = sunriseTime
	}
	if sunsetTime, ok := responseMap["sunset_time"].(string); ok {
		cacheData.SunsetTime = sunsetTime
	}

	// Convert events if present
	if events, ok := responseMap["events"].([]interface{}); ok {
		cacheData.Events = make([]cache.Event, len(events))
		for i, event := range events {
			if eventMap, ok := event.(map[string]interface{}); ok {
				cacheEvent := cache.Event{}
				if name, ok := eventMap["name"].(string); ok {
					cacheEvent.Name = name
				}
				if time, ok := eventMap["time"].(string); ok {
					cacheEvent.Time = time
				}
				if eventType, ok := eventMap["event_type"].(string); ok {
					cacheEvent.EventType = eventType
				}
				cacheData.Events[i] = cacheEvent
			}
		}
	}

	return cacheData
}

// convertCacheToResponse converts cache format to response format
func convertCacheToResponse(data *cache.PanchangamCacheData) interface{} {
	response := map[string]interface{}{
		"date":         data.Date,
		"tithi":        data.Tithi,
		"nakshatra":    data.Nakshatra,
		"yoga":         data.Yoga,
		"karana":       data.Karana,
		"sunrise_time": data.SunriseTime,
		"sunset_time":  data.SunsetTime,
		"events":       make([]map[string]interface{}, len(data.Events)),
	}

	for i, event := range data.Events {
		response["events"].([]map[string]interface{})[i] = map[string]interface{}{
			"name":       event.Name,
			"time":       event.Time,
			"event_type": event.EventType,
		}
	}

	return response
}

// handleCacheHealth handles cache health check requests
func (g *GatewayServer) handleCacheHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if g.cache == nil {
			writeErrorResponse(w, r, http.StatusServiceUnavailable, "CACHE_DISABLED", "Cache is not enabled", nil)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := g.cache.HealthCheck(ctx); err != nil {
			writeErrorResponse(w, r, http.StatusServiceUnavailable, "CACHE_UNHEALTHY", "Cache health check failed", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"status": "healthy",
			"timestamp": "%s",
			"service": "redis-cache"
		}`, time.Now().UTC().Format(time.RFC3339))
	}
}

// handleCacheStats handles cache statistics requests
func (g *GatewayServer) handleCacheStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if g.cache == nil {
			writeErrorResponse(w, r, http.StatusServiceUnavailable, "CACHE_DISABLED", "Cache is not enabled", nil)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		stats, err := g.cache.GetStats(ctx)
		if err != nil {
			writeErrorResponse(w, r, http.StatusInternalServerError, "STATS_ERROR", "Failed to get cache statistics", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			logger.Error("Failed to encode cache stats", "error", err)
			writeErrorResponse(w, r, http.StatusInternalServerError, "ENCODING_ERROR", "Failed to encode stats", nil)
		}
	}
}