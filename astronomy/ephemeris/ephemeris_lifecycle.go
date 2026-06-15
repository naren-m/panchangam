package ephemeris

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// tryProvider attempts to get data from a provider with observability.
func (m *Manager) tryProvider(ctx context.Context, provider EphemerisProvider, providerType string, operation func(EphemerisProvider) (interface{}, error)) (interface{}, error) {
	if provider == nil {
		return nil, fmt.Errorf("%s provider is nil", providerType)
	}

	_, span := m.observer.CreateSpan(ctx, fmt.Sprintf("ephemeris.try_%s_provider", providerType))
	defer span.End()

	span.SetAttributes(
		attribute.String("provider_type", providerType),
		attribute.String("provider_name", provider.GetProviderName()),
		attribute.String("provider_version", provider.GetVersion()),
	)

	start := time.Now()
	result, err := operation(provider)
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Int64("response_time_ms", duration.Milliseconds()),
		attribute.Bool("success", err == nil),
	)

	if err != nil {
		span.RecordError(err)
		span.AddEvent("Provider operation failed")
		return nil, err
	}

	span.AddEvent("Provider operation succeeded")
	return result, nil
}

// GetHealthStatus returns the health status of all providers.
func (m *Manager) GetHealthStatus(ctx context.Context) (map[string]*HealthStatus, error) {
	ctx, span := m.observer.CreateSpan(ctx, "ephemeris.GetHealthStatus")
	defer span.End()

	status := make(map[string]*HealthStatus)

	if m.primary != nil {
		if health, err := m.primary.GetHealthStatus(ctx); err == nil {
			status["primary"] = health
		}
	}

	if m.fallback != nil {
		if health, err := m.fallback.GetHealthStatus(ctx); err == nil {
			status["fallback"] = health
		}
	}

	span.SetAttributes(attribute.Int("provider_count", len(status)))
	span.AddEvent("Health status retrieved for all providers")

	return status, nil
}

// Close closes all providers and releases resources.
func (m *Manager) Close() error {
	var closeErrors []error

	if m.primary != nil {
		if err := m.primary.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("primary provider close error: %w", err))
		}
	}

	if m.fallback != nil {
		if err := m.fallback.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("fallback provider close error: %w", err))
		}
	}

	if m.cache != nil {
		if err := m.cache.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("cache close error: %w", err))
		}
	}

	if m.healthChecker != nil {
		m.healthChecker.Stop()
	}

	if len(closeErrors) > 0 {
		return fmt.Errorf("errors during close: %w", errors.Join(closeErrors...))
	}

	return nil
}
