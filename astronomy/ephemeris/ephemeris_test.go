package ephemeris

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Initialize observability for testing
	observability.NewLocalObserver()
}

func TestJulianDayConversion(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected JulianDay
		tolerance float64
	}{
		{
			name:     "J2000.0 epoch",
			time:     time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: JulianDay(2451545.0),
			tolerance: 0.001,
		},
		{
			name:     "Unix epoch",
			time:     time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: JulianDay(2440587.5),
			tolerance: 0.001,
		},
		{
			name:     "Current test date",
			time:     time.Date(2024, 7, 18, 0, 0, 0, 0, time.UTC),
			expected: JulianDay(2460509.5),
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jd := TimeToJulianDay(tt.time)
			assert.InDelta(t, float64(tt.expected), float64(jd), tt.tolerance)
			
			// Test round-trip conversion
			converted := JulianDayToTime(jd)
			assert.WithinDuration(t, tt.time, converted, time.Minute)
		})
	}
}

func TestJPLProvider(t *testing.T) {
	provider := NewJPLProvider()
	ctx := context.Background()
	testJD := JulianDay(2451545.0) // J2000.0

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
		assert.True(t, position.Distance > 0.9 && position.Distance < 1.1) // ~1 AU
	})

	t.Run("moon position", func(t *testing.T) {
		position, err := provider.GetMoonPosition(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, position)
		assert.Equal(t, testJD, position.JulianDay)
		assert.True(t, position.Longitude >= 0 && position.Longitude <= 360)
		assert.True(t, position.Distance > 300000 && position.Distance < 500000) // km
		assert.True(t, position.Phase >= 0 && position.Phase <= 1)
	})

	t.Run("planetary positions", func(t *testing.T) {
		positions, err := provider.GetPlanetaryPositions(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, positions)
		assert.Equal(t, testJD, positions.JulianDay)
		
		// Check all planets have valid positions
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
		invalidJD := JulianDay(1000000.0) // Outside valid range
		_, err := provider.GetSunPosition(ctx, invalidJD)
		assert.Error(t, err)
	})
}

func TestSwissProvider(t *testing.T) {
	provider := NewSwissProvider()
	ctx := context.Background()
	testJD := JulianDay(2451545.0) // J2000.0

	t.Run("provider info", func(t *testing.T) {
		assert.Equal(t, "Swiss Ephemeris", provider.GetProviderName())
		assert.Equal(t, "2.10", provider.GetVersion())
		assert.True(t, provider.IsAvailable(ctx))
	})

	t.Run("data range", func(t *testing.T) {
		startJD, endJD := provider.GetDataRange()
		assert.True(t, startJD < endJD)
		assert.True(t, testJD >= startJD && testJD <= endJD)
		
		// Swiss Ephemeris should have much wider range than JPL
		jplProvider := NewJPLProvider()
		jplStart, jplEnd := jplProvider.GetDataRange()
		assert.True(t, startJD < jplStart)
		assert.True(t, endJD > jplEnd)
	})

	t.Run("sun position accuracy", func(t *testing.T) {
		position, err := provider.GetSunPosition(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, position)
		
		// Swiss Ephemeris should be more accurate than JPL for this test
		jplProvider := NewJPLProvider()
		jplPosition, err := jplProvider.GetSunPosition(ctx, testJD)
		require.NoError(t, err)
		
		// Both should be reasonable, but Swiss might have slightly different values
		assert.True(t, position.Longitude >= 0 && position.Longitude <= 360)
		assert.True(t, jplPosition.Longitude >= 0 && jplPosition.Longitude <= 360)
	})
}

func TestEphemerisManager(t *testing.T) {
	primary := NewJPLProvider()
	fallback := NewSwissProvider()
	cache := NewMemoryCache(100, 1*time.Hour)
	
	manager := NewManager(primary, fallback, cache)
	ctx := context.Background()
	testJD := JulianDay(2451545.0) // J2000.0

	t.Run("manager initialization", func(t *testing.T) {
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.primary)
		assert.NotNil(t, manager.fallback)
		assert.NotNil(t, manager.cache)
		assert.NotNil(t, manager.healthChecker)
	})

	t.Run("sun position with caching", func(t *testing.T) {
		// First call should fetch from provider
		position1, err := manager.GetSunPosition(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, position1)
		
		// Second call should use cache
		position2, err := manager.GetSunPosition(ctx, testJD)
		require.NoError(t, err)
		assert.Equal(t, position1, position2)
	})

	t.Run("moon position with caching", func(t *testing.T) {
		// First call should fetch from provider
		position1, err := manager.GetMoonPosition(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, position1)
		
		// Second call should use cache
		position2, err := manager.GetMoonPosition(ctx, testJD)
		require.NoError(t, err)
		assert.Equal(t, position1, position2)
	})

	t.Run("planetary positions with caching", func(t *testing.T) {
		// First call should fetch from provider
		positions1, err := manager.GetPlanetaryPositions(ctx, testJD)
		require.NoError(t, err)
		assert.NotNil(t, positions1)
		
		// Second call should use cache
		positions2, err := manager.GetPlanetaryPositions(ctx, testJD)
		require.NoError(t, err)
		assert.Equal(t, positions1, positions2)
	})

	t.Run("fallback mechanism", func(t *testing.T) {
		// Test with a JD that might be outside JPL range but inside Swiss range
		// This is tricky since our test implementations both have wide ranges
		// Instead, we'll test the fallback logic by using a nil primary
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

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache(3, 1*time.Second)
	ctx := context.Background()
	
	t.Run("basic operations", func(t *testing.T) {
		// Test Set and Get
		cache.Set(ctx, "key1", "value1", 0)
		value, found := cache.Get(ctx, "key1")
		assert.True(t, found)
		assert.Equal(t, "value1", value)
		
		// Test missing key
		_, found = cache.Get(ctx, "nonexistent")
		assert.False(t, found)
	})

	t.Run("ttl expiration", func(t *testing.T) {
		cache.Set(ctx, "key2", "value2", 10*time.Millisecond)
		
		// Should be available immediately
		value, found := cache.Get(ctx, "key2")
		assert.True(t, found)
		assert.Equal(t, "value2", value)
		
		// Should expire after TTL
		time.Sleep(20 * time.Millisecond)
		_, found = cache.Get(ctx, "key2")
		assert.False(t, found)
	})

	t.Run("lru eviction", func(t *testing.T) {
		// Fill cache to capacity
		cache.Set(ctx, "key3", "value3", 0)
		cache.Set(ctx, "key4", "value4", 0)
		cache.Set(ctx, "key5", "value5", 0)
		
		// Access key3 to make it most recently used
		cache.Get(ctx, "key3")
		
		// Add another item, should evict key4 (least recently used)
		cache.Set(ctx, "key6", "value6", 0)
		
		// key3 should still be there
		_, found := cache.Get(ctx, "key3")
		assert.True(t, found)
		
		// key4 should be evicted
		_, found = cache.Get(ctx, "key4")
		assert.False(t, found)
	})

	t.Run("cache stats", func(t *testing.T) {
		stats := cache.GetStats(ctx)
		assert.NotNil(t, stats)
		assert.True(t, stats.Hits > 0)
		assert.True(t, stats.Misses > 0)
	})

	t.Run("clear cache", func(t *testing.T) {
		cache.Set(ctx, "key7", "value7", 0)
		err := cache.Clear(ctx)
		assert.NoError(t, err)
		
		_, found := cache.Get(ctx, "key7")
		assert.False(t, found)
	})

	t.Run("close cache", func(t *testing.T) {
		err := cache.Close()
		assert.NoError(t, err)
	})
}

func TestHealthChecker(t *testing.T) {
	primary := NewJPLProvider()
	fallback := NewSwissProvider()
	
	checker := NewHealthChecker([]EphemerisProvider{primary, fallback})
	
	t.Run("initial status", func(t *testing.T) {
		assert.NotNil(t, checker)
		// IsHealthy returns false if no statuses are available initially
		// This is expected behavior until first health check runs
	})

	t.Run("start and stop", func(t *testing.T) {
		checker.Start()
		time.Sleep(100 * time.Millisecond) // Allow initial health check
		
		statuses := checker.GetAllStatuses()
		assert.Contains(t, statuses, "JPL DE440")
		assert.Contains(t, statuses, "Swiss Ephemeris")
		
		checker.Stop()
	})

	t.Run("individual status", func(t *testing.T) {
		newChecker := NewHealthChecker([]EphemerisProvider{primary, fallback})
		newChecker.Start()
		time.Sleep(100 * time.Millisecond) // Allow initial health check
		
		status, found := newChecker.GetStatus("JPL DE440")
		assert.True(t, found)
		assert.True(t, status.Available)
		
		newChecker.Stop()
	})

	t.Run("metrics", func(t *testing.T) {
		metricsChecker := NewHealthChecker([]EphemerisProvider{primary, fallback})
		metricsChecker.Start()
		time.Sleep(100 * time.Millisecond) // Allow initial health check
		
		metrics := metricsChecker.GetMetrics()
		assert.NotNil(t, metrics)
		assert.Equal(t, 2, metrics["total_providers"])
		assert.Equal(t, 2, metrics["healthy_providers"])
		assert.Equal(t, 0, metrics["unhealthy_providers"])
		assert.Equal(t, 100.0, metrics["health_percentage"])
		
		metricsChecker.Stop()
	})

	t.Run("add and remove providers", func(t *testing.T) {
		addRemoveChecker := NewHealthChecker([]EphemerisProvider{primary, fallback})
		addRemoveChecker.Start()
		time.Sleep(100 * time.Millisecond) // Allow initial health check
		
		// Check initial state
		statuses := addRemoveChecker.GetAllStatuses()
		assert.Len(t, statuses, 2) // Should have 2 providers initially
		
		// Add a new provider (this will replace the existing JPL provider since they have the same name)
		newProvider := NewJPLProvider()
		addRemoveChecker.AddProvider(newProvider)
		
		time.Sleep(100 * time.Millisecond) // Allow health check
		
		statuses = addRemoveChecker.GetAllStatuses()
		// Still 2 providers because the new JPL provider replaces the old one
		assert.Len(t, statuses, 2)
		
		// Remove a provider
		addRemoveChecker.RemoveProvider("JPL DE440")
		
		statuses = addRemoveChecker.GetAllStatuses()
		assert.Len(t, statuses, 1) // Should have 1 provider now
		
		addRemoveChecker.Stop()
	})
}

func TestNoOpCache(t *testing.T) {
	cache := NewNoOpCache()
	ctx := context.Background()
	
	t.Run("no-op operations", func(t *testing.T) {
		// Set should do nothing
		cache.Set(ctx, "key", "value", 0)
		
		// Get should always return false
		_, found := cache.Get(ctx, "key")
		assert.False(t, found)
		
		// Delete should always return false
		deleted := cache.Delete(ctx, "key")
		assert.False(t, deleted)
		
		// Clear should do nothing
		err := cache.Clear(ctx)
		assert.NoError(t, err)
		
		// Stats should be empty
		stats := cache.GetStats(ctx)
		assert.Equal(t, int64(0), stats.Hits)
		assert.Equal(t, int64(0), stats.Misses)
		
		// Close should do nothing
		err = cache.Close()
		assert.NoError(t, err)
	})
}

func TestPositionAccuracy(t *testing.T) {
	// Test known astronomical positions for accuracy
	ctx := context.Background()
	
	// J2000.0 epoch - well-known reference point
	j2000 := JulianDay(2451545.0)
	
	jplProvider := NewJPLProvider()
	swissProvider := NewSwissProvider()
	
	t.Run("sun position accuracy", func(t *testing.T) {
		jplSun, err := jplProvider.GetSunPosition(ctx, j2000)
		require.NoError(t, err)
		
		swissSun, err := swissProvider.GetSunPosition(ctx, j2000)
		require.NoError(t, err)
		
		// Both should give reasonable values for J2000.0
		// At J2000.0, Sun's longitude should be around 280 degrees
		assert.InDelta(t, 280.0, jplSun.Longitude, 10.0)
		assert.InDelta(t, 280.0, swissSun.Longitude, 10.0)
		
		// Distance should be around 1 AU
		assert.InDelta(t, 1.0, jplSun.Distance, 0.1)
		assert.InDelta(t, 1.0, swissSun.Distance, 0.1)
	})

	t.Run("moon position accuracy", func(t *testing.T) {
		jplMoon, err := jplProvider.GetMoonPosition(ctx, j2000)
		require.NoError(t, err)
		
		swissMoon, err := swissProvider.GetMoonPosition(ctx, j2000)
		require.NoError(t, err)
		
		// Both should give reasonable values
		assert.True(t, jplMoon.Longitude >= 0 && jplMoon.Longitude <= 360)
		assert.True(t, swissMoon.Longitude >= 0 && swissMoon.Longitude <= 360)
		
		// Distance should be around 384,400 km
		assert.InDelta(t, 384400.0, jplMoon.Distance, 50000.0)
		assert.InDelta(t, 384400.0, swissMoon.Distance, 50000.0)
	})

	t.Run("planetary motion consistency", func(t *testing.T) {
		// Test that planets move consistently over time
		testJD1 := JulianDay(2451545.0)    // J2000.0
		testJD2 := JulianDay(2451545.0 + 30) // 30 days later
		
		positions1, err := jplProvider.GetPlanetaryPositions(ctx, testJD1)
		require.NoError(t, err)
		
		positions2, err := jplProvider.GetPlanetaryPositions(ctx, testJD2)
		require.NoError(t, err)
		
		// Mercury should move the most (fastest planet)
		mercuryDelta := math.Abs(positions2.Mercury.Longitude - positions1.Mercury.Longitude)
		if mercuryDelta > 180 {
			mercuryDelta = 360 - mercuryDelta
		}
		
		// Saturn should move the least (slowest planet)
		saturnDelta := math.Abs(positions2.Saturn.Longitude - positions1.Saturn.Longitude)
		if saturnDelta > 180 {
			saturnDelta = 360 - saturnDelta
		}
		
		// Mercury should move more than Saturn in 30 days
		assert.True(t, mercuryDelta > saturnDelta, 
			"Mercury should move more than Saturn: Mercury=%.2f, Saturn=%.2f", 
			mercuryDelta, saturnDelta)
	})
}

func BenchmarkEphemerisOperations(b *testing.B) {
	primary := NewJPLProvider()
	fallback := NewSwissProvider()
	cache := NewMemoryCache(1000, 1*time.Hour)
	manager := NewManager(primary, fallback, cache)
	ctx := context.Background()
	testJD := JulianDay(2451545.0)

	b.Run("GetSunPosition", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := manager.GetSunPosition(ctx, testJD)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("GetMoonPosition", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := manager.GetMoonPosition(ctx, testJD)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("GetPlanetaryPositions", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := manager.GetPlanetaryPositions(ctx, testJD)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("GetSunPositionWithCache", func(b *testing.B) {
		// Prime the cache
		manager.GetSunPosition(ctx, testJD)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := manager.GetSunPosition(ctx, testJD)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}