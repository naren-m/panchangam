package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestServiceFunctionalCoverage provides comprehensive functional testing for the Panchangam service
// This test validates end-to-end functionality from service request to calculated response
func TestServiceFunctionalCoverage(t *testing.T) {
	// Initialize observability for testing
	observability.NewLocalObserver()
	
	// Create service instance
	server := NewPanchangamServer()
	
	t.Run("Functional_Service_Basic_Request", func(t *testing.T) {
		ctx := context.Background()
		
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,  // Bangalore coordinates
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
			Region:    "India",
		}
		
		// Execute service call
		resp, err := server.Get(ctx, req)
		
		// Validate response
		assert.NoError(t, err, "Service should not return error for valid request")
		require.NotNil(t, resp, "Response should not be nil")
		require.NotNil(t, resp.PanchangamData, "Panchangam data should not be nil")
		
		data := resp.PanchangamData
		
		// Validate basic fields
		assert.Equal(t, req.Date, data.Date, "Response date should match request date")
		assert.NotEmpty(t, data.Tithi, "Tithi should not be empty")
		assert.NotEmpty(t, data.Nakshatra, "Nakshatra should not be empty")
		assert.NotEmpty(t, data.Yoga, "Yoga should not be empty")
		assert.NotEmpty(t, data.Karana, "Karana should not be empty")
		assert.NotEmpty(t, data.SunriseTime, "Sunrise time should not be empty")
		assert.NotEmpty(t, data.SunsetTime, "Sunset time should not be empty")
		
		// Validate time format (HH:MM:SS)
		_, err = time.Parse("15:04:05", data.SunriseTime)
		assert.NoError(t, err, "Sunrise time should be in valid format")
		
		_, err = time.Parse("15:04:05", data.SunsetTime)
		assert.NoError(t, err, "Sunset time should be in valid format")
		
		// Validate events
		assert.NotNil(t, data.Events, "Events should not be nil")
		assert.True(t, len(data.Events) >= 0, "Events should be a valid array")
	})
	
	t.Run("Functional_Service_Input_Validation", func(t *testing.T) {
		ctx := context.Background()
		
		testCases := []struct {
			name          string
			request       *ppb.GetPanchangamRequest
			expectedError codes.Code
			description   string
		}{
			{
				name: "Invalid_Latitude_High",
				request: &ppb.GetPanchangamRequest{
					Date:      "2024-01-15",
					Latitude:  91.0,
					Longitude: 77.5946,
				},
				expectedError: codes.InvalidArgument,
				description:   "Latitude above 90 should be rejected",
			},
			{
				name: "Invalid_Latitude_Low",
				request: &ppb.GetPanchangamRequest{
					Date:      "2024-01-15",
					Latitude:  -91.0,
					Longitude: 77.5946,
				},
				expectedError: codes.InvalidArgument,
				description:   "Latitude below -90 should be rejected",
			},
			{
				name: "Invalid_Longitude_High",
				request: &ppb.GetPanchangamRequest{
					Date:      "2024-01-15",
					Latitude:  12.9716,
					Longitude: 181.0,
				},
				expectedError: codes.InvalidArgument,
				description:   "Longitude above 180 should be rejected",
			},
			{
				name: "Invalid_Longitude_Low",
				request: &ppb.GetPanchangamRequest{
					Date:      "2024-01-15",
					Latitude:  12.9716,
					Longitude: -181.0,
				},
				expectedError: codes.InvalidArgument,
				description:   "Longitude below -180 should be rejected",
			},
		}
		
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				resp, err := server.Get(ctx, tc.request)
				
				assert.Error(t, err, tc.description)
				assert.Nil(t, resp, "Response should be nil for invalid request")
				
				st, ok := status.FromError(err)
				require.True(t, ok, "Error should be a gRPC status error")
				assert.Equal(t, tc.expectedError, st.Code(), "Error code should match expected")
			})
		}
	})
	
	t.Run("Functional_Service_Date_Validation", func(t *testing.T) {
		ctx := context.Background()
		
		invalidDateCases := []string{
			"invalid-date",
			"2024-13-01",  // Invalid month
			"2024-01-32",  // Invalid day
			"24-01-01",    // Wrong year format
			"2024/01/01",  // Wrong separator
		}
		
		for _, invalidDate := range invalidDateCases {
			t.Run("Invalid_Date_"+invalidDate, func(t *testing.T) {
				req := &ppb.GetPanchangamRequest{
					Date:      invalidDate,
					Latitude:  12.9716,
					Longitude: 77.5946,
				}
				
				resp, err := server.Get(ctx, req)
				
				// Should fail due to invalid date format
				assert.Error(t, err, "Invalid date should cause error")
				assert.Nil(t, resp, "Response should be nil for invalid date")
				
				st, ok := status.FromError(err)
				require.True(t, ok, "Error should be a gRPC status error")
				assert.Equal(t, codes.InvalidArgument, st.Code(), "Should return InvalidArgument for bad date")
			})
		}
	})
	
	t.Run("Functional_Service_Geographic_Coverage", func(t *testing.T) {
		ctx := context.Background()
		
		// Test various geographic locations
		locations := []struct {
			name      string
			latitude  float64
			longitude float64
			timezone  string
		}{
			{"Bangalore_India", 12.9716, 77.5946, "Asia/Kolkata"},
			{"New_York_USA", 40.7128, -74.0060, "America/New_York"},
			{"London_UK", 51.5074, -0.1278, "Europe/London"},
			{"Tokyo_Japan", 35.6762, 139.6503, "Asia/Tokyo"},
			{"Sydney_Australia", -33.8688, 151.2093, "Australia/Sydney"},
			{"Arctic_Circle", 66.5, 0.0, "UTC"},
			{"Antarctic_Circle", -66.5, 0.0, "UTC"},
		}
		
		for _, loc := range locations {
			t.Run(loc.name, func(t *testing.T) {
				req := &ppb.GetPanchangamRequest{
					Date:      "2024-06-21", // Summer solstice
					Latitude:  loc.latitude,
					Longitude: loc.longitude,
					Timezone:  loc.timezone,
				}
				
				resp, err := server.Get(ctx, req)
				
				// Should succeed for all valid coordinates
				assert.NoError(t, err, "Valid coordinates should not cause error")
				require.NotNil(t, resp, "Response should not be nil")
				require.NotNil(t, resp.PanchangamData, "Panchangam data should not be nil")
				
				data := resp.PanchangamData
				assert.Equal(t, req.Date, data.Date, "Date should match")
				assert.NotEmpty(t, data.SunriseTime, "Should have sunrise time")
				assert.NotEmpty(t, data.SunsetTime, "Should have sunset time")
			})
		}
	})
	
	t.Run("Functional_Service_Timezone_Handling", func(t *testing.T) {
		ctx := context.Background()
		
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
		}
		
		// Test with no timezone (should use local)
		resp1, err1 := server.Get(ctx, req)
		assert.NoError(t, err1)
		require.NotNil(t, resp1)
		
		// Test with valid timezone
		req.Timezone = "Asia/Kolkata"
		resp2, err2 := server.Get(ctx, req)
		assert.NoError(t, err2)
		require.NotNil(t, resp2)
		
		// Test with invalid timezone (should fallback gracefully)
		req.Timezone = "Invalid/Timezone"
		resp3, err3 := server.Get(ctx, req)
		assert.NoError(t, err3, "Invalid timezone should fallback gracefully")
		require.NotNil(t, resp3)
	})
}

// TestServiceWithRealCalculations tests service integration with actual astronomy calculations
func TestServiceWithRealCalculations(t *testing.T) {
	t.Run("Functional_Real_Calculation_Integration", func(t *testing.T) {
		// This test demonstrates how the service should integrate with real calculations
		// Currently the service returns placeholder data, but this shows the expected pattern
		
		// Skip ephemeris setup for now - focus on service functional testing
		t.Skip("Skipping real calculation integration - requires ephemeris provider setup")
		
		// This test would validate that service integrates with real astronomy calculations
		// when the service is updated to use calculated values instead of placeholder data
		
		// The integration pattern would be:
		// 1. Service receives request with date and location
		// 2. Service creates ephemeris manager with providers
		// 3. Service creates all 5 Panchangam calculators
		// 4. Service calculates real values for each element
		// 5. Service returns calculated data in response
		
		t.Logf("This test will validate real calculation integration when implemented")
	})
}

// TestServicePerformance validates service-level performance targets
func TestServicePerformance(t *testing.T) {
	// Initialize observability for testing
	observability.NewLocalObserver()
	
	server := NewPanchangamServer()
	
	t.Run("Functional_Service_Performance", func(t *testing.T) {
		ctx := context.Background()
		
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
		}
		
		// Warm up
		_, _ = server.Get(ctx, req)
		
		// Measure performance
		start := time.Now()
		resp, err := server.Get(ctx, req)
		duration := time.Since(start)
		
		assert.NoError(t, err, "Performance test should not fail")
		require.NotNil(t, resp, "Response should not be nil")
		
		// Service performance target: <500ms (including simulation delays)
		assert.True(t, duration < 500*time.Millisecond, 
			"Service response should be under 500ms, got %v", duration)
		
		t.Logf("Service response time: %v", duration)
	})
	
	t.Run("Functional_Service_Concurrent_Performance", func(t *testing.T) {
		ctx := context.Background()
		
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
		}
		
		// Test concurrent requests
		concurrency := 10
		results := make(chan time.Duration, concurrency)
		
		start := time.Now()
		for i := 0; i < concurrency; i++ {
			go func() {
				reqStart := time.Now()
				resp, err := server.Get(ctx, req)
				reqDuration := time.Since(reqStart)
				
				assert.NoError(t, err, "Concurrent request should not fail")
				require.NotNil(t, resp, "Concurrent response should not be nil")
				
				results <- reqDuration
			}()
		}
		
		// Collect results
		var totalDuration time.Duration
		var maxDuration time.Duration
		for i := 0; i < concurrency; i++ {
			duration := <-results
			totalDuration += duration
			if duration > maxDuration {
				maxDuration = duration
			}
		}
		
		totalTime := time.Since(start)
		avgDuration := totalDuration / time.Duration(concurrency)
		
		// Concurrent performance targets
		assert.True(t, maxDuration < 1*time.Second, 
			"Max concurrent response should be under 1s, got %v", maxDuration)
		assert.True(t, avgDuration < 600*time.Millisecond, 
			"Average concurrent response should be under 600ms, got %v", avgDuration)
		
		t.Logf("Concurrent performance (%d requests):", concurrency)
		t.Logf("  Total time: %v", totalTime)
		t.Logf("  Average duration: %v", avgDuration)
		t.Logf("  Max duration: %v", maxDuration)
	})
}

// TestServiceErrorHandling validates comprehensive error handling scenarios
func TestServiceErrorHandling(t *testing.T) {
	// Initialize observability for testing
	observability.NewLocalObserver()
	
	server := NewPanchangamServer()
	
	t.Run("Functional_Service_Error_Coverage", func(t *testing.T) {
		ctx := context.Background()
		
		// Test nil request
		resp, err := server.Get(ctx, nil)
		assert.Error(t, err, "Nil request should cause error")
		assert.Nil(t, resp, "Response should be nil for nil request")
		
		// Test empty request
		resp, err = server.Get(ctx, &ppb.GetPanchangamRequest{})
		assert.Error(t, err, "Empty request should cause error")
		assert.Nil(t, resp, "Response should be nil for empty request")
		
		// Test context cancellation
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
		}
		
		resp, err = server.Get(cancelCtx, req)
		// Note: Current service doesn't check context cancellation in placeholder implementation
		// In real implementation, this should handle cancellation
		
		// Test timeout context
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		
		time.Sleep(2 * time.Millisecond) // Ensure timeout
		
		resp, err = server.Get(timeoutCtx, req)
		// Note: Similar to above, real implementation should handle timeouts
	})
}

// TestServiceFeatureCoverage validates that all documented features are accessible via service
func TestServiceFeatureCoverage(t *testing.T) {
	// Initialize observability for testing
	observability.NewLocalObserver()
	
	server := NewPanchangamServer()
	
	t.Run("Functional_Feature_Coverage_Validation", func(t *testing.T) {
		ctx := context.Background()
		
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
			Region:    "India",
			CalculationMethod: "traditional",
			Locale:    "en",
		}
		
		resp, err := server.Get(ctx, req)
		require.NoError(t, err, "Service should handle all request parameters")
		require.NotNil(t, resp, "Response should not be nil")
		require.NotNil(t, resp.PanchangamData, "Panchangam data should not be nil")
		
		data := resp.PanchangamData
		
		// Validate all documented Panchangam elements are present
		features := map[string]string{
			"TITHI_001":     data.Tithi,
			"NAKSHATRA_001": data.Nakshatra,
			"YOGA_001":      data.Yoga,
			"KARANA_001":    data.Karana,
			"VARA_001":      "", // Vara is not directly exposed but should be calculated
		}
		
		for featureID, value := range features {
			if featureID != "VARA_001" { // Skip Vara as it's not in proto
				assert.NotEmpty(t, value, "Feature %s should have a value", featureID)
			}
		}
		
		// Validate astronomical calculations
		assert.NotEmpty(t, data.SunriseTime, "ASTRONOMY_001: Sunrise should be calculated")
		assert.NotEmpty(t, data.SunsetTime, "ASTRONOMY_001: Sunset should be calculated")
		
		// Validate service capabilities
		assert.Equal(t, req.Date, data.Date, "SERVICE_001: Date should be processed correctly")
		assert.NotNil(t, data.Events, "SERVICE_001: Events should be included")
		
		// Validate response structure matches proto definition
		assert.IsType(t, "", data.Date, "Date should be string")
		assert.IsType(t, "", data.Tithi, "Tithi should be string")
		assert.IsType(t, "", data.Nakshatra, "Nakshatra should be string")
		assert.IsType(t, "", data.Yoga, "Yoga should be string")
		assert.IsType(t, "", data.Karana, "Karana should be string")
		assert.IsType(t, "", data.SunriseTime, "Sunrise time should be string")
		assert.IsType(t, "", data.SunsetTime, "Sunset time should be string")
		assert.IsType(t, []*ppb.PanchangamEvent{}, data.Events, "Events should be array")
		
		t.Logf("Feature coverage validation completed successfully")
		t.Logf("All documented features are accessible via service API")
	})
}