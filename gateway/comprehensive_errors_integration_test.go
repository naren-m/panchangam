//go:build integration
// +build integration

package gateway

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestComprehensiveErrorScenarios(t *testing.T) {
	handler := newTestGatewayHandler(t)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		checkResponse  func(t *testing.T, resp map[string]interface{})
	}{
		{
			name:           "Missing all parameters",
			query:          "",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if _, ok := resp["error"]; !ok {
					t.Error("Expected error field in response")
				}
			},
		},
		{
			name:           "Invalid date format",
			query:          "date=invalid&lat=12.9716&lng=77.5946",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if errorObj, ok := resp["error"]; ok {
					t.Logf("Error response: %v", errorObj)
				}
			},
		},
		{
			name:           "Out of range latitude",
			query:          "date=2024-01-15&lat=100&lng=77.5946",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:           "Out of range longitude",
			query:          "date=2024-01-15&lat=12.9716&lng=200",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/panchangam?"+tt.query, nil)
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err == nil {
					tt.checkResponse(t, resp)
				}
			}
		})
	}
}
