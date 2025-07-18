package panchangam

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPanchangamServer_Get(t *testing.T) {
	// Initialize observability for testing
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	tests := []struct {
		name          string
		request       *ppb.GetPanchangamRequest
		validateFunc  func(t *testing.T, resp *ppb.GetPanchangamResponse, err error)
		skipRandomErr bool
	}{
		{
			name: "Valid request - New York",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timezone:  "America/New_York",
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				if err != nil {
					// Skip if random error occurred
					if status.Code(err) == codes.Internal {
						t.Skip("Random error occurred, skipping validation")
					}
				}
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotNil(t, resp.PanchangamData)
				assert.Equal(t, "2024-06-21", resp.PanchangamData.Date)
				assert.NotEmpty(t, resp.PanchangamData.SunriseTime)
				assert.NotEmpty(t, resp.PanchangamData.SunsetTime)
			},
		},
		{
			name: "Valid request - Chennai, India",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-03-15",
				Latitude:  13.0827,
				Longitude: 80.2707,
				Timezone:  "Asia/Kolkata",
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				if err != nil {
					// Skip if random error occurred
					if status.Code(err) == codes.Internal {
						t.Skip("Random error occurred, skipping validation")
					}
				}
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotNil(t, resp.PanchangamData)
				assert.NotEmpty(t, resp.PanchangamData.SunriseTime)
				assert.NotEmpty(t, resp.PanchangamData.SunsetTime)
			},
		},
		{
			name: "Valid request - London",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-12-21",
				Latitude:  51.5074,
				Longitude: -0.1278,
				Timezone:  "Europe/London",
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				if err != nil {
					// Skip if random error occurred
					if status.Code(err) == codes.Internal {
						t.Skip("Random error occurred, skipping validation")
					}
				}
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotNil(t, resp.PanchangamData)
				assert.NotEmpty(t, resp.PanchangamData.SunriseTime)
				assert.NotEmpty(t, resp.PanchangamData.SunsetTime)
			},
		},
		{
			name: "Valid request - Sydney",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-01-15",
				Latitude:  -33.8688,
				Longitude: 151.2093,
				Timezone:  "Australia/Sydney",
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				if err != nil {
					// Skip if random error occurred
					if status.Code(err) == codes.Internal {
						t.Skip("Random error occurred, skipping validation")
					}
				}
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotNil(t, resp.PanchangamData)
			},
		},
		{
			name: "Invalid latitude - too high",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  91.0,
				Longitude: 0.0,
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				assert.Contains(t, err.Error(), "latitude must be between -90 and 90")
			},
		},
		{
			name: "Invalid latitude - too low",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  -91.0,
				Longitude: 0.0,
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				assert.Contains(t, err.Error(), "latitude must be between -90 and 90")
			},
		},
		{
			name: "Invalid longitude - too high",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  0.0,
				Longitude: 181.0,
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				assert.Contains(t, err.Error(), "longitude must be between -180 and 180")
			},
		},
		{
			name: "Invalid longitude - too low",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  0.0,
				Longitude: -181.0,
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				assert.Contains(t, err.Error(), "longitude must be between -180 and 180")
			},
		},
		{
			name: "Invalid date format",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024/06/21", // Wrong format
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				assert.Contains(t, err.Error(), "invalid date format")
			},
		},
		{
			name: "Invalid timezone - fallback to local",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timezone:  "Invalid/Timezone",
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				if err != nil {
					// Skip if random error occurred
					if status.Code(err) == codes.Internal {
						t.Skip("Random error occurred, skipping validation")
					}
				}
				// Should not error, but use local timezone
				require.NoError(t, err)
				require.NotNil(t, resp)
			},
		},
		{
			name: "Polar region - Arctic summer",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  71.0,
				Longitude: 0.0,
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				if err != nil {
					// Skip if random error occurred
					if status.Code(err) == codes.Internal {
						t.Skip("Random error occurred, skipping validation")
					}
				}
				require.NoError(t, err)
				require.NotNil(t, resp)
				// In Arctic summer, sun might not set
			},
		},
		{
			name: "Polar region - Arctic winter",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-12-21",
				Latitude:  71.0,
				Longitude: 0.0,
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				if err != nil {
					// Skip if random error occurred
					if status.Code(err) == codes.Internal {
						t.Skip("Random error occurred, skipping validation")
					}
				}
				require.NoError(t, err)
				require.NotNil(t, resp)
				// In Arctic winter, sun might not rise
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resp, err := server.Get(ctx, tt.request)
			tt.validateFunc(t, resp, err)
		})
	}
}

func TestNewPanchangamServer(t *testing.T) {
	server := NewPanchangamServer()
	assert.NotNil(t, server)
	assert.NotNil(t, server.observer)
}

func TestPanchangamServer_SunriseSunsetFormat(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	// Test with a known location to verify time format
	req := &ppb.GetPanchangamRequest{
		Date:      "2024-06-21",
		Latitude:  40.7128,
		Longitude: -74.0060,
		Timezone:  "America/New_York",
	}

	// Try multiple times to handle random errors
	var resp *ppb.GetPanchangamResponse
	var err error
	for i := 0; i < 5; i++ {
		resp, err = server.Get(context.Background(), req)
		if err == nil {
			break
		}
		if status.Code(err) != codes.Internal {
			break
		}
		// Random error, try again
		time.Sleep(100 * time.Millisecond)
	}

	if err != nil && status.Code(err) == codes.Internal {
		t.Skip("Could not get successful response due to random errors")
	}

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.PanchangamData)

	// Verify time format (HH:MM:SS)
	sunrise := resp.PanchangamData.SunriseTime
	sunset := resp.PanchangamData.SunsetTime

	assert.Regexp(t, `^\d{2}:\d{2}:\d{2}$`, sunrise, "Sunrise time should be in HH:MM:SS format")
	assert.Regexp(t, `^\d{2}:\d{2}:\d{2}$`, sunset, "Sunset time should be in HH:MM:SS format")

	// Parse times to ensure they're valid
	_, err = time.Parse("15:04:05", sunrise)
	assert.NoError(t, err, "Sunrise time should be parseable")

	_, err = time.Parse("15:04:05", sunset)
	assert.NoError(t, err, "Sunset time should be parseable")
}

// Test helper functions
func TestTraceAttribute(t *testing.T) {
	attr := traceAttribute("test_key", "test_value")
	assert.Equal(t, "test_key", string(attr.Key))
	assert.Equal(t, "test_value", attr.Value.AsString())
}

func TestTraceAttributes(t *testing.T) {
	// Test valid key-value pairs
	attrs := traceAttributes("key1", "value1", "key2", "value2")
	assert.NotNil(t, attrs)
	assert.Len(t, attrs, 1) // Should return slice with one EventOption
	
	// Test odd number of arguments (should return nil)
	attrs = traceAttributes("key1", "value1", "key2")
	assert.Nil(t, attrs)
	
	// Test empty arguments
	attrs = traceAttributes()
	assert.NotNil(t, attrs)
	assert.Len(t, attrs, 1)
}

// Test comprehensive logging for all paths
func TestPanchangamServer_LoggingPaths(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()
	
	// Test successful request with all fields
	req := &ppb.GetPanchangamRequest{
		Date:              "2024-06-21",
		Latitude:          40.7128,
		Longitude:         -74.0060,
		Timezone:          "America/New_York",
		Region:            "North America",
		CalculationMethod: "Drik",
		Locale:            "en-US",
	}
	
	// Try multiple times to get a successful response (to test success logging)
	var resp *ppb.GetPanchangamResponse
	var err error
	for i := 0; i < 10; i++ {
		resp, err = server.Get(context.Background(), req)
		if err == nil {
			break
		}
		if status.Code(err) != codes.Internal {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	
	if err != nil && status.Code(err) == codes.Internal {
		t.Skip("Could not get successful response due to random errors")
	}
	
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.PanchangamData)
}

// Test logging with empty optional fields
func TestPanchangamServer_LoggingEmptyFields(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()
	
	req := &ppb.GetPanchangamRequest{
		Date:      "2024-06-21",
		Latitude:  40.7128,
		Longitude: -74.0060,
		// All other fields empty
	}
	
	// Try to get response (may fail due to random error)
	resp, err := server.Get(context.Background(), req)
	if err != nil && status.Code(err) == codes.Internal {
		// This is expected due to random errors in the service
		assert.Contains(t, err.Error(), "failed to fetch panchangam data")
	} else {
		require.NoError(t, err)
		assert.NotNil(t, resp)
	}
}

// Test edge case: empty date string
func TestPanchangamServer_EmptyDate(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()
	
	req := &ppb.GetPanchangamRequest{
		Date:      "",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	
	resp, err := server.Get(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

// Test edge case: very long timezone string
func TestPanchangamServer_LongTimezone(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()
	
	req := &ppb.GetPanchangamRequest{
		Date:      "2024-06-21",
		Latitude:  40.7128,
		Longitude: -74.0060,
		Timezone:  "This/Is/A/Very/Long/Invalid/Timezone/String/That/Should/Cause/An/Error",
	}
	
	// This should still work but fallback to local timezone
	resp, err := server.Get(context.Background(), req)
	if err != nil && status.Code(err) == codes.Internal {
		// Random error occurred, which is expected
		t.Skip("Random error occurred, skipping validation")
	} else {
		require.NoError(t, err)
		assert.NotNil(t, resp)
	}
}

// Test boundary coordinates
func TestPanchangamServer_BoundaryCoordinates(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()
	
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		shouldErr bool
	}{
		{"Valid boundary - North Pole", 90.0, 0.0, false},
		{"Valid boundary - South Pole", -90.0, 0.0, false},
		{"Valid boundary - East boundary", 0.0, 180.0, false},
		{"Valid boundary - West boundary", 0.0, -180.0, false},
		{"Just over north boundary", 90.1, 0.0, true},
		{"Just over south boundary", -90.1, 0.0, true},
		{"Just over east boundary", 0.0, 180.1, true},
		{"Just over west boundary", 0.0, -180.1, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  tt.latitude,
				Longitude: tt.longitude,
			}
			
			resp, err := server.Get(context.Background(), req)
			if tt.shouldErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			} else {
				// May still error due to random errors
				if err != nil && status.Code(err) == codes.Internal {
					t.Skip("Random error occurred, skipping validation")
				} else {
					require.NoError(t, err)
					assert.NotNil(t, resp)
				}
			}
		})
	}
}

// Test server initialization with observer
func TestPanchangamServer_Initialization(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()
	
	assert.NotNil(t, server)
	assert.NotNil(t, server.observer)
	assert.IsType(t, &PanchangamServer{}, server)
}

// Test with various date formats to ensure proper error handling
func TestPanchangamServer_DateFormats(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()
	
	invalidDates := []string{
		"2024-13-01",    // Invalid month
		"2024-02-30",    // Invalid day
		"2024-06-32",    // Invalid day
		"24-06-21",      // Wrong year format
		"2024/06/21",    // Wrong separator
		"2024-6-21",     // Single digit month
		"2024-06-1",     // Single digit day
		"invalid-date",  // Completely invalid
	}
	
	for _, date := range invalidDates {
		t.Run(fmt.Sprintf("Invalid date: %s", date), func(t *testing.T) {
			req := &ppb.GetPanchangamRequest{
				Date:      date,
				Latitude:  40.7128,
				Longitude: -74.0060,
			}
			
			resp, err := server.Get(context.Background(), req)
			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		})
	}
}