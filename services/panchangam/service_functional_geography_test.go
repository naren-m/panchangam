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

func TestServiceFunctionalGeographicCoverage(t *testing.T) {
	observability.NewLocalObserver()

	server := NewPanchangamServer()
	ctx := context.Background()

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
				Date:      "2024-06-21",
				Latitude:  loc.latitude,
				Longitude: loc.longitude,
				Timezone:  loc.timezone,
			}

			resp, err := server.Get(ctx, req)

			assert.NoError(t, err, "Valid coordinates should not cause error")
			require.NotNil(t, resp, "Response should not be nil")
			require.NotNil(t, resp.PanchangamData, "Panchangam data should not be nil")

			data := resp.PanchangamData
			assert.Equal(t, req.Date, data.Date, "Date should match")
			assert.NotEmpty(t, data.SunriseTime, "Should have sunrise time")
			assert.NotEmpty(t, data.SunsetTime, "Should have sunset time")
		})
	}
}

func TestServiceFunctionalTimezoneHandling(t *testing.T) {
	observability.NewLocalObserver()

	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	resp1, err1 := server.Get(ctx, req)
	assert.NoError(t, err1)
	require.NotNil(t, resp1)

	req.Timezone = "Asia/Kolkata"
	resp2, err2 := server.Get(ctx, req)
	assert.NoError(t, err2)
	require.NotNil(t, resp2)

	req.Timezone = "Invalid/Timezone"
	_, err3 := server.Get(ctx, req)
	require.Error(t, err3)
	assert.Equal(t, codes.InvalidArgument, status.Code(err3))
	assert.Contains(t, err3.Error(), "invalid timezone")
}
