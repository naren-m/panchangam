package astronomy

import (
	"context"
	"math"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	// Lunar constants
	LunarParallax      = 0.125  // degrees, average lunar parallax
	LunarSemidiameter  = 0.25   // degrees, average lunar semidiameter
	LunarDepressionAngle = 0.125 + 0.25 // parallax + semidiameter
	
	// Lunar orbital constants
	LunarMeanDistance  = 384400.0 // km, mean distance to moon
	LunarSynodicMonth  = 29.530588853 // days, synodic month
	LunarEccentricity  = 0.0549 // orbital eccentricity
	
	// Lunar phase constants
	NewMoon    = 0.0
	FirstQuarter = 0.25
	FullMoon   = 0.5
	LastQuarter = 0.75
)

// LunarTimes holds moonrise and moonset times
type LunarTimes struct {
	Moonrise time.Time
	Moonset  time.Time
	IsVisible bool // whether moon is visible (not below horizon all day)
}

// LunarPosition represents the moon's position
type LunarPosition struct {
	RightAscension float64 // degrees
	Declination    float64 // degrees
	Distance       float64 // km
	Phase          float64 // 0.0 = new, 0.5 = full
	Illumination   float64 // percentage illuminated
}

// LunarPhase represents moon phase information
type LunarPhase struct {
	Phase        float64   // 0.0-1.0, where 0=new moon, 0.5=full moon
	Illumination float64   // percentage illuminated (0-100)
	Name         string    // phase name
	Age          float64   // days since new moon
	NextPhase    time.Time // time of next major phase
}

// CalculateLunarTimes calculates moonrise and moonset times for a given location and date
func CalculateLunarTimes(loc Location, date time.Time) (*LunarTimes, error) {
	return CalculateLunarTimesWithContext(context.Background(), loc, date)
}

// CalculateLunarTimesWithContext calculates moonrise and moonset times with OpenTelemetry tracing
func CalculateLunarTimesWithContext(ctx context.Context, loc Location, date time.Time) (*LunarTimes, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "CalculateLunarTimes")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
		attribute.String("timezone", date.Location().String()),
	)
	
	year, month, day := date.Date()
	
	// Convert to Julian day number
	ctx, julianSpan := observer.CreateSpan(ctx, "calculateJulianDay")
	jd := julianDayNumber(year, int(month), day)
	julianSpan.SetAttributes(attribute.Float64("julian_day", jd))
	julianSpan.End()
	
	// Calculate lunar position for the day
	ctx, positionSpan := observer.CreateSpan(ctx, "calculateLunarPosition") 
	lunarPos := calculateLunarPositionJD(ctx, jd)
	positionSpan.SetAttributes(
		attribute.Float64("lunar.right_ascension", lunarPos.RightAscension),
		attribute.Float64("lunar.declination", lunarPos.Declination),
		attribute.Float64("lunar.distance", lunarPos.Distance),
		attribute.Float64("lunar.phase", lunarPos.Phase),
	)
	positionSpan.End()
	
	// Calculate hour angle for moonrise/moonset
	ctx, hourAngleSpan := observer.CreateSpan(ctx, "calculateLunarHourAngle")
	latRad := loc.Latitude * DegToRad
	declRad := lunarPos.Declination * DegToRad
	
	// Hour angle calculation with lunar depression angle
	cosH := (math.Cos((90.0 + LunarDepressionAngle) * DegToRad) - 
		math.Sin(latRad)*math.Sin(declRad)) / 
		(math.Cos(latRad) * math.Cos(declRad))
	
	hourAngleSpan.SetAttributes(
		attribute.Float64("hour_angle.cos_h", cosH),
		attribute.Float64("latitude_rad", latRad),
		attribute.Float64("declination_rad", declRad),
		attribute.Float64("lunar_depression_angle", LunarDepressionAngle),
	)
	
	// Check for circumstances where moon doesn't rise or set
	if cosH > 1 {
		// Moon never rises (always below horizon)
		hourAngleSpan.SetAttributes(attribute.String("condition", "moon_never_rises"))
		hourAngleSpan.AddEvent("Moon never rises - always below horizon")
		hourAngleSpan.End()
		
		result := &LunarTimes{
			Moonrise:  time.Date(year, month, day, 12, 0, 0, 0, date.Location()),
			Moonset:   time.Date(year, month, day, 12, 0, 0, 0, date.Location()),
			IsVisible: false,
		}
		span.SetAttributes(attribute.String("result_type", "moon_never_rises"))
		return result, nil
	} else if cosH < -1 {
		// Moon never sets (always above horizon)
		hourAngleSpan.SetAttributes(attribute.String("condition", "moon_never_sets"))
		hourAngleSpan.AddEvent("Moon never sets - always above horizon")
		hourAngleSpan.End()
		
		result := &LunarTimes{
			Moonrise:  time.Date(year, month, day, 0, 0, 0, 0, date.Location()),
			Moonset:   time.Date(year, month, day, 23, 59, 59, 0, date.Location()),
			IsVisible: true,
		}
		span.SetAttributes(attribute.String("result_type", "moon_never_sets"))
		return result, nil
	}
	
	// Hour angle in degrees
	H := math.Acos(cosH) * RadToDeg
	hourAngleSpan.SetAttributes(
		attribute.Float64("hour_angle_degrees", H),
		attribute.String("condition", "normal"),
	)
	hourAngleSpan.End()
	
	// Calculate lunar noon (when moon crosses meridian)
	ctx, lunarNoonSpan := observer.CreateSpan(ctx, "calculateLunarNoon")
	
	// Convert lunar right ascension to hour angle
	// This is an approximation - more precise calculation would use sidereal time
	lunarNoon := 12.0 + (lunarPos.RightAscension - loc.Longitude) / 15.0
	lunarNoon = math.Mod(lunarNoon + 24, 24) // normalize to 0-24 hours
	
	// Moonrise and moonset times (in decimal hours local time)
	moonriseDecimal := lunarNoon - H/15.0
	moonsetDecimal := lunarNoon + H/15.0
	
	// Normalize times to 0-24 hour range
	moonriseDecimal = math.Mod(moonriseDecimal + 24, 24)
	moonsetDecimal = math.Mod(moonsetDecimal + 24, 24)
	
	lunarNoonSpan.SetAttributes(
		attribute.Float64("lunar_noon", lunarNoon),
		attribute.Float64("moonrise_decimal", moonriseDecimal),
		attribute.Float64("moonset_decimal", moonsetDecimal),
	)
	lunarNoonSpan.End()
	
	// Convert to time objects
	_, conversionSpan := observer.CreateSpan(ctx, "convertToTime")
	moonriseTime := decimalHoursToTime(moonriseDecimal, year, month, day, date.Location())
	moonsetTime := decimalHoursToTime(moonsetDecimal, year, month, day, date.Location())
	
	// Adjust for the fact that lunar events might occur on adjacent days
	// This is a simplified adjustment - more precise calculations would track across days
	if moonriseDecimal > moonsetDecimal {
		// Moonrise is after moonset, so moonrise is likely the next day
		moonriseTime = moonriseTime.Add(24 * time.Hour)
	}
	
	conversionSpan.SetAttributes(
		attribute.String("moonrise_time", moonriseTime.Format("15:04:05")),
		attribute.String("moonset_time", moonsetTime.Format("15:04:05")),
	)
	conversionSpan.End()
	
	result := &LunarTimes{
		Moonrise:  moonriseTime,
		Moonset:   moonsetTime,
		IsVisible: true,
	}
	
	span.SetAttributes(
		attribute.String("result_type", "normal"),
		attribute.String("final_moonrise", result.Moonrise.Format("15:04:05")),
		attribute.String("final_moonset", result.Moonset.Format("15:04:05")),
		attribute.Bool("is_visible", result.IsVisible),
	)
	span.AddEvent("Lunar calculation completed", trace.WithAttributes(
		attribute.String("moonrise", result.Moonrise.Format("15:04:05")),
		attribute.String("moonset", result.Moonset.Format("15:04:05")),
		attribute.Bool("is_visible", result.IsVisible),
	))
	
	return result, nil
}

// calculateLunarPositionJD calculates the moon's position for a given Julian day
func calculateLunarPositionJD(ctx context.Context, jd float64) *LunarPosition {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "calculateLunarPositionJD")
	defer span.End()
	
	span.SetAttributes(attribute.Float64("julian_day", jd))
	
	// Days since J2000.0
	T := (jd - 2451545.0) / 36525.0
	
	// Moon's mean longitude (degrees)
	L := math.Mod(218.3164477 + 481267.88123421*T - 0.0015786*T*T + T*T*T/538841.0 - T*T*T*T/65194000.0, 360.0)
	
	// Mean elongation of the Moon from the Sun (degrees)
	D := math.Mod(297.8501921 + 445267.1114034*T - 0.0018819*T*T + T*T*T/545868.0 - T*T*T*T/113065000.0, 360.0)
	
	// Sun's mean anomaly (degrees)
	M := math.Mod(357.5291092 + 35999.0502909*T - 0.0001536*T*T + T*T*T/24490000.0, 360.0)
	
	// Moon's mean anomaly (degrees)
	MPrime := math.Mod(134.9633964 + 477198.8675055*T + 0.0087414*T*T + T*T*T/69699.0 - T*T*T*T/14712000.0, 360.0)
	
	// Moon's argument of latitude (degrees)
	F := math.Mod(93.2720950 + 483202.0175233*T - 0.0036539*T*T - T*T*T/3526000.0 + T*T*T*T/863310000.0, 360.0)
	
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
	lambda := L + lonCorrection  // Ecliptic longitude
	beta := latCorrection        // Ecliptic latitude
	delta := 385000.56 + distCorrection  // Distance in km
	
	// Convert to equatorial coordinates
	epsilon := 23.4392911 * DegToRad  // Obliquity of ecliptic (simplified)
	lambdaRad := lambda * DegToRad
	betaRad := beta * DegToRad
	
	// Right ascension and declination
	alpha := math.Atan2(math.Sin(lambdaRad)*math.Cos(epsilon) - math.Tan(betaRad)*math.Sin(epsilon), math.Cos(lambdaRad)) * RadToDeg
	if alpha < 0 {
		alpha += 360
	}
	
	delta_eq := math.Asin(math.Sin(betaRad)*math.Cos(epsilon) + math.Cos(betaRad)*math.Sin(epsilon)*math.Sin(lambdaRad)) * RadToDeg
	
	// Calculate lunar phase
	// Phase angle between Sun and Moon as seen from Earth
	elongation := math.Mod(D, 360.0) * DegToRad
	phaseAngle := math.Pi - elongation
	if phaseAngle < 0 {
		phaseAngle += 2 * math.Pi
	}
	
	phase := (1 - math.Cos(phaseAngle)) / 2  // 0 = new moon, 1 = full moon
	illumination := phase * 100  // percentage
	
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

// CalculateLunarPhase calculates detailed lunar phase information
func CalculateLunarPhase(date time.Time) (*LunarPhase, error) {
	return CalculateLunarPhaseWithContext(context.Background(), date)
}

// CalculateLunarPhaseWithContext calculates lunar phase with OpenTelemetry tracing
func CalculateLunarPhaseWithContext(ctx context.Context, date time.Time) (*LunarPhase, error) {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "CalculateLunarPhase")
	defer span.End()
	
	span.SetAttributes(attribute.String("date", date.Format("2006-01-02")))
	
	year, month, day := date.Date()
	jd := julianDayNumber(year, int(month), day)
	
	// Calculate days since known new moon (January 6, 2000)
	knownNewMoon := julianDayNumber(2000, 1, 6)  // JD 2451549.5
	daysSinceKnownNew := jd - knownNewMoon
	
	// Calculate the current lunation number
	lunationNumber := daysSinceKnownNew / LunarSynodicMonth
	currentLunation := math.Floor(lunationNumber)
	
	// Days into current lunar cycle
	daysIntoCycle := (lunationNumber - currentLunation) * LunarSynodicMonth
	
	// Calculate phase (0.0 = new moon, 0.5 = full moon)
	phase := daysIntoCycle / LunarSynodicMonth
	
	// Calculate illumination percentage
	phaseAngle := phase * 2 * math.Pi
	illumination := (1 - math.Cos(phaseAngle)) / 2 * 100
	
	// Determine phase name
	var phaseName string
	if phase < 0.125 || phase >= 0.875 {
		phaseName = "New Moon"
	} else if phase < 0.375 {
		phaseName = "Waxing Crescent"
	} else if phase < 0.625 {
		phaseName = "Full Moon"
	} else {
		phaseName = "Waning Crescent"
	}
	
	// More precise phase names
	if phase >= 0.125 && phase < 0.25 {
		phaseName = "Waxing Crescent"
	} else if phase >= 0.25 && phase < 0.375 {
		phaseName = "First Quarter"
	} else if phase >= 0.375 && phase < 0.5 {
		phaseName = "Waxing Gibbous"
	} else if phase >= 0.5 && phase < 0.625 {
		phaseName = "Full Moon"
	} else if phase >= 0.625 && phase < 0.75 {
		phaseName = "Waning Gibbous"
	} else if phase >= 0.75 && phase < 0.875 {
		phaseName = "Last Quarter"
	}
	
	// Calculate next major phase
	var nextPhaseJD float64
	var nextPhaseType float64
	
	if phase < 0.25 {
		nextPhaseType = 0.25 // First Quarter
		nextPhaseJD = knownNewMoon + (currentLunation * LunarSynodicMonth) + (0.25 * LunarSynodicMonth)
	} else if phase < 0.5 {
		nextPhaseType = 0.5 // Full Moon
		nextPhaseJD = knownNewMoon + (currentLunation * LunarSynodicMonth) + (0.5 * LunarSynodicMonth)
	} else if phase < 0.75 {
		nextPhaseType = 0.75 // Last Quarter
		nextPhaseJD = knownNewMoon + (currentLunation * LunarSynodicMonth) + (0.75 * LunarSynodicMonth)
	} else {
		nextPhaseType = 0.0 // New Moon (next cycle)
		nextPhaseJD = knownNewMoon + ((currentLunation + 1) * LunarSynodicMonth)
	}
	
	// Convert Julian day to time
	nextPhaseTime := julianDayToTime(nextPhaseJD, date.Location())
	
	result := &LunarPhase{
		Phase:        phase,
		Illumination: illumination,
		Name:         phaseName,
		Age:          daysIntoCycle,
		NextPhase:    nextPhaseTime,
	}
	
	span.SetAttributes(
		attribute.Float64("julian_day", jd),
		attribute.Float64("days_since_known_new", daysSinceKnownNew),
		attribute.Float64("lunation_number", lunationNumber),
		attribute.Float64("current_lunation", currentLunation),
		attribute.Float64("days_into_cycle", daysIntoCycle),
		attribute.Float64("phase", phase),
		attribute.Float64("illumination", illumination),
		attribute.String("phase_name", phaseName),
		attribute.Float64("age_days", daysIntoCycle),
		attribute.Float64("next_phase_type", nextPhaseType),
		attribute.String("next_phase_time", nextPhaseTime.Format("2006-01-02 15:04:05")),
	)
	
	span.AddEvent("Lunar phase calculated", trace.WithAttributes(
		attribute.Float64("phase", phase),
		attribute.Float64("illumination", illumination),
		attribute.String("phase_name", phaseName),
		attribute.Float64("age_days", daysIntoCycle),
	))
	
	return result, nil
}

// julianDayToTime converts Julian day to time.Time
func julianDayToTime(jd float64, loc *time.Location) time.Time {
	// Add 0.5 to adjust for noon-based Julian day
	jd += 0.5
	
	// Integer part is the number of days since Julian epoch
	z := int(jd)
	f := jd - float64(z)
	
	var a int
	if z < 2299161 {
		a = z
	} else {
		alpha := int((float64(z) - 1867216.25) / 36524.25)
		a = z + 1 + alpha - alpha/4
	}
	
	b := a + 1524
	c := int((float64(b) - 122.1) / 365.25)
	d := int(365.25 * float64(c))
	e := int(float64(b-d) / 30.6001)
	
	// Day of month
	day := b - d - int(30.6001*float64(e))
	
	// Month
	var month int
	if e < 14 {
		month = e - 1
	} else {
		month = e - 13
	}
	
	// Year
	var year int
	if month > 2 {
		year = c - 4716
	} else {
		year = c - 4715
	}
	
	// Hours, minutes, seconds from fractional part
	dayFraction := f
	hours := int(dayFraction * 24)
	minuteFraction := (dayFraction*24 - float64(hours)) * 60
	minutes := int(minuteFraction)
	seconds := int((minuteFraction - float64(minutes)) * 60)
	
	return time.Date(year, time.Month(month), day, hours, minutes, seconds, 0, loc)
}

// GetMoonriseTime returns just the moonrise time for a location and date
func GetMoonriseTime(loc Location, date time.Time) (time.Time, error) {
	return GetMoonriseTimeWithContext(context.Background(), loc, date)
}

// GetMoonriseTimeWithContext returns just the moonrise time with tracing
func GetMoonriseTimeWithContext(ctx context.Context, loc Location, date time.Time) (time.Time, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "GetMoonriseTime")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
	)
	
	lunarTimes, err := CalculateLunarTimesWithContext(ctx, loc, date)
	if err != nil {
		span.RecordError(err)
		return time.Time{}, err
	}
	
	span.SetAttributes(
		attribute.String("moonrise", lunarTimes.Moonrise.Format("15:04:05")),
		attribute.Bool("is_visible", lunarTimes.IsVisible),
	)
	span.AddEvent("Moonrise time extracted", trace.WithAttributes(
		attribute.String("moonrise", lunarTimes.Moonrise.Format("15:04:05")),
	))
	
	return lunarTimes.Moonrise, nil
}

// GetMoonsetTime returns just the moonset time for a location and date
func GetMoonsetTime(loc Location, date time.Time) (time.Time, error) {
	return GetMoonsetTimeWithContext(context.Background(), loc, date)
}

// GetMoonsetTimeWithContext returns just the moonset time with tracing
func GetMoonsetTimeWithContext(ctx context.Context, loc Location, date time.Time) (time.Time, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "GetMoonsetTime")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
	)
	
	lunarTimes, err := CalculateLunarTimesWithContext(ctx, loc, date)
	if err != nil {
		span.RecordError(err)
		return time.Time{}, err
	}
	
	span.SetAttributes(
		attribute.String("moonset", lunarTimes.Moonset.Format("15:04:05")),
		attribute.Bool("is_visible", lunarTimes.IsVisible),
	)
	span.AddEvent("Moonset time extracted", trace.WithAttributes(
		attribute.String("moonset", lunarTimes.Moonset.Format("15:04:05")),
	))
	
	return lunarTimes.Moonset, nil
}