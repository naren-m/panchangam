package ephemeris

import (
	"context"
	"testing"
	"time"
)

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
		_, err := manager.GetSunPosition(ctx, testJD)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := manager.GetSunPosition(ctx, testJD)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
