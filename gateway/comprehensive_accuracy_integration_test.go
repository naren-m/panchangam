//go:build integration
// +build integration

package gateway

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDataAccuracyValidation(t *testing.T) {
	handler := newTestGatewayHandler(t)

	tests := []struct {
		name         string
		date         string
		lat          float64
		lng          float64
		tz           string
		validateFunc func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "New Moon - January 2024",
			date: "2024-01-11",
			lat:  12.9716,
			lng:  77.5946,
			tz:   "Asia/Kolkata",
			validateFunc: func(t *testing.T, data map[string]interface{}) {
				if tithi, ok := data["tithi"].(string); ok && tithi == "" {
					t.Error("Tithi should not be empty for New Moon date")
				}
			},
		},
		{
			name: "Summer Solstice - London",
			date: "2024-06-20",
			lat:  51.5074,
			lng:  -0.1278,
			tz:   "Europe/London",
			validateFunc: func(t *testing.T, data map[string]interface{}) {
				sunriseTime, sunriseOk := data["sunrise_time"].(string)
				sunsetTime, sunsetOk := data["sunset_time"].(string)

				if sunriseOk && sunsetOk && sunriseTime != "" && sunsetTime != "" {
					sunrise, _ := time.Parse("15:04:05", sunriseTime)
					sunset, _ := time.Parse("15:04:05", sunsetTime)
					dayLength := sunset.Sub(sunrise)

					if dayLength < 15*time.Hour {
						t.Logf("Summer solstice day length: %v (expected > 15h)", dayLength)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := buildQueryString(tt.date, tt.lat, tt.lng, tt.tz)
			req := httptest.NewRequest("GET", "/api/v1/panchangam?"+query, nil)
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("Expected status 200, got %d", w.Code)
			}

			var data map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, data)
			}
		})
	}
}

func TestDataConsistency(t *testing.T) {
	handler := newTestGatewayHandler(t)
	query := buildQueryString("2024-01-15", 12.9716, 77.5946, "Asia/Kolkata")

	req1 := httptest.NewRequest("GET", "/api/v1/panchangam?"+query, nil)
	w1 := httptest.NewRecorder()
	handler(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("First request failed with status %d", w1.Code)
	}

	var data1 map[string]interface{}
	if err := json.Unmarshal(w1.Body.Bytes(), &data1); err != nil {
		t.Fatalf("Failed to parse first response: %v", err)
	}

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/v1/panchangam?"+query, nil)
		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Request %d failed with status %d", i+1, w.Code)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Fatalf("Failed to parse response %d: %v", i+1, err)
		}

		if !compareData(data1, data) {
			t.Errorf("Request %d returned different data", i+1)
		}
	}
}

func TestMultipleLocationsDataFlow(t *testing.T) {
	handler := newTestGatewayHandler(t)

	locations := []struct {
		name string
		lat  float64
		lng  float64
		tz   string
	}{
		{"Bangalore", 12.9716, 77.5946, "Asia/Kolkata"},
		{"Mumbai", 19.0760, 72.8777, "Asia/Kolkata"},
		{"New York", 40.7128, -74.0060, "America/New_York"},
		{"London", 51.5074, -0.1278, "Europe/London"},
	}

	date := "2024-01-15"

	for _, loc := range locations {
		t.Run(loc.name, func(t *testing.T) {
			query := buildQueryString(date, loc.lat, loc.lng, loc.tz)
			req := httptest.NewRequest("GET", "/api/v1/panchangam?"+query, nil)
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("Request for %s failed with status %d", loc.name, w.Code)
			}

			var data map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
				t.Fatalf("Failed to parse response for %s: %v", loc.name, err)
			}

			requiredFields := []string{"date", "tithi", "nakshatra", "sunrise_time", "sunset_time"}
			for _, field := range requiredFields {
				if _, ok := data[field]; !ok {
					t.Errorf("Missing required field '%s' for %s", field, loc.name)
				}
			}
		})
	}
}
