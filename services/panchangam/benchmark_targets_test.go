package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/stretchr/testify/require"
)

func BenchmarkFeaturePerformanceTargets(b *testing.B) {
	b.Run("Target_Individual_Calculator_50ms", func(b *testing.B) {
		benchmarkIndividualCalculatorTarget(b)
	})

	b.Run("Target_All_Elements_100ms", func(b *testing.B) {
		benchmarkAllElementsTarget(b)
	})

	b.Run("Target_Service_Response_500ms", func(b *testing.B) {
		benchmarkServiceResponseTarget(b)
	})

	b.Run("Target_End_to_End_500ms", func(b *testing.B) {
		benchmarkEndToEndTarget(b)
	})
}

func benchmarkIndividualCalculatorTarget(b *testing.B) {
	ctx := context.Background()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()

		_, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		require.NoError(b, err)

		duration := time.Since(start)
		if duration > 50*time.Millisecond {
			b.Errorf("Individual calculator exceeded 50ms target: %v", duration)
		}
	}
}

func benchmarkAllElementsTarget(b *testing.B) {
	ctx := context.Background()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()

		_, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		require.NoError(b, err)

		mockCalculateAllElements(b, ctx, testDate)

		duration := time.Since(start)
		if duration > 100*time.Millisecond {
			b.Errorf("All elements calculation exceeded 100ms target: %v", duration)
		}
	}
}
