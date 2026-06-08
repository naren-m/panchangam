package ephemeris

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// GetRetrogradePlanets returns all planets currently in retrograde motion.
func (rd *RetrogradeDetector) GetRetrogradePlanets(ctx context.Context, jd JulianDay) ([]string, error) {
	ctx, span := rd.observer.CreateSpan(ctx, "retrograde.GetRetrogradePlanets")
	defer span.End()

	planets := []string{"mercury", "venus", "mars", "jupiter", "saturn", "uranus", "neptune", "pluto"}
	retrogradePlanets := make([]string, 0)

	for _, planet := range planets {
		motion, err := rd.DetectRetrogradeMotion(ctx, jd, planet)
		if err != nil {
			continue
		}

		if motion == MotionRetrograde {
			retrogradePlanets = append(retrogradePlanets, planet)
		}
	}

	span.SetAttributes(
		attribute.Int("retrograde_count", len(retrogradePlanets)),
		attribute.StringSlice("retrograde_planets", retrogradePlanets),
	)

	return retrogradePlanets, nil
}

// ValidateKnownRetrograde validates detection against known retrograde periods.
func (rd *RetrogradeDetector) ValidateKnownRetrograde(ctx context.Context, planet string, knownStartJD, knownEndJD JulianDay) (bool, error) {
	ctx, span := rd.observer.CreateSpan(ctx, "retrograde.ValidateKnownRetrograde")
	defer span.End()

	midJD := (knownStartJD + knownEndJD) / 2
	motion, err := rd.DetectRetrogradeMotion(ctx, midJD, planet)
	if err != nil {
		span.RecordError(err)
		return false, err
	}

	isValid := motion == MotionRetrograde

	span.SetAttributes(
		attribute.String("planet", planet),
		attribute.Float64("known_start_jd", float64(knownStartJD)),
		attribute.Float64("known_end_jd", float64(knownEndJD)),
		attribute.Bool("validation_result", isValid),
	)

	return isValid, nil
}
