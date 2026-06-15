package ephemeris

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// updateHealthStatus updates the health status of the provider.
func (s *SwissProvider) updateHealthStatus(ctx context.Context) {
	ctx, span := s.observer.CreateSpan(ctx, "swiss.updateHealthStatus")
	defer span.End()

	start := time.Now()

	// Simple health check: verify we can perform basic calculations.
	testJD := JulianDay(2451545.0) // J2000.0
	available := true
	var errorMessage string

	if testJD < s.dataStartJD || testJD > s.dataEndJD {
		available = false
		errorMessage = "Test Julian day outside valid range"
	} else {
		_ = s.calculateSunPosition(ctx, testJD)
	}

	responseTime := time.Since(start)
	now := time.Now()

	s.healthStatus = &HealthStatus{
		Available:    available,
		LastCheck:    now,
		DataStartJD:  float64(s.dataStartJD),
		DataEndJD:    float64(s.dataEndJD),
		ResponseTime: responseTime,
		ErrorMessage: errorMessage,
		Version:      s.version,
		Source:       s.name,
	}
	s.lastHealthCheck = now

	span.SetAttributes(
		attribute.Bool("available", available),
		attribute.Int64("response_time_ms", responseTime.Milliseconds()),
		attribute.String("error_message", errorMessage),
	)
	span.AddEvent("Health status updated")
}
