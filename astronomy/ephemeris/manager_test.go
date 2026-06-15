package ephemeris

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingCloseProvider struct {
	name     string
	closeErr error
}

func (p failingCloseProvider) GetPlanetaryPositions(ctx context.Context, jd JulianDay) (*PlanetaryPositions, error) {
	return &PlanetaryPositions{JulianDay: jd}, nil
}

func (p failingCloseProvider) GetSunPosition(ctx context.Context, jd JulianDay) (*SolarPosition, error) {
	return &SolarPosition{JulianDay: jd}, nil
}

func (p failingCloseProvider) GetMoonPosition(ctx context.Context, jd JulianDay) (*LunarPosition, error) {
	return &LunarPosition{JulianDay: jd}, nil
}

func (p failingCloseProvider) IsAvailable(ctx context.Context) bool {
	return true
}

func (p failingCloseProvider) GetDataRange() (startJD, endJD JulianDay) {
	return 0, 0
}

func (p failingCloseProvider) GetHealthStatus(ctx context.Context) (*HealthStatus, error) {
	return &HealthStatus{Available: true, Source: p.name}, nil
}

func (p failingCloseProvider) GetProviderName() string {
	return p.name
}

func (p failingCloseProvider) GetVersion() string {
	return "test"
}

func (p failingCloseProvider) Close() error {
	return p.closeErr
}

func TestEphemerisManager(t *testing.T) {
	primary := NewJPLProvider()
	fallback := NewSwissProvider()
	cache := NewMemoryCache(100, 1*time.Hour)

	manager := NewManager(primary, fallback, cache)
	ctx := context.Background()
	testJD := JulianDay(2451545.0)

	t.Run("manager initialization", func(t *testing.T) {
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.primary)
		assert.NotNil(t, manager.fallback)
		assert.NotNil(t, manager.cache)
		assert.NotNil(t, manager.healthChecker)
	})

	t.Run("sun position with caching", func(t *testing.T) {
		position1, err := manager.GetSunPosition(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, position1)

		position2, err := manager.GetSunPosition(ctx, testJD)
		require.NoError(t, err)
		assert.Equal(t, position1, position2)
	})

	t.Run("moon position with caching", func(t *testing.T) {
		position1, err := manager.GetMoonPosition(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, position1)

		position2, err := manager.GetMoonPosition(ctx, testJD)
		require.NoError(t, err)
		assert.Equal(t, position1, position2)
	})

	t.Run("planetary positions with caching", func(t *testing.T) {
		positions1, err := manager.GetPlanetaryPositions(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, positions1)

		positions2, err := manager.GetPlanetaryPositions(ctx, testJD)
		require.NoError(t, err)
		assert.Equal(t, positions1, positions2)
	})

	t.Run("fallback mechanism", func(t *testing.T) {
		nilPrimary := NewManager(nil, fallback, cache)

		position, err := nilPrimary.GetSunPosition(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, position)
	})

	t.Run("health status", func(t *testing.T) {
		statuses, err := manager.GetHealthStatus(ctx)
		require.NoError(t, err)
		assert.Contains(t, statuses, "primary")
		assert.Contains(t, statuses, "fallback")
		assert.True(t, statuses["primary"].Available)
		assert.True(t, statuses["fallback"].Available)
	})

	t.Run("close manager", func(t *testing.T) {
		err := manager.Close()
		assert.NoError(t, err)
	})
}

func TestEphemerisManagerClosePreservesErrorCauses(t *testing.T) {
	primaryErr := errors.New("primary close failed")
	fallbackErr := errors.New("fallback close failed")
	manager := NewManager(
		failingCloseProvider{name: "primary", closeErr: primaryErr},
		failingCloseProvider{name: "fallback", closeErr: fallbackErr},
		nil,
	)

	err := manager.Close()
	require.Error(t, err)
	assert.True(t, errors.Is(err, primaryErr), "expected close error to preserve primary cause")
	assert.True(t, errors.Is(err, fallbackErr), "expected close error to preserve fallback cause")
}
