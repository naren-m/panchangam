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

func TestPanchangamServer_Get(t *testing.T) {
	// Initialize observability for testing
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	tests := []struct {
		name         string
		request      *ppb.GetPanchangamRequest
		validateFunc func(t *testing.T, resp *ppb.GetPanchangamResponse, err error)
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
			name: "Invalid timezone - returns clear error",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timezone:  "Invalid/Timezone",
			},
			validateFunc: func(t *testing.T, resp *ppb.GetPanchangamResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				assert.Contains(t, err.Error(), "invalid timezone")
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
