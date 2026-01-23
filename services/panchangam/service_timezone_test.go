package panchangam

import (
	"context"
	"testing"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPanchangamServer_TimezoneSupport(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	tests := []struct {
		name           string
		timezone       string
		expectedError  bool
		validateResult func(*testing.T, *ppb.PanchangamData)
	}{
		{
			name:          "IANA timezone Asia/Kolkata",
			timezone:      "Asia/Kolkata",
			expectedError: false,
			validateResult: func(t *testing.T, data *ppb.PanchangamData) {
				assert.Equal(t, "Asia/Kolkata", data.Timezone)
				assert.Equal(t, "+05:30", data.TimezoneOffset)
				assert.False(t, data.IsDst) // India doesn't use DST
			},
		},
		{
			name:          "IANA timezone America/New_York",
			timezone:      "America/New_York",
			expectedError: false,
			validateResult: func(t *testing.T, data *ppb.PanchangamData) {
				assert.Equal(t, "America/New_York", data.Timezone)
				// Offset depends on DST, so just check it's not empty
				assert.NotEmpty(t, data.TimezoneOffset)
			},
		},
		{
			name:          "UTC offset +05:30",
			timezone:      "+05:30",
			expectedError: false,
			validateResult: func(t *testing.T, data *ppb.PanchangamData) {
				assert.Equal(t, "UTC+05:30", data.Timezone)
				assert.Equal(t, "+05:30", data.TimezoneOffset)
				assert.False(t, data.IsDst) // Fixed offset has no DST
			},
		},
		{
			name:          "UTC offset -08:00",
			timezone:      "-08:00",
			expectedError: false,
			validateResult: func(t *testing.T, data *ppb.PanchangamData) {
				assert.Equal(t, "UTC-08:00", data.Timezone)
				assert.Equal(t, "-08:00", data.TimezoneOffset)
				assert.False(t, data.IsDst) // Fixed offset has no DST
			},
		},
		{
			name:          "UTC timezone",
			timezone:      "UTC",
			expectedError: false,
			validateResult: func(t *testing.T, data *ppb.PanchangamData) {
				assert.Equal(t, "UTC", data.Timezone)
				assert.Equal(t, "+00:00", data.TimezoneOffset)
				assert.False(t, data.IsDst)
			},
		},
		{
			name:          "Empty timezone (defaults to UTC)",
			timezone:      "",
			expectedError: false,
			validateResult: func(t *testing.T, data *ppb.PanchangamData) {
				assert.Equal(t, "UTC", data.Timezone)
				assert.Equal(t, "+00:00", data.TimezoneOffset)
				assert.False(t, data.IsDst)
			},
		},
		{
			name:          "Invalid timezone",
			timezone:      "Invalid/Timezone",
			expectedError: true,
			validateResult: func(t *testing.T, data *ppb.PanchangamData) {
				// Should not be called
			},
		},
		{
			name:          "Invalid UTC offset",
			timezone:      "+99:99",
			expectedError: true,
			validateResult: func(t *testing.T, data *ppb.PanchangamData) {
				// Should not be called
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ppb.GetPanchangamRequest{
				Date:      "2024-01-15",
				Latitude:  12.9716,
				Longitude: 77.5946,
				Timezone:  tt.timezone,
				Region:    "Karnataka",
			}

			resp, err := server.Get(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.NotNil(t, resp.PanchangamData)

				// Validate timezone fields are present
				assert.NotEmpty(t, resp.PanchangamData.Timezone)
				assert.NotEmpty(t, resp.PanchangamData.TimezoneOffset)

				// Run custom validation if provided
				if tt.validateResult != nil {
					tt.validateResult(t, resp.PanchangamData)
				}
			}
		})
	}
}

func TestPanchangamServer_DSTTransitions(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	tests := []struct {
		name     string
		date     string
		timezone string
		location struct {
			lat  float64
			long float64
		}
	}{
		{
			name:     "New York summer (DST active)",
			date:     "2024-07-01",
			timezone: "America/New_York",
			location: struct {
				lat  float64
				long float64
			}{lat: 40.7128, long: -74.0060},
		},
		{
			name:     "New York winter (DST inactive)",
			date:     "2024-01-01",
			timezone: "America/New_York",
			location: struct {
				lat  float64
				long float64
			}{lat: 40.7128, long: -74.0060},
		},
		{
			name:     "London summer (BST)",
			date:     "2024-07-01",
			timezone: "Europe/London",
			location: struct {
				lat  float64
				long float64
			}{lat: 51.5074, long: -0.1278},
		},
		{
			name:     "London winter (GMT)",
			date:     "2024-01-01",
			timezone: "Europe/London",
			location: struct {
				lat  float64
				long float64
			}{lat: 51.5074, long: -0.1278},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ppb.GetPanchangamRequest{
				Date:      tt.date,
				Latitude:  tt.location.lat,
				Longitude: tt.location.long,
				Timezone:  tt.timezone,
			}

			resp, err := server.Get(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.PanchangamData)

			// Verify timezone info is included
			assert.NotEmpty(t, resp.PanchangamData.Timezone)
			assert.NotEmpty(t, resp.PanchangamData.TimezoneOffset)

			t.Logf("Date: %s, Timezone: %s, Offset: %s, DST: %v",
				tt.date,
				resp.PanchangamData.Timezone,
				resp.PanchangamData.TimezoneOffset,
				resp.PanchangamData.IsDst,
			)
		})
	}
}

func TestPanchangamServer_TimezoneValidation(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	// Test timezone that doesn't match location (should warn but still work)
	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  40.7128,  // New York coordinates
		Longitude: -74.0060, // New York coordinates
		Timezone:  "Asia/Kolkata", // But using India timezone
	}

	resp, err := server.Get(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.PanchangamData)

	// Should still work, just with a warning in logs
	assert.Equal(t, "Asia/Kolkata", resp.PanchangamData.Timezone)
	assert.Equal(t, "+05:30", resp.PanchangamData.TimezoneOffset)
}

func TestPanchangamServer_UTCOffsetFormats(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	tests := []struct {
		name           string
		timezone       string
		expectedOffset string
	}{
		{
			name:           "Positive offset with prefix",
			timezone:       "UTC+05:30",
			expectedOffset: "+05:30",
		},
		{
			name:           "Negative offset with prefix",
			timezone:       "GMT-08:00",
			expectedOffset: "-08:00",
		},
		{
			name:           "Positive offset without prefix",
			timezone:       "+12:00",
			expectedOffset: "+12:00",
		},
		{
			name:           "Negative offset without prefix",
			timezone:       "-05:00",
			expectedOffset: "-05:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ppb.GetPanchangamRequest{
				Date:      "2024-01-15",
				Latitude:  0.0,
				Longitude: 0.0,
				Timezone:  tt.timezone,
			}

			resp, err := server.Get(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.PanchangamData)

			assert.Equal(t, tt.expectedOffset, resp.PanchangamData.TimezoneOffset)
			assert.False(t, resp.PanchangamData.IsDst) // Fixed offsets don't have DST
		})
	}
}
