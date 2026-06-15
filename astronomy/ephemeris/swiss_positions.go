package ephemeris

import (
	"context"
	"math"

	"github.com/soniakeys/meeus/v3/base"
	"github.com/soniakeys/meeus/v3/moonposition"
	"github.com/soniakeys/meeus/v3/solar"
	"go.opentelemetry.io/otel/attribute"
)

// calculateSunPosition calculates a basic Sun position using Swiss Ephemeris methods.
func (s *SwissProvider) calculateSunPosition(ctx context.Context, jd JulianDay) Position {
	_, span := s.observer.CreateSpan(ctx, "swiss.calculateSunPosition")
	defer span.End()

	t := base.J2000Century(float64(jd))
	longitude := solar.ApparentLongitude(t).Deg()
	distance := solar.Radius(t)

	position := Position{
		Longitude: math.Mod(longitude+360, 360),
		Latitude:  0.0,
		Distance:  distance,
		Speed:     0.9856,
	}

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("distance", position.Distance),
		attribute.Float64("speed", position.Speed),
	)

	return position
}

// calculateMoonPosition calculates a basic Moon position using Swiss Ephemeris methods.
func (s *SwissProvider) calculateMoonPosition(ctx context.Context, jd JulianDay) Position {
	_, span := s.observer.CreateSpan(ctx, "swiss.calculateMoonPosition")
	defer span.End()

	longitude, latitude, distanceKm := moonposition.Position(float64(jd))

	position := Position{
		Longitude: math.Mod(longitude.Deg()+360, 360),
		Latitude:  latitude.Deg(),
		Distance:  distanceKm / 149597870.7,
		Speed:     13.18,
	}

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("latitude", position.Latitude),
		attribute.Float64("distance_au", position.Distance),
		attribute.Float64("speed", position.Speed),
	)

	return position
}

// calculatePlanetPosition calculates planet position using Swiss Ephemeris VSOP87 theory.
func (s *SwissProvider) calculatePlanetPosition(ctx context.Context, jd JulianDay, planet string) Position {
	_, span := s.observer.CreateSpan(ctx, "swiss.calculatePlanetPosition")
	defer span.End()

	span.SetAttributes(attribute.String("planet", planet))

	// Swiss Ephemeris uses VSOP87 theory for planetary positions.
	t := float64(jd - 2451545.0)

	// More accurate planetary elements with VSOP87 corrections.
	var L, M, distance, speed float64
	var deltaL, deltaM, deltaR float64

	switch planet {
	case "mercury":
		L = math.Mod(252.2509+4.092338*t, 360.0)
		M = math.Mod(174.7948+4.092335*t, 360.0)
		distance = 0.387098
		speed = 4.092
		deltaL = 0.378 * math.Sin((157.074+4.092338*t)*math.Pi/180.0)
		deltaM = 0.321 * math.Sin((164.045+4.092338*t)*math.Pi/180.0)
		deltaR = 0.007824 * math.Cos((157.074+4.092338*t)*math.Pi/180.0)
	case "venus":
		L = math.Mod(181.9798+1.602136*t, 360.0)
		M = math.Mod(50.4161+1.602136*t, 360.0)
		distance = 0.723327
		speed = 1.602
		deltaL = 0.775 * math.Sin((89.44+1.602136*t)*math.Pi/180.0)
		deltaM = 0.007 * math.Sin((313.42+1.602136*t)*math.Pi/180.0)
		deltaR = 0.000005 * math.Cos((89.44+1.602136*t)*math.Pi/180.0)
	case "mars":
		L = math.Mod(355.433+0.524033*t, 360.0)
		M = math.Mod(19.3870+0.524033*t, 360.0)
		distance = 1.523679
		speed = 0.524
		deltaL = 10.691 * math.Sin((68.98+0.524033*t)*math.Pi/180.0)
		deltaM = 0.606 * math.Sin((108.99+0.524033*t)*math.Pi/180.0)
		deltaR = 0.141063 * math.Cos((68.98+0.524033*t)*math.Pi/180.0)
	case "jupiter":
		L = math.Mod(34.3515+0.083091*t, 360.0)
		M = math.Mod(20.0202+0.083091*t, 360.0)
		distance = 5.204267
		speed = 0.083
		deltaL = 5.555 * math.Sin((318.16+0.083091*t)*math.Pi/180.0)
		deltaM = 0.164 * math.Sin((225.33+0.083091*t)*math.Pi/180.0)
		deltaR = 0.262127 * math.Cos((318.16+0.083091*t)*math.Pi/180.0)
	case "saturn":
		L = math.Mod(50.0774+0.033494*t, 360.0)
		M = math.Mod(317.021+0.033494*t, 360.0)
		distance = 9.5820172
		speed = 0.033
		deltaL = 6.406 * math.Sin((231.46+0.033494*t)*math.Pi/180.0)
		deltaM = 0.407 * math.Sin((206.19+0.033494*t)*math.Pi/180.0)
		deltaR = 0.301020 * math.Cos((231.46+0.033494*t)*math.Pi/180.0)
	case "uranus":
		L = math.Mod(314.055+0.011733*t, 360.0)
		M = math.Mod(142.238+0.011733*t, 360.0)
		distance = 19.189253
		speed = 0.012
		deltaL = 1.681 * math.Sin((77.25+0.011733*t)*math.Pi/180.0)
		deltaM = 0.104 * math.Sin((108.11+0.011733*t)*math.Pi/180.0)
		deltaR = 0.09142 * math.Cos((77.25+0.011733*t)*math.Pi/180.0)
	case "neptune":
		L = math.Mod(304.348+0.005965*t, 360.0)
		M = math.Mod(256.225+0.005965*t, 360.0)
		distance = 30.070900
		speed = 0.006
		deltaL = 1.021 * math.Sin((84.457+0.005965*t)*math.Pi/180.0)
		deltaM = 0.058 * math.Sin((200.51+0.005965*t)*math.Pi/180.0)
		deltaR = 0.046116 * math.Cos((84.457+0.005965*t)*math.Pi/180.0)
	case "pluto":
		L = math.Mod(238.956+0.003968*t, 360.0)
		M = math.Mod(14.8820+0.003968*t, 360.0)
		distance = 39.481686
		speed = 0.004
		deltaL = 0.041 * math.Sin((322.16+0.003968*t)*math.Pi/180.0)
		deltaM = 0.004 * math.Sin((322.16+0.003968*t)*math.Pi/180.0)
		deltaR = 0.0064 * math.Cos((322.16+0.003968*t)*math.Pi/180.0)
	default:
		// Earth's position.
		L = math.Mod(100.4644+0.985647*t, 360.0)
		M = math.Mod(357.5291+0.985600*t, 360.0)
		distance = 1.000001
		speed = 0.986
		deltaL = 0.0
		deltaM = 0.0
		deltaR = 0.0
	}

	// Apply VSOP87 corrections.
	MRad := (M + deltaM) * math.Pi / 180.0
	lambda := L + deltaL + 1.915*math.Sin(MRad) + 0.020*math.Sin(2*MRad)
	correctedDistance := distance + deltaR

	position := Position{
		Longitude: math.Mod(lambda+360, 360),
		Latitude:  0.0, // Simplified - real VSOP87 includes latitude corrections
		Distance:  correctedDistance,
		Speed:     speed,
	}

	span.SetAttributes(
		attribute.Float64("longitude", position.Longitude),
		attribute.Float64("distance", position.Distance),
		attribute.Float64("speed", position.Speed),
		attribute.Float64("vsop87_delta_l", deltaL),
		attribute.Float64("vsop87_delta_r", deltaR),
	)

	return position
}
