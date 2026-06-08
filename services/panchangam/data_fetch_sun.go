package panchangam

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func calculateSunTimes(ctx context.Context, span trace.Span, date time.Time, location astronomy.Location) (*astronomy.SunTimes, error) {
	logger.InfoContext(ctx, "Calculating sun times",
		"operation", "CalculateSunTimesWithContext",
		"date", date.Format("2006-01-02"),
		"location", fmt.Sprintf("%.4f,%.4f", location.Latitude, location.Longitude))

	observability.RecordCalculationStart(ctx, "sun_times_calculation", map[string]interface{}{
		"date":      date.Format("2006-01-02"),
		"latitude":  location.Latitude,
		"longitude": location.Longitude,
		"timezone":  date.Location().String(),
	})

	calcStart := time.Now()
	sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, date)
	calcDuration := time.Since(calcStart)
	if err != nil {
		grpcErr := status.Error(codes.Internal, fmt.Sprintf("failed to calculate sun times: %v", err))

		observability.RecordCalculationEnd(ctx, "sun_times_calculation", false, calcDuration, nil)
		observability.RecordError(ctx, grpcErr, observability.ErrorContext{
			Severity:  observability.SeverityHigh,
			Category:  observability.CategoryCalculation,
			Operation: "sun_times_calculation",
			Component: "astronomy_service",
			Additional: map[string]interface{}{
				"calculation_type": "sun_times",
				"latitude":         location.Latitude,
				"longitude":        location.Longitude,
				"date":             date.Format("2006-01-02"),
				"duration_ms":      calcDuration.Milliseconds(),
				"original_error":   err.Error(),
			},
			Retryable:   true,
			ExpectedErr: false,
		})

		observability.RecordEvent(ctx, "Astronomical calculation failed", map[string]interface{}{
			"calculation_type": "sun_times",
			"error_type":       "calculation_failure",
			"location":         fmt.Sprintf("%.4f,%.4f", location.Latitude, location.Longitude),
			"date":             date.Format("2006-01-02"),
			"duration_ms":      calcDuration.Milliseconds(),
		})

		logger.ErrorContext(ctx, "Astronomical calculation failed",
			"operation", "CalculateSunTimesWithContext",
			"error", grpcErr,
			"location", fmt.Sprintf("%.4f,%.4f", location.Latitude, location.Longitude),
			"date", date.Format("2006-01-02"))
		span.RecordError(grpcErr)
		return nil, grpcErr
	}

	observability.RecordCalculationEnd(ctx, "sun_times_calculation", true, calcDuration, map[string]interface{}{
		"sunrise_time": sunTimes.Sunrise.Format("15:04:05"),
		"sunset_time":  sunTimes.Sunset.Format("15:04:05"),
	})
	logger.DebugContext(ctx, "Sun times calculated successfully",
		"sunrise", sunTimes.Sunrise.Format("15:04:05"),
		"sunset", sunTimes.Sunset.Format("15:04:05"))

	return sunTimes, nil
}
