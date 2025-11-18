package ephemeris

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

// RetrogradeMotion indicates whether a planet is in retrograde motion
type RetrogradeMotion string

const (
	// MotionDirect indicates normal forward motion
	MotionDirect RetrogradeMotion = "direct"

	// MotionRetrograde indicates backward (retrograde) motion
	MotionRetrograde RetrogradeMotion = "retrograde"

	// MotionStationary indicates the planet is at a stationary point
	MotionStationary RetrogradeMotion = "stationary"
)

// PlanetaryStation represents a stationary point where planet changes direction
type PlanetaryStation struct {
	Planet     string           // Planet name
	JulianDay  JulianDay        // JD when station occurs
	Time       time.Time        // Station time
	Longitude  float64          // Ecliptic longitude at station (degrees)
	StationType StationType     // Type of station (direct to retrograde, or vice versa)
	Speed      float64          // Speed at station (should be near zero)
}

// StationType indicates the type of planetary station
type StationType string

const (
	// StationRetrograde indicates planet is becoming retrograde (direct to retrograde)
	StationRetrograde StationType = "station_retrograde"

	// StationDirect indicates planet is becoming direct (retrograde to direct)
	StationDirect StationType = "station_direct"
)

// RetrogradePeriod represents a period of retrograde motion
type RetrogradePeriod struct {
	Planet           string        // Planet name
	StartJD          JulianDay     // Start of retrograde period
	EndJD            JulianDay     // End of retrograde period
	StartTime        time.Time     // Start time
	EndTime          time.Time     // End time
	StartLongitude   float64       // Longitude at start (degrees)
	EndLongitude     float64       // Longitude at end (degrees)
	Duration         time.Duration // Duration of retrograde period
	MaxRetroDistance float64       // Maximum retrograde distance traveled (degrees)
}

// MotionAnalysis provides comprehensive analysis of planetary motion
type MotionAnalysis struct {
	JulianDay        JulianDay        // Current JD
	Planet           string           // Planet name
	Motion           RetrogradeMotion // Current motion state
	Speed            float64          // Current speed (degrees/day)
	Longitude        float64          // Current longitude (degrees)
	IsNearStation    bool             // Whether near a stationary point
	NextStation      *PlanetaryStation // Next upcoming station (if any)
	CurrentPeriod    *RetrogradePeriod // Current retrograde period (if retrograde)
	RecentStations   []PlanetaryStation // Recent stations (last 6 months)
}

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

// FindPlanetaryStation finds the next stationary point for a planet
func (rd *RetrogradeDetector) FindPlanetaryStation(ctx context.Context, startJD JulianDay, planet string, searchDays int) (*PlanetaryStation, error) {
	ctx, span := rd.observer.CreateSpan(ctx, "retrograde.FindPlanetaryStation")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("start_jd", float64(startJD)),
		attribute.String("planet", planet),
		attribute.Int("search_days", searchDays),
	)

	// Sample positions at regular intervals
	const sampleInterval = 0.25 // 6 hours
	maxSamples := int(float64(searchDays) / sampleInterval)

	var prevSpeed float64
	var prevJD JulianDay

	// Get initial speed
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

	prevSpeed = pos.Speed
	prevJD = startJD

	// Search for speed sign change (indicates station)
	for i := 1; i < maxSamples; i++ {
		currentJD := JulianDay(float64(startJD) + float64(i)*sampleInterval)

		positions, err := rd.manager.GetPlanetaryPositions(ctx, currentJD)
		if err != nil {
			continue // Skip errors, keep searching
		}

		pos, err := rd.interpolator.extractPlanetPosition(positions, planet)
		if err != nil {
			continue
		}

		// Check for speed sign change
		if prevSpeed*pos.Speed < 0 || math.Abs(pos.Speed) < 0.01 {
			// Found a station! Refine the exact time using bisection
			stationJD, err := rd.refineStation(ctx, prevJD, currentJD, planet)
			if err != nil {
				span.RecordError(err)
				return nil, err
			}

			// Get position at station
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

			// Determine station type
			var stationType StationType
			if prevSpeed > 0 && pos.Speed < 0 {
				stationType = StationRetrograde
			} else {
				stationType = StationDirect
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
				attribute.Float64("station_jd", float64(stationJD)),
				attribute.String("station_type", string(stationType)),
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

// refineStation uses bisection to find the exact JD of a stationary point
func (rd *RetrogradeDetector) refineStation(ctx context.Context, jd1, jd2 JulianDay, planet string) (JulianDay, error) {
	const tolerance = 0.001 // ~1.4 minutes
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

		// Get speed at jd1
		positions1, err := rd.manager.GetPlanetaryPositions(ctx, jd1)
		if err != nil {
			return 0, err
		}

		pos1, err := rd.interpolator.extractPlanetPosition(positions1, planet)
		if err != nil {
			return 0, err
		}

		// Check which half contains the station
		if pos1.Speed*pos.Speed < 0 {
			jd2 = midJD
		} else {
			jd1 = midJD
		}
	}

	return (jd1 + jd2) / 2, nil
}

// FindRetrogradePeriod finds the complete retrograde period containing the given JD
func (rd *RetrogradeDetector) FindRetrogradePeriod(ctx context.Context, jd JulianDay, planet string) (*RetrogradePeriod, error) {
	ctx, span := rd.observer.CreateSpan(ctx, "retrograde.FindRetrogradePeriod")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("planet", planet),
	)

	// First, check if planet is retrograde at given JD
	motion, err := rd.DetectRetrogradeMotion(ctx, jd, planet)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	if motion != MotionRetrograde {
		return nil, fmt.Errorf("planet %s is not retrograde at JD %f", planet, jd)
	}

	// Search backward for station retrograde
	startStation, err := rd.findStationBackward(ctx, jd, planet, 200) // Search up to 200 days back
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to find start station: %w", err)
	}

	// Search forward for station direct
	endStation, err := rd.findStationForward(ctx, jd, planet, 200) // Search up to 200 days forward
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

	// Calculate max retrograde distance
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

// findStationBackward searches backward for a station
func (rd *RetrogradeDetector) findStationBackward(ctx context.Context, jd JulianDay, planet string, maxDays int) (*PlanetaryStation, error) {
	const stepSize = -1.0 // 1 day backward
	for i := 0; i < maxDays; i++ {
		searchJD := JulianDay(float64(jd) + float64(i)*stepSize)

		motion, err := rd.DetectRetrogradeMotion(ctx, searchJD, planet)
		if err != nil {
			continue
		}

		if motion == MotionDirect || motion == MotionStationary {
			// Found transition point, refine it
			return rd.FindPlanetaryStation(ctx, searchJD, planet, 10)
		}
	}

	return nil, fmt.Errorf("no station found in %d days backward search", maxDays)
}

// findStationForward searches forward for a station
func (rd *RetrogradeDetector) findStationForward(ctx context.Context, jd JulianDay, planet string, maxDays int) (*PlanetaryStation, error) {
	const stepSize = 1.0 // 1 day forward
	for i := 0; i < maxDays; i++ {
		searchJD := JulianDay(float64(jd) + float64(i)*stepSize)

		motion, err := rd.DetectRetrogradeMotion(ctx, searchJD, planet)
		if err != nil {
			continue
		}

		if motion == MotionDirect || motion == MotionStationary {
			// Found transition point, refine it
			return rd.FindPlanetaryStation(ctx, searchJD, planet, 10)
		}
	}

	return nil, fmt.Errorf("no station found in %d days forward search", maxDays)
}

// AnalyzeMotion provides comprehensive analysis of planetary motion
func (rd *RetrogradeDetector) AnalyzeMotion(ctx context.Context, jd JulianDay, planet string) (*MotionAnalysis, error) {
	ctx, span := rd.observer.CreateSpan(ctx, "retrograde.AnalyzeMotion")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("planet", planet),
	)

	// Get current motion
	motion, err := rd.DetectRetrogradeMotion(ctx, jd, planet)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Get current position and speed
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
		JulianDay: jd,
		Planet:    planet,
		Motion:    motion,
		Speed:     pos.Speed,
		Longitude: pos.Longitude,
	}

	// Check if near station (speed < 0.05 deg/day)
	analysis.IsNearStation = math.Abs(pos.Speed) < 0.05

	// Find next station
	nextStation, err := rd.FindPlanetaryStation(ctx, jd, planet, 400)
	if err == nil {
		analysis.NextStation = nextStation
	}

	// If retrograde, find current period
	if motion == MotionRetrograde {
		period, err := rd.FindRetrogradePeriod(ctx, jd, planet)
		if err == nil {
			analysis.CurrentPeriod = period
		}
	}

	// Find recent stations (last 6 months)
	recentStations := rd.findRecentStations(ctx, jd, planet, 180)
	analysis.RecentStations = recentStations

	span.SetAttributes(
		attribute.String("motion", string(motion)),
		attribute.Bool("is_near_station", analysis.IsNearStation),
		attribute.Bool("has_next_station", analysis.NextStation != nil),
	)

	return analysis, nil
}

// findRecentStations finds stations in the past N days
func (rd *RetrogradeDetector) findRecentStations(ctx context.Context, jd JulianDay, planet string, days int) []PlanetaryStation {
	stations := make([]PlanetaryStation, 0)

	// Search backward in chunks
	searchJD := jd
	const chunkSize = 30 // Search in 30-day chunks

	for i := 0; i < days/chunkSize; i++ {
		searchJD = JulianDay(float64(searchJD) - float64(chunkSize))

		station, err := rd.FindPlanetaryStation(ctx, searchJD, planet, chunkSize)
		if err == nil && station != nil {
			stations = append(stations, *station)
		}
	}

	return stations
}

// GetRetrogradePlanets returns all planets currently in retrograde motion
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

// ValidateKnownRetrograde validates detection against known retrograde periods
// This is useful for testing accuracy
func (rd *RetrogradeDetector) ValidateKnownRetrograde(ctx context.Context, planet string, knownStartJD, knownEndJD JulianDay) (bool, error) {
	ctx, span := rd.observer.CreateSpan(ctx, "retrograde.ValidateKnownRetrograde")
	defer span.End()

	// Check if planet is retrograde at midpoint
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
