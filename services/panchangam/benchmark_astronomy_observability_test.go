package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

func benchmarkAstronomyCalculations(b *testing.B) {
	ctx := context.Background()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		start := time.Now()
		sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		duration := time.Since(start)

		require.NoError(b, err)
		require.NotNil(b, sunTimes)
		require.True(b, sunTimes.Sunrise.Before(sunTimes.Sunset))

		if i == 0 {
			b.Logf("Astronomy calculation time: %v", duration)
		}
	}
}

func benchmarkObservability(b *testing.B) {
	observer := observability.Observer()
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		start := time.Now()

		spanCtx, span := observer.CreateSpan(ctx, "benchmark_span")
		_ = spanCtx

		span.SetAttributes(attribute.String("benchmark", "test"))
		span.AddEvent("benchmark_event")
		span.End()

		duration := time.Since(start)
		if i == 0 {
			b.Logf("Observability overhead: %v", duration)
		}
	}
}
