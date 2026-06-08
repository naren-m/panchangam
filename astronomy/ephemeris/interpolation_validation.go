package ephemeris

import (
	"context"
	"fmt"
	"math"

	"go.opentelemetry.io/otel/attribute"
)

// ValidateInterpolation validates interpolation accuracy against direct calculation.
func (i *Interpolator) ValidateInterpolation(ctx context.Context, jd JulianDay, planet string) (float64, error) {
	ctx, span := i.observer.CreateSpan(ctx, "interpolator.ValidateInterpolation")
	defer span.End()

	interpPos, err := i.InterpolatePlanetaryPosition(ctx, jd, planet)
	if err != nil {
		return 0, fmt.Errorf("interpolation failed: %w", err)
	}

	actualPositions, err := i.manager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		return 0, fmt.Errorf("failed to get actual positions: %w", err)
	}

	actualPos, err := i.extractPlanetPosition(actualPositions, planet)
	if err != nil {
		return 0, err
	}

	lonError := math.Abs(interpPos.Longitude - actualPos.Longitude)
	if lonError > 180 {
		lonError = 360 - lonError
	}

	latError := math.Abs(interpPos.Latitude - actualPos.Latitude)
	distError := math.Abs(interpPos.Distance-actualPos.Distance) / actualPos.Distance * 100
	totalError := lonError + latError*0.5 + distError*0.1

	span.SetAttributes(
		attribute.Float64("longitude_error_deg", lonError),
		attribute.Float64("latitude_error_deg", latError),
		attribute.Float64("distance_error_percent", distError),
		attribute.Float64("total_error", totalError),
		attribute.Bool("within_tolerance", totalError <= i.config.Tolerance),
	)

	return totalError, nil
}
