package ephemeris

import "github.com/naren-m/panchangam/observability"

// InterpolationMethod defines the type of interpolation to use.
type InterpolationMethod string

const (
	// InterpolationLinear uses simple linear interpolation.
	InterpolationLinear InterpolationMethod = "linear"

	// InterpolationLagrange uses Lagrange polynomial interpolation.
	InterpolationLagrange InterpolationMethod = "lagrange"

	// InterpolationCubicSpline uses cubic spline interpolation.
	InterpolationCubicSpline InterpolationMethod = "cubic_spline"
)

// InterpolationConfig holds configuration for interpolation operations.
type InterpolationConfig struct {
	Method       InterpolationMethod
	Order        int
	Tolerance    float64
	MaxCacheSize int
}

// Interpolator provides interpolation methods for planetary positions.
type Interpolator struct {
	manager  *Manager
	config   InterpolationConfig
	observer observability.ObserverInterface
	cache    map[string]*interpolationCache
}

// interpolationCache stores data points for interpolation.
type interpolationCache struct {
	jdPoints       []float64
	positionPoints []Position
	maxSize        int
}

// dataPoint holds a single data point for interpolation.
type dataPoint struct {
	jd       float64
	position Position
}

// DefaultInterpolationConfig returns the default configuration.
func DefaultInterpolationConfig() InterpolationConfig {
	return InterpolationConfig{
		Method:       InterpolationCubicSpline,
		Order:        5,
		Tolerance:    0.0001, // about 0.36 arcseconds
		MaxCacheSize: 100,
	}
}
