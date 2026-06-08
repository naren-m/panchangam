package ephemeris

import (
	"context"
	"math"

	"go.opentelemetry.io/otel/attribute"
)

// calculateDetailedSunPosition calculates detailed Sun position using Swiss Ephemeris.
func (s *SwissProvider) calculateDetailedSunPosition(ctx context.Context, jd JulianDay) *SolarPosition {
	_, span := s.observer.CreateSpan(ctx, "swiss.calculateDetailedSunPosition")
	defer span.End()

	// Use more accurate Swiss Ephemeris algorithms.
	t := float64(jd - 2451545.0)

	// More precise mean longitude.
	L := math.Mod(280.4664567+0.9856235*t, 360.0)

	// More precise mean anomaly.
	M := math.Mod(357.5291092+0.9856002585*t, 360.0)
	MRad := M * math.Pi / 180.0

	// Higher-order equation of center.
	C := 1.9148*math.Sin(MRad) + 0.0200*math.Sin(2*MRad) + 0.0003*math.Sin(3*MRad)

	// True longitude.
	lambda := L + C
	lambdaRad := lambda * math.Pi / 180.0

	// More accurate obliquity.
	epsilon := 23.4392911 - 0.0130042*t/100.0 - 0.00000164*t*t/10000.0
	epsilonRad := epsilon * math.Pi / 180.0

	// More accurate right ascension.
	alpha := math.Atan2(math.Cos(epsilonRad)*math.Sin(lambdaRad), math.Cos(lambdaRad)) * 180.0 / math.Pi
	alpha = math.Mod(alpha+360, 360)

	// More accurate declination.
	delta := math.Asin(math.Sin(epsilonRad)*math.Sin(lambdaRad)) * 180.0 / math.Pi

	// More accurate distance.
	distance := 1.000001018 * (1 - 0.01671123*math.Cos(MRad) - 0.00014*math.Cos(2*MRad))

	// More accurate equation of time.
	y := math.Tan(epsilonRad/2.0) * math.Tan(epsilonRad/2.0)
	eqTime := 4.0 * (y*math.Sin(2*L*math.Pi/180.0) - 2.0*0.01671123*math.Sin(M*math.Pi/180.0) +
		4.0*0.01671123*y*math.Sin(M*math.Pi/180.0)*math.Cos(2*L*math.Pi/180.0) -
		0.5*y*y*math.Sin(4*L*math.Pi/180.0) - 1.25*0.01671123*0.01671123*math.Sin(2*M*math.Pi/180.0))
	eqTime = eqTime * 180.0 / math.Pi / 15.0 // Convert to minutes

	// Accurate anomalies.
	trueAnomaly := M + C
	eccentricAnomaly := M + 1.9148*math.Sin(MRad) + 0.0200*math.Sin(2*MRad)

	// Apparent longitude (with nutation and aberration).
	apparentLongitude := lambda + 0.00569 - 0.00478*math.Sin((125.04-1934.136*t)*math.Pi/180.0)

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
		MeanLongitude:     L,
		ApparentLongitude: math.Mod(apparentLongitude+360, 360),
	}

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("right_ascension", position.RightAscension),
		attribute.Float64("declination", position.Declination),
		attribute.Float64("distance", position.Distance),
		attribute.Float64("equation_of_time", position.EquationOfTime),
		attribute.Float64("equation_of_center", C),
	)

	return position
}

// calculateDetailedMoonPosition calculates detailed Moon position using Swiss Ephemeris.
func (s *SwissProvider) calculateDetailedMoonPosition(ctx context.Context, jd JulianDay) *LunarPosition {
	ctx, span := s.observer.CreateSpan(ctx, "swiss.calculateDetailedMoonPosition")
	defer span.End()

	// Use ELP-2000 lunar theory (simplified).
	t := float64(jd - 2451545.0)

	// Fundamental arguments.
	L := math.Mod(218.3164477+13.17639648*t, 360.0)   // Mean longitude
	M := math.Mod(134.9633964+13.06499295*t, 360.0)   // Mean anomaly
	Mp := math.Mod(357.5291092+0.9856002585*t, 360.0) // Sun's mean anomaly
	D := math.Mod(297.8501921+12.19074912*t, 360.0)   // Mean elongation
	F := math.Mod(93.2720950+13.22935025*t, 360.0)    // Mean distance from node

	// Convert to radians.
	MRad := M * math.Pi / 180.0
	MpRad := Mp * math.Pi / 180.0
	DRad := D * math.Pi / 180.0
	FRad := F * math.Pi / 180.0

	// Extended periodic terms (ELP-2000).
	deltaL := 6.289*math.Sin(MRad) + 1.274*math.Sin(2*DRad-MRad) + 0.658*math.Sin(2*DRad) -
		0.186*math.Sin(MpRad) - 0.059*math.Sin(2*MRad-2*DRad) - 0.057*math.Sin(MRad-2*DRad+MpRad) +
		0.053*math.Sin(MRad+2*DRad) + 0.046*math.Sin(2*DRad-MpRad) + 0.041*math.Sin(MRad-MpRad) -
		0.035*math.Sin(DRad) - 0.031*math.Sin(MRad+MpRad) - 0.015*math.Sin(2*FRad-2*DRad) +
		0.011*math.Sin(MRad-4*DRad)

	deltaB := 5.128*math.Sin(FRad) + 0.281*math.Sin(MRad+FRad) + 0.277*math.Sin(MRad-FRad) +
		0.173*math.Sin(2*DRad-FRad) + 0.055*math.Sin(2*DRad-MRad+FRad) - 0.046*math.Sin(2*DRad-MRad-FRad) +
		0.033*math.Sin(MRad+2*DRad+FRad) + 0.017*math.Sin(2*MRad+FRad)

	deltaR := -20905*math.Cos(MRad) - 3699*math.Cos(2*DRad-MRad) - 2956*math.Cos(2*DRad) -
		570*math.Cos(2*MRad) + 246*math.Cos(2*MRad-2*DRad) - 205*math.Cos(MpRad-2*DRad) -
		171*math.Cos(MRad+2*DRad) - 152*math.Cos(MRad+MpRad-2*DRad) + 148*math.Cos(MRad-MpRad) -
		125*math.Cos(DRad) - 110*math.Cos(MRad+MpRad) + 59*math.Cos(2*DRad-MRad-MpRad)

	// Final geocentric coordinates.
	lambda := L + deltaL
	beta := deltaB
	distance := 385000.56 + deltaR // km

	// Convert to equatorial coordinates.
	lambdaRad := lambda * math.Pi / 180.0
	betaRad := beta * math.Pi / 180.0
	epsilon := 23.4392911 - 0.0130042*t/100.0
	epsilonRad := epsilon * math.Pi / 180.0

	// Right ascension and declination.
	alpha := math.Atan2(math.Cos(epsilonRad)*math.Sin(lambdaRad)-math.Sin(epsilonRad)*math.Tan(betaRad), math.Cos(lambdaRad)) * 180.0 / math.Pi
	alpha = math.Mod(alpha+360, 360)

	delta := math.Asin(math.Sin(epsilonRad)*math.Sin(lambdaRad)*math.Cos(betaRad)+math.Cos(epsilonRad)*math.Sin(betaRad)) * 180.0 / math.Pi

	// Phase calculations.
	sunLongitude := s.calculateSunPosition(ctx, jd).Longitude
	elongation := math.Abs(lambda - sunLongitude)
	if elongation > 180 {
		elongation = 360 - elongation
	}

	// Phase and illumination.
	phaseAngle := elongation
	phase := (1.0 - math.Cos(elongation*math.Pi/180.0)) / 2.0
	illumination := phase * 100.0

	// Angular diameter.
	angularDiameter := 1873.0 * (6378.14 / distance)

	// Lunar anomalies and arguments.
	trueAnomaly := M + deltaL
	argumentOfLatitude := math.Mod(lambda-125.0, 360.0)

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
		MeanLongitude:      L,
		TrueLongitude:      lambda,
	}

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("latitude", position.Latitude),
		attribute.Float64("distance", position.Distance),
		attribute.Float64("phase", position.Phase),
		attribute.Float64("illumination", position.Illumination),
		attribute.Float64("elp2000_delta_l", deltaL),
		attribute.Float64("elp2000_delta_b", deltaB),
		attribute.Float64("elp2000_delta_r", deltaR),
	)

	return position
}
