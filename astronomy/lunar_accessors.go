package astronomy

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

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
