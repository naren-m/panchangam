package ephemeris

import (
	"context"
	"fmt"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

// NewInterpolator creates a new interpolator
func NewInterpolator(manager *Manager, config InterpolationConfig) *Interpolator {
	return &Interpolator{
		manager:  manager,
		config:   config,
		observer: observability.Observer(),
		cache:    make(map[string]*interpolationCache),
	}
}

// InterpolatePlanetaryPosition calculates planetary position at a specific JD using interpolation
func (i *Interpolator) InterpolatePlanetaryPosition(ctx context.Context, jd JulianDay, planet string) (*Position, error) {
	ctx, span := i.observer.CreateSpan(ctx, "interpolator.InterpolatePlanetaryPosition")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("planet", planet),
		attribute.String("method", string(i.config.Method)),
	)

	// Get surrounding data points
	points, err := i.getDataPoints(ctx, jd, planet)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get data points: %w", err)
	}

	// Perform interpolation based on method
	var position Position
	switch i.config.Method {
	case InterpolationLinear:
		position, err = i.linearInterpolation(points, float64(jd))
	case InterpolationLagrange:
		position, err = i.lagrangeInterpolation(points, float64(jd))
	case InterpolationCubicSpline:
		position, err = i.cubicSplineInterpolation(points, float64(jd))
	default:
		return nil, fmt.Errorf("unsupported interpolation method: %s", i.config.Method)
	}

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Normalize longitude to 0-360 range
	position.Longitude = normalizeAngle(position.Longitude)

	span.SetAttributes(
		attribute.Float64("interpolated_longitude", position.Longitude),
		attribute.Float64("interpolated_latitude", position.Latitude),
		attribute.Bool("success", true),
	)

	return &position, nil
}

// GetInterpolationMethod returns the current interpolation method
func (i *Interpolator) GetInterpolationMethod() InterpolationMethod {
	return i.config.Method
}

// SetInterpolationMethod sets the interpolation method
func (i *Interpolator) SetInterpolationMethod(method InterpolationMethod) {
	i.config.Method = method
}

// GetInterpolationConfig returns the current configuration
func (i *Interpolator) GetInterpolationConfig() InterpolationConfig {
	return i.config
}
