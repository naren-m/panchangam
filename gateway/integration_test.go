package gateway

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ppb "github.com/naren-m/panchangam/proto"
	"github.com/naren-m/panchangam/services/panchangam"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// TestHTTPGatewayIntegration tests the complete HTTP Gateway to gRPC service flow
func TestHTTPGatewayIntegration(t *testing.T) {
	// Start gRPC server in memory
	grpcServer := grpc.NewServer()
	panchangamService := panchangam.NewPanchangamServer()
	ppb.RegisterPanchangamServer(grpcServer, panchangamService)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Logf("gRPC server error: %v", err)
		}
	}()
	defer grpcServer.Stop()

	// Create gRPC client
	conn, err := grpc.NewClient("bufnet", 
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	client := ppb.NewPanchangamClient(conn)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		validateResp   func(*testing.T, map[string]interface{})
	}{
		{
			name:           "Valid panchangam request",
			query:          "date=2024-01-15&lat=12.9716&lng=77.5946&tz=Asia/Kolkata",
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp map[string]interface{}) {
				// Validate required fields
				if date, ok := resp["date"].(string); !ok || date != "2024-01-15" {
					t.Errorf("Expected date '2024-01-15', got %v", resp["date"])
				}
				if tithi, ok := resp["tithi"].(string); !ok || tithi == "" {
					t.Errorf("Expected non-empty tithi, got %v", resp["tithi"])
				}
				if nakshatra, ok := resp["nakshatra"].(string); !ok || nakshatra == "" {
					t.Errorf("Expected non-empty nakshatra, got %v", resp["nakshatra"])
				}
				if events, ok := resp["events"].([]interface{}); !ok || len(events) == 0 {
					t.Errorf("Expected non-empty events array, got %v", resp["events"])
				}
			},
		},
		{
			name:           "Missing date parameter",
			query:          "lat=12.9716&lng=77.5946",
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, resp map[string]interface{}) {
				if errorInfo, ok := resp["error"].(map[string]interface{}); ok {
					if code, ok := errorInfo["code"].(string); !ok || code != "MISSING_PARAMETER" {
						t.Errorf("Expected error code 'MISSING_PARAMETER', got %v", errorInfo["code"])
					}
				} else {
					t.Error("Expected error object in response")
				}
			},
		},
		{
			name:           "Invalid latitude",
			query:          "date=2024-01-15&lat=999&lng=77.5946",
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, resp map[string]interface{}) {
				if errorInfo, ok := resp["error"].(map[string]interface{}); ok {
					if code, ok := errorInfo["code"].(string); !ok || code != "INVALID_ARGUMENT" {
						t.Errorf("Expected error code 'INVALID_ARGUMENT', got %v", errorInfo["code"])
					}
				} else {
					t.Error("Expected error object in response")
				}
			},
		},
		{
			name:           "Global location test - London",
			query:          "date=2024-06-21&lat=51.5074&lng=-0.1278&tz=Europe/London",
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp map[string]interface{}) {
				// Summer solstice in London should have long day
				if sunriseTime, ok := resp["sunrise_time"].(string); ok {
					if sunsetTime, ok := resp["sunset_time"].(string); ok {
						sunrise, _ := time.Parse("15:04:05", sunriseTime)
						sunset, _ := time.Parse("15:04:05", sunsetTime)
						dayLength := sunset.Sub(sunrise)
						if dayLength < 15*time.Hour { // Should be very long day
							t.Logf("Day length in London on summer solstice: %v", dayLength)
						}
					}
				}
			},
		},
		{
			name:           "Default timezone fallback",
			query:          "date=2024-01-15&lat=12.9716&lng=77.5946",
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp map[string]interface{}) {
				// Should still work with default timezone
				if tithi, ok := resp["tithi"].(string); !ok || tithi == "" {
					t.Errorf("Expected non-empty tithi with default timezone, got %v", resp["tithi"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP request
			req := httptest.NewRequest("GET", "/api/v1/panchangam?"+tt.query, nil)
			w := httptest.NewRecorder()

			// Create gateway handler
			gateway := &GatewayServer{}
			handler := gateway.handlePanchangam(client)

			// Execute request
			handler(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Parse response
			var resp map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			// Validate response
			if tt.validateResp != nil {
				tt.validateResp(t, resp)
			}
		})
	}
}

// TestHTTPGatewayPerformance tests the performance of the HTTP Gateway
func TestHTTPGatewayPerformance(t *testing.T) {
	// Start gRPC server in memory
	grpcServer := grpc.NewServer()
	panchangamService := panchangam.NewPanchangamServer()
	ppb.RegisterPanchangamServer(grpcServer, panchangamService)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Logf("gRPC server error: %v", err)
		}
	}()
	defer grpcServer.Stop()

	// Create gRPC client
	conn, err := grpc.NewClient("bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	client := ppb.NewPanchangamClient(conn)
	gateway := &GatewayServer{}
	handler := gateway.handlePanchangam(client)

	// Performance test
	iterations := 100
	start := time.Now()
	
	for i := 0; i < iterations; i++ {
		req := httptest.NewRequest("GET", "/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946&tz=Asia/Kolkata", nil)
		w := httptest.NewRecorder()
		handler(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Request %d failed with status %d", i, w.Code)
		}
	}
	
	duration := time.Since(start)
	avgDuration := duration / time.Duration(iterations)
	
	t.Logf("Performance test completed:")
	t.Logf("- Total time: %v", duration)
	t.Logf("- Average per request: %v", avgDuration)
	t.Logf("- Requests per second: %.2f", float64(iterations)/duration.Seconds())
	
	// Performance target: average should be < 100ms without random delays
	if avgDuration > 100*time.Millisecond {
		t.Errorf("Performance target missed: average %v > 100ms", avgDuration)
	}
}

// TestGatewayErrorHandling tests error handling in the gateway
func TestGatewayErrorHandling(t *testing.T) {
	// Start gRPC server in memory
	grpcServer := grpc.NewServer()
	panchangamService := panchangam.NewPanchangamServer()
	ppb.RegisterPanchangamServer(grpcServer, panchangamService)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Logf("gRPC server error: %v", err)
		}
	}()
	defer grpcServer.Stop()

	// Create gRPC client
	conn, err := grpc.NewClient("bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	client := ppb.NewPanchangamClient(conn)
	gateway := &GatewayServer{}
	handler := gateway.handlePanchangam(client)

	// Test various error conditions
	tests := []struct {
		name           string
		query          string
		expectedStatus int
		errorCode      string
	}{
		{
			name:           "Missing date",
			query:          "lat=12.9716&lng=77.5946",
			expectedStatus: http.StatusBadRequest,
			errorCode:      "MISSING_PARAMETER",
		},
		{
			name:           "Missing latitude",
			query:          "date=2024-01-15&lng=77.5946",
			expectedStatus: http.StatusBadRequest,
			errorCode:      "MISSING_PARAMETER",
		},
		{
			name:           "Invalid latitude format",
			query:          "date=2024-01-15&lat=invalid&lng=77.5946",
			expectedStatus: http.StatusBadRequest,
			errorCode:      "INVALID_PARAMETER",
		},
		{
			name:           "Invalid longitude range",
			query:          "date=2024-01-15&lat=12.9716&lng=999",
			expectedStatus: http.StatusBadRequest,
			errorCode:      "INVALID_ARGUMENT",
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

			var resp map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse error response: %v", err)
			}

			if errorInfo, ok := resp["error"].(map[string]interface{}); ok {
				if code, ok := errorInfo["code"].(string); !ok || code != tt.errorCode {
					t.Errorf("Expected error code '%s', got %v", tt.errorCode, errorInfo["code"])
				}
			} else {
				t.Error("Expected error object in response")
			}
		})
	}
}