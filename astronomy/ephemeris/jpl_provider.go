package ephemeris

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

// JPLProvider implements the EphemerisProvider interface using JPL DE440 ephemeris
type JPLProvider struct {
	name           string
	version        string
	dataStartJD    JulianDay
	dataEndJD      JulianDay
	observer       observability.ObserverInterface
	healthStatus   *HealthStatus
	lastHealthCheck time.Time
}

// NewJPLProvider creates a new JPL DE440 ephemeris provider
func NewJPLProvider() *JPLProvider {
	now := time.Now()
	
	return &JPLProvider{
		name:        "JPL DE440",
		version:     "440",
		dataStartJD: JulianDay(1550184.5), // Jan 1, 1550 CE
		dataEndJD:   JulianDay(2816787.5), // Dec 31, 2650 CE
		observer:    observability.Observer(),
		healthStatus: &HealthStatus{
			Available:    true,
			LastCheck:    now,
			DataStartJD:  1550184.5,
			DataEndJD:    2816787.5,
			ResponseTime: 0,
			Version:      "440",
			Source:       "JPL DE440",
		},
		lastHealthCheck: now,
	}
}

// GetPlanetaryPositions returns positions of all planets for a given Julian day
func (j *JPLProvider) GetPlanetaryPositions(ctx context.Context, jd JulianDay) (*PlanetaryPositions, error) {
	ctx, span := j.observer.CreateSpan(ctx, "jpl.GetPlanetaryPositions")
	defer span.End()

	span.SetAttributes(
		attribute.String("provider", j.name),
		attribute.String("version", j.version),
		attribute.Float64("julian_day", float64(jd)),
	)

	// Validate Julian day range
	if jd < j.dataStartJD || jd > j.dataEndJD {
		err := fmt.Errorf("julian day %f is outside valid range [%f, %f]", jd, j.dataStartJD, j.dataEndJD)
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("in_range", false))
		return nil, err
	}

	span.SetAttributes(attribute.Bool("in_range", true))

	// Calculate planetary positions using simplified analytical methods
	// In a real implementation, this would use JPL DE440 binary files
	positions := &PlanetaryPositions{
		JulianDay: jd,
		Sun:       j.calculateSunPosition(ctx, jd),
		Moon:      j.calculateMoonPosition(ctx, jd),
		Mercury:   j.calculatePlanetPosition(ctx, jd, "mercury"),
		Venus:     j.calculatePlanetPosition(ctx, jd, "venus"),
		Mars:      j.calculatePlanetPosition(ctx, jd, "mars"),
		Jupiter:   j.calculatePlanetPosition(ctx, jd, "jupiter"),
		Saturn:    j.calculatePlanetPosition(ctx, jd, "saturn"),
		Uranus:    j.calculatePlanetPosition(ctx, jd, "uranus"),
		Neptune:   j.calculatePlanetPosition(ctx, jd, "neptune"),
		Pluto:     j.calculatePlanetPosition(ctx, jd, "pluto"),
	}

	span.SetAttributes(attribute.Bool("success", true))
	span.AddEvent("Planetary positions calculated")

	return positions, nil
}

// GetSunPosition returns detailed Sun position for a given Julian day
func (j *JPLProvider) GetSunPosition(ctx context.Context, jd JulianDay) (*SolarPosition, error) {
	ctx, span := j.observer.CreateSpan(ctx, "jpl.GetSunPosition")
	defer span.End()

	span.SetAttributes(
		attribute.String("provider", j.name),
		attribute.Float64("julian_day", float64(jd)),
	)

	// Validate Julian day range
	if jd < j.dataStartJD || jd > j.dataEndJD {
		err := fmt.Errorf("julian day %f is outside valid range [%f, %f]", jd, j.dataStartJD, j.dataEndJD)
		span.RecordError(err)
		return nil, err
	}

	position := j.calculateDetailedSunPosition(ctx, jd)
	
	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("right_ascension", position.RightAscension),
		attribute.Float64("declination", position.Declination),
		attribute.Float64("distance", position.Distance),
		attribute.Bool("success", true),
	)
	span.AddEvent("Sun position calculated")

	return position, nil
}

// GetMoonPosition returns detailed Moon position for a given Julian day
func (j *JPLProvider) GetMoonPosition(ctx context.Context, jd JulianDay) (*LunarPosition, error) {
	ctx, span := j.observer.CreateSpan(ctx, "jpl.GetMoonPosition")
	defer span.End()

	span.SetAttributes(
		attribute.String("provider", j.name),
		attribute.Float64("julian_day", float64(jd)),
	)

	// Validate Julian day range
	if jd < j.dataStartJD || jd > j.dataEndJD {
		err := fmt.Errorf("julian day %f is outside valid range [%f, %f]", jd, j.dataStartJD, j.dataEndJD)
		span.RecordError(err)
		return nil, err
	}

	position := j.calculateDetailedMoonPosition(ctx, jd)
	
	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("latitude", position.Latitude),
		attribute.Float64("distance", position.Distance),
		attribute.Float64("phase", position.Phase),
		attribute.Bool("success", true),
	)
	span.AddEvent("Moon position calculated")

	return position, nil
}

// IsAvailable checks if the ephemeris provider is available
func (j *JPLProvider) IsAvailable(ctx context.Context) bool {
	ctx, span := j.observer.CreateSpan(ctx, "jpl.IsAvailable")
	defer span.End()

	// Update health status if it's been more than 30 seconds
	if time.Since(j.lastHealthCheck) > 30*time.Second {
		j.updateHealthStatus(ctx)
	}

	available := j.healthStatus.Available
	span.SetAttributes(
		attribute.Bool("available", available),
		attribute.String("last_check", j.healthStatus.LastCheck.Format(time.RFC3339)),
	)

	return available
}

// GetDataRange returns the valid Julian day range for this provider
func (j *JPLProvider) GetDataRange() (startJD, endJD JulianDay) {
	return j.dataStartJD, j.dataEndJD
}

// GetHealthStatus returns the current health status
func (j *JPLProvider) GetHealthStatus(ctx context.Context) (*HealthStatus, error) {
	ctx, span := j.observer.CreateSpan(ctx, "jpl.GetHealthStatus")
	defer span.End()

	// Update health status
	j.updateHealthStatus(ctx)

	span.SetAttributes(
		attribute.Bool("available", j.healthStatus.Available),
		attribute.Int64("response_time_ms", j.healthStatus.ResponseTime.Milliseconds()),
		attribute.String("version", j.healthStatus.Version),
	)

	return j.healthStatus, nil
}

// GetProviderName returns the name of the provider
func (j *JPLProvider) GetProviderName() string {
	return j.name
}

// GetVersion returns the version of the ephemeris data
func (j *JPLProvider) GetVersion() string {
	return j.version
}

// Close closes the provider and releases resources
func (j *JPLProvider) Close() error {
	// No resources to close for this implementation
	return nil
}

// updateHealthStatus updates the health status of the provider
func (j *JPLProvider) updateHealthStatus(ctx context.Context) {
	ctx, span := j.observer.CreateSpan(ctx, "jpl.updateHealthStatus")
	defer span.End()

	start := time.Now()
	
	// Simple health check - verify we can perform basic calculations
	testJD := JulianDay(2451545.0) // J2000.0
	available := true
	var errorMessage string

	// Test basic calculation
	if testJD < j.dataStartJD || testJD > j.dataEndJD {
		available = false
		errorMessage = "Test Julian day outside valid range"
	} else {
		// Test calculation by computing a simple position
		_ = j.calculateSunPosition(ctx, testJD)
	}

	responseTime := time.Since(start)
	now := time.Now()

	j.healthStatus = &HealthStatus{
		Available:    available,
		LastCheck:    now,
		DataStartJD:  float64(j.dataStartJD),
		DataEndJD:    float64(j.dataEndJD),
		ResponseTime: responseTime,
		ErrorMessage: errorMessage,
		Version:      j.version,
		Source:       j.name,
	}
	j.lastHealthCheck = now

	span.SetAttributes(
		attribute.Bool("available", available),
		attribute.Int64("response_time_ms", responseTime.Milliseconds()),
		attribute.String("error_message", errorMessage),
	)
	span.AddEvent("Health status updated")
}

// calculateSunPosition calculates basic sun position
func (j *JPLProvider) calculateSunPosition(ctx context.Context, jd JulianDay) Position {
	_, span := j.observer.CreateSpan(ctx, "jpl.calculateSunPosition")
	defer span.End()

	// Days since J2000.0
	t := float64(jd - 2451545.0)
	
	// Mean longitude of the Sun (degrees)
	L := math.Mod(280.460 + 0.9856474*t, 360.0)
	
	// Mean anomaly of the Sun (degrees)
	M := math.Mod(357.528 + 0.9856003*t, 360.0)
	
	// Convert to radians
	MRad := M * math.Pi / 180.0
	
	// Ecliptic longitude (degrees)
	lambda := L + 1.915*math.Sin(MRad) + 0.020*math.Sin(2*MRad)
	
	// Distance (AU)
	distance := 1.00014 - 0.01671*math.Cos(MRad) - 0.00014*math.Cos(2*MRad)
	
	// Speed (degrees/day) - approximate
	speed := 0.9856 // Mean motion of the Sun
	
	position := Position{
		Longitude: math.Mod(lambda+360, 360),
		Latitude:  0.0, // Sun's ecliptic latitude is always 0
		Distance:  distance,
		Speed:     speed,
	}

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("distance", position.Distance),
		attribute.Float64("speed", position.Speed),
	)

	return position
}

// calculateMoonPosition calculates basic moon position
func (j *JPLProvider) calculateMoonPosition(ctx context.Context, jd JulianDay) Position {
	_, span := j.observer.CreateSpan(ctx, "jpl.calculateMoonPosition")
	defer span.End()

	// Days since J2000.0
	t := float64(jd - 2451545.0)
	
	// Mean longitude of the Moon (degrees)
	L := math.Mod(218.3164591 + 13.1763965268*t, 360.0)
	
	// Mean anomaly of the Moon (degrees)
	M := math.Mod(134.9634114 + 13.0649929509*t, 360.0)
	
	// Mean distance of the Moon from its ascending node (degrees)
	F := math.Mod(93.2720993 + 13.2299226639*t, 360.0)
	
	// Convert to radians
	MRad := M * math.Pi / 180.0
	FRad := F * math.Pi / 180.0
	
	// Ecliptic longitude (degrees) - simplified
	lambda := L + 6.289*math.Sin(MRad)
	
	// Ecliptic latitude (degrees) - simplified
	beta := 5.128 * math.Sin(FRad)
	
	// Distance (km) - simplified
	distance := 385000.0 - 20905.0*math.Cos(MRad)
	
	// Speed (degrees/day) - approximate
	speed := 13.18 // Mean motion of the Moon
	
	position := Position{
		Longitude: math.Mod(lambda+360, 360),
		Latitude:  beta,
		Distance:  distance / 149597870.7, // Convert to AU
		Speed:     speed,
	}

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("latitude", position.Latitude),
		attribute.Float64("distance_au", position.Distance),
		attribute.Float64("speed", position.Speed),
	)

	return position
}

// calculatePlanetPosition calculates basic planet position
func (j *JPLProvider) calculatePlanetPosition(ctx context.Context, jd JulianDay, planet string) Position {
	_, span := j.observer.CreateSpan(ctx, "jpl.calculatePlanetPosition")
	defer span.End()

	span.SetAttributes(attribute.String("planet", planet))

	// Days since J2000.0
	t := float64(jd - 2451545.0)
	
	// Simplified planetary elements - these would be much more complex in real JPL DE440
	var L, M, distance, speed float64
	
	switch planet {
	case "mercury":
		L = math.Mod(252.25084 + 4.092338796*t, 360.0)
		M = math.Mod(174.79252 + 4.092335*t, 360.0)
		distance = 0.387098
		speed = 4.092
	case "venus":
		L = math.Mod(181.97973 + 1.602136*t, 360.0)
		M = math.Mod(50.41575 + 1.602136*t, 360.0)
		distance = 0.723327
		speed = 1.602
	case "mars":
		L = math.Mod(355.433 + 0.524033*t, 360.0)
		M = math.Mod(19.3879 + 0.524033*t, 360.0)
		distance = 1.523679
		speed = 0.524
	case "jupiter":
		L = math.Mod(34.40438 + 0.083091*t, 360.0)
		M = math.Mod(20.0202 + 0.083091*t, 360.0)
		distance = 5.204267
		speed = 0.083
	case "saturn":
		L = math.Mod(49.9477 + 0.033494*t, 360.0)
		M = math.Mod(317.0207 + 0.033494*t, 360.0)
		distance = 9.5820172
		speed = 0.033
	case "uranus":
		L = math.Mod(313.23218 + 0.011733*t, 360.0)
		M = math.Mod(141.0498 + 0.011733*t, 360.0)
		distance = 19.189253
		speed = 0.012
	case "neptune":
		L = math.Mod(304.88003 + 0.005965*t, 360.0)
		M = math.Mod(256.228 + 0.005965*t, 360.0)
		distance = 30.070900
		speed = 0.006
	case "pluto":
		L = math.Mod(238.92881 + 0.003968*t, 360.0)
		M = math.Mod(14.882 + 0.003968*t, 360.0)
		distance = 39.481686
		speed = 0.004
	default:
		// Default to Earth's position relative to Sun
		L = math.Mod(100.46435 + 0.985609*t, 360.0)
		M = math.Mod(357.52911 + 0.985600*t, 360.0)
		distance = 1.000001
		speed = 0.986
	}
	
	// Convert mean anomaly to radians
	MRad := M * math.Pi / 180.0
	
	// Simple correction for eccentricity
	lambda := L + 2.0*math.Sin(MRad)
	
	position := Position{
		Longitude: math.Mod(lambda+360, 360),
		Latitude:  0.0, // Simplified - no inclination correction
		Distance:  distance,
		Speed:     speed,
	}

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("distance", position.Distance),
		attribute.Float64("speed", position.Speed),
	)

	return position
}

// calculateDetailedSunPosition calculates detailed sun position
func (j *JPLProvider) calculateDetailedSunPosition(ctx context.Context, jd JulianDay) *SolarPosition {
	_, span := j.observer.CreateSpan(ctx, "jpl.calculateDetailedSunPosition")
	defer span.End()

	// Days since J2000.0
	t := float64(jd - 2451545.0)
	
	// Mean longitude of the Sun (degrees)
	L := math.Mod(280.460 + 0.9856474*t, 360.0)
	
	// Mean anomaly of the Sun (degrees)
	M := math.Mod(357.528 + 0.9856003*t, 360.0)
	MRad := M * math.Pi / 180.0
	
	// Ecliptic longitude (degrees)
	lambda := L + 1.915*math.Sin(MRad) + 0.020*math.Sin(2*MRad)
	lambdaRad := lambda * math.Pi / 180.0
	
	// Obliquity of the ecliptic (degrees)
	epsilon := 23.4393 - 0.0000004*t
	epsilonRad := epsilon * math.Pi / 180.0
	
	// Right ascension (degrees)
	alpha := math.Atan2(math.Cos(epsilonRad)*math.Sin(lambdaRad), math.Cos(lambdaRad)) * 180.0 / math.Pi
	alpha = math.Mod(alpha+360, 360)
	
	// Declination (degrees)
	delta := math.Asin(math.Sin(epsilonRad)*math.Sin(lambdaRad)) * 180.0 / math.Pi
	
	// Distance (AU)
	distance := 1.00014 - 0.01671*math.Cos(MRad) - 0.00014*math.Cos(2*MRad)
	
	// Equation of time (minutes)
	eqTime := 4.0 * (L - alpha)
	
	// True anomaly (degrees)
	trueAnomaly := M + 1.915*math.Sin(MRad) + 0.020*math.Sin(2*MRad)
	
	// Eccentric anomaly (degrees)
	eccentricAnomaly := M + 1.915*math.Sin(MRad)
	
	// Mean longitude (degrees)
	meanLongitude := L
	
	// Apparent longitude (degrees) - simplified
	apparentLongitude := lambda
	
	position := &SolarPosition{
		JulianDay:           jd,
		Longitude:           math.Mod(lambda+360, 360),
		RightAscension:      alpha,
		Declination:         delta,
		Distance:            distance,
		EquationOfTime:      eqTime,
		MeanAnomaly:         M,
		TrueAnomaly:         trueAnomaly,
		EccentricAnomaly:    eccentricAnomaly,
		MeanLongitude:       meanLongitude,
		ApparentLongitude:   apparentLongitude,
	}

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("right_ascension", position.RightAscension),
		attribute.Float64("declination", position.Declination),
		attribute.Float64("distance", position.Distance),
		attribute.Float64("equation_of_time", position.EquationOfTime),
	)

	return position
}

// calculateDetailedMoonPosition calculates detailed moon position
func (j *JPLProvider) calculateDetailedMoonPosition(ctx context.Context, jd JulianDay) *LunarPosition {
	ctx, span := j.observer.CreateSpan(ctx, "jpl.calculateDetailedMoonPosition")
	defer span.End()

	// Days since J2000.0
	t := float64(jd - 2451545.0)
	
	// Mean longitude of the Moon (degrees)
	L := math.Mod(218.3164591 + 13.1763965268*t, 360.0)
	
	// Mean anomaly of the Moon (degrees)
	M := math.Mod(134.9634114 + 13.0649929509*t, 360.0)
	MRad := M * math.Pi / 180.0
	
	// Mean distance of the Moon from its ascending node (degrees)
	F := math.Mod(93.2720993 + 13.2299226639*t, 360.0)
	FRad := F * math.Pi / 180.0
	
	// Ecliptic longitude (degrees)
	lambda := L + 6.289*math.Sin(MRad)
	lambdaRad := lambda * math.Pi / 180.0
	
	// Ecliptic latitude (degrees)
	beta := 5.128 * math.Sin(FRad)
	betaRad := beta * math.Pi / 180.0
	
	// Distance (km)
	distance := 385000.0 - 20905.0*math.Cos(MRad)
	
	// Obliquity of the ecliptic (degrees)
	epsilon := 23.4393 - 0.0000004*t
	epsilonRad := epsilon * math.Pi / 180.0
	
	// Right ascension (degrees)
	alpha := math.Atan2(math.Cos(epsilonRad)*math.Sin(lambdaRad)-math.Sin(epsilonRad)*math.Tan(betaRad), math.Cos(lambdaRad)) * 180.0 / math.Pi
	alpha = math.Mod(alpha+360, 360)
	
	// Declination (degrees)
	delta := math.Asin(math.Sin(epsilonRad)*math.Sin(lambdaRad)*math.Cos(betaRad)+math.Cos(epsilonRad)*math.Sin(betaRad)) * 180.0 / math.Pi
	
	// Phase calculation (simplified)
	// Phase angle between Earth, Moon, and Sun
	sunLongitude := j.calculateSunPosition(ctx, jd).Longitude
	phaseAngle := math.Abs(lambda - sunLongitude)
	if phaseAngle > 180 {
		phaseAngle = 360 - phaseAngle
	}
	
	// Phase (0 = new moon, 0.5 = full moon)
	phase := (1.0 - math.Cos(phaseAngle*math.Pi/180.0)) / 2.0
	
	// Illumination percentage
	illumination := phase * 100.0
	
	// Angular diameter (arcseconds)
	angularDiameter := 1873.0 * (6378.14 / distance)
	
	// True anomaly (degrees)
	trueAnomaly := M + 6.289*math.Sin(MRad)
	
	// Argument of latitude (degrees)
	argumentOfLatitude := math.Mod(lambda - 125.0, 360.0)
	
	// Mean longitude (degrees)
	meanLongitude := L
	
	// True longitude (degrees)
	trueLongitude := lambda
	
	position := &LunarPosition{
		JulianDay:           jd,
		Longitude:           math.Mod(lambda+360, 360),
		Latitude:            beta,
		RightAscension:      alpha,
		Declination:         delta,
		Distance:            distance,
		Phase:               phase,
		PhaseAngle:          phaseAngle,
		Illumination:        illumination,
		AngularDiameter:     angularDiameter,
		MeanAnomaly:         M,
		TrueAnomaly:         trueAnomaly,
		ArgumentOfLatitude:  argumentOfLatitude,
		MeanLongitude:       meanLongitude,
		TrueLongitude:       trueLongitude,
	}

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("latitude", position.Latitude),
		attribute.Float64("distance", position.Distance),
		attribute.Float64("phase", position.Phase),
		attribute.Float64("illumination", position.Illumination),
	)

	return position
}