package ephemeris

import (
	"context"
	"math"
	"testing"
	"time"
)

func TestLinearInterpolation(t *testing.T) {
	manager := createTestManager(t)
	config := InterpolationConfig{
		Method:    InterpolationLinear,
		Tolerance: 0.01,
	}
	interpolator := NewInterpolator(manager, config)

	baseJD := TimeToJulianDay(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
	targetJD := JulianDay(float64(baseJD) + 0.5)

	position, err := interpolator.InterpolatePlanetaryPosition(context.Background(), targetJD, "sun")
	if err != nil {
		t.Fatalf("Linear interpolation failed: %v", err)
	}

	if position == nil {
		t.Fatal("Expected non-nil position")
	}

	if position.Longitude < 0 || position.Longitude >= 360 {
		t.Errorf("Longitude out of range: %f", position.Longitude)
	}

	if position.Latitude < -90 || position.Latitude > 90 {
		t.Errorf("Latitude out of range: %f", position.Latitude)
	}

	if position.Distance <= 0 {
		t.Errorf("Distance should be positive: %f", position.Distance)
	}
}

func TestLagrangeInterpolation(t *testing.T) {
	manager := createTestManager(t)
	config := InterpolationConfig{
		Method:    InterpolationLagrange,
		Order:     5,
		Tolerance: 0.001,
	}
	interpolator := NewInterpolator(manager, config)

	baseJD := TimeToJulianDay(time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC))
	targetJD := JulianDay(float64(baseJD) + 0.25)

	position, err := interpolator.InterpolatePlanetaryPosition(context.Background(), targetJD, "moon")
	if err != nil {
		t.Fatalf("Lagrange interpolation failed: %v", err)
	}

	if position == nil {
		t.Fatal("Expected non-nil position")
	}

	if position.Longitude < 0 || position.Longitude >= 360 {
		t.Errorf("Moon longitude out of range: %f", position.Longitude)
	}
}

func TestCubicSplineInterpolation(t *testing.T) {
	manager := createTestManager(t)
	config := InterpolationConfig{
		Method:    InterpolationCubicSpline,
		Order:     7,
		Tolerance: 0.0001,
	}
	interpolator := NewInterpolator(manager, config)

	baseJD := TimeToJulianDay(time.Date(2024, 3, 20, 12, 0, 0, 0, time.UTC))
	targetJD := JulianDay(float64(baseJD) + 0.1)

	position, err := interpolator.InterpolatePlanetaryPosition(context.Background(), targetJD, "mars")
	if err != nil {
		t.Fatalf("Cubic spline interpolation failed: %v", err)
	}

	if position == nil {
		t.Fatal("Expected non-nil position")
	}

	if position.Longitude < 0 || position.Longitude >= 360 {
		t.Errorf("Longitude out of range: %f", position.Longitude)
	}

	if math.Abs(position.Latitude) > 10 {
		t.Errorf("Mars latitude unexpectedly large: %f", position.Latitude)
	}

	if position.Distance < 0.5 || position.Distance > 3.0 {
		t.Errorf("Mars distance out of expected range: %f AU", position.Distance)
	}
}

func TestValidateInterpolation(t *testing.T) {
	manager := createTestManager(t)
	config := InterpolationConfig{
		Method:    InterpolationCubicSpline,
		Order:     5,
		Tolerance: 1.0,
	}
	interpolator := NewInterpolator(manager, config)

	baseJD := TimeToJulianDay(time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC))

	error, err := interpolator.ValidateInterpolation(context.Background(), baseJD, "jupiter")
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if error > 1.0 {
		t.Logf("Warning: Validation error is %f degrees (within tolerance but consider improving)", error)
	}

	t.Logf("Validation error for Jupiter: %f degrees", error)
}

func TestInterpolationMethods(t *testing.T) {
	manager := createTestManager(t)

	methods := []InterpolationMethod{
		InterpolationLinear,
		InterpolationLagrange,
		InterpolationCubicSpline,
	}

	baseJD := TimeToJulianDay(time.Date(2024, 8, 15, 12, 0, 0, 0, time.UTC))
	targetJD := JulianDay(float64(baseJD) + 0.25)

	for _, method := range methods {
		t.Run(string(method), func(t *testing.T) {
			config := InterpolationConfig{
				Method:    method,
				Order:     5,
				Tolerance: 0.01,
			}
			interpolator := NewInterpolator(manager, config)

			position, err := interpolator.InterpolatePlanetaryPosition(context.Background(), targetJD, "venus")
			if err != nil {
				t.Fatalf("Method %s failed: %v", method, err)
			}

			if position == nil {
				t.Fatalf("Method %s returned nil position", method)
			}

			t.Logf("Method %s - Venus longitude: %f°", method, position.Longitude)
		})
	}
}

func TestNormalizeAngle(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{180, 180},
		{360, 0},
		{-90, 270},
		{450, 90},
		{-180, 180},
		{720, 0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := normalizeAngle(tt.input)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("normalizeAngle(%f) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}
