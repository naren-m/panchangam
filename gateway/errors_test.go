package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestConvertGRPCError(t *testing.T) {
	tests := []struct {
		name           string
		grpcError      error
		expectedStatus int
		expectedCode   string
		expectedMsg    string
	}{
		{
			name:           "Invalid argument",
			grpcError:      status.Error(codes.InvalidArgument, "latitude must be between -90 and 90"),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_PARAMETERS",
			expectedMsg:    "Latitude must be between -90 and 90 degrees",
		},
		{
			name:           "Not found",
			grpcError:      status.Error(codes.NotFound, "resource not found"),
			expectedStatus: http.StatusNotFound,
			expectedCode:   "RESOURCE_NOT_FOUND",
			expectedMsg:    "The requested resource was not found",
		},
		{
			name:           "Already exists",
			grpcError:      status.Error(codes.AlreadyExists, "resource exists"),
			expectedStatus: http.StatusConflict,
			expectedCode:   "RESOURCE_EXISTS",
			expectedMsg:    "The resource already exists",
		},
		{
			name:           "Permission denied",
			grpcError:      status.Error(codes.PermissionDenied, "access denied"),
			expectedStatus: http.StatusForbidden,
			expectedCode:   "ACCESS_DENIED",
			expectedMsg:    "Permission denied",
		},
		{
			name:           "Unauthenticated",
			grpcError:      status.Error(codes.Unauthenticated, "not authenticated"),
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "AUTHENTICATION_REQUIRED",
			expectedMsg:    "Authentication required",
		},
		{
			name:           "Resource exhausted",
			grpcError:      status.Error(codes.ResourceExhausted, "too many requests"),
			expectedStatus: http.StatusTooManyRequests,
			expectedCode:   "RATE_LIMITED",
			expectedMsg:    "Too many requests, please try again later",
		},
		{
			name:           "Internal error",
			grpcError:      status.Error(codes.Internal, "internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
			expectedMsg:    "An internal server error occurred",
		},
		{
			name:           "Unavailable",
			grpcError:      status.Error(codes.Unavailable, "service unavailable"),
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   "SERVICE_UNAVAILABLE",
			expectedMsg:    "Service is temporarily unavailable",
		},
		{
			name:           "Deadline exceeded",
			grpcError:      status.Error(codes.DeadlineExceeded, "timeout"),
			expectedStatus: http.StatusGatewayTimeout,
			expectedCode:   "REQUEST_TIMEOUT",
			expectedMsg:    "Request timed out",
		},
		{
			name:           "Non-gRPC error",
			grpcError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
			expectedMsg:    "An internal server error occurred",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpStatus, apiError := convertGRPCError(tt.grpcError, "test-123", "/api/v1/test")
			
			assert.Equal(t, tt.expectedStatus, httpStatus)
			assert.Equal(t, tt.expectedCode, apiError.Error.Code)
			assert.Equal(t, tt.expectedMsg, apiError.Error.Message)
			assert.Equal(t, "test-123", apiError.Error.RequestID)
			assert.Equal(t, "/api/v1/test", apiError.Error.Path)
			assert.NotEmpty(t, apiError.Error.Timestamp)
		})
	}
}

func TestEnhanceValidationMessage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "latitude is invalid",
			expected: "Latitude must be between -90 and 90 degrees",
		},
		{
			input:    "longitude out of range",
			expected: "Longitude must be between -180 and 180 degrees",
		},
		{
			input:    "invalid date format",
			expected: "Date must be in YYYY-MM-DD format",
		},
		{
			input:    "unknown timezone",
			expected: "Invalid timezone identifier",
		},
		{
			input:    "some other error",
			expected: "some other error",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := enhanceValidationMessage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCustomErrorHandler(t *testing.T) {
	// Create a test request
	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	req.Header.Set("X-Request-Id", "test-request-123")
	
	// Create a test response writer
	w := httptest.NewRecorder()
	
	// Create test error
	testErr := status.Error(codes.InvalidArgument, "test error")
	
	// Call custom error handler
	customErrorHandler(context.Background(), runtime.NewServeMux(), nil, w, req, testErr)
	
	// Check response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "test-request-123", w.Header().Get("X-Request-Id"))
	
	// Parse response body
	var errResp APIError
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_PARAMETERS", errResp.Error.Code)
	assert.Equal(t, "/api/v1/test", errResp.Error.Path)
	assert.Equal(t, "test-request-123", errResp.Error.RequestID)
}

func TestWriteErrorResponse(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	w := httptest.NewRecorder()
	
	details := map[string]interface{}{
		"field": "date",
		"value": "invalid",
	}
	
	writeErrorResponse(w, req, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid date format", details)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.NotEmpty(t, w.Header().Get("X-Request-Id"))
	
	var errResp APIError
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_ERROR", errResp.Error.Code)
	assert.Equal(t, "Invalid date format", errResp.Error.Message)
	assert.Equal(t, "date", errResp.Error.Details["field"])
	assert.Equal(t, "invalid", errResp.Error.Details["value"])
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		str      string
		substr   string
		expected bool
	}{
		{"latitude is invalid", "latitude", true},
		{"LATITUDE is invalid", "latitude", true},
		{"Latitude is invalid", "LATITUDE", true},
		{"something else", "latitude", false},
		{"lat", "latitude", false},
		{"", "latitude", false},
		{"latitude", "", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.str+"-"+tt.substr, func(t *testing.T) {
			result := contains(tt.str, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandleGRPCError(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/panchangam", nil)
	req.Header.Set("X-Request-Id", "test-456")
	w := httptest.NewRecorder()
	
	grpcErr := status.Error(codes.NotFound, "panchangam data not found")
	
	handleGRPCError(w, req, grpcErr)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "test-456", w.Header().Get("X-Request-Id"))
	
	var errResp APIError
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Equal(t, "RESOURCE_NOT_FOUND", errResp.Error.Code)
	assert.Equal(t, "test-456", errResp.Error.RequestID)
}

// Benchmark tests
func BenchmarkConvertGRPCError(b *testing.B) {
	err := status.Error(codes.InvalidArgument, "test error")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		convertGRPCError(err, "test-123", "/api/v1/test")
	}
}

func BenchmarkEnhanceValidationMessage(b *testing.B) {
	messages := []string{
		"latitude is invalid",
		"longitude out of range",
		"invalid date format",
		"unknown timezone",
		"some other error",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enhanceValidationMessage(messages[i%len(messages)])
	}
}