package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandlePanchangam_Success(t *testing.T) {
	mockClient := new(MockPanchangamClient)

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

	mockClient.On("Get", mock.Anything, mock.MatchedBy(func(req *ppb.GetPanchangamRequest) bool {
		return req.Date == "2024-01-15" &&
			req.Latitude == 12.9716 &&
			req.Longitude == 77.5946 &&
			req.Timezone == "Asia/Kolkata"
	})).Return(expectedResponse, nil)

	server := &GatewayServer{}
	handler := server.handlePanchangam(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946&tz=Asia/Kolkata", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

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
