package ephemeris

import (
	"context"
	"fmt"
	"math"

	"go.opentelemetry.io/otel/attribute"
)

// FindPlanetaryStation finds the next stationary point for a planet.
func (rd *RetrogradeDetector) FindPlanetaryStation(ctx context.Context, startJD JulianDay, planet string, searchDays int) (*PlanetaryStation, error) {
	ctx, span := rd.observer.CreateSpan(ctx, "retrograde.FindPlanetaryStation")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("start_jd", float64(startJD)),
		attribute.String("planet", planet),
		attribute.Int("search_days", searchDays),
	)

	const sampleInterval = 0.25 // 6 hours
	maxSamples := int(float64(searchDays) / sampleInterval)

	positions, err := rd.manager.GetPlanetaryPositions(ctx, startJD)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get initial positions: %w", err)
	}

	pos, err := rd.interpolator.extractPlanetPosition(positions, planet)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	prevSpeed := pos.Speed
	prevJD := startJD

	for i := 1; i < maxSamples; i++ {
		currentJD := JulianDay(float64(startJD) + float64(i)*sampleInterval)

		positions, err := rd.manager.GetPlanetaryPositions(ctx, currentJD)
		if err != nil {
			continue
		}

		pos, err := rd.interpolator.extractPlanetPosition(positions, planet)
		if err != nil {
			continue
		}

		if prevSpeed*pos.Speed < 0 || math.Abs(pos.Speed) < 0.01 {
			stationJD, err := rd.refineStation(ctx, prevJD, currentJD, planet)
			if err != nil {
				span.RecordError(err)
				return nil, err
			}

			stationPos, err := rd.manager.GetPlanetaryPositions(ctx, stationJD)
			if err != nil {
				span.RecordError(err)
				return nil, err
			}

			stationPlanetPos, err := rd.interpolator.extractPlanetPosition(stationPos, planet)
			if err != nil {
				span.RecordError(err)
				return nil, err
			}

			stationType := StationDirect
			if prevSpeed > 0 && pos.Speed < 0 {
				stationType = StationRetrograde
			}

			station := &PlanetaryStation{
				Planet:      planet,
				JulianDay:   stationJD,
				Time:        JulianDayToTime(stationJD),
				Longitude:   stationPlanetPos.Longitude,
				StationType: stationType,
				Speed:       stationPlanetPos.Speed,
			}

			span.SetAttributes(
				attribute.Float64("station_jd", float64(station.JulianDay)),
				attribute.String("station_type", string(station.StationType)),
				attribute.Bool("found", true),
			)

			return station, nil
		}

		prevSpeed = pos.Speed
		prevJD = currentJD
	}

	span.SetAttributes(attribute.Bool("found", false))
	return nil, fmt.Errorf("no station found within %d days", searchDays)
}

// refineStation uses bisection to find the exact JD of a stationary point.
func (rd *RetrogradeDetector) refineStation(ctx context.Context, jd1, jd2 JulianDay, planet string) (JulianDay, error) {
	const tolerance = 0.001 // about 1.4 minutes
	const maxIterations = 20

	for i := 0; i < maxIterations; i++ {
		if float64(jd2-jd1) < tolerance {
			return (jd1 + jd2) / 2, nil
		}

		midJD := (jd1 + jd2) / 2

		positions, err := rd.manager.GetPlanetaryPositions(ctx, midJD)
		if err != nil {
			return 0, err
		}

		pos, err := rd.interpolator.extractPlanetPosition(positions, planet)
		if err != nil {
			return 0, err
		}

		positions1, err := rd.manager.GetPlanetaryPositions(ctx, jd1)
		if err != nil {
			return 0, err
		}

		pos1, err := rd.interpolator.extractPlanetPosition(positions1, planet)
		if err != nil {
			return 0, err
		}

		if pos1.Speed*pos.Speed < 0 {
			jd2 = midJD
		} else {
			jd1 = midJD
		}
	}

	return (jd1 + jd2) / 2, nil
}
