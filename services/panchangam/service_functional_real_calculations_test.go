package panchangam

import (
	"context"
	"testing"

	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/require"
)

func TestServiceWithRealCalculations(t *testing.T) {
	t.Run("Functional_Real_Calculation_Integration", func(t *testing.T) {
		server := NewPanchangamServer()

		resp, err := server.Get(context.Background(), &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
			Region:    "India",
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.PanchangamData)

		data := resp.PanchangamData
		require.Equal(t, "2024-01-15", data.Date)
		require.NotEmpty(t, data.Tithi)
		require.NotEmpty(t, data.Nakshatra)
		require.NotEmpty(t, data.Yoga)
		require.NotEmpty(t, data.Karana)
		require.NotEmpty(t, data.SunriseTime)
		require.NotEmpty(t, data.SunsetTime)
	})
}
