package ephemeris

import (
	"context"
	"testing"
	"time"
)

func BenchmarkLinearInterpolation(b *testing.B) {
	manager := createTestManager(b)
	config := InterpolationConfig{
		Method:    InterpolationLinear,
		Tolerance: 0.01,
	}
	interpolator := NewInterpolator(manager, config)

	baseJD := TimeToJulianDay(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))
	targetJD := JulianDay(float64(baseJD) + 0.5)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := interpolator.InterpolatePlanetaryPosition(ctx, targetJD, "sun")
		if err != nil {
			b.Fatalf("Interpolation failed: %v", err)
		}
	}
}

func BenchmarkLagrangeInterpolation(b *testing.B) {
	manager := createTestManager(b)
	config := InterpolationConfig{
		Method:    InterpolationLagrange,
		Order:     5,
		Tolerance: 0.001,
	}
	interpolator := NewInterpolator(manager, config)

	baseJD := TimeToJulianDay(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))
	targetJD := JulianDay(float64(baseJD) + 0.5)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := interpolator.InterpolatePlanetaryPosition(ctx, targetJD, "mars")
		if err != nil {
			b.Fatalf("Interpolation failed: %v", err)
		}
	}
}

func BenchmarkCubicSplineInterpolation(b *testing.B) {
	manager := createTestManager(b)
	config := DefaultInterpolationConfig()
	interpolator := NewInterpolator(manager, config)

	baseJD := TimeToJulianDay(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))
	targetJD := JulianDay(float64(baseJD) + 0.5)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := interpolator.InterpolatePlanetaryPosition(ctx, targetJD, "jupiter")
		if err != nil {
			b.Fatalf("Interpolation failed: %v", err)
		}
	}
}

func BenchmarkInterpolateAllPlanets(b *testing.B) {
	manager := createTestManager(b)
	config := DefaultInterpolationConfig()
	interpolator := NewInterpolator(manager, config)

	baseJD := TimeToJulianDay(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))
	targetJD := JulianDay(float64(baseJD) + 0.5)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := interpolator.InterpolatePlanetaryPositions(ctx, targetJD)
		if err != nil {
			b.Fatalf("Interpolation failed: %v", err)
		}
	}
}
