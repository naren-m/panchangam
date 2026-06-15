package ephemeris

import (
	"context"
	"math"

	"go.opentelemetry.io/otel/attribute"
)

// calculateSunPosition calculates a basic Sun position.
func (j *JPLProvider) calculateSunPosition(ctx context.Context, jd JulianDay) Position {
	_, span := j.observer.CreateSpan(ctx, "jpl.calculateSunPosition")
	defer span.End()

	// Days since J2000.0
	t := float64(jd - 2451545.0)

	// Mean longitude of the Sun (degrees)
	L := math.Mod(280.460+0.9856474*t, 360.0)

	// Mean anomaly of the Sun (degrees)
	M := math.Mod(357.528+0.9856003*t, 360.0)

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

// calculateMoonPosition calculates a basic Moon position.
func (j *JPLProvider) calculateMoonPosition(ctx context.Context, jd JulianDay) Position {
	_, span := j.observer.CreateSpan(ctx, "jpl.calculateMoonPosition")
	defer span.End()

	// Days since J2000.0
	t := float64(jd - 2451545.0)

	// Mean longitude of the Moon (degrees)
	L := math.Mod(218.3164591+13.1763965268*t, 360.0)

	// Mean anomaly of the Moon (degrees)
	M := math.Mod(134.9634114+13.0649929509*t, 360.0)

	// Mean distance of the Moon from its ascending node (degrees)
	F := math.Mod(93.2720993+13.2299226639*t, 360.0)

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

// calculatePlanetPosition calculates a basic planet position.
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
		L = math.Mod(252.25084+4.092338796*t, 360.0)
		M = math.Mod(174.79252+4.092335*t, 360.0)
		distance = 0.387098
		speed = 4.092
	case "venus":
		L = math.Mod(181.97973+1.602136*t, 360.0)
		M = math.Mod(50.41575+1.602136*t, 360.0)
		distance = 0.723327
		speed = 1.602
	case "mars":
		L = math.Mod(355.433+0.524033*t, 360.0)
		M = math.Mod(19.3879+0.524033*t, 360.0)
		distance = 1.523679
		speed = 0.524
	case "jupiter":
		L = math.Mod(34.40438+0.083091*t, 360.0)
		M = math.Mod(20.0202+0.083091*t, 360.0)
		distance = 5.204267
		speed = 0.083
	case "saturn":
		L = math.Mod(49.9477+0.033494*t, 360.0)
		M = math.Mod(317.0207+0.033494*t, 360.0)
		distance = 9.5820172
		speed = 0.033
	case "uranus":
		L = math.Mod(313.23218+0.011733*t, 360.0)
		M = math.Mod(141.0498+0.011733*t, 360.0)
		distance = 19.189253
		speed = 0.012
	case "neptune":
		L = math.Mod(304.88003+0.005965*t, 360.0)
		M = math.Mod(256.228+0.005965*t, 360.0)
		distance = 30.070900
		speed = 0.006
	case "pluto":
		L = math.Mod(238.92881+0.003968*t, 360.0)
		M = math.Mod(14.882+0.003968*t, 360.0)
		distance = 39.481686
		speed = 0.004
	default:
		// Default to Earth's position relative to Sun
		L = math.Mod(100.46435+0.985609*t, 360.0)
		M = math.Mod(357.52911+0.985600*t, 360.0)
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
