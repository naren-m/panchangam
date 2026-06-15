package ephemeris

import (
	"context"
	"fmt"
	"math"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

// RetrogradeDetector detects retrograde motion and stationary points
type RetrogradeDetector struct {
	manager      *Manager
	interpolator *Interpolator
	observer     observability.ObserverInterface
}

// NewRetrogradeDetector creates a new retrograde detector
func NewRetrogradeDetector(manager *Manager) *RetrogradeDetector {
	config := DefaultInterpolationConfig()
	interpolator := NewInterpolator(manager, config)

	return &RetrogradeDetector{
		manager:      manager,
		interpolator: interpolator,
		observer:     observability.Observer(),
	}
}

// DetectRetrogradeMotion determines if a planet is in retrograde motion
func (rd *RetrogradeDetector) DetectRetrogradeMotion(ctx context.Context, jd JulianDay, planet string) (RetrogradeMotion, error) {
	ctx, span := rd.observer.CreateSpan(ctx, "retrograde.DetectRetrogradeMotion")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("planet", planet),
	)

	// Get planetary position
	positions, err := rd.manager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to get positions: %w", err)
	}

	pos, err := rd.interpolator.extractPlanetPosition(positions, planet)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	// Check speed to determine motion
	// Negative speed indicates retrograde motion
	// Speed near zero indicates stationary
	const stationaryThreshold = 0.01 // degrees per day

	var motion RetrogradeMotion
	if math.Abs(pos.Speed) < stationaryThreshold {
		motion = MotionStationary
	} else if pos.Speed < 0 {
		motion = MotionRetrograde
	} else {
		motion = MotionDirect
	}

	span.SetAttributes(
		attribute.String("motion", string(motion)),
		attribute.Float64("speed", pos.Speed),
		attribute.Float64("longitude", pos.Longitude),
	)

	return motion, nil
}
