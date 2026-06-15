package ephemeris

import "testing"

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
