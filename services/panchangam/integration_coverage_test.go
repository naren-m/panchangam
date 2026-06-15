package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceCalculationIntegration(t *testing.T) {
	observability.NewLocalObserver()

	t.Run("Integration_Service_With_Real_Calculations", func(t *testing.T) {
		ctx := context.Background()
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		location := astronomy.Location{
			Latitude:  12.9716,
			Longitude: 77.5946,
		}

		validateCalculationModulesWork(t, ctx, testDate, location)
		validateServiceStructureReady(t, ctx)
		validateIntegrationReadiness(t)
	})

	t.Run("Integration_Mock_Real_Service_Flow", func(t *testing.T) {
		ctx := context.Background()

		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
			Region:    "India",
		}

		mockRealServiceFlow(t, ctx, req)
	})

	t.Run("Integration_Performance_With_Calculations", func(t *testing.T) {
		ctx := context.Background()
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

		start := time.Now()

		sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		require.NoError(t, err, "Astronomy calculation should work")
		require.NotNil(t, sunTimes, "Sun times should be calculated")

		astronomyDuration := time.Since(start)

		assert.True(t, astronomyDuration < 100*time.Millisecond,
			"Combined calculations should be <100ms, got %v", astronomyDuration)

		t.Logf("Integration Performance: Astronomy calculation %v", astronomyDuration)
	})
}
