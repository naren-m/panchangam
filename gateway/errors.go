package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// APIError represents a structured API error response
type APIError struct {
	Error ErrorDetails `json:"error"`
}

// ErrorDetails contains detailed error information
type ErrorDetails struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"requestId"`
	Timestamp string                 `json:"timestamp"`
	Path      string                 `json:"path"`
}

// customErrorHandler handles gRPC errors and converts them to HTTP responses
func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	// Extract request ID from headers
	requestID := r.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = generateRequestID()
	}

	// Set common headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-Id", requestID)

	// Convert gRPC error to HTTP status and error response
	httpStatus, apiError := convertGRPCError(err, requestID, r.URL.Path)
	
	// Set HTTP status code
	w.WriteHeader(httpStatus)

	// Marshal and write error response
	if err := json.NewEncoder(w).Encode(apiError); err != nil {
		logger.Error("Failed to encode error response", "error", err, "request_id", requestID)
		// Fallback to plain text error
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Internal server error")
	}

	// Log the error
	logger.Error("API error",
		"http_status", httpStatus,
		"grpc_code", status.Code(err),
		"error_code", apiError.Error.Code,
		"message", apiError.Error.Message,
		"request_id", requestID,
		"path", r.URL.Path,
		"method", r.Method,
	)
}

// convertGRPCError converts a gRPC error to HTTP status code and API error
func convertGRPCError(err error, requestID, path string) (int, *APIError) {
	// Extract gRPC status
	s, ok := status.FromError(err)
	if !ok {
		// Not a gRPC error, treat as internal error
		return http.StatusInternalServerError, &APIError{
			Error: ErrorDetails{
				Code:      "INTERNAL_ERROR",
				Message:   "An internal server error occurred",
				RequestID: requestID,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      path,
			},
		}
	}

	// Map gRPC codes to HTTP status codes and error details
	httpStatus, errorCode, message, details := mapGRPCCodeToHTTP(s)

	return httpStatus, &APIError{
		Error: ErrorDetails{
			Code:      errorCode,
			Message:   message,
			Details:   details,
			RequestID: requestID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Path:      path,
		},
	}
}

// mapGRPCCodeToHTTP maps gRPC status codes to HTTP status codes and error details
func mapGRPCCodeToHTTP(s *status.Status) (int, string, string, map[string]interface{}) {
	var details map[string]interface{}

	switch s.Code() {
	case codes.OK:
		return http.StatusOK, "SUCCESS", "Request completed successfully", nil

	case codes.InvalidArgument:
		details = map[string]interface{}{
			"validation": "Request parameters are invalid",
			"grpc_code":  "INVALID_ARGUMENT",
		}
		return http.StatusBadRequest, "INVALID_PARAMETERS", enhanceValidationMessage(s.Message()), details

	case codes.NotFound:
		details = map[string]interface{}{
			"resource": "The requested resource was not found",
			"grpc_code": "NOT_FOUND",
		}
		return http.StatusNotFound, "RESOURCE_NOT_FOUND", "The requested resource was not found", details

	case codes.AlreadyExists:
		details = map[string]interface{}{
			"conflict": "Resource already exists",
			"grpc_code": "ALREADY_EXISTS",
		}
		return http.StatusConflict, "RESOURCE_EXISTS", "The resource already exists", details

	case codes.PermissionDenied:
		details = map[string]interface{}{
			"authorization": "Insufficient permissions",
			"grpc_code":     "PERMISSION_DENIED",
		}
		return http.StatusForbidden, "ACCESS_DENIED", "Permission denied", details

	case codes.Unauthenticated:
		details = map[string]interface{}{
			"authentication": "Authentication required",
			"grpc_code":      "UNAUTHENTICATED",
		}
		return http.StatusUnauthorized, "AUTHENTICATION_REQUIRED", "Authentication required", details

	case codes.ResourceExhausted:
		details = map[string]interface{}{
			"rate_limiting": "Too many requests",
			"grpc_code":     "RESOURCE_EXHAUSTED",
			"retry_after":   30,
		}
		return http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests, please try again later", details

	case codes.FailedPrecondition:
		details = map[string]interface{}{
			"precondition": "Request precondition failed",
			"grpc_code":    "FAILED_PRECONDITION",
		}
		return http.StatusPreconditionFailed, "PRECONDITION_FAILED", "Request precondition failed", details

	case codes.OutOfRange:
		details = map[string]interface{}{
			"range": "Parameter out of valid range",
			"grpc_code": "OUT_OF_RANGE",
		}
		return http.StatusBadRequest, "PARAMETER_OUT_OF_RANGE", "Parameter value is out of valid range", details

	case codes.Unimplemented:
		details = map[string]interface{}{
			"feature": "Feature not implemented",
			"grpc_code": "UNIMPLEMENTED",
		}
		return http.StatusNotImplemented, "NOT_IMPLEMENTED", "This feature is not yet implemented", details

	case codes.Internal:
		details = map[string]interface{}{
			"server": "Internal server error",
			"grpc_code": "INTERNAL",
		}
		return http.StatusInternalServerError, "INTERNAL_ERROR", "An internal server error occurred", details

	case codes.Unavailable:
		details = map[string]interface{}{
			"service": "Service temporarily unavailable",
			"grpc_code": "UNAVAILABLE",
			"retry_after": 60,
		}
		return http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Service is temporarily unavailable", details

	case codes.DataLoss:
		details = map[string]interface{}{
			"data": "Data loss detected",
			"grpc_code": "DATA_LOSS",
		}
		return http.StatusInternalServerError, "DATA_LOSS", "Data loss detected", details

	case codes.DeadlineExceeded:
		details = map[string]interface{}{
			"timeout": "Request timeout",
			"grpc_code": "DEADLINE_EXCEEDED",
		}
		return http.StatusGatewayTimeout, "REQUEST_TIMEOUT", "Request timed out", details

	case codes.Canceled:
		details = map[string]interface{}{
			"cancellation": "Request was cancelled",
			"grpc_code": "CANCELLED",
		}
		return http.StatusRequestTimeout, "REQUEST_CANCELLED", "Request was cancelled", details

	default:
		details = map[string]interface{}{
			"unknown": "Unknown error occurred",
			"grpc_code": s.Code().String(),
		}
		return http.StatusInternalServerError, "UNKNOWN_ERROR", fmt.Sprintf("Unknown error: %s", s.Message()), details
	}
}

// enhanceValidationMessage provides more specific validation error messages
func enhanceValidationMessage(original string) string {
	switch {
	case contains(original, "latitude"):
		return "Latitude must be between -90 and 90 degrees"
	case contains(original, "longitude"):
		return "Longitude must be between -180 and 180 degrees"
	case contains(original, "date"):
		return "Date must be in YYYY-MM-DD format"
	case contains(original, "timezone"):
		return "Invalid timezone identifier"
	default:
		return original
	}
}

// contains checks if a string contains a substring (case-insensitive)
func contains(str, substr string) bool {
	return len(str) >= len(substr) && 
		   (str == substr || 
		    (len(str) > len(substr) && 
		     containsIgnoreCase(str, substr)))
}

func containsIgnoreCase(str, substr string) bool {
	str = toLower(str)
	substr = toLower(substr)
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			result[i] = s[i] + 32
		} else {
			result[i] = s[i]
		}
	}
	return string(result)
}