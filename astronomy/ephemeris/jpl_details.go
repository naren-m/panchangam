package ephemeris

import (
	"context"
	"math"

	"go.opentelemetry.io/otel/attribute"
)

// calculateDetailedSunPosition calculates detailed Sun position.
func (j *JPLProvider) calculateDetailedSunPosition(ctx context.Context, jd JulianDay) *SolarPosition {
	_, span := j.observer.CreateSpan(ctx, "jpl.calculateDetailedSunPosition")
	defer span.End()

	// Days since J2000.0
	t := float64(jd - 2451545.0)

	// Mean longitude of the Sun (degrees)
	L := math.Mod(280.460+0.9856474*t, 360.0)

	// Mean anomaly of the Sun (degrees)
	M := math.Mod(357.528+0.9856003*t, 360.0)
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
		JulianDay:         jd,
		Longitude:         math.Mod(lambda+360, 360),
		RightAscension:    alpha,
		Declination:       delta,
		Distance:          distance,
		EquationOfTime:    eqTime,
		MeanAnomaly:       M,
		TrueAnomaly:       trueAnomaly,
		EccentricAnomaly:  eccentricAnomaly,
		MeanLongitude:     meanLongitude,
		ApparentLongitude: apparentLongitude,
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

// calculateDetailedMoonPosition calculates detailed Moon position.
func (j *JPLProvider) calculateDetailedMoonPosition(ctx context.Context, jd JulianDay) *LunarPosition {
	ctx, span := j.observer.CreateSpan(ctx, "jpl.calculateDetailedMoonPosition")
	defer span.End()

	// Days since J2000.0
	t := float64(jd - 2451545.0)

	// Mean longitude of the Moon (degrees)
	L := math.Mod(218.3164591+13.1763965268*t, 360.0)

	// Mean anomaly of the Moon (degrees)
	M := math.Mod(134.9634114+13.0649929509*t, 360.0)
	MRad := M * math.Pi / 180.0

	// Mean distance of the Moon from its ascending node (degrees)
	F := math.Mod(93.2720993+13.2299226639*t, 360.0)
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
	argumentOfLatitude := math.Mod(lambda-125.0, 360.0)

	// Mean longitude (degrees)
	meanLongitude := L

	// True longitude (degrees)
	trueLongitude := lambda

	position := &LunarPosition{
		JulianDay:          jd,
		Longitude:          math.Mod(lambda+360, 360),
		Latitude:           beta,
		RightAscension:     alpha,
		Declination:        delta,
		Distance:           distance,
		Phase:              phase,
		PhaseAngle:         phaseAngle,
		Illumination:       illumination,
		AngularDiameter:    angularDiameter,
		MeanAnomaly:        M,
		TrueAnomaly:        trueAnomaly,
		ArgumentOfLatitude: argumentOfLatitude,
		MeanLongitude:      meanLongitude,
		TrueLongitude:      trueLongitude,
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
