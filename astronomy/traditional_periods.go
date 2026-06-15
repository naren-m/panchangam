package astronomy

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// calculateRahuKalam calculates Rahu Kalam based on the day of the week
// Rahu Kalam is considered inauspicious and varies by day of the week
func calculateRahuKalam(ctx context.Context, sunrise, sunset time.Time, date time.Time) *TimePeriod {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "calculateRahuKalam")
	defer span.End()

	dayOfWeek := date.Weekday()
	span.SetAttributes(attribute.String("day_of_week", dayOfWeek.String()))

	// Calculate day length and divide into 8 parts
	dayLength := sunset.Sub(sunrise)
	partDuration := dayLength / 8

	span.SetAttributes(
		attribute.Float64("day_length_hours", dayLength.Hours()),
		attribute.Float64("part_duration_minutes", partDuration.Minutes()),
	)

	// Rahu Kalam timing based on day of the week (traditional calculation)
	var rahuPart int
	switch dayOfWeek {
	case time.Sunday:
		rahuPart = 7 // 7th part (4:30 PM - 6:00 PM traditionally)
	case time.Monday:
		rahuPart = 1 // 1st part (7:30 AM - 9:00 AM traditionally)
	case time.Tuesday:
		rahuPart = 6 // 6th part (3:00 PM - 4:30 PM traditionally)
	case time.Wednesday:
		rahuPart = 4 // 4th part (12:00 PM - 1:30 PM traditionally)
	case time.Thursday:
		rahuPart = 3 // 3rd part (10:30 AM - 12:00 PM traditionally)
	case time.Friday:
		rahuPart = 2 // 2nd part (9:00 AM - 10:30 AM traditionally)
	case time.Saturday:
		rahuPart = 5 // 5th part (1:30 PM - 3:00 PM traditionally)
	}

	// Calculate start and end times
	start := sunrise.Add(time.Duration(rahuPart-1) * partDuration)
	end := sunrise.Add(time.Duration(rahuPart) * partDuration)

	span.SetAttributes(
		attribute.Int("rahu_part", rahuPart),
		attribute.String("start_time", start.Format("15:04:05")),
		attribute.String("end_time", end.Format("15:04:05")),
	)

	result := &TimePeriod{
		Start:       start,
		End:         end,
		Duration:    int(partDuration.Minutes()),
		Description: "Rahu Kalam - Inauspicious period ruled by Rahu",
		Auspicious:  false,
	}

	span.AddEvent("Rahu Kalam calculated", trace.WithAttributes(
		attribute.String("time_period", start.Format("15:04:05")+" - "+end.Format("15:04:05")),
		attribute.Int("duration_minutes", result.Duration),
	))

	return result
}

// calculateYamagandam calculates Yamagandam based on the day of the week
// Yamagandam is another inauspicious period
func calculateYamagandam(ctx context.Context, sunrise, sunset time.Time, date time.Time) *TimePeriod {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "calculateYamagandam")
	defer span.End()

	dayOfWeek := date.Weekday()
	span.SetAttributes(attribute.String("day_of_week", dayOfWeek.String()))

	// Calculate day length and divide into 8 parts
	dayLength := sunset.Sub(sunrise)
	partDuration := dayLength / 8

	// Yamagandam timing based on day of the week
	var yamaPart int
	switch dayOfWeek {
	case time.Sunday:
		yamaPart = 4 // 4th part
	case time.Monday:
		yamaPart = 7 // 7th part
	case time.Tuesday:
		yamaPart = 2 // 2nd part
	case time.Wednesday:
		yamaPart = 5 // 5th part
	case time.Thursday:
		yamaPart = 8 // 8th part
	case time.Friday:
		yamaPart = 6 // 6th part
	case time.Saturday:
		yamaPart = 3 // 3rd part
	}

	// Calculate start and end times
	start := sunrise.Add(time.Duration(yamaPart-1) * partDuration)
	end := sunrise.Add(time.Duration(yamaPart) * partDuration)

	span.SetAttributes(
		attribute.Int("yama_part", yamaPart),
		attribute.String("start_time", start.Format("15:04:05")),
		attribute.String("end_time", end.Format("15:04:05")),
	)

	result := &TimePeriod{
		Start:       start,
		End:         end,
		Duration:    int(partDuration.Minutes()),
		Description: "Yamagandam - Inauspicious period ruled by Yama",
		Auspicious:  false,
	}

	span.AddEvent("Yamagandam calculated", trace.WithAttributes(
		attribute.String("time_period", start.Format("15:04:05")+" - "+end.Format("15:04:05")),
		attribute.Int("duration_minutes", result.Duration),
	))

	return result
}

// calculateGulikaKalam calculates Gulika Kalam based on the day of the week
// Gulika Kalam is also considered inauspicious
func calculateGulikaKalam(ctx context.Context, sunrise, sunset time.Time, date time.Time) *TimePeriod {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "calculateGulikaKalam")
	defer span.End()

	dayOfWeek := date.Weekday()
	span.SetAttributes(attribute.String("day_of_week", dayOfWeek.String()))

	// Calculate day length and divide into 8 parts
	dayLength := sunset.Sub(sunrise)
	partDuration := dayLength / 8

	// Gulika Kalam timing based on day of the week
	var gulikaPart int
	switch dayOfWeek {
	case time.Sunday:
		gulikaPart = 6 // 6th part
	case time.Monday:
		gulikaPart = 8 // 8th part
	case time.Tuesday:
		gulikaPart = 4 // 4th part
	case time.Wednesday:
		gulikaPart = 7 // 7th part
	case time.Thursday:
		gulikaPart = 2 // 2nd part
	case time.Friday:
		gulikaPart = 5 // 5th part
	case time.Saturday:
		gulikaPart = 1 // 1st part
	}

	// Calculate start and end times
	start := sunrise.Add(time.Duration(gulikaPart-1) * partDuration)
	end := sunrise.Add(time.Duration(gulikaPart) * partDuration)

	span.SetAttributes(
		attribute.Int("gulika_part", gulikaPart),
		attribute.String("start_time", start.Format("15:04:05")),
		attribute.String("end_time", end.Format("15:04:05")),
	)

	result := &TimePeriod{
		Start:       start,
		End:         end,
		Duration:    int(partDuration.Minutes()),
		Description: "Gulika Kalam - Inauspicious period ruled by Gulika",
		Auspicious:  false,
	}

	span.AddEvent("Gulika Kalam calculated", trace.WithAttributes(
		attribute.String("time_period", start.Format("15:04:05")+" - "+end.Format("15:04:05")),
		attribute.Int("duration_minutes", result.Duration),
	))

	return result
}

// calculateAbhijitMuhurta calculates Abhijit Muhurta - the most auspicious time of the day
// Abhijit Muhurta is approximately the 8th muhurta of the day (around midday)
func calculateAbhijitMuhurta(ctx context.Context, sunrise, sunset time.Time) *TimePeriod {
	observer := observability.Observer()
	_, span := observer.CreateSpan(ctx, "calculateAbhijitMuhurta")
	defer span.End()

	// Abhijit Muhurta is the 48-minute window centered on local solar noon.
	dayLength := sunset.Sub(sunrise)
	duration := 48 * time.Minute
	solarNoon := sunrise.Add(dayLength / 2)
	start := solarNoon.Add(-duration / 2)
	end := solarNoon.Add(duration / 2)

	span.SetAttributes(
		attribute.Float64("day_length_hours", dayLength.Hours()),
		attribute.Float64("duration_minutes", duration.Minutes()),
		attribute.String("solar_noon", solarNoon.Format("15:04:05")),
	)

	span.SetAttributes(
		attribute.Bool("abhijit_valid", true),
		attribute.String("start_time", start.Format("15:04:05")),
		attribute.String("end_time", end.Format("15:04:05")),
	)

	result := &TimePeriod{
		Start:       start,
		End:         end,
		Duration:    int(duration.Minutes()),
		Description: "Abhijit Muhurta - Most auspicious period of the day",
		Auspicious:  true,
	}

	span.AddEvent("Abhijit Muhurta calculated", trace.WithAttributes(
		attribute.String("time_period", start.Format("15:04:05")+" - "+end.Format("15:04:05")),
		attribute.Int("duration_minutes", result.Duration),
		attribute.Bool("is_valid", true),
	))

	return result
}
