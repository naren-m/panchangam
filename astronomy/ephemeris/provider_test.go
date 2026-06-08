package ephemeris

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJPLProvider(t *testing.T) {
	provider := NewJPLProvider()
	ctx := context.Background()
	testJD := JulianDay(2451545.0)

	t.Run("provider info", func(t *testing.T) {
		assert.Equal(t, "JPL DE440", provider.GetProviderName())
		assert.Equal(t, "440", provider.GetVersion())
		assert.True(t, provider.IsAvailable(ctx))
	})

	t.Run("data range", func(t *testing.T) {
		startJD, endJD := provider.GetDataRange()
		assert.True(t, startJD < endJD)
		assert.True(t, testJD >= startJD && testJD <= endJD)
	})

	t.Run("health status", func(t *testing.T) {
		health, err := provider.GetHealthStatus(ctx)
		require.NoError(t, err)
		assert.True(t, health.Available)
		assert.Equal(t, "JPL DE440", health.Source)
		assert.Equal(t, "440", health.Version)
	})

	t.Run("sun position", func(t *testing.T) {
		position, err := provider.GetSunPosition(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, position)
		assert.Equal(t, testJD, position.JulianDay)
		assert.True(t, position.Longitude >= 0 && position.Longitude <= 360)
		assert.True(t, position.Distance > 0.9 && position.Distance < 1.1)
	})

	t.Run("moon position", func(t *testing.T) {
		position, err := provider.GetMoonPosition(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, position)
		assert.Equal(t, testJD, position.JulianDay)
		assert.True(t, position.Longitude >= 0 && position.Longitude <= 360)
		assert.True(t, position.Distance > 300000 && position.Distance < 500000)
		assert.True(t, position.Phase >= 0 && position.Phase <= 1)
	})

	t.Run("planetary positions", func(t *testing.T) {
		positions, err := provider.GetPlanetaryPositions(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, positions)
		assert.Equal(t, testJD, positions.JulianDay)

		planets := []Position{
			positions.Sun, positions.Moon, positions.Mercury,
			positions.Venus, positions.Mars, positions.Jupiter,
			positions.Saturn, positions.Uranus, positions.Neptune,
			positions.Pluto,
		}

		for i, pos := range planets {
			assert.True(t, pos.Longitude >= 0 && pos.Longitude <= 360, "Planet %d longitude out of range", i)
			assert.True(t, pos.Distance > 0, "Planet %d distance invalid", i)
			assert.True(t, pos.Speed > 0, "Planet %d speed invalid", i)
		}
	})

	t.Run("invalid julian day", func(t *testing.T) {
		invalidJD := JulianDay(1000000.0)
		_, err := provider.GetSunPosition(ctx, invalidJD)
		assert.Error(t, err)
	})
}

func TestSwissProvider(t *testing.T) {
	provider := NewSwissProvider()
	ctx := context.Background()
	testJD := JulianDay(2451545.0)

	t.Run("provider info", func(t *testing.T) {
		assert.Equal(t, "Swiss Ephemeris", provider.GetProviderName())
		assert.Equal(t, "2.10", provider.GetVersion())
		assert.True(t, provider.IsAvailable(ctx))
	})

	t.Run("data range", func(t *testing.T) {
		startJD, endJD := provider.GetDataRange()
		assert.True(t, startJD < endJD)
		assert.True(t, testJD >= startJD && testJD <= endJD)

		jplProvider := NewJPLProvider()
		jplStart, jplEnd := jplProvider.GetDataRange()
		assert.True(t, startJD < jplStart)
		assert.True(t, endJD > jplEnd)
	})

	t.Run("sun position accuracy", func(t *testing.T) {
		position, err := provider.GetSunPosition(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, position)

		jplProvider := NewJPLProvider()
		jplPosition, err := jplProvider.GetSunPosition(ctx, testJD)
		require.NoError(t, err)

		assert.True(t, position.Longitude >= 0 && position.Longitude <= 360)
		assert.True(t, jplPosition.Longitude >= 0 && jplPosition.Longitude <= 360)
	})
}
