package ephemeris

import (
	"context"
	"fmt"
	"math"

	"go.opentelemetry.io/otel/attribute"
)

// FindRetrogradePeriod finds the complete retrograde period containing the given JD.
func (rd *RetrogradeDetector) FindRetrogradePeriod(ctx context.Context, jd JulianDay, planet string) (*RetrogradePeriod, error) {
	ctx, span := rd.observer.CreateSpan(ctx, "retrograde.FindRetrogradePeriod")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("planet", planet),
	)

	motion, err := rd.DetectRetrogradeMotion(ctx, jd, planet)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	if motion != MotionRetrograde {
		return nil, fmt.Errorf("planet %s is not retrograde at JD %f", planet, jd)
	}

	startStation, err := rd.findStationBackward(ctx, jd, planet, 200)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to find start station: %w", err)
	}

	endStation, err := rd.findStationForward(ctx, jd, planet, 200)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to find end station: %w", err)
	}

	period := &RetrogradePeriod{
		Planet:         planet,
		StartJD:        startStation.JulianDay,
		EndJD:          endStation.JulianDay,
		StartTime:      startStation.Time,
		EndTime:        endStation.Time,
		StartLongitude: startStation.Longitude,
		EndLongitude:   endStation.Longitude,
		Duration:       endStation.Time.Sub(startStation.Time),
	}

	period.MaxRetroDistance = math.Abs(endStation.Longitude - startStation.Longitude)
	if period.MaxRetroDistance > 180 {
		period.MaxRetroDistance = 360 - period.MaxRetroDistance
	}

	span.SetAttributes(
		attribute.Float64("period_start_jd", float64(period.StartJD)),
		attribute.Float64("period_end_jd", float64(period.EndJD)),
		attribute.Float64("duration_days", period.Duration.Hours()/24),
	)

	return period, nil
}

// findStationBackward searches backward for a station.
func (rd *RetrogradeDetector) findStationBackward(ctx context.Context, jd JulianDay, planet string, maxDays int) (*PlanetaryStation, error) {
	const stepSize = -1.0 // 1 day backward

	for i := 0; i < maxDays; i++ {
		searchJD := JulianDay(float64(jd) + float64(i)*stepSize)

		motion, err := rd.DetectRetrogradeMotion(ctx, searchJD, planet)
		if err != nil {
			continue
		}

		if motion == MotionDirect || motion == MotionStationary {
			return rd.FindPlanetaryStation(ctx, searchJD, planet, 10)
		}
	}

	return nil, fmt.Errorf("no station found in %d days backward search", maxDays)
}

// findStationForward searches forward for a station.
func (rd *RetrogradeDetector) findStationForward(ctx context.Context, jd JulianDay, planet string, maxDays int) (*PlanetaryStation, error) {
	const stepSize = 1.0 // 1 day forward

	for i := 0; i < maxDays; i++ {
		searchJD := JulianDay(float64(jd) + float64(i)*stepSize)

		motion, err := rd.DetectRetrogradeMotion(ctx, searchJD, planet)
		if err != nil {
			continue
		}

		if motion == MotionDirect || motion == MotionStationary {
			return rd.FindPlanetaryStation(ctx, searchJD, planet, 10)
		}
	}

	return nil, fmt.Errorf("no station found in %d days forward search", maxDays)
}
