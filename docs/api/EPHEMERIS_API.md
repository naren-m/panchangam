# Ephemeris Integration API Documentation

## Overview

The Ephemeris API provides precise astronomical calculations for planetary positions, lunar phases, and solar phenomena. This documentation covers the complete interface, data structures, and integration patterns for astronomical calculations in the Panchangam system.

## Core Architecture

### EphemerisProvider Interface

The foundation of all astronomical calculations:

```go
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
    
    // GetHealth returns the health status of the provider
    GetHealth(ctx context.Context) (*HealthStatus, error)
    
    // GetProviderInfo returns information about the provider
    GetProviderInfo() ProviderInfo
}
```

### Data Structures

#### Planetary Positions
```go
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
```

#### Position Structure
```go
type Position struct {
    Longitude float64 `json:"longitude"` // Ecliptic longitude in degrees
    Latitude  float64 `json:"latitude"`  // Ecliptic latitude in degrees
    Distance  float64 `json:"distance"`  // Distance from Earth in AU
    Speed     float64 `json:"speed"`     // Speed in degrees per day
}
```

## Provider Implementations

### Swiss Ephemeris Provider

High-precision astronomical calculations using Swiss Ephemeris library.

#### Features
- **Precision**: ±0.001 arcsecond accuracy
- **Time Range**: 30,000-year span (historical and future)
- **Data Size**: 99MB compressed from 2.8GB NASA JPL data
- **Performance**: ~10-50ms per calculation

#### Configuration
```go
type SwissConfig struct {
    DataPath     string        `json:"data_path"`
    CacheSize    int           `json:"cache_size"`
    Timeout      time.Duration `json:"timeout"`
    EnableCache  bool          `json:"enable_cache"`
    LogLevel     string        `json:"log_level"`
}
```

#### Usage Example
```go
// Initialize Swiss Ephemeris provider
swissProvider, err := ephemeris.NewSwissProvider(ephemeris.SwissConfig{
    DataPath:    "/path/to/ephemeris/data",
    CacheSize:   1000,
    Timeout:     30 * time.Second,
    EnableCache: true,
})

// Get planetary positions
jd := ephemeris.TimeToJulianDay(time.Now())
positions, err := swissProvider.GetPlanetaryPositions(ctx, jd)
```

### JPL Ephemeris Provider

NASA JPL planetary ephemeris integration for cross-validation.

#### Features
- **Source**: NASA Jet Propulsion Laboratory
- **Accuracy**: Reference standard for astronomical calculations
- **Integration**: Fallback and validation provider
- **Data Format**: Binary ephemeris files

#### Configuration
```go
type JPLConfig struct {
    EphemerisFile string        `json:"ephemeris_file"`
    LeapSeconds   string        `json:"leap_seconds"`
    CacheSize     int           `json:"cache_size"`
    Timeout       time.Duration `json:"timeout"`
}
```

## Manager Interface

### Ephemeris Manager

Orchestrates multiple providers with fallback and caching.

```go
type Manager struct {
    primaryProvider   EphemerisProvider
    fallbackProvider  EphemerisProvider
    cache            CacheInterface
    observer         observability.ObserverInterface
    config           ManagerConfig
}
```

#### Manager Configuration
```go
type ManagerConfig struct {
    PrimaryProvider    string        `json:"primary_provider"`
    FallbackProvider   string        `json:"fallback_provider"`
    EnableCache        bool          `json:"enable_cache"`
    CacheTTL          time.Duration `json:"cache_ttl"`
    HealthCheckInterval time.Duration `json:"health_check_interval"`
    MaxRetries        int           `json:"max_retries"`
    RetryDelay        time.Duration `json:"retry_delay"`
}
```

#### Initialization
```go
// Create manager with providers
manager := ephemeris.NewManager(ephemeris.ManagerConfig{
    PrimaryProvider:     "swiss",
    FallbackProvider:    "jpl",
    EnableCache:         true,
    CacheTTL:           1 * time.Hour,
    HealthCheckInterval: 5 * time.Minute,
    MaxRetries:         3,
    RetryDelay:         1 * time.Second,
})

// Add providers
manager.RegisterProvider("swiss", swissProvider)
manager.RegisterProvider("jpl", jplProvider)

// Initialize
err := manager.Initialize(ctx)
```

## API Methods

### GetPlanetaryPositions

Retrieve positions for all major planets.

```go
func (m *Manager) GetPlanetaryPositions(ctx context.Context, jd JulianDay) (*PlanetaryPositions, error)
```

**Parameters:**
- `ctx`: Context for cancellation and tracing
- `jd`: Julian Day number for calculation

**Returns:**
- `*PlanetaryPositions`: Complete planetary position data
- `error`: Any calculation or provider errors

**Example:**
```go
jd := ephemeris.TimeToJulianDay(time.Date(2025, 7, 20, 12, 0, 0, 0, time.UTC))
positions, err := manager.GetPlanetaryPositions(ctx, jd)
if err != nil {
    return fmt.Errorf("failed to get positions: %w", err)
}

fmt.Printf("Sun longitude: %.6f°\n", positions.Sun.Longitude)
fmt.Printf("Moon longitude: %.6f°\n", positions.Moon.Longitude)
```

### GetSunPosition

Detailed solar position and timing information.

```go
func (m *Manager) GetSunPosition(ctx context.Context, jd JulianDay) (*SolarPosition, error)
```

**Solar Position Structure:**
```go
type SolarPosition struct {
    JulianDay           JulianDay `json:"julian_day"`
    Longitude           float64   `json:"longitude"`           // Ecliptic longitude
    RightAscension      float64   `json:"right_ascension"`     // Right ascension
    Declination         float64   `json:"declination"`         // Declination
    Distance            float64   `json:"distance"`            // Distance in AU
    EquationOfTime      float64   `json:"equation_of_time"`    // Minutes
    MeanAnomaly         float64   `json:"mean_anomaly"`        // Degrees
    TrueAnomaly         float64   `json:"true_anomaly"`        // Degrees
    EccentricAnomaly    float64   `json:"eccentric_anomaly"`   // Degrees
    MeanLongitude       float64   `json:"mean_longitude"`      // Degrees
    ApparentLongitude   float64   `json:"apparent_longitude"`  // Degrees
}
```

### GetMoonPosition

Comprehensive lunar position and phase information.

```go
func (m *Manager) GetMoonPosition(ctx context.Context, jd JulianDay) (*LunarPosition, error)
```

**Lunar Position Structure:**
```go
type LunarPosition struct {
    JulianDay         JulianDay `json:"julian_day"`
    Longitude         float64   `json:"longitude"`          // Ecliptic longitude
    Latitude          float64   `json:"latitude"`           // Ecliptic latitude
    RightAscension    float64   `json:"right_ascension"`    // Right ascension
    Declination       float64   `json:"declination"`        // Declination
    Distance          float64   `json:"distance"`           // Distance in km
    Phase             float64   `json:"phase"`              // 0-1 (0=new, 0.5=full)
    PhaseAngle        float64   `json:"phase_angle"`        // Phase angle in degrees
    Illumination      float64   `json:"illumination"`       // Percentage
    AngularDiameter   float64   `json:"angular_diameter"`   // Arcseconds
    MeanAnomaly       float64   `json:"mean_anomaly"`       // Degrees
    TrueAnomaly       float64   `json:"true_anomaly"`       // Degrees
    ArgumentOfLatitude float64  `json:"argument_of_latitude"` // Degrees
    MeanLongitude     float64   `json:"mean_longitude"`     // Degrees
    TrueLongitude     float64   `json:"true_longitude"`     // Degrees
}
```

## Utility Functions

### Time Conversion

#### TimeToJulianDay
```go
func TimeToJulianDay(t time.Time) JulianDay
```

Converts Go time.Time to Julian Day number.

**Example:**
```go
// Current time to Julian Day
jd := ephemeris.TimeToJulianDay(time.Now())

// Specific date
date := time.Date(2025, 7, 20, 12, 0, 0, 0, time.UTC)
jd := ephemeris.TimeToJulianDay(date)
```

#### JulianDayToTime
```go
func JulianDayToTime(jd JulianDay) time.Time
```

Converts Julian Day number back to time.Time.

**Example:**
```go
jd := JulianDay(2460509.0) // J2000.0 epoch
t := ephemeris.JulianDayToTime(jd)
```

## Error Handling

### Error Types

```go
var (
    ErrProviderUnavailable = errors.New("ephemeris provider unavailable")
    ErrInvalidJulianDay   = errors.New("invalid Julian Day value")
    ErrDataOutOfRange     = errors.New("date outside ephemeris data range")
    ErrCalculationFailed  = errors.New("astronomical calculation failed")
    ErrProviderTimeout    = errors.New("provider request timeout")
)
```

### Error Handling Pattern

```go
positions, err := manager.GetPlanetaryPositions(ctx, jd)
if err != nil {
    switch {
    case errors.Is(err, ephemeris.ErrProviderUnavailable):
        // Handle provider unavailability
        log.Warn("Primary provider unavailable, using fallback")
    case errors.Is(err, ephemeris.ErrDataOutOfRange):
        // Handle out-of-range dates
        return fmt.Errorf("date %v outside ephemeris range", date)
    case errors.Is(err, ephemeris.ErrCalculationFailed):
        // Handle calculation errors
        log.Error("Calculation failed", "error", err)
    default:
        // Handle other errors
        return fmt.Errorf("ephemeris error: %w", err)
    }
}
```

## Health Monitoring

### Health Status

```go
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
```

### Health Check

```go
// Check provider health
health, err := manager.GetHealth(ctx)
if err != nil {
    log.Error("Health check failed", "error", err)
}

if !health.Available {
    log.Warn("Provider unavailable", "message", health.ErrorMessage)
}
```

## Caching Strategy

### Cache Interface

```go
type CacheInterface interface {
    Get(ctx context.Context, key string) (interface{}, bool)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Clear(ctx context.Context) error
}
```

### Cache Keys

```go
// Planetary positions: "positions:{jd}"
// Solar position: "sun:{jd}"
// Lunar position: "moon:{jd}"
```

### Cache Configuration

```go
type CacheConfig struct {
    TTL           time.Duration `json:"ttl"`
    MaxSize       int           `json:"max_size"`
    EvictionPolicy string       `json:"eviction_policy"`
    EnableMetrics bool          `json:"enable_metrics"`
}
```

## Performance Optimization

### Batch Operations

```go
// Batch planetary positions for multiple dates
func (m *Manager) GetPlanetaryPositionsBatch(ctx context.Context, julianDays []JulianDay) ([]*PlanetaryPositions, error)
```

### Pre-computation

```go
// Pre-compute positions for a date range
func (m *Manager) PrecomputeRange(ctx context.Context, startJD, endJD JulianDay, interval time.Duration) error
```

### Memory Management

```go
// Memory optimization for long-running services
type MemoryConfig struct {
    MaxCacheSize    int           `json:"max_cache_size"`
    GCInterval      time.Duration `json:"gc_interval"`
    MemoryThreshold int64         `json:"memory_threshold"`
}
```

## Observability Integration

### OpenTelemetry Tracing

```go
// Automatic span creation for all operations
ctx, span := tracer.Start(ctx, "ephemeris.GetPlanetaryPositions")
defer span.End()

// Comprehensive attribute logging
span.SetAttributes(
    attribute.Float64("julian_day", float64(jd)),
    attribute.String("provider", "swiss"),
    attribute.Float64("response_time_ms", responseTime.Milliseconds()),
)
```

### Metrics

```go
// Performance metrics
- ephemeris_requests_total
- ephemeris_request_duration_seconds
- ephemeris_cache_hits_total
- ephemeris_cache_misses_total
- ephemeris_provider_health
```

## Integration Examples

### Tithi Calculation Integration

```go
func (tc *TithiCalculator) calculatePositions(ctx context.Context, date time.Time) error {
    jd := ephemeris.TimeToJulianDay(date)
    
    positions, err := tc.ephemerisManager.GetPlanetaryPositions(ctx, jd)
    if err != nil {
        return fmt.Errorf("failed to get positions: %w", err)
    }
    
    // Use positions for tithi calculation
    sunLong := positions.Sun.Longitude
    moonLong := positions.Moon.Longitude
    
    return tc.calculateTithi(sunLong, moonLong)
}
```

### Sunrise Calculation Integration

```go
func (sc *SunriseCalculator) getSolarPosition(ctx context.Context, date time.Time) (*SolarPosition, error) {
    jd := ephemeris.TimeToJulianDay(date)
    
    return sc.ephemerisManager.GetSunPosition(ctx, jd)
}
```

## Configuration Examples

### Production Configuration

```yaml
ephemeris:
  manager:
    primary_provider: "swiss"
    fallback_provider: "jpl"
    enable_cache: true
    cache_ttl: "1h"
    health_check_interval: "5m"
    max_retries: 3
    retry_delay: "1s"
  
  swiss:
    data_path: "/data/ephemeris/swiss"
    cache_size: 1000
    timeout: "30s"
    enable_cache: true
    log_level: "info"
  
  jpl:
    ephemeris_file: "/data/ephemeris/jpl/de440.bsp"
    leap_seconds: "/data/ephemeris/jpl/leap_seconds.dat"
    cache_size: 500
    timeout: "60s"
  
  cache:
    ttl: "1h"
    max_size: 10000
    eviction_policy: "lru"
    enable_metrics: true
```

### Development Configuration

```yaml
ephemeris:
  manager:
    primary_provider: "swiss"
    enable_cache: false
    health_check_interval: "1m"
    max_retries: 1
  
  swiss:
    data_path: "./testdata/ephemeris"
    cache_size: 100
    timeout: "10s"
    log_level: "debug"
```

## Testing

### Unit Tests

```go
func TestEphemerisManager_GetPlanetaryPositions(t *testing.T) {
    manager := setupTestManager(t)
    
    jd := ephemeris.TimeToJulianDay(time.Date(2025, 7, 20, 12, 0, 0, 0, time.UTC))
    
    positions, err := manager.GetPlanetaryPositions(context.Background(), jd)
    require.NoError(t, err)
    require.NotNil(t, positions)
    
    // Validate position ranges
    assert.True(t, positions.Sun.Longitude >= 0 && positions.Sun.Longitude < 360)
    assert.True(t, positions.Moon.Longitude >= 0 && positions.Moon.Longitude < 360)
}
```

### Integration Tests

```go
func TestEphemerisIntegration(t *testing.T) {
    // Test with real ephemeris data
    manager := ephemeris.NewManager(ephemeris.ManagerConfig{
        PrimaryProvider: "swiss",
    })
    
    // Test historical date
    historicalDate := time.Date(1900, 1, 1, 12, 0, 0, 0, time.UTC)
    jd := ephemeris.TimeToJulianDay(historicalDate)
    
    positions, err := manager.GetPlanetaryPositions(context.Background(), jd)
    require.NoError(t, err)
    
    // Verify against known values
    assertPositionAccuracy(t, positions)
}
```

## References

### External Resources
- [Swiss Ephemeris Documentation](https://www.astro.com/swisseph/)
- [JPL Planetary Ephemeris](https://ssd.jpl.nasa.gov/horizons/)
- [Astronomical Algorithms by Jean Meeus](https://www.willbell.com/math/MC1.HTM)

### Standards
- IAU 2000 Resolutions
- IERS Conventions 2010
- JPL DE440 Ephemeris Standard

---

*Last updated: July 2025*
*API Version: 1.0.0*
*Maintainer: Panchangam Development Team*