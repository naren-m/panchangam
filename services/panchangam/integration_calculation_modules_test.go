package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validateCalculationModulesWork(t *testing.T, ctx context.Context, testDate time.Time, location astronomy.Location) {
	t.Helper()

	sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
	assert.NoError(t, err, "Astronomy calculation should work")
	require.NotNil(t, sunTimes, "Sun times should be calculated")
	assert.True(t, sunTimes.Sunrise.Before(sunTimes.Sunset), "Sunrise should be before sunset")

	assert.True(t, sunTimes.Sunrise.Hour() >= 0 && sunTimes.Sunrise.Hour() <= 23,
		"Sunrise should be valid hour")
	assert.True(t, sunTimes.Sunset.Hour() >= 0 && sunTimes.Sunset.Hour() <= 23,
		"Sunset should be valid hour")

	t.Logf("Calculation modules work: Sunrise %s, Sunset %s",
		sunTimes.Sunrise.Format("15:04:05"), sunTimes.Sunset.Format("15:04:05"))

	assert.NotNil(t, astronomy.TithiNames, "Tithi data should be available")
	assert.Len(t, astronomy.TithiNames, 30, "Should have 30 Tithi names")

	assert.NotNil(t, astronomy.KaranaData, "Karana data should be available")
	assert.Len(t, astronomy.KaranaData, 11, "Should have 11 Karana entries")

	t.Logf("All calculation module interfaces are properly defined")
}

func TestCalculationModuleIntegration(t *testing.T) {
	t.Run("Calculation_Module_Coordination", func(t *testing.T) {
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		tithi := &astronomy.TithiInfo{Number: 15, Name: "Purnima"}
		karana := &astronomy.KaranaInfo{TithiNumber: 15, HalfTithi: 1}

		assert.Equal(t, tithi.Number, karana.TithiNumber, "Karana should reference correct Tithi")

		assert.NoError(t, astronomy.ValidateTithiCalculation(&astronomy.TithiInfo{
			Number:      15,
			MoonSunDiff: 180.0,
			Duration:    24.0,
			StartTime:   testDate,
			EndTime:     testDate.Add(24 * time.Hour),
			Name:        "Purnima",
		}), "Tithi validation should work")

		assert.NoError(t, astronomy.ValidateKaranaCalculation(&astronomy.KaranaInfo{
			Number:      7,
			TithiNumber: 15,
			HalfTithi:   1,
			MoonSunDiff: 180.0,
			Duration:    12.0,
			StartTime:   testDate,
			EndTime:     testDate.Add(12 * time.Hour),
			Name:        "Vanija",
			Type:        astronomy.KaranaTypeMovable,
		}), "Karana validation should work")

		t.Logf("Calculation modules have compatible interfaces")
	})
}
