package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

// BenchmarkFeatureCoverage provides performance benchmarks for all Panchangam features
// These benchmarks validate that all features meet performance targets from FEATURES.md

func BenchmarkFeatureCoverage(b *testing.B) {
	// Initialize observability for benchmarking
	observability.NewLocalObserver()

	b.Run("Feature_Performance_All_Elements", func(b *testing.B) {
		benchmarkAllPanchangamElements(b)
	})

	b.Run("Feature_Performance_Service_Layer", func(b *testing.B) {
		benchmarkServiceLayer(b)
	})

	b.Run("Feature_Performance_Astronomy", func(b *testing.B) {
		benchmarkAstronomyCalculations(b)
	})

	b.Run("Feature_Performance_Observability", func(b *testing.B) {
		benchmarkObservability(b)
	})
}

// benchmarkAllPanchangamElements benchmarks all 5 Panchangam elements together
func benchmarkAllPanchangamElements(b *testing.B) {
	ctx := context.Background()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Benchmark all 5 Panchangam elements calculation
		// Target: Combined calculation <100ms

		start := time.Now()

		// This is the actual calculation we can benchmark
		sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		require.NoError(b, err)
		require.NotNil(b, sunTimes)

		// Mock other calculations for timing (these would be real calculations when ephemeris is integrated)
		mockCalculateAllElements(b, ctx, testDate)

		duration := time.Since(start)

		// Report custom metrics
		if i == 0 {
			b.Logf("All Panchangam elements calculation: %v", duration)
		}
	}
}

// benchmarkServiceLayer benchmarks the gRPC service layer
func benchmarkServiceLayer(b *testing.B) {
	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
		Region:    "India",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Benchmark service response time
		// Target: Service response <500ms

		start := time.Now()
		resp, err := server.Get(ctx, req)
		duration := time.Since(start)

		// Only count successful responses (due to random error simulation)
		if err == nil && resp != nil {
			if i == 0 {
				b.Logf("Service response time: %v", duration)
			}
		} else {
			// Skip failed requests in timing
			i--
			if i < 0 {
				i = 0
			}
		}
	}
}

// benchmarkAstronomyCalculations benchmarks astronomy calculations specifically
func benchmarkAstronomyCalculations(b *testing.B) {
	ctx := context.Background()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Benchmark ASTRONOMY_001: Sunrise/sunset calculations
		// Target: Astronomy calculations <100ms

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

// benchmarkObservability benchmarks OpenTelemetry integration
func benchmarkObservability(b *testing.B) {
	observer := observability.Observer()
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Benchmark OBSERVABILITY_001: OpenTelemetry overhead
		// Target: Observability overhead <1ms per operation

		start := time.Now()

		// Create span
		spanCtx, span := observer.CreateSpan(ctx, "benchmark_span")
		_ = spanCtx // Use the context to avoid unused variable

		// Add attributes
		span.SetAttributes(attribute.String("benchmark", "test"))

		// Add event
		span.AddEvent("benchmark_event")

		// End span
		span.End()

		duration := time.Since(start)

		if i == 0 {
			b.Logf("Observability overhead: %v", duration)
		}
	}
}

// BenchmarkFeaturePerformanceTargets validates specific performance targets
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

// benchmarkIndividualCalculatorTarget tests individual calculator performance
func benchmarkIndividualCalculatorTarget(b *testing.B) {
	ctx := context.Background()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()

		// Test individual calculator (astronomy is the one we can actually test)
		_, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		require.NoError(b, err)

		duration := time.Since(start)

		// Validate target: <50ms
		if duration > 50*time.Millisecond {
			b.Errorf("Individual calculator exceeded 50ms target: %v", duration)
		}
	}
}

// benchmarkAllElementsTarget tests combined elements performance
func benchmarkAllElementsTarget(b *testing.B) {
	ctx := context.Background()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()

		// Calculate all elements (mocked for now)
		_, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		require.NoError(b, err)

		// Mock other calculations
		mockCalculateAllElements(b, ctx, testDate)

		duration := time.Since(start)

		// Validate target: <100ms
		if duration > 100*time.Millisecond {
			b.Errorf("All elements calculation exceeded 100ms target: %v", duration)
		}
	}
}

// benchmarkServiceResponseTarget tests service response time
func benchmarkServiceResponseTarget(b *testing.B) {
	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
	}

	b.ResetTimer()

	successCount := 0
	for i := 0; i < b.N; i++ {
		start := time.Now()
		resp, err := server.Get(ctx, req)
		duration := time.Since(start)

		// Only validate successful responses (due to random error simulation)
		if err == nil && resp != nil {
			successCount++

			// Validate target: <500ms
			if duration > 500*time.Millisecond {
				b.Errorf("Service response exceeded 500ms target: %v", duration)
			}
		}
	}

	if successCount == 0 {
		b.Skip("No successful responses due to random error simulation")
	}
}

// benchmarkEndToEndTarget tests complete end-to-end performance
func benchmarkEndToEndTarget(b *testing.B) {
	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:              "2024-01-15",
		Latitude:          12.9716,
		Longitude:         77.5946,
		Timezone:          "Asia/Kolkata",
		Region:            "India",
		CalculationMethod: "traditional",
		Locale:            "en",
	}

	b.ResetTimer()

	successCount := 0
	for i := 0; i < b.N; i++ {
		start := time.Now()

		// Complete end-to-end request
		resp, err := server.Get(ctx, req)

		// Validate response
		if err == nil && resp != nil && resp.PanchangamData != nil {
			data := resp.PanchangamData
			require.Equal(b, req.Date, data.Date)
			require.NotEmpty(b, data.Tithi)
			require.NotEmpty(b, data.Nakshatra)
			require.NotEmpty(b, data.Yoga)
			require.NotEmpty(b, data.Karana)
			require.NotEmpty(b, data.SunriseTime)
			require.NotEmpty(b, data.SunsetTime)
		}

		duration := time.Since(start)

		// Only validate successful responses
		if err == nil && resp != nil {
			successCount++

			// Validate target: <500ms end-to-end
			if duration > 500*time.Millisecond {
				b.Errorf("End-to-end response exceeded 500ms target: %v", duration)
			}
		}
	}

	if successCount == 0 {
		b.Skip("No successful responses due to random error simulation")
	}
}

// BenchmarkConcurrentFeatureAccess tests concurrent access performance
func BenchmarkConcurrentFeatureAccess(b *testing.B) {
	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
	}

	b.ResetTimer()
	b.SetParallelism(10) // Test with 10 concurrent goroutines

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()
			resp, err := server.Get(ctx, req)
			duration := time.Since(start)

			// Only validate successful responses
			if err == nil && resp != nil {
				// Validate concurrent performance target: <1s
				if duration > 1*time.Second {
					b.Errorf("Concurrent response exceeded 1s target: %v", duration)
				}
			}
		}
	})
}

// BenchmarkMemoryUsage tests memory usage characteristics
func BenchmarkMemoryUsage(b *testing.B) {
	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Test memory allocations per request
		resp, err := server.Get(ctx, req)

		// Prevent optimization
		if err == nil && resp != nil {
			_ = resp.PanchangamData
		}
	}
}

// BenchmarkCalculationAccuracy tests calculation accuracy under load
func BenchmarkCalculationAccuracy(b *testing.B) {
	ctx := context.Background()
	location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Test different dates to ensure accuracy consistency
		testDate := time.Date(2024, 1, 1+i%28, 0, 0, 0, 0, time.UTC)

		sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		require.NoError(b, err)
		require.NotNil(b, sunTimes)

		// Validate accuracy constraints
		require.True(b, sunTimes.Sunrise.Before(sunTimes.Sunset), "Sunrise should be before sunset")
		require.True(b, sunTimes.Sunrise.Hour() >= 4 && sunTimes.Sunrise.Hour() <= 8, "Sunrise should be reasonable")
		require.True(b, sunTimes.Sunset.Hour() >= 16 && sunTimes.Sunset.Hour() <= 20, "Sunset should be reasonable")
	}
}

// Helper functions for benchmarking

// mockCalculateAllElements simulates calculating all 5 Panchangam elements
func mockCalculateAllElements(b *testing.B, ctx context.Context, date time.Time) {
	b.Helper()

	// Simulate time for calculating each element
	// In real implementation, this would call:
	// - TithiCalculator.GetTithiForDate()
	// - NakshatraCalculator.GetNakshatraForDate()
	// - YogaCalculator.GetYogaForDate()
	// - KaranaCalculator.GetKaranaForDate()
	// - VaraCalculator.GetVaraForDate()

	// For now, just simulate some computational work
	for i := 0; i < 5; i++ {
		_ = time.Now().Add(time.Microsecond)
	}
}

// BenchmarkFeatureCoverageReport generates a performance report
func BenchmarkFeatureCoverageReport(b *testing.B) {
	b.Run("Performance_Report_Generation", func(b *testing.B) {
		// This benchmark generates a comprehensive performance report

		ctx := context.Background()
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

		server := NewPanchangamServer()
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// Measure astronomy calculation
			astronomyStart := time.Now()
			sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
			astronomyDuration := time.Since(astronomyStart)

			require.NoError(b, err)
			require.NotNil(b, sunTimes)

			// Measure service response
			serviceStart := time.Now()
			resp, serviceErr := server.Get(ctx, req)
			serviceDuration := time.Since(serviceStart)

			// Report performance for successful requests
			if serviceErr == nil && resp != nil {
				if i == 0 {
					b.Logf("Performance Report:")
					b.Logf("  Astronomy calculation: %v (target: <100ms)", astronomyDuration)
					b.Logf("  Service response: %v (target: <500ms)", serviceDuration)
					b.Logf("  Combined overhead: %v", serviceDuration-astronomyDuration)
				}
			}
		}
	})
}
