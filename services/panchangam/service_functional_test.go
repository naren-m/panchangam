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

func TestServiceFunctionalBasicRequest(t *testing.T) {
	observability.NewLocalObserver()

	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
		Region:    "India",
	}

	resp, err := server.Get(ctx, req)

	assert.NoError(t, err, "Service should not return error for valid request")
	require.NotNil(t, resp, "Response should not be nil")
	require.NotNil(t, resp.PanchangamData, "Panchangam data should not be nil")

	data := resp.PanchangamData
	assert.Equal(t, req.Date, data.Date, "Response date should match request date")
	assert.NotEmpty(t, data.Tithi, "Tithi should not be empty")
	assert.NotEmpty(t, data.Nakshatra, "Nakshatra should not be empty")
	assert.NotEmpty(t, data.Yoga, "Yoga should not be empty")
	assert.NotEmpty(t, data.Karana, "Karana should not be empty")
	assert.NotEmpty(t, data.SunriseTime, "Sunrise time should not be empty")
	assert.NotEmpty(t, data.SunsetTime, "Sunset time should not be empty")

	_, err = time.Parse("15:04:05", data.SunriseTime)
	assert.NoError(t, err, "Sunrise time should be in valid format")

	_, err = time.Parse("15:04:05", data.SunsetTime)
	assert.NoError(t, err, "Sunset time should be in valid format")

	assert.NotNil(t, data.Events, "Events should not be nil")
}
