package astronomy

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CalculateTraditionalPeriods calculates all traditional time periods for a given location and date
func CalculateTraditionalPeriods(loc Location, date time.Time) (*TraditionalPeriods, error) {
	return CalculateTraditionalPeriodsWithContext(context.Background(), loc, date)
}

// CalculateTraditionalPeriodsWithContext calculates traditional periods with OpenTelemetry tracing
func CalculateTraditionalPeriodsWithContext(ctx context.Context, loc Location, date time.Time) (*TraditionalPeriods, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "CalculateTraditionalPeriods")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
		attribute.String("timezone", date.Location().String()),
	)

	// First, get sunrise and sunset times
	ctx, sunTimesSpan := observer.CreateSpan(ctx, "getSunTimes")
	sunTimes, err := CalculateSunTimesWithContext(ctx, loc, date)
	if err != nil {
		sunTimesSpan.RecordError(err)
		span.RecordError(err)
		return nil, err
	}
	sunrise := sunTimes.Sunrise.In(date.Location())
	sunset := sunTimes.Sunset.In(date.Location())
	if !sunset.After(sunrise) {
		sunset = sunset.Add(24 * time.Hour)
	}
	sunTimesSpan.SetAttributes(
		attribute.String("sunrise", sunrise.Format("15:04:05")),
		attribute.String("sunset", sunset.Format("15:04:05")),
	)
	sunTimesSpan.End()

	// Calculate day length in minutes
	dayLength := sunset.Sub(sunrise)
	dayLengthMinutes := int(dayLength.Minutes())

	span.SetAttributes(
		attribute.Float64("day_length_hours", dayLength.Hours()),
		attribute.Int("day_length_minutes", dayLengthMinutes),
	)

	// Calculate Rahu Kalam
	ctx, rahuSpan := observer.CreateSpan(ctx, "calculateRahuKalam")
	rahuKalam := calculateRahuKalam(ctx, sunrise, sunset, date)
	rahuSpan.SetAttributes(
		attribute.String("rahu_kalam_start", rahuKalam.Start.Format("15:04:05")),
		attribute.String("rahu_kalam_end", rahuKalam.End.Format("15:04:05")),
		attribute.Int("rahu_kalam_duration", rahuKalam.Duration),
	)
	rahuSpan.End()

	// Calculate Yamagandam
	ctx, yamaSpan := observer.CreateSpan(ctx, "calculateYamagandam")
	yamagandam := calculateYamagandam(ctx, sunrise, sunset, date)
	yamaSpan.SetAttributes(
		attribute.String("yamagandam_start", yamagandam.Start.Format("15:04:05")),
		attribute.String("yamagandam_end", yamagandam.End.Format("15:04:05")),
		attribute.Int("yamagandam_duration", yamagandam.Duration),
	)
	yamaSpan.End()

	// Calculate Gulika Kalam
	ctx, gulikaSpan := observer.CreateSpan(ctx, "calculateGulikaKalam")
	gulikaKalam := calculateGulikaKalam(ctx, sunrise, sunset, date)
	gulikaSpan.SetAttributes(
		attribute.String("gulika_kalam_start", gulikaKalam.Start.Format("15:04:05")),
		attribute.String("gulika_kalam_end", gulikaKalam.End.Format("15:04:05")),
		attribute.Int("gulika_kalam_duration", gulikaKalam.Duration),
	)
	gulikaSpan.End()

	// Calculate Abhijit Muhurta
	ctx, abhijitSpan := observer.CreateSpan(ctx, "calculateAbhijitMuhurta")
	abhijitMuhurta := calculateAbhijitMuhurta(ctx, sunrise, sunset)
	abhijitSpan.SetAttributes(
		attribute.String("abhijit_start", abhijitMuhurta.Start.Format("15:04:05")),
		attribute.String("abhijit_end", abhijitMuhurta.End.Format("15:04:05")),
		attribute.Int("abhijit_duration", abhijitMuhurta.Duration),
	)
	abhijitSpan.End()

	result := &TraditionalPeriods{
		RahuKalam:      rahuKalam,
		Yamagandam:     yamagandam,
		GulikaKalam:    gulikaKalam,
		AbhijitMuhurta: abhijitMuhurta,
	}

	span.AddEvent("Traditional periods calculated", trace.WithAttributes(
		attribute.String("rahu_kalam", rahuKalam.Start.Format("15:04:05")+" - "+rahuKalam.End.Format("15:04:05")),
		attribute.String("yamagandam", yamagandam.Start.Format("15:04:05")+" - "+yamagandam.End.Format("15:04:05")),
		attribute.String("gulika_kalam", gulikaKalam.Start.Format("15:04:05")+" - "+gulikaKalam.End.Format("15:04:05")),
		attribute.String("abhijit_muhurta", abhijitMuhurta.Start.Format("15:04:05")+" - "+abhijitMuhurta.End.Format("15:04:05")),
	))

	return result, nil
}
