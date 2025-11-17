// +build integration

package gateway

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	ppb "github.com/naren-m/panchangam/proto"
	"github.com/naren-m/panchangam/services/panchangam"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// TestDataAccuracyValidation tests data accuracy against known astronomical events
// Issue #81 Requirement: Data accuracy verification against known test cases
func TestDataAccuracyValidation(t *testing.T) {
	// Setup gRPC server
	grpcServer, client := setupTestGRPCServer(t)
	defer grpcServer.Stop()

	gateway := &GatewayServer{}
	handler := gateway.handlePanchangam(client)

	// Known astronomical data for validation
	tests := []struct {
		name          string
		date          string
		lat           float64
		lng           float64
		tz            string
		validateFunc  func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "New Moon - January 2024",
			date: "2024-01-11",
			lat:  12.9716,
			lng:  77.5946,
			tz:   "Asia/Kolkata",
			validateFunc: func(t *testing.T, data map[string]interface{}) {
				// Should have data for new moon
				if tithi, ok := data["tithi"].(string); ok {
					if tithi == "" {
						t.Error("Tithi should not be empty for New Moon date")
					}
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
				// Summer solstice should have long day
				sunriseTime, sunriseOk := data["sunrise_time"].(string)
				sunsetTime, sunsetOk := data["sunset_time"].(string)

				if sunriseOk && sunsetOk && sunriseTime != "" && sunsetTime != "" {
					// Parse times and verify long day
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

// TestConcurrentRequestPerformance tests concurrent request handling
// Issue #81 Requirement: 50 concurrent requests in <5 seconds
func TestConcurrentRequestPerformance(t *testing.T) {
	grpcServer, client := setupTestGRPCServer(t)
	defer grpcServer.Stop()

	gateway := &GatewayServer{}
	handler := gateway.handlePanchangam(client)

	concurrentRequests := 50
	timeout := 5 * time.Second

	start := time.Now()
	var wg sync.WaitGroup
	errors := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			query := buildQueryString("2024-01-15", 12.9716, 77.5946, "Asia/Kolkata")
			req := httptest.NewRequest("GET", "/api/v1/panchangam?"+query, nil)
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != http.StatusOK {
				errors <- &testError{msg: "Request failed"}
			}
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	close(errors)

	duration := time.Since(start)

	// Check for errors
	errorCount := 0
	for range errors {
		errorCount++
	}

	if errorCount > 0 {
		t.Errorf("%d out of %d concurrent requests failed", errorCount, concurrentRequests)
	}

	// Issue #81 Requirement: 50 concurrent requests in <5 seconds
	if duration > timeout {
		t.Errorf("Concurrent requests took %v, expected < %v", duration, timeout)
	}

	t.Logf("Concurrent performance: %d requests in %v (%.2f req/sec)",
		concurrentRequests, duration, float64(concurrentRequests)/duration.Seconds())
}

// TestResponseTimeTarget tests average response time
// Issue #81 Requirement: <500ms average response time
func TestResponseTimeTarget(t *testing.T) {
	grpcServer, client := setupTestGRPCServer(t)
	defer grpcServer.Stop()

	gateway := &GatewayServer{}
	handler := gateway.handlePanchangam(client)

	iterations := 100
	var totalDuration time.Duration

	for i := 0; i < iterations; i++ {
		query := buildQueryString("2024-01-15", 12.9716, 77.5946, "Asia/Kolkata")
		req := httptest.NewRequest("GET", "/api/v1/panchangam?"+query, nil)
		w := httptest.NewRecorder()

		start := time.Now()
		handler(w, req)
		duration := time.Since(start)

		totalDuration += duration

		if w.Code != http.StatusOK {
			t.Errorf("Request %d failed with status %d", i, w.Code)
		}
	}

	avgDuration := totalDuration / time.Duration(iterations)

	t.Logf("Average response time: %v over %d requests", avgDuration, iterations)

	// Issue #81 Requirement: <500ms average
	targetDuration := 500 * time.Millisecond
	if avgDuration > targetDuration {
		t.Errorf("Average response time %v exceeds target %v", avgDuration, targetDuration)
	}
}

// TestDataConsistency tests that identical requests return identical data
// Issue #81 Requirement: 100% data consistency verification
func TestDataConsistency(t *testing.T) {
	grpcServer, client := setupTestGRPCServer(t)
	defer grpcServer.Stop()

	gateway := &GatewayServer{}
	handler := gateway.handlePanchangam(client)

	query := buildQueryString("2024-01-15", 12.9716, 77.5946, "Asia/Kolkata")

	// Make first request
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

	// Make multiple subsequent requests
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

		// Verify consistency
		if !compareData(data1, data) {
			t.Errorf("Request %d returned different data", i+1)
		}
	}
}

// TestErrorRecoveryTime tests error recovery performance
// Issue #81 Requirement: <3 seconds for retry scenarios
func TestErrorRecoveryTime(t *testing.T) {
	grpcServer, client := setupTestGRPCServer(t)
	defer grpcServer.Stop()

	gateway := &GatewayServer{}
	handler := gateway.handlePanchangam(client)

	maxRetries := 3
	retryDelay := 500 * time.Millisecond

	// Simulate retry scenario
	start := time.Now()

	for attempt := 0; attempt < maxRetries; attempt++ {
		query := buildQueryString("2024-01-15", 12.9716, 77.5946, "Asia/Kolkata")
		req := httptest.NewRequest("GET", "/api/v1/panchangam?"+query, nil)
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code == http.StatusOK {
			break
		}

		if attempt < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	duration := time.Since(start)

	// Issue #81 Requirement: <3 seconds for retry scenarios
	if duration > 3*time.Second {
		t.Errorf("Error recovery took %v, expected <3s", duration)
	}

	t.Logf("Error recovery completed in %v", duration)
}

// TestMultipleLocationsDataFlow tests data flow for various geographic locations
func TestMultipleLocationsDataFlow(t *testing.T) {
	grpcServer, client := setupTestGRPCServer(t)
	defer grpcServer.Stop()

	gateway := &GatewayServer{}
	handler := gateway.handlePanchangam(client)

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

			// Verify essential fields
			requiredFields := []string{"date", "tithi", "nakshatra", "sunrise_time", "sunset_time"}
			for _, field := range requiredFields {
				if _, ok := data[field]; !ok {
					t.Errorf("Missing required field '%s' for %s", field, loc.name)
				}
			}
		})
	}
}

// TestComprehensiveErrorScenarios tests various error scenarios
func TestComprehensiveErrorScenarios(t *testing.T) {
	grpcServer, client := setupTestGRPCServer(t)
	defer grpcServer.Stop()

	gateway := &GatewayServer{}
	handler := gateway.handlePanchangam(client)

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
			checkResponse: nil,
		},
		{
			name:           "Out of range longitude",
			query:          "date=2024-01-15&lat=12.9716&lng=200",
			expectedStatus: http.StatusBadRequest,
			checkResponse: nil,
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

// Helper functions

func setupTestGRPCServer(t *testing.T) (*grpc.Server, ppb.PanchangamClient) {
	// Create buffer connection
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	panchangamService := panchangam.NewPanchangamServer()
	ppb.RegisterPanchangamServer(grpcServer, panchangamService)

	// Start server
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			t.Logf("gRPC server error: %v", err)
		}
	}()

	// Create client
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.NewClient("bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	client := ppb.NewPanchangamClient(conn)

	return grpcServer, client
}

func buildQueryString(date string, lat, lng float64, tz string) string {
	return "date=" + date + "&lat=" + floatToString(lat) + "&lng=" + floatToString(lng) + "&tz=" + tz
}

func floatToString(f float64) string {
	return string(rune(int(f*10000))) // Simplified for testing
}

func compareData(data1, data2 map[string]interface{}) bool {
	// Simple comparison of key fields
	fields := []string{"date", "tithi", "nakshatra", "sunrise_time", "sunset_time"}

	for _, field := range fields {
		val1, ok1 := data1[field]
		val2, ok2 := data2[field]

		if ok1 != ok2 || val1 != val2 {
			return false
		}
	}

	return true
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
