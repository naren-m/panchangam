package ephemeris

import (
	"context"
	"math"

	"go.opentelemetry.io/otel/attribute"
)

// AnalyzeMotion provides comprehensive analysis of planetary motion.
func (rd *RetrogradeDetector) AnalyzeMotion(ctx context.Context, jd JulianDay, planet string) (*MotionAnalysis, error) {
	ctx, span := rd.observer.CreateSpan(ctx, "retrograde.AnalyzeMotion")
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

	positions, err := rd.manager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	pos, err := rd.interpolator.extractPlanetPosition(positions, planet)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	analysis := &MotionAnalysis{
		JulianDay:     jd,
		Planet:        planet,
		Motion:        motion,
		Speed:         pos.Speed,
		Longitude:     pos.Longitude,
		IsNearStation: math.Abs(pos.Speed) < 0.05,
	}

	nextStation, err := rd.FindPlanetaryStation(ctx, jd, planet, 400)
	if err == nil {
		analysis.NextStation = nextStation
	}

	if motion == MotionRetrograde {
		period, err := rd.FindRetrogradePeriod(ctx, jd, planet)
		if err == nil {
			analysis.CurrentPeriod = period
		}
	}

	analysis.RecentStations = rd.findRecentStations(ctx, jd, planet, 180)

	span.SetAttributes(
		attribute.String("motion", string(motion)),
		attribute.Bool("is_near_station", analysis.IsNearStation),
		attribute.Bool("has_next_station", analysis.NextStation != nil),
	)

	return analysis, nil
}

// findRecentStations finds stations in the past N days.
func (rd *RetrogradeDetector) findRecentStations(ctx context.Context, jd JulianDay, planet string, days int) []PlanetaryStation {
	stations := make([]PlanetaryStation, 0)
	searchJD := jd
	const chunkSize = 30

	for i := 0; i < days/chunkSize; i++ {
		searchJD = JulianDay(float64(searchJD) - float64(chunkSize))

		station, err := rd.FindPlanetaryStation(ctx, searchJD, planet, chunkSize)
		if err == nil && station != nil {
			stations = append(stations, *station)
		}
	}

	return stations
}
