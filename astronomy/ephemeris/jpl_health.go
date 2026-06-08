package ephemeris

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// updateHealthStatus updates the health status of the provider.
func (j *JPLProvider) updateHealthStatus(ctx context.Context) {
	ctx, span := j.observer.CreateSpan(ctx, "jpl.updateHealthStatus")
	defer span.End()

	start := time.Now()

	// Simple health check: verify we can perform basic calculations.
	testJD := JulianDay(2451545.0) // J2000.0
	available := true
	var errorMessage string

	if testJD < j.dataStartJD || testJD > j.dataEndJD {
		available = false
		errorMessage = "Test Julian day outside valid range"
	} else {
		_ = j.calculateSunPosition(ctx, testJD)
	}

	responseTime := time.Since(start)
	now := time.Now()

	j.healthStatus = &HealthStatus{
		Available:    available,
		LastCheck:    now,
		DataStartJD:  float64(j.dataStartJD),
		DataEndJD:    float64(j.dataEndJD),
		ResponseTime: responseTime,
		ErrorMessage: errorMessage,
		Version:      j.version,
		Source:       j.name,
	}
	j.lastHealthCheck = now

	span.SetAttributes(
		attribute.Bool("available", available),
		attribute.Int64("response_time_ms", responseTime.Milliseconds()),
		attribute.String("error_message", errorMessage),
	)
	span.AddEvent("Health status updated")
}
