package ephemeris

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPositionAccuracy(t *testing.T) {
	ctx := context.Background()
	j2000 := JulianDay(2451545.0)

	jplProvider := NewJPLProvider()
	swissProvider := NewSwissProvider()

	t.Run("sun position accuracy", func(t *testing.T) {
		jplSun, err := jplProvider.GetSunPosition(ctx, j2000)
		require.NoError(t, err)

		swissSun, err := swissProvider.GetSunPosition(ctx, j2000)
		require.NoError(t, err)

		assert.InDelta(t, 280.0, jplSun.Longitude, 10.0)
		assert.InDelta(t, 280.0, swissSun.Longitude, 10.0)
		assert.InDelta(t, 1.0, jplSun.Distance, 0.1)
		assert.InDelta(t, 1.0, swissSun.Distance, 0.1)
	})

	t.Run("moon position accuracy", func(t *testing.T) {
		jplMoon, err := jplProvider.GetMoonPosition(ctx, j2000)
		require.NoError(t, err)

		swissMoon, err := swissProvider.GetMoonPosition(ctx, j2000)
		require.NoError(t, err)

		assert.True(t, jplMoon.Longitude >= 0 && jplMoon.Longitude <= 360)
		assert.True(t, swissMoon.Longitude >= 0 && swissMoon.Longitude <= 360)
		assert.InDelta(t, 384400.0, jplMoon.Distance, 50000.0)
		assert.InDelta(t, 384400.0, swissMoon.Distance, 50000.0)
	})

	t.Run("planetary motion consistency", func(t *testing.T) {
		testJD1 := JulianDay(2451545.0)
		testJD2 := JulianDay(2451545.0 + 30)

		positions1, err := jplProvider.GetPlanetaryPositions(ctx, testJD1)
		require.NoError(t, err)

		positions2, err := jplProvider.GetPlanetaryPositions(ctx, testJD2)
		require.NoError(t, err)

		mercuryDelta := math.Abs(positions2.Mercury.Longitude - positions1.Mercury.Longitude)
		if mercuryDelta > 180 {
			mercuryDelta = 360 - mercuryDelta
		}

		saturnDelta := math.Abs(positions2.Saturn.Longitude - positions1.Saturn.Longitude)
		if saturnDelta > 180 {
			saturnDelta = 360 - saturnDelta
		}

		assert.True(t, mercuryDelta > saturnDelta,
			"Mercury should move more than Saturn: Mercury=%.2f, Saturn=%.2f",
			mercuryDelta, saturnDelta)
	})
}
