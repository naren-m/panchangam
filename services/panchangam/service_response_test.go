package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPanchangamServer_SunriseSunsetFormat(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-06-21",
		Latitude:  40.7128,
		Longitude: -74.0060,
		Timezone:  "America/New_York",
	}

	resp, err := server.Get(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.PanchangamData)

	sunrise := resp.PanchangamData.SunriseTime
	sunset := resp.PanchangamData.SunsetTime

	assert.Regexp(t, `^\d{2}:\d{2}:\d{2}$`, sunrise, "Sunrise time should be in HH:MM:SS format")
	assert.Regexp(t, `^\d{2}:\d{2}:\d{2}$`, sunset, "Sunset time should be in HH:MM:SS format")

	_, err = time.Parse("15:04:05", sunrise)
	assert.NoError(t, err, "Sunrise time should be parseable")

	_, err = time.Parse("15:04:05", sunset)
	assert.NoError(t, err, "Sunset time should be parseable")
}
