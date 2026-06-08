package ephemeris

import (
	"context"
	"testing"
	"time"
)

func TestInterpolationEdgeCases(t *testing.T) {
	manager := createTestManager(t)
	config := DefaultInterpolationConfig()
	interpolator := NewInterpolator(manager, config)

	t.Run("very_small_time_interval", func(t *testing.T) {
		baseJD := TimeToJulianDay(time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC))
		targetJD := JulianDay(float64(baseJD) + 1.0/1440.0)

		position, err := interpolator.InterpolatePlanetaryPosition(context.Background(), targetJD, "sun")
		if err != nil {
			t.Fatalf("Failed on small interval: %v", err)
		}

		if position == nil {
			t.Fatal("Expected non-nil position")
		}
	})

	t.Run("large_time_interval", func(t *testing.T) {
		baseJD := TimeToJulianDay(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		targetJD := JulianDay(float64(baseJD) + 0.5)

		position, err := interpolator.InterpolatePlanetaryPosition(context.Background(), targetJD, "saturn")
		if err != nil {
			t.Fatalf("Failed on large interval: %v", err)
		}

		if position == nil {
			t.Fatal("Expected non-nil position")
		}
	})

	t.Run("angle_wrapping_near_zero", func(t *testing.T) {
		baseJD := TimeToJulianDay(time.Date(2024, 12, 21, 12, 0, 0, 0, time.UTC))
		targetJD := JulianDay(float64(baseJD) + 0.2)

		position, err := interpolator.InterpolatePlanetaryPosition(context.Background(), targetJD, "sun")
		if err != nil {
			t.Fatalf("Failed with angle wrapping: %v", err)
		}

		if position.Longitude < 0 || position.Longitude >= 360 {
			t.Errorf("Angle wrapping failed: %f", position.Longitude)
		}
	})
}
