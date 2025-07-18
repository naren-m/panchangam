package ephemeris

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

// JulianDay represents a Julian day number
type JulianDay float64

// PlanetaryPositions holds the positions of all planets
type PlanetaryPositions struct {
	JulianDay JulianDay `json:"julian_day"`
	Sun       Position  `json:"sun"`
	Moon      Position  `json:"moon"`
	Mercury   Position  `json:"mercury"`
	Venus     Position  `json:"venus"`
	Mars      Position  `json:"mars"`
	Jupiter   Position  `json:"jupiter"`
	Saturn    Position  `json:"saturn"`
	Uranus    Position  `json:"uranus"`
	Neptune   Position  `json:"neptune"`
	Pluto     Position  `json:"pluto"`
}

// Position represents a celestial body's position
type Position struct {
	Longitude float64 `json:"longitude"` // Ecliptic longitude in degrees
	Latitude  float64 `json:"latitude"`  // Ecliptic latitude in degrees
	Distance  float64 `json:"distance"`  // Distance from Earth in AU
	Speed     float64 `json:"speed"`     // Speed in degrees per day
}

// SolarPosition represents the Sun's position
type SolarPosition struct {
	JulianDay           JulianDay `json:"julian_day"`
	Longitude           float64   `json:"longitude"`           // Ecliptic longitude in degrees
	RightAscension      float64   `json:"right_ascension"`     // Right ascension in degrees
	Declination         float64   `json:"declination"`         // Declination in degrees
	Distance            float64   `json:"distance"`            // Distance from Earth in AU
	EquationOfTime      float64   `json:"equation_of_time"`    // Equation of time in minutes
	MeanAnomaly         float64   `json:"mean_anomaly"`        // Mean anomaly in degrees
	TrueAnomaly         float64   `json:"true_anomaly"`        // True anomaly in degrees
	EccentricAnomaly    float64   `json:"eccentric_anomaly"`   // Eccentric anomaly in degrees
	MeanLongitude       float64   `json:"mean_longitude"`      // Mean longitude in degrees
	ApparentLongitude   float64   `json:"apparent_longitude"`  // Apparent longitude in degrees
}

// LunarPosition represents the Moon's position
type LunarPosition struct {
	JulianDay         JulianDay `json:"julian_day"`
	Longitude         float64   `json:"longitude"`          // Ecliptic longitude in degrees
	Latitude          float64   `json:"latitude"`           // Ecliptic latitude in degrees
	RightAscension    float64   `json:"right_ascension"`    // Right ascension in degrees
	Declination       float64   `json:"declination"`        // Declination in degrees
	Distance          float64   `json:"distance"`           // Distance from Earth in km
	Phase             float64   `json:"phase"`              // Phase (0-1, 0=new, 0.5=full)
	PhaseAngle        float64   `json:"phase_angle"`        // Phase angle in degrees
	Illumination      float64   `json:"illumination"`       // Illumination percentage
	AngularDiameter   float64   `json:"angular_diameter"`   // Angular diameter in arcseconds
	MeanAnomaly       float64   `json:"mean_anomaly"`       // Mean anomaly in degrees
	TrueAnomaly       float64   `json:"true_anomaly"`       // True anomaly in degrees
	ArgumentOfLatitude float64  `json:"argument_of_latitude"` // Argument of latitude in degrees
	MeanLongitude     float64   `json:"mean_longitude"`     // Mean longitude in degrees
	TrueLongitude     float64   `json:"true_longitude"`     // True longitude in degrees
}

// HealthStatus represents the health status of an ephemeris provider
type HealthStatus struct {
	Available     bool      `json:"available"`
	LastCheck     time.Time `json:"last_check"`
	DataStartJD   float64   `json:"data_start_jd"`
	DataEndJD     float64   `json:"data_end_jd"`
	ResponseTime  time.Duration `json:"response_time"`
	ErrorMessage  string    `json:"error_message,omitempty"`
	Version       string    `json:"version,omitempty"`
	Source        string    `json:"source,omitempty"`
}

// EphemerisProvider defines the interface for ephemeris data providers
type EphemerisProvider interface {
	// GetPlanetaryPositions returns positions of all planets for a given Julian day
	GetPlanetaryPositions(ctx context.Context, jd JulianDay) (*PlanetaryPositions, error)
	
	// GetSunPosition returns detailed Sun position for a given Julian day
	GetSunPosition(ctx context.Context, jd JulianDay) (*SolarPosition, error)
	
	// GetMoonPosition returns detailed Moon position for a given Julian day
	GetMoonPosition(ctx context.Context, jd JulianDay) (*LunarPosition, error)
	
	// IsAvailable checks if the ephemeris provider is available
	IsAvailable(ctx context.Context) bool
	
	// GetDataRange returns the valid Julian day range for this provider
	GetDataRange() (startJD, endJD JulianDay)
	
	// GetHealthStatus returns the current health status
	GetHealthStatus(ctx context.Context) (*HealthStatus, error)
	
	// GetProviderName returns the name of the provider
	GetProviderName() string
	
	// GetVersion returns the version of the ephemeris data
	GetVersion() string
	
	// Close closes the provider and releases resources
	Close() error
}

// Manager manages multiple ephemeris providers with fallback and caching
type Manager struct {
	primary       EphemerisProvider
	fallback      EphemerisProvider
	cache         Cache
	observer      observability.ObserverInterface
	healthChecker *HealthChecker
}

// NewManager creates a new ephemeris manager
func NewManager(primary, fallback EphemerisProvider, cache Cache) *Manager {
	manager := &Manager{
		primary:  primary,
		fallback: fallback,
		cache:    cache,
		observer: observability.Observer(),
	}
	
	// Initialize health checker
	manager.healthChecker = NewHealthChecker([]EphemerisProvider{primary, fallback})
	
	return manager
}

// GetPlanetaryPositions retrieves planetary positions with caching and fallback
func (m *Manager) GetPlanetaryPositions(ctx context.Context, jd JulianDay) (*PlanetaryPositions, error) {
	ctx, span := m.observer.CreateSpan(ctx, "ephemeris.GetPlanetaryPositions")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("operation", "get_planetary_positions"),
	)
	
	// Check cache first
	cacheKey := fmt.Sprintf("planetary_positions_%f", jd)
	if cached, found := m.cache.Get(ctx, cacheKey); found {
		span.SetAttributes(attribute.Bool("cache_hit", true))
		span.AddEvent("Cache hit for planetary positions")
		if positions, ok := cached.(*PlanetaryPositions); ok {
			return positions, nil
		}
	}
	
	span.SetAttributes(attribute.Bool("cache_hit", false))
	
	// Try primary provider
	result, err := m.tryProvider(ctx, m.primary, "primary", func(provider EphemerisProvider) (interface{}, error) {
		return provider.GetPlanetaryPositions(ctx, jd)
	})
	
	var positions *PlanetaryPositions
	if err == nil {
		positions = result.(*PlanetaryPositions)
	} else {
		span.AddEvent("Primary provider failed, trying fallback")
		
		// Try fallback provider
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
	
	// Cache the result
	m.cache.Set(ctx, cacheKey, positions, 1*time.Hour)
	span.SetAttributes(attribute.Bool("success", true))
	span.AddEvent("Successfully retrieved planetary positions")
	
	return positions, nil
}

// GetSunPosition retrieves Sun position with caching and fallback
func (m *Manager) GetSunPosition(ctx context.Context, jd JulianDay) (*SolarPosition, error) {
	ctx, span := m.observer.CreateSpan(ctx, "ephemeris.GetSunPosition")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("operation", "get_sun_position"),
	)
	
	// Check cache first
	cacheKey := fmt.Sprintf("sun_position_%f", jd)
	if cached, found := m.cache.Get(ctx, cacheKey); found {
		span.SetAttributes(attribute.Bool("cache_hit", true))
		span.AddEvent("Cache hit for sun position")
		if position, ok := cached.(*SolarPosition); ok {
			return position, nil
		}
	}
	
	span.SetAttributes(attribute.Bool("cache_hit", false))
	
	// Try primary provider
	result, err := m.tryProvider(ctx, m.primary, "primary", func(provider EphemerisProvider) (interface{}, error) {
		return provider.GetSunPosition(ctx, jd)
	})
	
	var position *SolarPosition
	if err == nil {
		position = result.(*SolarPosition)
	} else {
		span.AddEvent("Primary provider failed, trying fallback")
		
		// Try fallback provider
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
	
	// Cache the result
	m.cache.Set(ctx, cacheKey, position, 1*time.Hour)
	span.SetAttributes(attribute.Bool("success", true))
	span.AddEvent("Successfully retrieved sun position")
	
	return position, nil
}

// GetMoonPosition retrieves Moon position with caching and fallback
func (m *Manager) GetMoonPosition(ctx context.Context, jd JulianDay) (*LunarPosition, error) {
	ctx, span := m.observer.CreateSpan(ctx, "ephemeris.GetMoonPosition")
	defer span.End()
	
	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("operation", "get_moon_position"),
	)
	
	// Check cache first
	cacheKey := fmt.Sprintf("moon_position_%f", jd)
	if cached, found := m.cache.Get(ctx, cacheKey); found {
		span.SetAttributes(attribute.Bool("cache_hit", true))
		span.AddEvent("Cache hit for moon position")
		if position, ok := cached.(*LunarPosition); ok {
			return position, nil
		}
	}
	
	span.SetAttributes(attribute.Bool("cache_hit", false))
	
	// Try primary provider
	result, err := m.tryProvider(ctx, m.primary, "primary", func(provider EphemerisProvider) (interface{}, error) {
		return provider.GetMoonPosition(ctx, jd)
	})
	
	var position *LunarPosition
	if err == nil {
		position = result.(*LunarPosition)
	} else {
		span.AddEvent("Primary provider failed, trying fallback")
		
		// Try fallback provider
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
	
	// Cache the result
	m.cache.Set(ctx, cacheKey, position, 1*time.Hour)
	span.SetAttributes(attribute.Bool("success", true))
	span.AddEvent("Successfully retrieved moon position")
	
	return position, nil
}

// tryProvider attempts to get data from a provider with observability
func (m *Manager) tryProvider(ctx context.Context, provider EphemerisProvider, providerType string, operation func(EphemerisProvider) (interface{}, error)) (interface{}, error) {
	
	if provider == nil {
		return nil, fmt.Errorf("%s provider is nil", providerType)
	}
	
	ctx, span := m.observer.CreateSpan(ctx, fmt.Sprintf("ephemeris.try_%s_provider", providerType))
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

// GetHealthStatus returns the health status of all providers
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

// Close closes all providers and releases resources
func (m *Manager) Close() error {
	var errors []error
	
	if m.primary != nil {
		if err := m.primary.Close(); err != nil {
			errors = append(errors, fmt.Errorf("primary provider close error: %w", err))
		}
	}
	
	if m.fallback != nil {
		if err := m.fallback.Close(); err != nil {
			errors = append(errors, fmt.Errorf("fallback provider close error: %w", err))
		}
	}
	
	if m.cache != nil {
		if err := m.cache.Close(); err != nil {
			errors = append(errors, fmt.Errorf("cache close error: %w", err))
		}
	}
	
	if m.healthChecker != nil {
		m.healthChecker.Stop()
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("errors during close: %v", errors)
	}
	
	return nil
}

// TimeToJulianDay converts a time.Time to Julian day number
func TimeToJulianDay(t time.Time) JulianDay {
	// Convert to UTC
	utc := t.UTC()
	
	// Julian day calculation
	year := utc.Year()
	month := int(utc.Month())
	day := utc.Day()
	
	// Adjust for January and February
	if month <= 2 {
		year--
		month += 12
	}
	
	// Calculate Julian day number
	a := year / 100
	b := 2 - a + a/4
	
	jd := math.Floor(365.25*float64(year+4716)) +
		math.Floor(30.6001*float64(month+1)) +
		float64(day) + float64(b) - 1524.5
	
	// Add time of day (Julian day starts at noon)
	hour := float64(utc.Hour())
	minute := float64(utc.Minute())
	second := float64(utc.Second())
	
	jd += (hour-12.0)/24.0 + minute/1440.0 + second/86400.0
	
	return JulianDay(jd)
}

// JulianDayToTime converts a Julian day number to time.Time
func JulianDayToTime(jd JulianDay) time.Time {
	// Convert Julian day to calendar date
	z := math.Floor(float64(jd) + 0.5)
	f := float64(jd) + 0.5 - z
	
	var a float64
	if z < 2299161 {
		a = z
	} else {
		alpha := math.Floor((z - 1867216.25) / 36524.25)
		a = z + 1 + alpha - math.Floor(alpha/4)
	}
	
	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)
	
	day := int(b - d - math.Floor(30.6001*e) + f)
	var month int
	if e < 14 {
		month = int(e - 1)
	} else {
		month = int(e - 13)
	}
	
	var year int
	if month > 2 {
		year = int(c - 4716)
	} else {
		year = int(c - 4715)
	}
	
	// Calculate time of day
	dayFraction := f
	hours := dayFraction * 24
	hour := int(hours)
	minutes := (hours - float64(hour)) * 60
	minute := int(minutes)
	seconds := (minutes - float64(minute)) * 60
	second := int(seconds)
	nanosecond := int((seconds - float64(second)) * 1e9)
	
	return time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, time.UTC)
}