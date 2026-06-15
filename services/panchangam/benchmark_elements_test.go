package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/stretchr/testify/require"
)

func benchmarkAllPanchangamElements(b *testing.B) {
	ctx := context.Background()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		start := time.Now()

		sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		require.NoError(b, err)
		require.NotNil(b, sunTimes)

		mockCalculateAllElements(b, ctx, testDate)

		duration := time.Since(start)
		if i == 0 {
			b.Logf("All Panchangam elements calculation: %v", duration)
		}
	}
}

func BenchmarkCalculationAccuracy(b *testing.B) {
	ctx := context.Background()
	location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		testDate := time.Date(2024, 1, 1+i%28, 0, 0, 0, 0, time.UTC)

		sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		require.NoError(b, err)
		require.NotNil(b, sunTimes)

		require.True(b, sunTimes.Sunrise.Before(sunTimes.Sunset), "Sunrise should be before sunset")
		require.True(b, sunTimes.Sunrise.Hour() >= 4 && sunTimes.Sunrise.Hour() <= 8, "Sunrise should be reasonable")
		require.True(b, sunTimes.Sunset.Hour() >= 16 && sunTimes.Sunset.Hour() <= 20, "Sunset should be reasonable")
	}
}

func mockCalculateAllElements(b *testing.B, ctx context.Context, date time.Time) {
	b.Helper()

	for i := 0; i < 5; i++ {
		_ = time.Now().Add(time.Microsecond)
	}
}
