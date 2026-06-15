package ephemeris

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// GetPlanetaryPositions retrieves planetary positions with caching and fallback.
func (m *Manager) GetPlanetaryPositions(ctx context.Context, jd JulianDay) (*PlanetaryPositions, error) {
	ctx, span := m.observer.CreateSpan(ctx, "ephemeris.GetPlanetaryPositions")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("operation", "get_planetary_positions"),
	)

	cacheKey := fmt.Sprintf("planetary_positions_%f", jd)
	if cached, found := m.cache.Get(ctx, cacheKey); found {
		span.SetAttributes(attribute.Bool("cache_hit", true))
		span.AddEvent("Cache hit for planetary positions")
		if positions, ok := cached.(*PlanetaryPositions); ok {
			return positions, nil
		}
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))

	result, err := m.tryProvider(ctx, m.primary, "primary", func(provider EphemerisProvider) (interface{}, error) {
		return provider.GetPlanetaryPositions(ctx, jd)
	})

	var positions *PlanetaryPositions
	if err == nil {
		positions = result.(*PlanetaryPositions)
	} else {
		span.AddEvent("Primary provider failed, trying fallback")

		result, err = m.tryProvider(ctx, m.fallback, "fallback", func(provider EphemerisProvider) (interface{}, error) {
			return provider.GetPlanetaryPositions(ctx, jd)
		})
		if err == nil {
			positions = result.(*PlanetaryPositions)
		}
	}

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return nil, fmt.Errorf("failed to get planetary positions from all providers: %w", err)
	}

	m.cache.Set(ctx, cacheKey, positions, time.Hour)
	span.SetAttributes(attribute.Bool("success", true))
	span.AddEvent("Successfully retrieved planetary positions")

	return positions, nil
}

// GetSunPosition retrieves Sun position with caching and fallback.
func (m *Manager) GetSunPosition(ctx context.Context, jd JulianDay) (*SolarPosition, error) {
	ctx, span := m.observer.CreateSpan(ctx, "ephemeris.GetSunPosition")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("operation", "get_sun_position"),
	)

	cacheKey := fmt.Sprintf("sun_position_%f", jd)
	if cached, found := m.cache.Get(ctx, cacheKey); found {
		span.SetAttributes(attribute.Bool("cache_hit", true))
		span.AddEvent("Cache hit for sun position")
		if position, ok := cached.(*SolarPosition); ok {
			return position, nil
		}
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))

	result, err := m.tryProvider(ctx, m.primary, "primary", func(provider EphemerisProvider) (interface{}, error) {
		return provider.GetSunPosition(ctx, jd)
	})

	var position *SolarPosition
	if err == nil {
		position = result.(*SolarPosition)
	} else {
		span.AddEvent("Primary provider failed, trying fallback")

		result, err = m.tryProvider(ctx, m.fallback, "fallback", func(provider EphemerisProvider) (interface{}, error) {
			return provider.GetSunPosition(ctx, jd)
		})
		if err == nil {
			position = result.(*SolarPosition)
		}
	}

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return nil, fmt.Errorf("failed to get sun position from all providers: %w", err)
	}

	m.cache.Set(ctx, cacheKey, position, time.Hour)
	span.SetAttributes(attribute.Bool("success", true))
	span.AddEvent("Successfully retrieved sun position")

	return position, nil
}

// GetMoonPosition retrieves Moon position with caching and fallback.
func (m *Manager) GetMoonPosition(ctx context.Context, jd JulianDay) (*LunarPosition, error) {
	ctx, span := m.observer.CreateSpan(ctx, "ephemeris.GetMoonPosition")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("operation", "get_moon_position"),
	)

	cacheKey := fmt.Sprintf("moon_position_%f", jd)
	if cached, found := m.cache.Get(ctx, cacheKey); found {
		span.SetAttributes(attribute.Bool("cache_hit", true))
		span.AddEvent("Cache hit for moon position")
		if position, ok := cached.(*LunarPosition); ok {
			return position, nil
		}
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))

	result, err := m.tryProvider(ctx, m.primary, "primary", func(provider EphemerisProvider) (interface{}, error) {
		return provider.GetMoonPosition(ctx, jd)
	})

	var position *LunarPosition
	if err == nil {
		position = result.(*LunarPosition)
	} else {
		span.AddEvent("Primary provider failed, trying fallback")

		result, err = m.tryProvider(ctx, m.fallback, "fallback", func(provider EphemerisProvider) (interface{}, error) {
			return provider.GetMoonPosition(ctx, jd)
		})
		if err == nil {
			position = result.(*LunarPosition)
		}
	}

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return nil, fmt.Errorf("failed to get moon position from all providers: %w", err)
	}

	m.cache.Set(ctx, cacheKey, position, time.Hour)
	span.SetAttributes(attribute.Bool("success", true))
	span.AddEvent("Successfully retrieved moon position")

	return position, nil
}
