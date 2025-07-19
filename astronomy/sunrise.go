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
	// Solar depression angle for sunrise/sunset (geometric horizon)
	SolarDepressionAngle = 0.833

	// Degrees to radians conversion
	DegToRad = math.Pi / 180
	RadToDeg = 180 / math.Pi
)

// Location represents a geographic location with latitude and longitude
type Location struct {
	Latitude  float64
	Longitude float64
}

// SunTimes holds sunrise and sunset times
type SunTimes struct {
	Sunrise time.Time
	Sunset  time.Time
}

// CalculateSunTimes calculates sunrise and sunset times for a given location and date
func CalculateSunTimes(loc Location, date time.Time) (*SunTimes, error) {
	return CalculateSunTimesWithContext(context.Background(), loc, date)
}

// CalculateSunTimesWithContext calculates sunrise and sunset times for a given location and date with OpenTelemetry tracing
func CalculateSunTimesWithContext(ctx context.Context, loc Location, date time.Time) (*SunTimes, error) {
	// Create span for the entire calculation
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "CalculateSunTimes")
	defer span.End()
	
	// Add span attributes for tracing
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
	julianSpan.SetAttributes(
		attribute.Float64("julian_day", jd),
		attribute.Int("year", year),
		attribute.Int("month", int(month)),
		attribute.Int("day", day),
	)
	julianSpan.End()
	
	// Calculate centuries since J2000.0
	n := jd - 2451545.0
	span.SetAttributes(attribute.Float64("centuries_since_j2000", n))
	
	// Solar position calculations
	ctx, solarSpan := observer.CreateSpan(ctx, "calculateSolarPosition")
	
	// Mean longitude of the Sun
	L := math.Mod(280.460+0.9856474*n, 360.0)
	
	// Mean anomaly of the Sun
	g := math.Mod(357.528+0.9856003*n, 360.0) * DegToRad
	
	// Ecliptic longitude of the Sun
	lambda := L + 1.915*math.Sin(g) + 0.020*math.Sin(2*g)
	
	// Obliquity of the ecliptic
	epsilon := 23.439 - 0.0000004*n
	
	// Right ascension
	alpha := math.Atan2(math.Cos(epsilon*DegToRad)*math.Sin(lambda*DegToRad), math.Cos(lambda*DegToRad)) * RadToDeg
	
	// Declination
	delta := math.Asin(math.Sin(epsilon*DegToRad)*math.Sin(lambda*DegToRad)) * RadToDeg
	
	// Equation of time (in minutes)
	EqT := 4 * (L - alpha)
	
	solarSpan.SetAttributes(
		attribute.Float64("solar.mean_longitude", L),
		attribute.Float64("solar.mean_anomaly_rad", g),
		attribute.Float64("solar.ecliptic_longitude", lambda),
		attribute.Float64("solar.obliquity", epsilon),
		attribute.Float64("solar.right_ascension", alpha),
		attribute.Float64("solar.declination", delta),
		attribute.Float64("solar.equation_of_time", EqT),
	)
	solarSpan.End()
	
	// Hour angle for sunrise/sunset
	ctx, hourAngleSpan := observer.CreateSpan(ctx, "calculateHourAngle")
	latRad := loc.Latitude * DegToRad
	deltaRad := delta * DegToRad
	
	// Calculate hour angle
	cosH := (math.Cos(90.833*DegToRad) - math.Sin(latRad)*math.Sin(deltaRad)) / (math.Cos(latRad) * math.Cos(deltaRad))
	
	hourAngleSpan.SetAttributes(
		attribute.Float64("hour_angle.cos_h", cosH),
		attribute.Float64("latitude_rad", latRad),
		attribute.Float64("declination_rad", deltaRad),
	)
	
	// Check for polar day or polar night
	if cosH > 1 {
		// Polar night - sun never rises
		hourAngleSpan.SetAttributes(attribute.String("condition", "polar_night"))
		hourAngleSpan.AddEvent("Polar night detected - sun never rises")
		hourAngleSpan.End()
		
		result := &SunTimes{
			Sunrise: time.Date(year, month, day, 12, 0, 0, 0, date.Location()),
			Sunset:  time.Date(year, month, day, 12, 0, 0, 0, date.Location()),
		}
		span.SetAttributes(attribute.String("result_type", "polar_night"))
		span.AddEvent("Calculation completed", trace.WithAttributes(
			attribute.String("sunrise", result.Sunrise.Format("15:04:05")),
			attribute.String("sunset", result.Sunset.Format("15:04:05")),
		))
		return result, nil
	} else if cosH < -1 {
		// Polar day - sun never sets
		hourAngleSpan.SetAttributes(attribute.String("condition", "polar_day"))
		hourAngleSpan.AddEvent("Polar day detected - sun never sets")
		hourAngleSpan.End()
		
		result := &SunTimes{
			Sunrise: time.Date(year, month, day, 0, 0, 0, 0, date.Location()),
			Sunset:  time.Date(year, month, day, 23, 59, 59, 0, date.Location()),
		}
		span.SetAttributes(attribute.String("result_type", "polar_day"))
		span.AddEvent("Calculation completed", trace.WithAttributes(
			attribute.String("sunrise", result.Sunrise.Format("15:04:05")),
			attribute.String("sunset", result.Sunset.Format("15:04:05")),
		))
		return result, nil
	}
	
	// Hour angle in degrees
	H := math.Acos(cosH) * RadToDeg
	hourAngleSpan.SetAttributes(
		attribute.Float64("hour_angle_degrees", H),
		attribute.String("condition", "normal"),
	)
	hourAngleSpan.End()
	
	// Solar noon calculations
	ctx, solarNoonSpan := observer.CreateSpan(ctx, "calculateSolarNoon")
	solarNoon := 12.0 - loc.Longitude/15.0 - EqT/60.0
	
	// Sunrise and sunset times (in decimal hours UTC)
	sunriseDecimal := solarNoon - H/15.0
	sunsetDecimal := solarNoon + H/15.0
	
	solarNoonSpan.SetAttributes(
		attribute.Float64("solar_noon", solarNoon),
		attribute.Float64("sunrise_decimal", sunriseDecimal),
		attribute.Float64("sunset_decimal", sunsetDecimal),
	)
	solarNoonSpan.End()
	
	// Convert to time
	ctx, conversionSpan := observer.CreateSpan(ctx, "convertToTime")
	sunriseTime := decimalHoursToTime(sunriseDecimal, year, month, day, time.UTC)
	sunsetTime := decimalHoursToTime(sunsetDecimal, year, month, day, time.UTC)
	
	conversionSpan.SetAttributes(
		attribute.String("sunrise_time", sunriseTime.Format("15:04:05")),
		attribute.String("sunset_time", sunsetTime.Format("15:04:05")),
	)
	conversionSpan.End()
	
	result := &SunTimes{
		Sunrise: sunriseTime,
		Sunset:  sunsetTime,
	}
	
	// Calculate day length for metrics
	dayLength := sunsetTime.Sub(sunriseTime)
	if dayLength < 0 {
		dayLength = dayLength + 24*time.Hour
	}
	
	span.SetAttributes(
		attribute.String("result_type", "normal"),
		attribute.Float64("day_length_hours", dayLength.Hours()),
		attribute.String("final_sunrise", result.Sunrise.Format("15:04:05")),
		attribute.String("final_sunset", result.Sunset.Format("15:04:05")),
	)
	span.AddEvent("Calculation completed", trace.WithAttributes(
		attribute.String("sunrise", result.Sunrise.Format("15:04:05")),
		attribute.String("sunset", result.Sunset.Format("15:04:05")),
		attribute.Float64("day_length_hours", dayLength.Hours()),
	))
	
	return result, nil
}

// julianDate calculates Julian date for noon of the given date
func julianDate(date time.Time) float64 {
	year := date.Year()
	month := int(date.Month())
	day := date.Day()
	
	if month <= 2 {
		year--
		month += 12
	}
	
	a := year / 100
	b := 2 - a + a/4
	
	jd := math.Floor(365.25*(float64(year)+4716)) + 
		math.Floor(30.6001*(float64(month)+1)) + 
		float64(day) + float64(b) - 1524.5
	
	return jd
}

// solarPosition calculates equation of time and solar declination
func solarPosition(jd float64) (float64, float64) {
	return solarPositionWithContext(context.Background(), jd)
}

// solarPositionWithContext calculates equation of time and solar declination with tracing
func solarPositionWithContext(ctx context.Context, jd float64) (float64, float64) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "solarPosition")
	defer span.End()
	
	span.SetAttributes(attribute.Float64("julian_day", jd))
	
	n := jd - 2451545.0
	L := math.Mod(280.460+0.9856474*n, 360)
	g := math.Mod(357.528+0.9856003*n, 360) * DegToRad
	
	// Solar ecliptic longitude
	lambda := L + 1.915*math.Sin(g) + 0.020*math.Sin(2*g)
	
	// Right ascension
	ra := math.Atan2(math.Cos(23.44*DegToRad)*math.Sin(lambda*DegToRad), math.Cos(lambda*DegToRad)) * RadToDeg
	ra = math.Mod(ra+360, 360)
	
	// Equation of time (in minutes)
	eqTime := 4 * (L - ra)
	
	// Solar declination
	decl := math.Asin(math.Sin(23.44*DegToRad) * math.Sin(lambda*DegToRad))
	
	span.SetAttributes(
		attribute.Float64("centuries_since_j2000", n),
		attribute.Float64("mean_longitude", L),
		attribute.Float64("mean_anomaly_rad", g),
		attribute.Float64("ecliptic_longitude", lambda),
		attribute.Float64("right_ascension", ra),
		attribute.Float64("equation_of_time", eqTime),
		attribute.Float64("declination_rad", decl),
	)
	
	span.AddEvent("Solar position calculated", trace.WithAttributes(
		attribute.Float64("equation_of_time", eqTime),
		attribute.Float64("declination_rad", decl),
	))
	
	return eqTime, decl
}

// calculateRiseSet calculates sunrise and sunset times in minutes from midnight
func calculateRiseSet(latitude, longitude, jd, eqTime, decl float64) (float64, float64) {
	return calculateRiseSetWithContext(context.Background(), latitude, longitude, jd, eqTime, decl)
}

// calculateRiseSetWithContext calculates sunrise and sunset times in minutes from midnight with tracing
func calculateRiseSetWithContext(ctx context.Context, latitude, longitude, jd, eqTime, decl float64) (float64, float64) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "calculateRiseSet")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("latitude", latitude),
		attribute.Float64("longitude", longitude),
		attribute.Float64("julian_day", jd),
		attribute.Float64("equation_of_time", eqTime),
		attribute.Float64("declination", decl),
	)
	
	latRad := latitude * DegToRad
	
	// Hour angle with solar depression angle
	cosH := (math.Cos(SolarDepressionAngle*DegToRad) - math.Sin(latRad)*math.Sin(decl)) / 
		(math.Cos(latRad) * math.Cos(decl))
	
	span.SetAttributes(
		attribute.Float64("latitude_rad", latRad),
		attribute.Float64("cos_hour_angle", cosH),
		attribute.Float64("solar_depression_angle", SolarDepressionAngle),
	)
	
	// Check for polar day or polar night
	if cosH > 1 {
		// Polar night - sun never rises
		span.SetAttributes(attribute.String("condition", "polar_night"))
		span.AddEvent("Polar night - sun never rises")
		return 0, 0
	} else if cosH < -1 {
		// Polar day - sun never sets
		span.SetAttributes(attribute.String("condition", "polar_day"))
		span.AddEvent("Polar day - sun never sets")
		return 0, 24 * 60
	}
	
	H := math.Acos(cosH) * RadToDeg
	
	// Time corrections (equation of time and longitude correction)
	timeCorrection := eqTime + longitude*4
	
	// Sunrise and sunset times (minutes from midnight UTC)
	sunrise := 720 - 4*H - timeCorrection
	sunset := 720 + 4*H - timeCorrection
	
	// Ensure times are within 0-1440 minutes (24 hours)
	sunrise = math.Mod(sunrise+1440, 1440)
	sunset = math.Mod(sunset+1440, 1440)
	
	span.SetAttributes(
		attribute.String("condition", "normal"),
		attribute.Float64("hour_angle_degrees", H),
		attribute.Float64("time_correction", timeCorrection),
		attribute.Float64("sunrise_minutes", sunrise),
		attribute.Float64("sunset_minutes", sunset),
	)
	
	span.AddEvent("Rise/set times calculated", trace.WithAttributes(
		attribute.Float64("sunrise_minutes", sunrise),
		attribute.Float64("sunset_minutes", sunset),
	))
	
	return sunrise, sunset
}

// GetSunriseTime returns just the sunrise time for a location and date
func GetSunriseTime(loc Location, date time.Time) (time.Time, error) {
	return GetSunriseTimeWithContext(context.Background(), loc, date)
}

// GetSunriseTimeWithContext returns just the sunrise time for a location and date with tracing
func GetSunriseTimeWithContext(ctx context.Context, loc Location, date time.Time) (time.Time, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "GetSunriseTime")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
	)
	
	sunTimes, err := CalculateSunTimesWithContext(ctx, loc, date)
	if err != nil {
		span.RecordError(err)
		return time.Time{}, err
	}
	
	span.SetAttributes(attribute.String("sunrise", sunTimes.Sunrise.Format("15:04:05")))
	span.AddEvent("Sunrise time extracted", trace.WithAttributes(
		attribute.String("sunrise", sunTimes.Sunrise.Format("15:04:05")),
	))
	
	return sunTimes.Sunrise, nil
}

// GetSunsetTime returns just the sunset time for a location and date
func GetSunsetTime(loc Location, date time.Time) (time.Time, error) {
	return GetSunsetTimeWithContext(context.Background(), loc, date)
}

// GetSunsetTimeWithContext returns just the sunset time for a location and date with tracing
func GetSunsetTimeWithContext(ctx context.Context, loc Location, date time.Time) (time.Time, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "GetSunsetTime")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
	)
	
	sunTimes, err := CalculateSunTimesWithContext(ctx, loc, date)
	if err != nil {
		span.RecordError(err)
		return time.Time{}, err
	}
	
	span.SetAttributes(attribute.String("sunset", sunTimes.Sunset.Format("15:04:05")))
	span.AddEvent("Sunset time extracted", trace.WithAttributes(
		attribute.String("sunset", sunTimes.Sunset.Format("15:04:05")),
	))
	
	return sunTimes.Sunset, nil
}

// julianDayNumber calculates Julian day number for the given date
func julianDayNumber(year, month, day int) float64 {
	if month <= 2 {
		year--
		month += 12
	}
	
	a := year / 100
	b := 2 - a + a/4
	
	jd := math.Floor(365.25*(float64(year)+4716)) + 
		math.Floor(30.6001*(float64(month)+1)) + 
		float64(day) + float64(b) - 1524.5
	
	return jd
}

// decimalHoursToTime converts decimal hours to time.Time
func decimalHoursToTime(decimalHours float64, year int, month time.Month, day int, loc *time.Location) time.Time {
	// Ensure the decimal hours is within 0-24 range
	decimalHours = math.Mod(decimalHours+24, 24)
	
	hours := int(decimalHours)
	minutes := int((decimalHours - float64(hours)) * 60)
	seconds := int(((decimalHours - float64(hours)) * 60 - float64(minutes)) * 60)
	
	return time.Date(year, month, day, hours, minutes, seconds, 0, loc)
}

