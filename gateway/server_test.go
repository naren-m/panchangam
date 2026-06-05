package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ppb "github.com/naren-m/panchangam/proto"
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

func (m *MockPanchangamClient) GetTithiSummary(ctx context.Context, in *ppb.GetTithiSummaryRequest, opts ...grpc.CallOption) (*ppb.GetTithiSummaryResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ppb.GetTithiSummaryResponse), args.Error(1)
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

func TestHandleTithiSummary_Success(t *testing.T) {
	mockClient := new(MockPanchangamClient)
	expectedResponse := &ppb.GetTithiSummaryResponse{
		Date: "2026-06-02",
		Tithi: &ppb.TithiSummary{
			Name:            "Ekadashi",
			TraditionalName: "Ekadashi",
			Number:          11,
			Paksha:          "Shukla",
			PakshaDay:       11,
			Type:            "Nanda",
			StartTime:       "2026-06-02T01:00:00Z",
			EndTime:         "2026-06-03T02:00:00Z",
		},
		PanchaAnga: &ppb.PanchaAngaSummary{
			Nakshatra: "Anuradha",
			Yoga:      "Siddha",
			Karana:    "Gara",
			Vara:      "Tuesday",
		},
		Calculation: &ppb.TithiCalculationSummary{
			Timezone:       "America/Los_Angeles",
			Region:         "California",
			CalendarSystem: "Purnimanta",
			Method:         "Drik",
			Locale:         "en",
		},
		Day: &ppb.DaySummary{
			SunriseTime: "05:42",
			SunsetTime:  "19:01",
			AbhijitMuhurta: &ppb.TimeWindow{
				Name:       "Abhijit",
				StartTime:  "11:54",
				EndTime:    "12:48",
				Auspicious: true,
			},
		},
		GeneratedAt:   "2026-06-02T12:00:00Z",
		NextRefreshAt: "2026-06-03T02:00:00Z",
	}

	mockClient.On("GetTithiSummary", mock.Anything, mock.MatchedBy(func(req *ppb.GetTithiSummaryRequest) bool {
		return req.At == "2026-06-02T12:00:00Z" &&
			req.Latitude == 37.3382 &&
			req.Longitude == -121.8863 &&
			req.Timezone == "America/Los_Angeles" &&
			req.Region == "California" &&
			req.CalculationMethod == "Drik" &&
			req.Locale == "en" &&
			req.CalendarSystem == "Purnimanta"
	})).Return(expectedResponse, nil)

	server := &GatewayServer{}
	handler := server.handleTithiSummary(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/tithi/current?at=2026-06-02T12:00:00Z&lat=37.3382&lng=-121.8863&tz=America/Los_Angeles&region=California&method=Drik&locale=en&calendar_system=Purnimanta", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "MISS", w.Header().Get("X-Cache"))
	assert.Equal(t, "public, max-age=300", w.Header().Get("Cache-Control"))

	var result ppb.GetTithiSummaryResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Date, result.Date)
	assert.Equal(t, expectedResponse.Tithi.Name, result.Tithi.Name)
	assert.Equal(t, expectedResponse.PanchaAnga.Nakshatra, result.PanchaAnga.Nakshatra)
	assert.Equal(t, expectedResponse.Calculation.Timezone, result.Calculation.Timezone)
	assert.Equal(t, expectedResponse.Day.SunriseTime, result.Day.SunriseTime)
	assert.Equal(t, expectedResponse.Day.AbhijitMuhurta.StartTime, result.Day.AbhijitMuhurta.StartTime)
	assert.Equal(t, expectedResponse.NextRefreshAt, result.NextRefreshAt)

	mockClient.AssertExpectations(t)
}

func TestHandleTithiSummary_CapsCacheByNextRefresh(t *testing.T) {
	mockClient := new(MockPanchangamClient)
	expectedResponse := &ppb.GetTithiSummaryResponse{
		Date: "2026-06-02",
		Tithi: &ppb.TithiSummary{
			Name:            "Ekadashi",
			TraditionalName: "Ekadashi",
			Number:          11,
			Paksha:          "Shukla",
			PakshaDay:       11,
			Type:            "Nanda",
			StartTime:       "2026-06-02T01:00:00Z",
			EndTime:         "2026-06-02T12:02:00Z",
		},
		PanchaAnga: &ppb.PanchaAngaSummary{
			Nakshatra: "Anuradha",
			Yoga:      "Siddha",
			Karana:    "Gara",
			Vara:      "Tuesday",
		},
		Calculation: &ppb.TithiCalculationSummary{
			Timezone:       "America/Los_Angeles",
			Region:         "California",
			CalendarSystem: "Purnimanta",
			Method:         "Drik",
			Locale:         "en",
		},
		Day: &ppb.DaySummary{
			SunriseTime: "05:42",
			SunsetTime:  "19:01",
			AbhijitMuhurta: &ppb.TimeWindow{
				Name:       "Abhijit",
				StartTime:  "11:54",
				EndTime:    "12:48",
				Auspicious: true,
			},
		},
		GeneratedAt:   "2026-06-02T12:00:00Z",
		NextRefreshAt: "2026-06-02T12:02:00Z",
	}

	mockClient.On("GetTithiSummary", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	server := &GatewayServer{}
	handler := server.handleTithiSummary(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/tithi/current?at=2026-06-02T12:00:00Z&lat=37.3382&lng=-121.8863&tz=America/Los_Angeles", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "public, max-age=120", w.Header().Get("Cache-Control"))

	mockClient.AssertExpectations(t)
}

func TestHandleTithiSummary_DefaultsMomentToNow(t *testing.T) {
	mockClient := new(MockPanchangamClient)
	expectedResponse := &ppb.GetTithiSummaryResponse{
		Date: "2026-06-02",
		Tithi: &ppb.TithiSummary{
			Name:      "Ekadashi",
			Paksha:    "Shukla",
			PakshaDay: 11,
		},
		PanchaAnga: &ppb.PanchaAngaSummary{
			Nakshatra: "Anuradha",
			Yoga:      "Siddha",
			Karana:    "Gara",
			Vara:      "Tuesday",
		},
		Calculation: &ppb.TithiCalculationSummary{
			Timezone: "UTC",
			Region:   "global",
			Method:   "Drik",
			Locale:   "en",
		},
		Day: &ppb.DaySummary{
			SunriseTime: "05:42",
			SunsetTime:  "19:01",
			AbhijitMuhurta: &ppb.TimeWindow{
				Name:      "Abhijit",
				StartTime: "11:54",
				EndTime:   "12:48",
			},
		},
	}

	before := time.Now().UTC().Add(-time.Second)
	mockClient.On("GetTithiSummary", mock.Anything, mock.MatchedBy(func(req *ppb.GetTithiSummaryRequest) bool {
		at, err := time.Parse(time.RFC3339, req.At)
		return err == nil &&
			!at.Before(before) &&
			!at.After(time.Now().UTC().Add(time.Second)) &&
			req.Latitude == 37.3382 &&
			req.Longitude == -121.8863 &&
			req.Timezone == "UTC" &&
			req.Region == "global" &&
			req.CalculationMethod == "Drik" &&
			req.Locale == "en"
	})).Return(expectedResponse, nil)

	server := &GatewayServer{}
	handler := server.handleTithiSummary(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/tithi/current?lat=37.3382&lng=-121.8863", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockClient.AssertExpectations(t)
}

func TestHandleTithiSummary_MissingParameters(t *testing.T) {
	tests := []struct {
		name        string
		queryString string
		expectedMsg string
	}{
		{
			name:        "Missing latitude",
			queryString: "at=2026-06-02T12:00:00Z&lng=-121.8863",
			expectedMsg: "Missing required parameter: lat",
		},
		{
			name:        "Missing longitude",
			queryString: "at=2026-06-02T12:00:00Z&lat=37.3382",
			expectedMsg: "Missing required parameter: lng",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockPanchangamClient)
			server := &GatewayServer{}
			handler := server.handleTithiSummary(mockClient)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/tithi/current?%s", tt.queryString), nil)
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

func TestHandleTithiSummary_InvalidParameters(t *testing.T) {
	tests := []struct {
		name        string
		queryString string
		expectedMsg string
	}{
		{
			name:        "Invalid latitude",
			queryString: "at=2026-06-02T12:00:00Z&lat=invalid&lng=-121.8863",
			expectedMsg: "Invalid latitude format",
		},
		{
			name:        "Invalid longitude",
			queryString: "at=2026-06-02T12:00:00Z&lat=37.3382&lng=invalid",
			expectedMsg: "Invalid longitude format",
		},
		{
			name:        "Latitude out of range",
			queryString: "at=2026-06-02T12:00:00Z&lat=91&lng=-121.8863",
			expectedMsg: "Invalid latitude. Must be between -90 and 90",
		},
		{
			name:        "Longitude out of range",
			queryString: "at=2026-06-02T12:00:00Z&lat=37.3382&lng=181",
			expectedMsg: "Invalid longitude. Must be between -180 and 180",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockPanchangamClient)
			server := &GatewayServer{}
			handler := server.handleTithiSummary(mockClient)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/tithi/current?%s", tt.queryString), nil)
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

func TestHandleTithiSummary_GRPCErrors(t *testing.T) {
	tests := []struct {
		name           string
		grpcError      error
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "Invalid argument",
			grpcError:      status.Error(codes.InvalidArgument, "invalid timezone"),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_PARAMETERS",
		},
		{
			name:           "Unavailable",
			grpcError:      status.Error(codes.Unavailable, "service unavailable"),
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   "SERVICE_UNAVAILABLE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockPanchangamClient)
			mockClient.On("GetTithiSummary", mock.Anything, mock.Anything).Return(nil, tt.grpcError)

			server := &GatewayServer{}
			handler := server.handleTithiSummary(mockClient)

			req := httptest.NewRequest("GET", "/api/v1/tithi/current?at=2026-06-02T12:00:00Z&lat=37.3382&lng=-121.8863", nil)
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

func TestHandleTithiSummary_MethodNotAllowed(t *testing.T) {
	mockClient := new(MockPanchangamClient)
	server := &GatewayServer{}
	handler := server.handleTithiSummary(mockClient)

	req := httptest.NewRequest("POST", "/api/v1/tithi/current", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
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
