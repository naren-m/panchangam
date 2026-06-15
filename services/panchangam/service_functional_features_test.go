package panchangam

import (
	"context"
	"testing"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceFeatureCoverage(t *testing.T) {
	observability.NewLocalObserver()

	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:              "2024-01-15",
		Latitude:          12.9716,
		Longitude:         77.5946,
		Timezone:          "Asia/Kolkata",
		Region:            "India",
		CalculationMethod: "traditional",
		Locale:            "en",
	}

	resp, err := server.Get(ctx, req)
	require.NoError(t, err, "Service should handle all request parameters")
	require.NotNil(t, resp, "Response should not be nil")
	require.NotNil(t, resp.PanchangamData, "Panchangam data should not be nil")

	data := resp.PanchangamData
	features := map[string]string{
		"TITHI_001":     data.Tithi,
		"NAKSHATRA_001": data.Nakshatra,
		"YOGA_001":      data.Yoga,
		"KARANA_001":    data.Karana,
	}

	for featureID, value := range features {
		assert.NotEmpty(t, value, "Feature %s should have a value", featureID)
	}

	assert.NotEmpty(t, data.SunriseTime, "ASTRONOMY_001: Sunrise should be calculated")
	assert.NotEmpty(t, data.SunsetTime, "ASTRONOMY_001: Sunset should be calculated")

	assert.Equal(t, req.Date, data.Date, "SERVICE_001: Date should be processed correctly")
	assert.NotNil(t, data.Events, "SERVICE_001: Events should be included")

	assert.IsType(t, "", data.Date, "Date should be string")
	assert.IsType(t, "", data.Tithi, "Tithi should be string")
	assert.IsType(t, "", data.Nakshatra, "Nakshatra should be string")
	assert.IsType(t, "", data.Yoga, "Yoga should be string")
	assert.IsType(t, "", data.Karana, "Karana should be string")
	assert.IsType(t, "", data.SunriseTime, "Sunrise time should be string")
	assert.IsType(t, "", data.SunsetTime, "Sunset time should be string")
	assert.IsType(t, []*ppb.PanchangamEvent{}, data.Events, "Events should be array")
}
