package ephemeris

import (
	"context"
	"testing"
	"time"
)

func TestInterpolationAllPlanets(t *testing.T) {
	manager := createTestManager(t)
	config := DefaultInterpolationConfig()
	interpolator := NewInterpolator(manager, config)

	baseJD := TimeToJulianDay(time.Date(2024, 7, 4, 12, 0, 0, 0, time.UTC))
	targetJD := JulianDay(float64(baseJD) + 0.3)

	planets := []string{"sun", "moon", "mercury", "venus", "mars", "jupiter", "saturn"}

	for _, planet := range planets {
		t.Run(planet, func(t *testing.T) {
			position, err := interpolator.InterpolatePlanetaryPosition(context.Background(), targetJD, planet)
			if err != nil {
				t.Fatalf("Failed to interpolate %s: %v", planet, err)
			}

			if position == nil {
				t.Fatalf("Expected non-nil position for %s", planet)
			}

			if position.Longitude < 0 || position.Longitude >= 360 {
				t.Errorf("%s longitude out of range: %f", planet, position.Longitude)
			}

			if position.Distance <= 0 {
				t.Errorf("%s distance should be positive: %f", planet, position.Distance)
			}
		})
	}
}

func TestInterpolatePlanetaryPositions(t *testing.T) {
	manager := createTestManager(t)
	config := DefaultInterpolationConfig()
	interpolator := NewInterpolator(manager, config)

	baseJD := TimeToJulianDay(time.Date(2024, 9, 15, 12, 0, 0, 0, time.UTC))
	targetJD := JulianDay(float64(baseJD) + 0.4)

	positions, err := interpolator.InterpolatePlanetaryPositions(context.Background(), targetJD)
	if err != nil {
		t.Fatalf("Failed to interpolate all planetary positions: %v", err)
	}

	if positions == nil {
		t.Fatal("Expected non-nil positions")
	}

	if positions.JulianDay != targetJD {
		t.Errorf("JulianDay mismatch: got %f, want %f", positions.JulianDay, targetJD)
	}

	planetPositions := []Position{
		positions.Sun, positions.Moon, positions.Mercury, positions.Venus,
		positions.Mars, positions.Jupiter, positions.Saturn,
	}

	for i, pos := range planetPositions {
		if pos.Longitude < 0 || pos.Longitude >= 360 {
			t.Errorf("Planet %d longitude out of range: %f", i, pos.Longitude)
		}

		if pos.Distance <= 0 {
			t.Errorf("Planet %d distance should be positive: %f", i, pos.Distance)
		}
	}
}

func TestExtractPlanetPosition(t *testing.T) {
	manager := createTestManager(t)
	config := DefaultInterpolationConfig()
	interpolator := NewInterpolator(manager, config)

	baseJD := TimeToJulianDay(time.Date(2024, 2, 14, 12, 0, 0, 0, time.UTC))
	positions, err := manager.GetPlanetaryPositions(context.Background(), baseJD)
	if err != nil {
		t.Fatalf("Failed to get positions: %v", err)
	}

	planets := []string{"sun", "moon", "mercury", "venus", "mars", "jupiter", "saturn", "uranus", "neptune", "pluto"}

	for _, planet := range planets {
		t.Run(planet, func(t *testing.T) {
			pos, err := interpolator.extractPlanetPosition(positions, planet)
			if err != nil {
				t.Fatalf("Failed to extract %s position: %v", planet, err)
			}

			if pos == nil {
				t.Fatalf("Expected non-nil position for %s", planet)
			}
		})
	}

	_, err = interpolator.extractPlanetPosition(positions, "invalid")
	if err == nil {
		t.Error("Expected error for invalid planet")
	}
}
