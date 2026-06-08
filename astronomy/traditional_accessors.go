package astronomy

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetRahuKalam returns just the Rahu Kalam period for a location and date
func GetRahuKalam(loc Location, date time.Time) (*TimePeriod, error) {
	return GetRahuKalamWithContext(context.Background(), loc, date)
}

// GetRahuKalamWithContext returns Rahu Kalam with tracing
func GetRahuKalamWithContext(ctx context.Context, loc Location, date time.Time) (*TimePeriod, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "GetRahuKalam")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
	)

	periods, err := CalculateTraditionalPeriodsWithContext(ctx, loc, date)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("rahu_kalam", periods.RahuKalam.Start.Format("15:04:05")+" - "+periods.RahuKalam.End.Format("15:04:05")),
	)
	span.AddEvent("Rahu Kalam extracted", trace.WithAttributes(
		attribute.String("rahu_kalam", periods.RahuKalam.Start.Format("15:04:05")+" - "+periods.RahuKalam.End.Format("15:04:05")),
	))

	return periods.RahuKalam, nil
}

// GetAbhijitMuhurta returns just the Abhijit Muhurta period for a location and date
func GetAbhijitMuhurta(loc Location, date time.Time) (*TimePeriod, error) {
	return GetAbhijitMuhurtaWithContext(context.Background(), loc, date)
}

// GetAbhijitMuhurtaWithContext returns Abhijit Muhurta with tracing
func GetAbhijitMuhurtaWithContext(ctx context.Context, loc Location, date time.Time) (*TimePeriod, error) {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(ctx, "GetAbhijitMuhurta")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("location.latitude", loc.Latitude),
		attribute.Float64("location.longitude", loc.Longitude),
		attribute.String("date", date.Format("2006-01-02")),
	)

	periods, err := CalculateTraditionalPeriodsWithContext(ctx, loc, date)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("abhijit_muhurta", periods.AbhijitMuhurta.Start.Format("15:04:05")+" - "+periods.AbhijitMuhurta.End.Format("15:04:05")),
		attribute.Bool("abhijit_auspicious", periods.AbhijitMuhurta.Auspicious),
	)
	span.AddEvent("Abhijit Muhurta extracted", trace.WithAttributes(
		attribute.String("abhijit_muhurta", periods.AbhijitMuhurta.Start.Format("15:04:05")+" - "+periods.AbhijitMuhurta.End.Format("15:04:05")),
	))

	return periods.AbhijitMuhurta, nil
}
