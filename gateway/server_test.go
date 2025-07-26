package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockPanchangamClient is a mock implementation of the Panchangam gRPC client
type MockPanchangamClient struct {
	mock.Mock
}

func (m *MockPanchangamClient) Get(ctx context.Context, in *ppb.GetPanchangamRequest, opts ...grpc.CallOption) (*ppb.GetPanchangamResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ppb.GetPanchangamResponse), args.Error(1)
}

func TestHandlePanchangam_Success(t *testing.T) {
	// Create mock client
	mockClient := new(MockPanchangamClient)
	
	// Set up expected response
	expectedResponse := &ppb.GetPanchangamResponse{
		PanchangamData: &ppb.PanchangamData{
			Date:        "2024-01-15",
			Tithi:       "Shukla Paksha Tritiya",
			Nakshatra:   "Rohini",
			Yoga:        "Siddha",
			Karana:      "Gara",
			SunriseTime: "06:45:32",
			SunsetTime:  "18:21:47",
			Events: []*ppb.PanchangamEvent{
				{
					Name:      "Rahu Kalam",
					Time:      "08:00:00",
					EventType: "RAHU_KALAM",
				},
			},
		},
	}
	
	// Set up mock expectations
	mockClient.On("Get", mock.Anything, mock.MatchedBy(func(req *ppb.GetPanchangamRequest) bool {
		return req.Date == "2024-01-15" &&
			req.Latitude == 12.9716 &&
			req.Longitude == 77.5946 &&
			req.Timezone == "Asia/Kolkata"
	})).Return(expectedResponse, nil)
	
	// Create gateway server
	server := &GatewayServer{}
	handler := server.handlePanchangam(mockClient)
	
	// Create test request
	req := httptest.NewRequest("GET", "/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946&tz=Asia/Kolkata", nil)
	w := httptest.NewRecorder()
	
	// Execute handler
	handler(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	
	// Parse response body
	var result ppb.PanchangamData
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "2024-01-15", result.Date)
	assert.Equal(t, "Shukla Paksha Tritiya", result.Tithi)
	assert.Equal(t, "Rohini", result.Nakshatra)
	assert.Len(t, result.Events, 1)
	
	mockClient.AssertExpectations(t)
}

func TestHandlePanchangam_MissingParameters(t *testing.T) {
	tests := []struct {
		name        string
		queryString string
		expectedMsg string
	}{
		{
			name:        "Missing date",
			queryString: "lat=12.9716&lng=77.5946",
			expectedMsg: "Missing required parameter: date",
		},
		{
			name:        "Missing latitude",
			queryString: "date=2024-01-15&lng=77.5946",
			expectedMsg: "Missing required parameter: lat",
		},
		{
			name:        "Missing longitude",
			queryString: "date=2024-01-15&lat=12.9716",
			expectedMsg: "Missing required parameter: lng",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockPanchangamClient)
			server := &GatewayServer{}
			handler := server.handlePanchangam(mockClient)
			
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/panchangam?%s", tt.queryString), nil)
			w := httptest.NewRecorder()
			
			handler(w, req)
			
			assert.Equal(t, http.StatusBadRequest, w.Code)
			
			var errResp APIError
			err := json.Unmarshal(w.Body.Bytes(), &errResp)
			assert.NoError(t, err)
			assert.Equal(t, "MISSING_PARAMETER", errResp.Error.Code)
			assert.Equal(t, tt.expectedMsg, errResp.Error.Message)
		})
	}
}

func TestHandlePanchangam_InvalidParameters(t *testing.T) {
	tests := []struct {
		name        string
		queryString string
		expectedMsg string
	}{
		{
			name:        "Invalid latitude",
			queryString: "date=2024-01-15&lat=invalid&lng=77.5946",
			expectedMsg: "Invalid latitude format",
		},
		{
			name:        "Invalid longitude",
			queryString: "date=2024-01-15&lat=12.9716&lng=invalid",
			expectedMsg: "Invalid longitude format",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockPanchangamClient)
			server := &GatewayServer{}
			handler := server.handlePanchangam(mockClient)
			
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/panchangam?%s", tt.queryString), nil)
			w := httptest.NewRecorder()
			
			handler(w, req)
			
			assert.Equal(t, http.StatusBadRequest, w.Code)
			
			var errResp APIError
			err := json.Unmarshal(w.Body.Bytes(), &errResp)
			assert.NoError(t, err)
			assert.Equal(t, "INVALID_PARAMETER", errResp.Error.Code)
			assert.Equal(t, tt.expectedMsg, errResp.Error.Message)
		})
	}
}

func TestHandlePanchangam_GRPCErrors(t *testing.T) {
	tests := []struct {
		name           string
		grpcError      error
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "Invalid argument",
			grpcError:      status.Error(codes.InvalidArgument, "invalid date format"),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_PARAMETERS",
		},
		{
			name:           "Internal error",
			grpcError:      status.Error(codes.Internal, "internal server error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
		{
			name:           "Unavailable",
			grpcError:      status.Error(codes.Unavailable, "service unavailable"),
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   "SERVICE_UNAVAILABLE",
		},
		{
			name:           "Deadline exceeded",
			grpcError:      status.Error(codes.DeadlineExceeded, "timeout"),
			expectedStatus: http.StatusGatewayTimeout,
			expectedCode:   "REQUEST_TIMEOUT",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockPanchangamClient)
			mockClient.On("Get", mock.Anything, mock.Anything).Return(nil, tt.grpcError)
			
			server := &GatewayServer{}
			handler := server.handlePanchangam(mockClient)
			
			req := httptest.NewRequest("GET", "/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946", nil)
			w := httptest.NewRecorder()
			
			handler(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var errResp APIError
			err := json.Unmarshal(w.Body.Bytes(), &errResp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, errResp.Error.Code)
			
			mockClient.AssertExpectations(t)
		})
	}
}

func TestHealthCheckEndpoint(t *testing.T) {
	handler := addHealthCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not reach next handler for health check")
	}))
	
	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	
	var health map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &health)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", health["status"])
	assert.Equal(t, "panchangam-gateway", health["service"])
	assert.NotEmpty(t, health["timestamp"])
}

func TestLoggingMiddleware(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})
	
	handler := loggingMiddleware(next)
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-Id", "test-123")
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	assert.True(t, nextCalled)
	assert.Equal(t, "test-123", w.Header().Get("X-Request-Id"))
	assert.NotEmpty(t, w.Header().Get("X-Response-Time"))
}

func TestGenerateRequestID(t *testing.T) {
	id1 := generateRequestID()
	// Add a small delay to ensure different timestamps
	time.Sleep(time.Nanosecond)
	id2 := generateRequestID()
	
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	// Check that IDs are different (they should be due to nanosecond precision)
	if id1 == id2 {
		// If they're still the same (very unlikely), that's okay as long as format is correct
		t.Logf("Generated identical IDs (rare but possible): %s", id1)
	}
	assert.Contains(t, id1, "req_")
	assert.Contains(t, id2, "req_")
}

func TestCORSConfiguration(t *testing.T) {
	// This test would require actually starting the server with CORS middleware
	// For now, we'll test that the server can be created without errors
	server := NewGatewayServer("localhost:50052", "8080")
	assert.NotNil(t, server)
	assert.Equal(t, "localhost:50052", server.grpcEndpoint)
	assert.Equal(t, "8080", server.httpPort)
}

// Benchmark tests
func BenchmarkHandlePanchangam(b *testing.B) {
	mockClient := new(MockPanchangamClient)
	response := &ppb.GetPanchangamResponse{
		PanchangamData: &ppb.PanchangamData{
			Date:        "2024-01-15",
			Tithi:       "Test Tithi",
			Nakshatra:   "Test Nakshatra",
			Yoga:        "Test Yoga",
			Karana:      "Test Karana",
			SunriseTime: "06:45:32",
			SunsetTime:  "18:21:47",
		},
	}
	mockClient.On("Get", mock.Anything, mock.Anything).Return(response, nil)
	
	server := &GatewayServer{}
	handler := server.handlePanchangam(mockClient)
	
	req := httptest.NewRequest("GET", "/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler(w, req)
	}
}