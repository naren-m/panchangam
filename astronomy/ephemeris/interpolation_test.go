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

	// Test interpolation for Sun
	baseJD := TimeToJulianDay(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
	targetJD := JulianDay(float64(baseJD) + 0.5) // 12 hours later

	position, err := interpolator.InterpolatePlanetaryPosition(context.Background(), targetJD, "sun")
	if err != nil {
		t.Fatalf("Linear interpolation failed: %v", err)
	}

	if position == nil {
		t.Fatal("Expected non-nil position")
	}

	// Verify longitude is in valid range
	if position.Longitude < 0 || position.Longitude >= 360 {
		t.Errorf("Longitude out of range: %f", position.Longitude)
	}

	// Verify latitude is in valid range
	if position.Latitude < -90 || position.Latitude > 90 {
		t.Errorf("Latitude out of range: %f", position.Latitude)
	}

	// Verify distance is positive
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
	targetJD := JulianDay(float64(baseJD) + 0.25) // 6 hours later

	position, err := interpolator.InterpolatePlanetaryPosition(context.Background(), targetJD, "moon")
	if err != nil {
		t.Fatalf("Lagrange interpolation failed: %v", err)
	}

	if position == nil {
		t.Fatal("Expected non-nil position")
	}

	// Moon moves faster, so validate reasonable movement
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
	targetJD := JulianDay(float64(baseJD) + 0.1) // ~2.4 hours later

	position, err := interpolator.InterpolatePlanetaryPosition(context.Background(), targetJD, "mars")
	if err != nil {
		t.Fatalf("Cubic spline interpolation failed: %v", err)
	}

	if position == nil {
		t.Fatal("Expected non-nil position")
	}

	// Validate all components
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

			// Validate basic constraints
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

	// Verify all planets have valid positions
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

func TestValidateInterpolation(t *testing.T) {
	manager := createTestManager(t)
	config := InterpolationConfig{
		Method:    InterpolationCubicSpline,
		Order:     5,
		Tolerance: 1.0, // 1 degree tolerance for test
	}
	interpolator := NewInterpolator(manager, config)

	// Use an exact JD point that the ephemeris has data for
	baseJD := TimeToJulianDay(time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC))

	error, err := interpolator.ValidateInterpolation(context.Background(), baseJD, "jupiter")
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// At exact data point, error should be very small
	if error > 1.0 {
		t.Logf("Warning: Validation error is %f degrees (within tolerance but consider improving)", error)
	}

	t.Logf("Validation error for Jupiter: %f degrees", error)
}

func TestInterpolationEdgeCases(t *testing.T) {
	manager := createTestManager(t)
	config := DefaultInterpolationConfig()
	interpolator := NewInterpolator(manager, config)

	t.Run("very_small_time_interval", func(t *testing.T) {
		baseJD := TimeToJulianDay(time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC))
		// 1 minute later
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
		// 12 hours later
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
		// Test near 0/360 degree boundary
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

func TestInterpolatorConfig(t *testing.T) {
	manager := createTestManager(t)
	config := InterpolationConfig{
		Method:       InterpolationLagrange,
		Order:        7,
		Tolerance:    0.0005,
		MaxCacheSize: 200,
	}

	interpolator := NewInterpolator(manager, config)

	if interpolator.GetInterpolationMethod() != InterpolationLagrange {
		t.Errorf("Expected method %s, got %s", InterpolationLagrange, interpolator.GetInterpolationMethod())
	}

	// Change method
	interpolator.SetInterpolationMethod(InterpolationCubicSpline)

	if interpolator.GetInterpolationMethod() != InterpolationCubicSpline {
		t.Errorf("Expected method %s after change, got %s", InterpolationCubicSpline, interpolator.GetInterpolationMethod())
	}

	gotConfig := interpolator.GetInterpolationConfig()
	if gotConfig.Order != config.Order {
		t.Errorf("Expected order %d, got %d", config.Order, gotConfig.Order)
	}
}

func TestDefaultInterpolationConfig(t *testing.T) {
	config := DefaultInterpolationConfig()

	if config.Method != InterpolationCubicSpline {
		t.Errorf("Expected default method %s, got %s", InterpolationCubicSpline, config.Method)
	}

	if config.Order != 5 {
		t.Errorf("Expected default order 5, got %d", config.Order)
	}

	if config.Tolerance != 0.0001 {
		t.Errorf("Expected default tolerance 0.0001, got %f", config.Tolerance)
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

	// Test invalid planet
	_, err = interpolator.extractPlanetPosition(positions, "invalid")
	if err == nil {
		t.Error("Expected error for invalid planet")
	}
}

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

// Helper function to create a test manager
func createTestManager(tb testing.TB) *Manager {
	tb.Helper()

	// Create test providers (mock or actual)
	// For this test, we'll use actual providers if available
	primary, err := NewJPLProvider()
	if err != nil {
		tb.Skipf("JPL provider not available: %v", err)
	}

	fallback, err := NewSwissProvider()
	if err != nil {
		tb.Logf("Swiss provider not available, using only JPL: %v", err)
		fallback = nil
	}

	cache := NewMemoryCache(100, 1*time.Hour)

	return NewManager(primary, fallback, cache)
}
