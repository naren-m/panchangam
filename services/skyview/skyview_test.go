package skyview

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

// MockEphemerisProvider for testing
type MockEphemerisProvider struct {
	planetary *ephemeris.PlanetaryPositions
	solar     *ephemeris.SolarPosition
	lunar     *ephemeris.LunarPosition
	available bool
}

func (m *MockEphemerisProvider) GetPlanetaryPositions(ctx context.Context, jd ephemeris.JulianDay) (*ephemeris.PlanetaryPositions, error) {
	return m.planetary, nil
}

func (m *MockEphemerisProvider) GetSunPosition(ctx context.Context, jd ephemeris.JulianDay) (*ephemeris.SolarPosition, error) {
	return m.solar, nil
}

func (m *MockEphemerisProvider) GetMoonPosition(ctx context.Context, jd ephemeris.JulianDay) (*ephemeris.LunarPosition, error) {
	return m.lunar, nil
}

func (m *MockEphemerisProvider) IsAvailable(ctx context.Context) bool {
	return m.available
}

func (m *MockEphemerisProvider) GetDataRange() (startJD, endJD ephemeris.JulianDay) {
	return 2451545.0, 2488070.0 // J2000 to 2100
}

func (m *MockEphemerisProvider) GetHealthStatus(ctx context.Context) (*ephemeris.HealthStatus, error) {
	return &ephemeris.HealthStatus{
		Available: m.available,
		LastCheck: time.Now(),
	}, nil
}

func (m *MockEphemerisProvider) GetProviderName() string {
	return "MockEphemerisProvider"
}

func (m *MockEphemerisProvider) GetVersion() string {
	return "1.0.0-test"
}

func (m *MockEphemerisProvider) Close() error {
	return nil
}

// Test helper to create mock ephemeris data
func createMockEphemerisData() *ephemeris.PlanetaryPositions {
	return &ephemeris.PlanetaryPositions{
		JulianDay: 2451545.0, // J2000
		Sun: ephemeris.Position{
			Longitude: 280.0,
			Latitude:  0.0,
			Distance:  1.0,
			Speed:     0.9856,
		},
		Moon: ephemeris.Position{
			Longitude: 218.0,
			Latitude:  5.0,
			Distance:  0.00257,
			Speed:     13.176,
		},
		Mercury: ephemeris.Position{
			Longitude: 252.0,
			Latitude:  1.0,
			Distance:  0.5,
			Speed:     1.5,
		},
		Venus: ephemeris.Position{
			Longitude: 330.0,
			Latitude:  2.0,
			Distance:  0.7,
			Speed:     1.2,
		},
		Mars: ephemeris.Position{
			Longitude: 355.0,
			Latitude:  1.5,
			Distance:  1.5,
			Speed:     0.5,
		},
		Jupiter: ephemeris.Position{
			Longitude: 45.0,
			Latitude:  1.0,
			Distance:  5.2,
			Speed:     0.08,
		},
		Saturn: ephemeris.Position{
			Longitude: 180.0,
			Latitude:  2.0,
			Distance:  9.5,
			Speed:     0.03,
		},
		Uranus: ephemeris.Position{
			Longitude: 300.0,
			Latitude:  0.5,
			Distance:  19.2,
			Speed:     0.01,
		},
		Neptune: ephemeris.Position{
			Longitude: 270.0,
			Latitude:  1.0,
			Distance:  30.1,
			Speed:     0.006,
		},
	}
}

func createMockLunarPosition() *ephemeris.LunarPosition {
	return &ephemeris.LunarPosition{
		JulianDay:       2451545.0,
		Longitude:       218.0,
		Latitude:        5.0,
		Distance:        384400.0, // km
		Phase:           0.5,      // Half moon
		Illumination:    50.0,     // 50% illuminated
		PhaseAngle:      90.0,
		AngularDiameter: 1800.0, // arcseconds
	}
}
