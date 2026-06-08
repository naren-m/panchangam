package ephemeris

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
)

// InterpolatePlanetaryPositions interpolates all planetary positions at once.
func (i *Interpolator) InterpolatePlanetaryPositions(ctx context.Context, jd JulianDay) (*PlanetaryPositions, error) {
	ctx, span := i.observer.CreateSpan(ctx, "interpolator.InterpolatePlanetaryPositions")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("method", string(i.config.Method)),
	)

	planets := []string{"sun", "moon", "mercury", "venus", "mars", "jupiter", "saturn", "uranus", "neptune", "pluto"}
	positions := &PlanetaryPositions{
		JulianDay: jd,
	}

	for _, planet := range planets {
		pos, err := i.InterpolatePlanetaryPosition(ctx, jd, planet)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("failed to interpolate %s: %w", planet, err)
		}

		switch planet {
		case "sun":
			positions.Sun = *pos
		case "moon":
			positions.Moon = *pos
		case "mercury":
			positions.Mercury = *pos
		case "venus":
			positions.Venus = *pos
		case "mars":
			positions.Mars = *pos
		case "jupiter":
			positions.Jupiter = *pos
		case "saturn":
			positions.Saturn = *pos
		case "uranus":
			positions.Uranus = *pos
		case "neptune":
			positions.Neptune = *pos
		case "pluto":
			positions.Pluto = *pos
		}
	}

	span.SetAttributes(attribute.Bool("success", true))
	return positions, nil
}
