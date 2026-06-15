package astronomy

import (
	"context"
	"math"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

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
	cosH := (math.Cos((90.0+LunarDepressionAngle)*DegToRad) -
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
	lunarNoon := 12.0 + (lunarPos.RightAscension-loc.Longitude)/15.0
	lunarNoon = math.Mod(lunarNoon+24, 24) // normalize to 0-24 hours

	// Moonrise and moonset times (in decimal hours local time)
	moonriseDecimal := lunarNoon - H/15.0
	moonsetDecimal := lunarNoon + H/15.0

	// Normalize times to 0-24 hour range
	moonriseDecimal = math.Mod(moonriseDecimal+24, 24)
	moonsetDecimal = math.Mod(moonsetDecimal+24, 24)

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
