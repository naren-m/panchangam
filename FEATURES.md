# Panchangam Feature Specification

## Overview
The Panchangam system is a comprehensive Hindu calendar calculation engine that provides the five essential Panchangam elements (Panchang) along with astronomical calculations and timing data.

## Core Features

### 1. Panchangam Elements (पञ्चाङ्ग)

#### 1.1 Tithi Calculation (तिथि)
**Feature ID**: `TITHI_001`
**Description**: Lunar day calculation based on Moon-Sun longitude difference
**Status**: ✅ Implemented
**Coverage**: Unit + Integration + Performance Tests

**Components**:
- `TithiCalculator` in `astronomy/tithi.go`
- 30 Tithis per lunar month (1-15 Shukla Paksha, 16-30 Krishna Paksha)
- 5 Tithi types: Nanda, Bhadra, Jaya, Rikta, Purna
- Precise timing calculations with start/end times
- Duration calculation in hours

**API Endpoints**:
- `GetTithiForDate(ctx, date)` - Primary calculation method
- `GetTithiFromLongitudes(ctx, sunLong, moonLong, date)` - Direct calculation
- `ValidateTithiCalculation(tithi)` - Validation utility

**Test Coverage**:
- Unit tests: ✅ Complete
- Integration tests: ✅ Complete
- Performance tests: ✅ <100ms target met

#### 1.2 Nakshatra Calculation (नक्षत्र)
**Feature ID**: `NAKSHATRA_001`
**Description**: Lunar mansion calculation with 27 divisions of the zodiac
**Status**: ✅ Implemented
**Coverage**: Unit + Integration + Performance Tests

**Components**:
- `NakshatraCalculator` in `astronomy/nakshatra.go`
- 27 Nakshatras with Sanskrit names, deities, and planetary lords
- Pada calculation (4 quarters per Nakshatra)
- Symbol and characteristic information

**API Endpoints**:
- `GetNakshatraForDate(ctx, date)` - Primary calculation method
- `GetNakshatraFromLongitude(ctx, moonLong, date)` - Direct calculation
- `GetNakshatraInfo(number)` - Static information lookup

**Test Coverage**:
- Unit tests: ✅ Complete
- Integration tests: ✅ Complete
- Performance tests: ✅ <100ms target met

#### 1.3 Yoga Calculation (योग)
**Feature ID**: `YOGA_001`
**Description**: Auspicious combinations based on Sun+Moon longitude sum
**Status**: ✅ Implemented
**Coverage**: Unit + Integration + Performance Tests

**Components**:
- `YogaCalculator` in `astronomy/yoga.go`
- 27 Yogas with quality categorization (Auspicious/Inauspicious/Mixed)
- Sanskrit names and descriptions
- Quality assessment for planning activities

**API Endpoints**:
- `GetYogaForDate(ctx, date)` - Primary calculation method
- `GetYogaFromLongitudes(ctx, sunLong, moonLong, date)` - Direct calculation
- `IsAuspiciousYoga(yoga)` - Quality assessment utility

**Test Coverage**:
- Unit tests: ✅ Complete
- Integration tests: ✅ Complete
- Performance tests: ✅ <100ms target met

#### 1.4 Karana Calculation (करण)
**Feature ID**: `KARANA_001`
**Description**: Half-Tithi divisions with 11-Karana cycle
**Status**: ✅ Implemented (Refactored for code reuse)
**Coverage**: Unit + Integration + Performance Tests

**Components**:
- `KaranaCalculator` in `astronomy/karana.go`
- 11 Karanas: 7 Movable (Chara) + 4 Fixed (Sthira)
- Vishti (Bhadra) special handling for inauspicious periods
- Reuses TithiCalculator for astronomical calculations
- Activity recommendations based on Karana type

**API Endpoints**:
- `GetKaranaForDate(ctx, date)` - Primary calculation method
- `GetKaranaFromLongitudes(ctx, sunLong, moonLong, date)` - Direct calculation
- `IsAuspiciousKarana(karana)` - Auspiciousness assessment
- `GetKaranaRecommendations(karana)` - Activity suggestions

**Test Coverage**:
- Unit tests: ✅ Complete
- Integration tests: ✅ Complete
- Performance tests: ✅ <100ms target met

#### 1.5 Vara Calculation (वार)
**Feature ID**: `VARA_001`
**Description**: Weekday calculation with hora system and planetary lords
**Status**: ✅ Implemented
**Coverage**: Unit + Integration + Performance Tests

**Components**:
- `VaraCalculator` in `astronomy/vara.go`
- 7 Varas (weekdays) with Sanskrit names
- Planetary lord associations
- 24-hour Hora system with hourly planetary rulers
- Sunrise-based day assignment

**API Endpoints**:
- `GetVaraForDate(ctx, date)` - Primary calculation method
- `GetVaraFromGregorianDay(ctx, weekday, sunrise, nextSunrise, date)` - Direct calculation
- `GetHoraForTime(ctx, time, sunrise, sunset)` - Hora calculation

**Test Coverage**:
- Unit tests: ✅ Complete
- Integration tests: ✅ Complete
- Performance tests: ✅ <100ms target met

### 2. Astronomical Calculations

#### 2.1 Sunrise/Sunset Calculations
**Feature ID**: `ASTRONOMY_001`
**Description**: Precise sunrise and sunset time calculations
**Status**: ✅ Implemented
**Coverage**: Unit + Integration Tests

**Components**:
- `CalculateSunTimesWithContext(ctx, location, date)` in `astronomy/sunrise.go`
- Geographic location support (latitude/longitude)
- Timezone handling
- Atmospheric refraction corrections

**API Endpoints**:
- `CalculateSunTimesWithContext(ctx, location, date)` - Primary method
- Location-based calculations with coordinate validation

**Test Coverage**:
- Unit tests: ✅ Complete
- Integration tests: ✅ Complete
- Geographic variation testing: ✅ Complete

#### 2.2 Ephemeris Integration
**Feature ID**: `ASTRONOMY_002`
**Description**: Planetary position calculations via Swiss Ephemeris
**Status**: ✅ Implemented
**Coverage**: Unit + Integration Tests

**Components**:
- `ephemeris.Manager` for planetary position calculations
- Julian day conversion utilities
- Sun and Moon longitude calculations
- Integration with Swiss Ephemeris library

**API Endpoints**:
- `GetPlanetaryPositions(ctx, julianDay)` - Planetary positions
- `TimeToJulianDay(time)` - Time conversion utility

### 3. Service Layer

#### 3.1 gRPC Service
**Feature ID**: `SERVICE_001`
**Description**: High-performance gRPC service for Panchangam data
**Status**: ✅ Implemented
**Coverage**: Unit Tests (Service), Integration Tests Needed

**Components**:
- `PanchangamServer` in `services/panchangam/service.go`
- Protocol Buffers definitions in `proto/panchangam.proto`
- Request validation and error handling
- OpenTelemetry instrumentation
- Comprehensive logging

**API Endpoints**:
- `Get(ctx, GetPanchangamRequest) -> GetPanchangamResponse` - Main service endpoint

**Request Parameters**:
- `date` (string): Target date in YYYY-MM-DD format
- `latitude` (float64): Geographic latitude (-90 to 90)
- `longitude` (float64): Geographic longitude (-180 to 180)
- `timezone` (string): IANA timezone identifier
- `region` (string): Regional calculation preferences
- `calculation_method` (string): Calculation method selection
- `locale` (string): Localization preferences

**Response Data**:
- Complete Panchangam data for the requested date
- All 5 Panchangam elements with detailed information
- Sunrise/sunset times
- Additional events and recommendations

**Test Coverage**:
- Unit tests: ✅ Complete for service logic
- Integration tests: ⚠️ **Missing** - functional tests needed
- Performance tests: ⚠️ **Missing** - end-to-end performance validation needed

### 4. Observability & Monitoring

#### 4.1 OpenTelemetry Integration
**Feature ID**: `OBSERVABILITY_001`
**Description**: Comprehensive observability with tracing, logging, and metrics
**Status**: ✅ Implemented
**Coverage**: Integration Tests

**Components**:
- Distributed tracing with OpenTelemetry
- Structured logging with context propagation
- Error categorization and severity levels
- Performance monitoring and alerting
- Calculation timing and success rate tracking

**Features**:
- Span creation for all major operations
- Detailed error recording with context
- Event tracking for important milestones
- Attribute enrichment for debugging
- Custom error categories (Validation, Calculation, Internal)

### 5. Quality Assurance

#### 5.1 Test Infrastructure
**Feature ID**: `QA_001`
**Description**: Comprehensive testing framework
**Status**: ✅ Implemented
**Coverage**: Multiple test types

**Components**:
- Unit tests for all calculators
- Integration tests for end-to-end workflows
- Performance benchmarks with <100ms targets
- Historical validation against known dates
- Mock-based testing for external dependencies

**Test Types**:
- `simple_test.go`: Core functionality tests
- `performance_test.go`: Performance benchmarks
- `historical_validation_test.go`: Accuracy validation
- Individual calculator test files

**Coverage Metrics**:
- Overall: 93.0% (services/panchangam)
- Astronomy modules: High coverage across all calculators
- Service layer: Complete logic coverage

#### 5.2 Code Quality
**Feature ID**: `QA_002`
**Description**: Code quality standards and validation
**Status**: ✅ Implemented
**Coverage**: Automated validation

**Components**:
- Input validation for all public APIs
- Error handling with proper error types
- Documentation for all public interfaces
- Consistent coding patterns and conventions
- Thread-safe implementations

## Feature Gap Analysis

### ✅ Fully Implemented Features
1. All 5 Panchangam elements with comprehensive calculations
2. Astronomical calculations (sunrise/sunset, planetary positions)
3. gRPC service layer with protocol buffers
4. OpenTelemetry observability integration
5. Comprehensive unit and performance testing
6. Code quality standards and validation

### ⚠️ Partially Implemented Features
1. **Service Integration Testing**: Missing functional tests for service layer
2. **End-to-End Performance Testing**: Missing service-level performance validation
3. **Service Response Integration**: Service returns placeholder data instead of calculated values

### ❌ Missing Features (Not in current scope)
1. REST API endpoints (only gRPC implemented)
2. Database persistence layer
3. Multi-language localization
4. Advanced festival calculations
5. Astrological chart generation
6. Mobile application interfaces

## Priority Implementation Items

### High Priority - Service Integration
1. **Service Integration Tests**: Connect service layer to astronomy calculations
2. **End-to-End Functional Tests**: Complete request-response validation
3. **Service Performance Tests**: Validate service-level performance targets

### Medium Priority - Enhancement
1. **Response Data Integration**: Replace placeholder data with calculated values
2. **Error Handling Enhancement**: Improve service-level error responses
3. **Additional Validation**: Enhanced input validation and sanitization

### Low Priority - Future Features
1. **Multi-language Support**: Localization for different languages
2. **Additional Event Types**: Festival and special day calculations
3. **Advanced Timing**: More precise astronomical calculations

## Test Coverage Summary

| Component | Unit Tests | Integration Tests | Performance Tests | Functional Tests |
|-----------|------------|-------------------|-------------------|------------------|
| Tithi | ✅ Complete | ✅ Complete | ✅ Complete | ⚠️ Service Missing |
| Nakshatra | ✅ Complete | ✅ Complete | ✅ Complete | ⚠️ Service Missing |
| Yoga | ✅ Complete | ✅ Complete | ✅ Complete | ⚠️ Service Missing |
| Karana | ✅ Complete | ✅ Complete | ✅ Complete | ⚠️ Service Missing |
| Vara | ✅ Complete | ✅ Complete | ✅ Complete | ⚠️ Service Missing |
| Sunrise/Sunset | ✅ Complete | ✅ Complete | ✅ Complete | ⚠️ Service Missing |
| gRPC Service | ✅ Complete | ❌ Missing | ❌ Missing | ❌ Missing |
| Overall Coverage | 93.0% | High | High | **Missing** |

## Performance Targets

| Operation | Current Performance | Target | Status |
|-----------|---------------------|---------|---------|
| Individual Calculator | <10ms | <50ms | ✅ Exceeds |
| All 5 Elements | <50ms | <100ms | ✅ Meets |
| Service Response | Not measured | <200ms | ⚠️ Needs testing |
| End-to-End Request | Not measured | <500ms | ⚠️ Needs testing |

## Next Steps for Complete Feature Coverage

1. **Implement Service Functional Tests** - High priority missing component
2. **Integrate Service with Calculations** - Replace placeholder data
3. **Add End-to-End Performance Tests** - Service-level performance validation
4. **Complete Feature Coverage Report** - Comprehensive testing validation