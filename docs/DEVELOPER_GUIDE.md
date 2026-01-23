# Developer Guide

## Table of Contents

1. [Getting Started](#getting-started)
2. [Architecture Overview](#architecture-overview)
3. [Core Components](#core-components)
4. [API Extensions](#api-extensions)
5. [Testing](#testing)
6. [Contributing](#contributing)

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+ (for UI)
- Docker (optional)

### Quick Start

```bash
# Clone repository
git clone https://github.com/naren-m/panchangam.git
cd panchangam

# Build backend
make build-backend

# Build frontend
make build-frontend

# Run tests
make test

# Start services
make run
```

## Architecture Overview

```
panchangam/
├── astronomy/           # Core astronomical calculations
│   ├── ephemeris/      # Ephemeris integration
│   └── validation/     # Validation framework
├── api/                # Plugin system and extensions
├── services/           # gRPC services
├── ui/                 # React frontend
├── cmd/                # Command-line tools
├── proto/              # Protocol buffers
└── docs/               # Documentation
```

## Core Components

### Astronomical Calculations

#### Tithi Calculator

```go
import "github.com/naren-m/panchangam/astronomy"

// Create calculator
tithiCalc := astronomy.NewTithiCalculator(ephemerisManager, observer)

// Calculate tithi
tithi, err := tithiCalc.Calculate(ctx, date, location)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Tithi: %s\n", tithi.Name)
fmt.Printf("Start: %s\n", tithi.StartTime)
fmt.Printf("End: %s\n", tithi.EndTime)
```

#### Nakshatra Calculator

```go
nakshatraCalc := astronomy.NewNakshatraCalculator(ephemerisManager, observer)

nakshatra, err := nakshatraCalc.Calculate(ctx, date, location)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Nakshatra: %s\n", nakshatra.Name)
fmt.Printf("Pada: %d\n", nakshatra.Pada)
fmt.Printf("Lord: %s\n", nakshatra.Lord)
```

### Ephemeris System

#### Using the Ephemeris Manager

```go
import "github.com/naren-m/panchangam/astronomy/ephemeris"

// Create providers
primary, _ := ephemeris.NewJPLProvider()
fallback, _ := ephemeris.NewSwissProvider()

// Create cache
cache := ephemeris.NewMemoryCache(100, 1*time.Hour)

// Create manager
manager := ephemeris.NewManager(primary, fallback, cache)

// Get planetary positions
jd := ephemeris.TimeToJulianDay(time.Now())
positions, err := manager.GetPlanetaryPositions(ctx, jd)
```

#### Interpolation

```go
// Create interpolator with configuration
config := ephemeris.InterpolationConfig{
    Method:    ephemeris.InterpolationCubicSpline,
    Order:     5,
    Tolerance: 0.0001,
}

interpolator := ephemeris.NewInterpolator(manager, config)

// Interpolate position
position, err := interpolator.InterpolatePlanetaryPosition(ctx, jd, "mars")
```

#### Retrograde Detection

```go
detector := ephemeris.NewRetrogradeDetector(manager)

// Check if retrograde
motion, err := detector.DetectRetrogradeMotion(ctx, jd, "mercury")

// Find next station
station, err := detector.FindPlanetaryStation(ctx, jd, "mercury", 120)

// Get comprehensive analysis
analysis, err := detector.AnalyzeMotion(ctx, jd, "venus")
```

### Validation Framework

```go
import "github.com/naren-m/panchangam/astronomy/validation"

// Create validator
validator := validation.NewValidator(
    tithiCalc, nakshatraCalc, yogaCalc,
    karanaCalc, varaCalc, sunriseCalc,
    ephemerisManager,
)

// Prepare reference data
refData := []validation.ReferenceData{
    {
        Source:        "Drik Panchang",
        Date:          time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
        Location:      location,
        TithiName:     "Panchami",
        NakshatraName: "Rohini",
        Sunrise:       time.Date(2024, 1, 15, 6, 45, 0, 0, time.UTC),
    },
}

// Run validation suite
suite := validator.ValidateAgainstDrikPanchang(ctx, refData)

// Generate report
report := suite.GenerateReport()
fmt.Println(report)
```

## API Extensions

### Plugin System

#### Creating a Custom Plugin

```go
package main

import (
    "context"
    "github.com/naren-m/panchangam/api"
)

type CustomPlugin struct {
    name    string
    version string
}

func (p *CustomPlugin) ProcessPanchangamData(
    ctx context.Context,
    data *api.PanchangamData,
) (*api.ProcessedData, error) {
    // Your custom processing logic
    processed := &api.ProcessedData{
        OriginalData: data,
        Modifications: make(map[string]interface{}),
    }

    // Add custom calculations
    processed.Modifications["custom_field"] = "custom_value"

    return processed, nil
}

func (p *CustomPlugin) GetName() string {
    return p.name
}

func (p *CustomPlugin) GetVersion() api.Version {
    return api.Version{Major: 1, Minor: 0, Patch: 0}
}
```

#### Registering a Plugin

```go
// Create plugin manager
manager := api.NewPluginManager()

// Register plugin
customPlugin := &CustomPlugin{
    name:    "custom-calculator",
    version: "1.0.0",
}

manager.RegisterPlugin("custom", customPlugin)

// Process data with plugins
result := manager.ProcessPanchangamData(ctx, data)
```

### API Versioning

#### Version Structure

```go
type Version struct {
    Major int  // Breaking changes
    Minor int  // New features, backward compatible
    Patch int  // Bug fixes
}

// Create version
v := api.Version{Major: 2, Minor: 1, Patch: 3}
fmt.Println(v.String()) // "2.1.3"

// Compare versions
if v.GreaterThan(otherVersion) {
    // Handle version difference
}
```

#### Versioned API Endpoints

```proto
// proto/panchangam/v2/panchangam.proto
package panchangam.v2;

service PanchangamV2 {
    rpc Get(GetPanchangamRequest) returns (GetPanchangamResponse);

    // Version 2 additions
    rpc GetExtended(GetExtendedRequest) returns (GetExtendedResponse);
    rpc GetInterpolated(GetInterpolatedRequest) returns (GetInterpolatedResponse);
}
```

#### Handling API Versions in Code

```go
// Version-aware handler
func HandleRequest(ctx context.Context, req *Request) (*Response, error) {
    switch req.ApiVersion {
    case "v1":
        return handleV1(ctx, req)
    case "v2":
        return handleV2(ctx, req)
    default:
        return nil, fmt.Errorf("unsupported API version: %s", req.ApiVersion)
    }
}
```

### Extension Points

#### Custom Festival Plugin

```go
type FestivalPlugin struct{}

func (p *FestivalPlugin) GetFestivalDates(
    year int,
    region string,
    locale string,
) ([]Festival, error) {
    // Calculate festival dates for the year
    festivals := []Festival{
        {
            Name:   "Custom Festival",
            Date:   time.Date(year, 4, 14, 0, 0, 0, 0, time.UTC),
            Type:   "regional",
            Region: region,
        },
    }

    return festivals, nil
}

func (p *FestivalPlugin) ValidateFestivalDate(
    date time.Time,
    festival string,
    region string,
) (bool, error) {
    // Validation logic
    return true, nil
}
```

#### Custom Calculation Method Plugin

```go
type CalculationMethodPlugin struct{}

func (p *CalculationMethodPlugin) ProcessPanchangamData(
    ctx context.Context,
    data *api.PanchangamData,
) (*api.ProcessedData, error) {
    // Apply custom calculation method
    // E.g., custom ayanamsa, custom tithi calculation

    return &api.ProcessedData{
        OriginalData:  data,
        Modifications: modifications,
    }, nil
}
```

### gRPC Service Implementation

#### Custom Service

```go
type CustomPanchangamServer struct {
    panchangam.UnimplementedPanchangamServer
    manager *api.PluginManager
}

func (s *CustomPanchangamServer) Get(
    ctx context.Context,
    req *panchangam.GetPanchangamRequest,
) (*panchangam.GetPanchangamResponse, error) {
    // Validate request
    if err := validateRequest(req); err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
    }

    // Process with plugins
    data := convertRequestToData(req)
    processed := s.manager.ProcessPanchangamData(ctx, data)

    // Build response
    response := buildResponse(processed)

    return response, nil
}
```

## Testing

### Unit Testing

```go
func TestTithiCalculation(t *testing.T) {
    // Setup
    manager := createTestManager(t)
    calc := astronomy.NewTithiCalculator(manager, nil)

    // Test data
    date := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
    location := astronomy.Location{
        Latitude:  13.0827,
        Longitude: 80.2707,
        Name:      "Chennai",
    }

    // Execute
    tithi, err := calc.Calculate(context.Background(), date, location)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, tithi)
    assert.NotEmpty(t, tithi.Name)
}
```

### Integration Testing

```go
func TestPanchangamServiceIntegration(t *testing.T) {
    // Start test server
    server := startTestServer(t)
    defer server.Stop()

    // Create client
    conn, err := grpc.Dial(server.Address(), grpc.WithInsecure())
    require.NoError(t, err)
    defer conn.Close()

    client := panchangam.NewPanchangamClient(conn)

    // Make request
    req := &panchangam.GetPanchangamRequest{
        Date:      "2024-01-15",
        Latitude:  13.0827,
        Longitude: 80.2707,
        Timezone:  "Asia/Kolkata",
    }

    resp, err := client.Get(context.Background(), req)

    // Verify
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.NotNil(t, resp.Data)
}
```

### Validation Testing

```go
func TestValidationFramework(t *testing.T) {
    validator := createTestValidator(t)

    // Reference data from Drik Panchang
    refData := validation.ReferenceData{
        Source:        "Drik Panchang",
        Date:          time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
        Location:      testLocation,
        TithiName:     "Shashthi",
        NakshatraName: "Punarvasu",
    }

    // Validate
    result := validator.ValidateTithi(context.Background(), refData, 5.0)

    // Check result
    assert.True(t, result.Passed, "Tithi validation should pass")
    assert.LessOrEqual(t, result.Error, 5.0, "Error within tolerance")
}
```

### Benchmark Testing

```go
func BenchmarkTithiCalculation(b *testing.B) {
    manager := createTestManager(b)
    calc := astronomy.NewTithiCalculator(manager, nil)

    date := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
    location := testLocation
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := calc.Calculate(ctx, date, location)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Best Practices

### Error Handling

```go
// Use wrapped errors for context
if err != nil {
    return nil, fmt.Errorf("failed to calculate tithi: %w", err)
}

// Check for specific error types
if errors.Is(err, ephemeris.ErrDataNotAvailable) {
    // Handle specific error
}
```

### Logging and Observability

```go
// Use structured logging
logger.WithFields(log.Fields{
    "date":     date,
    "location": location.Name,
    "tithi":    tithi.Name,
}).Info("Tithi calculated successfully")

// Add OpenTelemetry spans
ctx, span := tracer.Start(ctx, "calculate_tithi")
defer span.End()

span.SetAttributes(
    attribute.String("location", location.Name),
    attribute.Float64("latitude", location.Latitude),
)
```

### Context Usage

```go
// Always pass context
func Calculate(ctx context.Context, date time.Time) error {
    // Check context cancellation
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // Pass context to sub-functions
    result, err := subFunction(ctx, date)

    return err
}
```

### Configuration Management

```go
// Use configuration structs
type Config struct {
    EphemerisProvider string `json:"ephemeris_provider"`
    CacheSize         int    `json:"cache_size"`
    Region            string `json:"region"`
    CalculationMethod string `json:"calculation_method"`
}

// Load from file
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

## Contributing

### Code Style

- Follow Go conventions
- Use `gofmt` for formatting
- Run `golint` before committing
- Keep functions small and focused
- Write descriptive comments

### Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Write tests for new features
4. Ensure all tests pass
5. Update documentation
6. Submit pull request

### Documentation

- Document public APIs
- Include examples
- Update CHANGELOG.md
- Add inline comments for complex logic

### Testing Requirements

- Unit test coverage > 90%
- Integration tests for new features
- Performance benchmarks for critical paths
- Validation tests against reference data

## Resources

### Internal Documentation

- [Astronomical Algorithms](./ASTRONOMICAL_ALGORITHMS.md)
- [Regional Variations](./REGIONAL_VARIATIONS.md)
- [API Reference](./API_REFERENCE.md)
- [Feature Documentation](./FEATURES.md)

### External Resources

- [Go Documentation](https://golang.org/doc/)
- [gRPC Guide](https://grpc.io/docs/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [OpenTelemetry](https://opentelemetry.io/docs/)

---

*Last Updated: 2025-11-18*
*Maintainer: Panchangam Development Team*
