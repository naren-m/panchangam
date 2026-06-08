package ephemeris

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

// SwissProvider implements the EphemerisProvider interface using Swiss Ephemeris
type SwissProvider struct {
	name            string
	version         string
	dataStartJD     JulianDay
	dataEndJD       JulianDay
	observer        observability.ObserverInterface
	healthStatus    *HealthStatus
	lastHealthCheck time.Time
}

// NewSwissProvider creates a new Swiss Ephemeris provider
func NewSwissProvider() *SwissProvider {
	now := time.Now()

	return &SwissProvider{
		name:        "Swiss Ephemeris",
		version:     "2.10",
		dataStartJD: JulianDay(-3027215.5), // 13201 BCE
		dataEndJD:   JulianDay(7857061.5),  // 17191 CE
		observer:    observability.Observer(),
		healthStatus: &HealthStatus{
			Available:    true,
			LastCheck:    now,
			DataStartJD:  -3027215.5,
			DataEndJD:    7857061.5,
			ResponseTime: 0,
			Version:      "2.10",
			Source:       "Swiss Ephemeris",
		},
		lastHealthCheck: now,
	}
}

// GetPlanetaryPositions returns positions of all planets for a given Julian day
func (s *SwissProvider) GetPlanetaryPositions(ctx context.Context, jd JulianDay) (*PlanetaryPositions, error) {
	ctx, span := s.observer.CreateSpan(ctx, "swiss.GetPlanetaryPositions")
	defer span.End()

	span.SetAttributes(
		attribute.String("provider", s.name),
		attribute.String("version", s.version),
		attribute.Float64("julian_day", float64(jd)),
	)

	// Validate Julian day range
	if jd < s.dataStartJD || jd > s.dataEndJD {
		err := fmt.Errorf("julian day %f is outside valid range [%f, %f]", jd, s.dataStartJD, s.dataEndJD)
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("in_range", false))
		return nil, err
	}

	span.SetAttributes(attribute.Bool("in_range", true))

	// Calculate planetary positions using Swiss Ephemeris algorithms
	// This is a simplified implementation - real Swiss Ephemeris uses complex algorithms
	positions := &PlanetaryPositions{
		JulianDay: jd,
		Sun:       s.calculateSunPosition(ctx, jd),
		Moon:      s.calculateMoonPosition(ctx, jd),
		Mercury:   s.calculatePlanetPosition(ctx, jd, "mercury"),
		Venus:     s.calculatePlanetPosition(ctx, jd, "venus"),
		Mars:      s.calculatePlanetPosition(ctx, jd, "mars"),
		Jupiter:   s.calculatePlanetPosition(ctx, jd, "jupiter"),
		Saturn:    s.calculatePlanetPosition(ctx, jd, "saturn"),
		Uranus:    s.calculatePlanetPosition(ctx, jd, "uranus"),
		Neptune:   s.calculatePlanetPosition(ctx, jd, "neptune"),
		Pluto:     s.calculatePlanetPosition(ctx, jd, "pluto"),
	}

	span.SetAttributes(attribute.Bool("success", true))
	span.AddEvent("Planetary positions calculated using Swiss Ephemeris")

	return positions, nil
}

// GetSunPosition returns detailed Sun position for a given Julian day
func (s *SwissProvider) GetSunPosition(ctx context.Context, jd JulianDay) (*SolarPosition, error) {
	ctx, span := s.observer.CreateSpan(ctx, "swiss.GetSunPosition")
	defer span.End()

	span.SetAttributes(
		attribute.String("provider", s.name),
		attribute.Float64("julian_day", float64(jd)),
	)

	// Validate Julian day range
	if jd < s.dataStartJD || jd > s.dataEndJD {
		err := fmt.Errorf("julian day %f is outside valid range [%f, %f]", jd, s.dataStartJD, s.dataEndJD)
		span.RecordError(err)
		return nil, err
	}

	position := s.calculateDetailedSunPosition(ctx, jd)

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("right_ascension", position.RightAscension),
		attribute.Float64("declination", position.Declination),
		attribute.Float64("distance", position.Distance),
		attribute.Bool("success", true),
	)
	span.AddEvent("Sun position calculated using Swiss Ephemeris")

	return position, nil
}

// GetMoonPosition returns detailed Moon position for a given Julian day
func (s *SwissProvider) GetMoonPosition(ctx context.Context, jd JulianDay) (*LunarPosition, error) {
	ctx, span := s.observer.CreateSpan(ctx, "swiss.GetMoonPosition")
	defer span.End()

	span.SetAttributes(
		attribute.String("provider", s.name),
		attribute.Float64("julian_day", float64(jd)),
	)

	// Validate Julian day range
	if jd < s.dataStartJD || jd > s.dataEndJD {
		err := fmt.Errorf("julian day %f is outside valid range [%f, %f]", jd, s.dataStartJD, s.dataEndJD)
		span.RecordError(err)
		return nil, err
	}

	position := s.calculateDetailedMoonPosition(ctx, jd)

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("latitude", position.Latitude),
		attribute.Float64("distance", position.Distance),
		attribute.Float64("phase", position.Phase),
		attribute.Bool("success", true),
	)
	span.AddEvent("Moon position calculated using Swiss Ephemeris")

	return position, nil
}

// IsAvailable checks if the ephemeris provider is available
func (s *SwissProvider) IsAvailable(ctx context.Context) bool {
	ctx, span := s.observer.CreateSpan(ctx, "swiss.IsAvailable")
	defer span.End()

	// Update health status if it's been more than 30 seconds
	if time.Since(s.lastHealthCheck) > 30*time.Second {
		s.updateHealthStatus(ctx)
	}

	available := s.healthStatus.Available
	span.SetAttributes(
		attribute.Bool("available", available),
		attribute.String("last_check", s.healthStatus.LastCheck.Format(time.RFC3339)),
	)

	return available
}

// GetDataRange returns the valid Julian day range for this provider
func (s *SwissProvider) GetDataRange() (startJD, endJD JulianDay) {
	return s.dataStartJD, s.dataEndJD
}

// GetHealthStatus returns the current health status
func (s *SwissProvider) GetHealthStatus(ctx context.Context) (*HealthStatus, error) {
	ctx, span := s.observer.CreateSpan(ctx, "swiss.GetHealthStatus")
	defer span.End()

	// Update health status
	s.updateHealthStatus(ctx)

	span.SetAttributes(
		attribute.Bool("available", s.healthStatus.Available),
		attribute.Int64("response_time_ms", s.healthStatus.ResponseTime.Milliseconds()),
		attribute.String("version", s.healthStatus.Version),
	)

	return s.healthStatus, nil
}

// GetProviderName returns the name of the provider
func (s *SwissProvider) GetProviderName() string {
	return s.name
}

// GetVersion returns the version of the ephemeris data
func (s *SwissProvider) GetVersion() string {
	return s.version
}

// Close closes the provider and releases resources
func (s *SwissProvider) Close() error {
	// No resources to close for this implementation
	return nil
}
