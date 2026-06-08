package astronomy

import (
	"context"
	"math"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// calculateLunarPositionJD calculates the moon's position for a given Julian day
func calculateLunarPositionJD(ctx context.Context, jd float64) *LunarPosition {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "calculateLunarPositionJD")
	defer span.End()

	span.SetAttributes(attribute.Float64("julian_day", jd))

	// Days since J2000.0
	T := (jd - 2451545.0) / 36525.0

	// Moon's mean longitude (degrees)
	L := math.Mod(218.3164477+481267.88123421*T-0.0015786*T*T+T*T*T/538841.0-T*T*T*T/65194000.0, 360.0)

	// Mean elongation of the Moon from the Sun (degrees)
	D := math.Mod(297.8501921+445267.1114034*T-0.0018819*T*T+T*T*T/545868.0-T*T*T*T/113065000.0, 360.0)

	// Sun's mean anomaly (degrees)
	M := math.Mod(357.5291092+35999.0502909*T-0.0001536*T*T+T*T*T/24490000.0, 360.0)

	// Moon's mean anomaly (degrees)
	MPrime := math.Mod(134.9633964+477198.8675055*T+0.0087414*T*T+T*T*T/69699.0-T*T*T*T/14712000.0, 360.0)

	// Moon's argument of latitude (degrees)
	F := math.Mod(93.2720950+483202.0175233*T-0.0036539*T*T-T*T*T/3526000.0+T*T*T*T/863310000.0, 360.0)

	// Convert to radians for calculations
	DRad := D * DegToRad
	MRad := M * DegToRad
	MPrimeRad := MPrime * DegToRad
	FRad := F * DegToRad

	// Primary lunar perturbations (simplified - full theory has hundreds of terms)
	// These are the most significant terms
	lonCorrection := 6.288774*math.Sin(MPrimeRad) +
		1.274027*math.Sin(2*DRad-MPrimeRad) +
		0.658314*math.Sin(2*DRad) +
		0.213618*math.Sin(2*MPrimeRad) -
		0.185116*math.Sin(MRad) -
		0.114332*math.Sin(2*FRad) +
		0.058793*math.Sin(2*(DRad-MPrimeRad)) +
		0.057066*math.Sin(2*DRad-MRad-MPrimeRad) +
		0.053322*math.Sin(2*DRad+MPrimeRad) +
		0.045758*math.Sin(2*DRad-MRad)

	latCorrection := 5.128122*math.Sin(FRad) +
		0.280602*math.Sin(MPrimeRad+FRad) +
		0.277693*math.Sin(MPrimeRad-FRad) +
		0.173237*math.Sin(2*DRad-FRad) +
		0.055413*math.Sin(2*DRad-MPrimeRad+FRad) +
		0.046271*math.Sin(2*DRad-MPrimeRad-FRad) +
		0.032573*math.Sin(2*DRad+FRad)

	distCorrection := -20905.355*math.Cos(MPrimeRad) -
		3699.111*math.Cos(2*DRad-MPrimeRad) -
		2955.968*math.Cos(2*DRad) -
		569.925*math.Cos(2*MPrimeRad) +
		246.158*math.Cos(MRad) -
		204.586*math.Cos(2*FRad) -
		170.733*math.Cos(2*(DRad-MPrimeRad)) -
		152.138*math.Cos(2*DRad-MRad-MPrimeRad)

	// Final lunar coordinates
	lambda := L + lonCorrection         // Ecliptic longitude
	beta := latCorrection               // Ecliptic latitude
	delta := 385000.56 + distCorrection // Distance in km

	// Convert to equatorial coordinates
	epsilon := 23.4392911 * DegToRad // Obliquity of ecliptic (simplified)
	lambdaRad := lambda * DegToRad
	betaRad := beta * DegToRad

	// Right ascension and declination
	alpha := math.Atan2(math.Sin(lambdaRad)*math.Cos(epsilon)-math.Tan(betaRad)*math.Sin(epsilon), math.Cos(lambdaRad)) * RadToDeg
	if alpha < 0 {
		alpha += 360
	}

	delta_eq := math.Asin(math.Sin(betaRad)*math.Cos(epsilon)+math.Cos(betaRad)*math.Sin(epsilon)*math.Sin(lambdaRad)) * RadToDeg

	// Calculate lunar phase
	// Phase angle between Sun and Moon as seen from Earth
	elongation := math.Mod(D, 360.0) * DegToRad
	phaseAngle := math.Pi - elongation
	if phaseAngle < 0 {
		phaseAngle += 2 * math.Pi
	}

	phase := (1 - math.Cos(phaseAngle)) / 2 // 0 = new moon, 1 = full moon
	illumination := phase * 100             // percentage

	result := &LunarPosition{
		RightAscension: alpha,
		Declination:    delta_eq,
		Distance:       delta,
		Phase:          phase,
		Illumination:   illumination,
	}

	span.SetAttributes(
		attribute.Float64("centuries_since_j2000", T),
		attribute.Float64("mean_longitude", L),
		attribute.Float64("mean_elongation", D),
		attribute.Float64("sun_mean_anomaly", M),
		attribute.Float64("moon_mean_anomaly", MPrime),
		attribute.Float64("argument_of_latitude", F),
		attribute.Float64("longitude_correction", lonCorrection),
		attribute.Float64("latitude_correction", latCorrection),
		attribute.Float64("distance_correction", distCorrection),
		attribute.Float64("ecliptic_longitude", lambda),
		attribute.Float64("ecliptic_latitude", beta),
		attribute.Float64("distance_km", delta),
		attribute.Float64("right_ascension", alpha),
		attribute.Float64("declination", delta_eq),
		attribute.Float64("phase", phase),
		attribute.Float64("illumination_percent", illumination),
	)

	span.AddEvent("Lunar position calculated", trace.WithAttributes(
		attribute.Float64("right_ascension", alpha),
		attribute.Float64("declination", delta_eq),
		attribute.Float64("distance_km", delta),
		attribute.Float64("phase", phase),
		attribute.Float64("illumination_percent", illumination),
	))

	return result
}
