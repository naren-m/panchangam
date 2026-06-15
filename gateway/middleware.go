package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// writeErrorResponse writes a standardized error response.
func writeErrorResponse(w http.ResponseWriter, r *http.Request, status int, code, message string, details map[string]interface{}) {
	requestID := requestIDFromRequest(r)

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

// handleGRPCError converts gRPC errors to HTTP responses.
func handleGRPCError(w http.ResponseWriter, r *http.Request, err error) {
	requestID := requestIDFromRequest(r)

	httpStatus, apiError := convertGRPCError(err, requestID, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-Id", requestID)
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(apiError); err != nil {
		logger.Error("Failed to encode gRPC error response", "error", err)
	}
}

// loggingMiddleware adds request logging.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := requestIDFromRequest(r)

		w.Header().Set("X-Request-Id", requestID)

		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)
		w.Header().Set("X-Response-Time", duration.String())

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

// addHealthCheck adds a health check endpoint.
func addHealthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/health" {
			if err := writeHealthCheckResponse(w, time.Now().UTC()); err != nil {
				logger.Error("Failed to encode health check response", "error", err)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

type healthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Service   string `json:"service"`
	Version   string `json:"version"`
}

func writeHealthCheckResponse(w http.ResponseWriter, now time.Time) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(healthCheckResponse{
		Status:    "healthy",
		Timestamp: now.Format(time.RFC3339),
		Service:   "panchangam-gateway",
		Version:   "1.0.0",
	})
}

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// generateRequestID generates a simple request ID.
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

func requestIDFromRequest(r *http.Request) string {
	requestID := r.Header.Get("X-Request-Id")
	if requestID == "" {
		return generateRequestID()
	}
	return requestID
}

// getCORSOrigins returns the configured CORS origins.
func getCORSOrigins() []string {
	defaultOrigins := []string{
		"http://localhost:5173",
		"http://localhost:3000",
		"http://localhost:8086",
	}

	return parseCORSOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"), defaultOrigins)
}

func parseCORSOrigins(value string, defaultOrigins []string) []string {
	if value == "" {
		return defaultOrigins
	}

	envOrigins := strings.Split(value, ",")
	origins := make([]string, 0, len(envOrigins))

	for _, origin := range envOrigins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	if len(origins) == 0 {
		logger.Warn("No valid CORS origins found in CORS_ALLOWED_ORIGINS, using defaults")
		return defaultOrigins
	}

	return origins
}
