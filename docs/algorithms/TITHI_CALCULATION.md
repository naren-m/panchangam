# Tithi Calculation Algorithm Documentation

## Overview

Tithi represents the lunar day in the Hindu calendar system, fundamental to Panchangam calculations. This document details the mathematical algorithms, implementation specifics, and validation methods used in the Panchangam project.

## Mathematical Foundation

### Core Formula
```
Tithi Number = floor((Moon_Longitude - Sun_Longitude) / 12°) + 1
```

### Key Principles
- **12-degree segments**: Each tithi spans exactly 12 degrees of lunar-solar separation
- **30 tithis per lunar month**: 15 for Shukla Paksha (waxing), 15 for Krishna Paksha (waning)
- **Variable duration**: 19h 59m to 26h 47m due to elliptical orbits
- **Continuous calculation**: Real-time precision using ephemeris data

## Implementation Architecture

### Core Components

#### TithiCalculator Structure
```go
type TithiCalculator struct {
    ephemerisManager *ephemeris.Manager
    observer         observability.ObserverInterface
}
```

#### TithiInfo Data Structure
```go
type TithiInfo struct {
    Number      int       // 1-15 (Shukla), 16-30 (Krishna)
    Name        string    // Sanskrit name
    Type        TithiType // Categorization
    StartTime   time.Time // Precise start time
    EndTime     time.Time // Precise end time
    Duration    float64   // Duration in hours
    IsShukla    bool      // Paksha identification
    MoonSunDiff float64   // Angular separation
}
```

### Calculation Process

#### 1. Astronomical Position Retrieval
```go
// Convert date to Julian Day for ephemeris lookup
jd := ephemeris.TimeToJulianDay(noonDate)

// Get precise planetary positions
positions, err := tc.ephemerisManager.GetPlanetaryPositions(ctx, jd)
sunLong := positions.Sun.Longitude
moonLong := positions.Moon.Longitude
```

#### 2. Angular Difference Calculation
```go
// Calculate lunar-solar separation
moonSunDiff := moonLong - sunLong

// Normalize to 0-360 degree range
if moonSunDiff < 0 {
    moonSunDiff += 360
}
if moonSunDiff >= 360 {
    moonSunDiff -= 360
}
```

#### 3. Tithi Number Determination
```go
// Convert angular difference to tithi number
tithiFloat := moonSunDiff / 12.0
tithiNumber := int(tithiFloat) + 1

// Validate range (1-30)
if tithiNumber > 30 { tithiNumber = 30 }
if tithiNumber < 1 { tithiNumber = 1 }
```

#### 4. Paksha Classification
```go
// Determine lunar phase
isShukla := tithiNumber <= 15  // Waxing moon
```

## Tithi Classification System

### Five Categories (Pancha Guna)

| Type | Numbers | Sanskrit | Quality | Usage |
|------|---------|----------|---------|-------|
| Nanda | 1,6,11 | नन्दा | Joyful | Celebrations, new beginnings |
| Bhadra | 2,7,12 | भद्रा | Auspicious | All beneficial activities |
| Jaya | 3,8,13 | जया | Victorious | Success-oriented tasks |
| Rikta | 4,9,14 | रिक्ता | Empty | Avoid new ventures |
| Purna | 5,10,15 | पूर्णा | Complete | Task completion, fulfillment |

### Implementation
```go
func getTithiType(tithiNumber int) TithiType {
    normalizedTithi := tithiNumber
    if normalizedTithi > 15 {
        normalizedTithi = normalizedTithi - 15
    }
    
    switch normalizedTithi {
    case 1, 6, 11: return TithiTypeNanda
    case 2, 7, 12: return TithiTypeBhadra
    case 3, 8, 13: return TithiTypeJaya
    case 4, 9, 14: return TithiTypeRikta
    case 5, 10, 15: return TithiTypePurna
    }
}
```

## Time Calculation Methods

### Duration Estimation
```go
// Average tithi duration: 24.79 hours (lunar month / 30)
avgTithiDuration := time.Duration(24.79 * float64(time.Hour))

// Calculate progress within current tithi
tithiProgress := tithiFloat - math.Floor(tithiFloat)
timeIntoTithi := time.Duration(tithiProgress * float64(avgTithiDuration))

// Estimate start and end times
startTime = referenceTime.Add(-timeIntoTithi)
endTime = startTime.Add(avgTithiDuration)
```

### Precision Considerations
- **Ephemeris accuracy**: ±0.001 arcsecond precision from Swiss Ephemeris
- **Timing precision**: ±30 minutes typical accuracy for tithi transitions
- **Iteration refinement**: Multiple ephemeris queries for exact transition moments

## Validation Framework

### Input Validation
```go
func ValidateTithiCalculation(tithi *TithiInfo) error {
    // Null check
    if tithi == nil {
        return fmt.Errorf("tithi cannot be nil")
    }
    
    // Range validation
    if tithi.Number < 1 || tithi.Number > 30 {
        return fmt.Errorf("invalid tithi number: %d", tithi.Number)
    }
    
    // Angular difference validation
    if tithi.MoonSunDiff < 0 || tithi.MoonSunDiff >= 360 {
        return fmt.Errorf("invalid moon-sun difference: %f", tithi.MoonSunDiff)
    }
    
    // Duration validation
    if tithi.Duration <= 0 || tithi.Duration > 48 {
        return fmt.Errorf("invalid duration: %f hours", tithi.Duration)
    }
    
    // Time sequence validation
    if tithi.EndTime.Before(tithi.StartTime) {
        return fmt.Errorf("end time before start time")
    }
    
    return nil
}
```

### Cross-Reference Validation
- **Historical panchangams**: Compare against established sources
- **Multiple ephemeris**: Swiss vs JPL validation
- **Regional variations**: Account for calculation method differences

## Error Handling

### Common Edge Cases
1. **Leap seconds**: UTC time adjustments
2. **Timezone transitions**: DST handling
3. **Ephemeris boundaries**: Data availability limits
4. **Precision limits**: Floating-point arithmetic considerations

### Observability Integration
```go
// OpenTelemetry tracing for debugging
ctx, span := tc.observer.CreateSpan(ctx, "TithiCalculator.GetTithiForDate")
defer span.End()

// Comprehensive attribute logging
span.SetAttributes(
    attribute.String("date", date.Format("2006-01-02")),
    attribute.Float64("sun_longitude", sunLong),
    attribute.Float64("moon_longitude", moonLong),
    attribute.Int("tithi_number", tithi.Number),
    attribute.String("tithi_type", string(tithi.Type)),
)
```

## Performance Optimization

### Caching Strategy
- **Position caching**: Cache planetary positions for repeated queries
- **Calculation memoization**: Store results for identical inputs
- **Batch processing**: Multiple date calculations in single ephemeris call

### Computational Complexity
- **Time complexity**: O(1) for single date calculation
- **Space complexity**: O(1) for data structures
- **Ephemeris calls**: Most expensive operation (~10-50ms)

## Integration Points

### API Integration
```go
// Primary interface
func (tc *TithiCalculator) GetTithiForDate(ctx context.Context, date time.Time) (*TithiInfo, error)

// Direct longitude input
func (tc *TithiCalculator) GetTithiFromLongitudes(ctx context.Context, sunLong, moonLong float64, date time.Time) (*TithiInfo, error)
```

### Dependencies
- **Ephemeris Manager**: Swiss/JPL ephemeris data
- **Observability**: OpenTelemetry tracing
- **Time Handling**: Go time package with timezone support

## Regional Variations

### Calculation Methods
- **Drik Ganita**: Modern observational astronomy (default)
- **Vakya Ganita**: Traditional verse-based calculations
- **Regional preferences**: South India vs North India methods

### Future Enhancements
- **Precise transition calculations**: Iterative refinement for exact timing
- **Multiple ayanamsa support**: Different zodiacal reference systems
- **Historical date validation**: Extended ephemeris range support

## Testing Strategy

### Unit Tests
- **Boundary conditions**: New moon, full moon edge cases
- **Mathematical accuracy**: Verify formula implementation
- **Type classification**: Ensure correct categorization

### Integration Tests
- **Ephemeris integration**: End-to-end calculation validation
- **Historical verification**: Cross-check with known panchangam dates
- **Performance benchmarks**: Response time measurements

### Validation Sources
- **Drik Panchang**: Online verification
- **Traditional panchangams**: Regional publication comparison
- **Astronomical software**: Swiss Ephemeris validation

## References

### Mathematical Sources
- Surya Siddhanta (ancient astronomical text)
- Siddhanta Shiromani by Bhaskaracharya
- Modern astronomical algorithms (Meeus)

### Implementation References
- Swiss Ephemeris documentation
- JPL planetary ephemeris
- Hindu calendar research papers

---

*Last updated: July 2025*
*Maintainer: Panchangam Development Team*