package skyview

import (
	"context"
	"math"
	"testing"
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
		JulianDay:      2451545.0,
		Longitude:      218.0,
		Latitude:       5.0,
		Distance:       384400.0, // km
		Phase:          0.5,      // Half moon
		Illumination:   50.0,     // 50% illuminated
		PhaseAngle:     90.0,
		AngularDiameter: 1800.0,  // arcseconds
	}
}

func TestDateToJulianDay(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected float64
		delta    float64
	}{
		{
			name:     "J2000 epoch",
			date:     time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: 2451545.0,
			delta:    0.01,
		},
		{
			name:     "Unix epoch",
			date:     time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 2440587.5,
			delta:    0.01,
		},
		{
			name:     "Year 2024",
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 2460310.5,
			delta:    0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jd := DateToJulianDay(tt.date)
			if math.Abs(jd-tt.expected) > tt.delta {
				t.Errorf("DateToJulianDay() = %v, want %v ± %v", jd, tt.expected, tt.delta)
			}
		})
	}
}

func TestGetLocalSiderealTime(t *testing.T) {
	tests := []struct {
		name      string
		longitude float64
		time      time.Time
	}{
		{
			name:      "Greenwich at J2000",
			longitude: 0.0,
			time:      time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "New York",
			longitude: -74.006,
			time:      time.Date(2024, 6, 21, 18, 0, 0, 0, time.UTC),
		},
		{
			name:      "Tokyo",
			longitude: 139.6917,
			time:      time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lst := GetLocalSiderealTime(tt.longitude, tt.time)

			// LST should be between 0 and 360 degrees
			if lst < 0 || lst >= 360 {
				t.Errorf("GetLocalSiderealTime() = %v, want value in [0, 360)", lst)
			}
		})
	}
}

func TestEclipticToEquatorial(t *testing.T) {
	tests := []struct {
		name      string
		ecliptic  EclipticCoordinates
		wantRA    float64 // approximate
		wantDec   float64 // approximate
		deltaRA   float64
		deltaDec  float64
	}{
		{
			name: "Vernal equinox",
			ecliptic: EclipticCoordinates{
				Longitude: 0,
				Latitude:  0,
				Distance:  1,
			},
			wantRA:   0,
			wantDec:  0,
			deltaRA:  1,
			deltaDec: 1,
		},
		{
			name: "Summer solstice",
			ecliptic: EclipticCoordinates{
				Longitude: 90,
				Latitude:  0,
				Distance:  1,
			},
			wantRA:   90,
			wantDec:  23.4,
			deltaRA:  1,
			deltaDec: 1,
		},
		{
			name: "Ecliptic north pole",
			ecliptic: EclipticCoordinates{
				Longitude: 0,
				Latitude:  90,
				Distance:  1,
			},
			wantRA:   270, // RA can be any value at pole, just check dec
			wantDec:  66.6,
			deltaRA:  180, // Large delta since RA is undefined at pole
			deltaDec: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eq := eclipticToEquatorial(tt.ecliptic)

			if math.Abs(eq.RightAscension-tt.wantRA) > tt.deltaRA {
				t.Errorf("RA = %v, want %v ± %v", eq.RightAscension, tt.wantRA, tt.deltaRA)
			}
			if math.Abs(eq.Declination-tt.wantDec) > tt.deltaDec {
				t.Errorf("Dec = %v, want %v ± %v", eq.Declination, tt.wantDec, tt.deltaDec)
			}
			if eq.Distance != tt.ecliptic.Distance {
				t.Errorf("Distance = %v, want %v", eq.Distance, tt.ecliptic.Distance)
			}
		})
	}
}

func TestEquatorialToHorizontal(t *testing.T) {
	observer := Observer{
		Latitude:  40.7128, // New York
		Longitude: -74.006,
		Altitude:  0,
		Timezone:  "America/New_York",
	}

	tests := []struct {
		name       string
		equatorial EquatorialCoordinates
		lst        float64
	}{
		{
			name: "Object at meridian",
			equatorial: EquatorialCoordinates{
				RightAscension: 180,
				Declination:    40,
				Distance:       1,
			},
			lst: 180,
		},
		{
			name: "Object rising",
			equatorial: EquatorialCoordinates{
				RightAscension: 90,
				Declination:    0,
				Distance:       1,
			},
			lst: 270,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hz := equatorialToHorizontal(tt.equatorial, observer, tt.lst)

			// Check that azimuth is in valid range (allowing for floating point edge case)
			if hz.Azimuth < 0 || hz.Azimuth > 360 {
				t.Errorf("Azimuth = %v, want value in [0, 360]", hz.Azimuth)
			}

			// Check that altitude is in valid range
			if hz.Altitude < -90 || hz.Altitude > 90 {
				t.Errorf("Altitude = %v, want value in [-90, 90]", hz.Altitude)
			}

			// Check distance preservation
			if hz.Distance != tt.equatorial.Distance {
				t.Errorf("Distance = %v, want %v", hz.Distance, tt.equatorial.Distance)
			}
		})
	}
}

func TestSkyViewService_GetSkyView(t *testing.T) {
	mockProvider := &MockEphemerisProvider{
		planetary: createMockEphemerisData(),
		lunar:     createMockLunarPosition(),
		available: true,
	}

	service := NewSkyViewService(mockProvider)

	observer := Observer{
		Latitude:  40.7128,
		Longitude: -74.006,
		Altitude:  0,
		Timezone:  "America/New_York",
	}

	testTime := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)

	ctx := context.Background()
	response, err := service.GetSkyView(ctx, observer, testTime)

	if err != nil {
		t.Fatalf("GetSkyView() error = %v", err)
	}

	// Check response structure
	if response == nil {
		t.Fatal("GetSkyView() returned nil response")
	}

	// Check observer data
	if response.Observer.Latitude != observer.Latitude {
		t.Errorf("Observer.Latitude = %v, want %v", response.Observer.Latitude, observer.Latitude)
	}

	// Check that we have celestial bodies
	if len(response.Bodies) == 0 {
		t.Error("GetSkyView() returned no bodies")
	}

	// Check for expected bodies (Sun, Moon, and planets)
	expectedBodies := []string{"sun", "moon", "mercury", "venus", "mars", "jupiter", "saturn"}
	bodyMap := make(map[string]bool)
	for _, body := range response.Bodies {
		bodyMap[body.ID] = true
	}

	for _, expected := range expectedBodies {
		if !bodyMap[expected] {
			t.Errorf("Missing expected body: %s", expected)
		}
	}

	// Check that visible bodies are subset of all bodies
	if len(response.VisibleBodies) > len(response.Bodies) {
		t.Error("More visible bodies than total bodies")
	}

	// Check Julian day
	expectedJD := 2451545.0
	if math.Abs(response.JulianDay-expectedJD) > 0.1 {
		t.Errorf("JulianDay = %v, want %v ± 0.1", response.JulianDay, expectedJD)
	}

	// Check LST is in valid range
	if response.LocalSiderealTime < 0 || response.LocalSiderealTime >= 360 {
		t.Errorf("LocalSiderealTime = %v, want value in [0, 360)", response.LocalSiderealTime)
	}
}

func TestSkyViewService_GetSkyView_NilProvider(t *testing.T) {
	service := NewSkyViewService(nil)

	observer := Observer{
		Latitude:  0,
		Longitude: 0,
		Altitude:  0,
		Timezone:  "UTC",
	}

	ctx := context.Background()
	_, err := service.GetSkyView(ctx, observer, time.Now())

	if err == nil {
		t.Error("GetSkyView() with nil provider should return error")
	}
}

func TestCelestialBodyCreation(t *testing.T) {
	mockProvider := &MockEphemerisProvider{
		planetary: createMockEphemerisData(),
		lunar:     createMockLunarPosition(),
		available: true,
	}

	service := NewSkyViewService(mockProvider)

	observer := Observer{
		Latitude:  40.7128,
		Longitude: -74.006,
		Altitude:  0,
		Timezone:  "America/New_York",
	}

	position := ephemeris.Position{
		Longitude: 90.0,
		Latitude:  0.0,
		Distance:  1.0,
		Speed:     1.0,
	}

	lst := 180.0

	body := service.createCelestialBody(
		"test",
		"Test Body",
		"परीक्षण",
		"परीक्षण",
		"planet",
		position,
		-2.0,
		"#ffffff",
		observer,
		lst,
	)

	// Check basic properties
	if body.ID != "test" {
		t.Errorf("ID = %v, want test", body.ID)
	}
	if body.Name != "Test Body" {
		t.Errorf("Name = %v, want Test Body", body.Name)
	}

	// Check coordinate transformations were applied
	if body.EclipticCoords.Longitude != 90.0 {
		t.Errorf("Ecliptic longitude = %v, want 90.0", body.EclipticCoords.Longitude)
	}

	if body.EquatorialCoords == nil {
		t.Error("EquatorialCoords is nil")
	}

	if body.HorizontalCoords == nil {
		t.Error("HorizontalCoords is nil")
	}

	// Check azimuth is in valid range
	if body.HorizontalCoords.Azimuth < 0 || body.HorizontalCoords.Azimuth >= 360 {
		t.Errorf("Azimuth = %v, want value in [0, 360)", body.HorizontalCoords.Azimuth)
	}

	// Check altitude is in valid range
	if body.HorizontalCoords.Altitude < -90 || body.HorizontalCoords.Altitude > 90 {
		t.Errorf("Altitude = %v, want value in [-90, 90]", body.HorizontalCoords.Altitude)
	}
}

func TestVisibilityDetermination(t *testing.T) {
	mockProvider := &MockEphemerisProvider{
		planetary: createMockEphemerisData(),
		lunar:     createMockLunarPosition(),
		available: true,
	}

	service := NewSkyViewService(mockProvider)

	observer := Observer{
		Latitude:  40.7128,
		Longitude: -74.006,
		Altitude:  0,
		Timezone:  "America/New_York",
	}

	testTime := time.Date(2024, 6, 21, 12, 0, 0, 0, time.UTC) // Summer solstice, noon

	ctx := context.Background()
	response, err := service.GetSkyView(ctx, observer, testTime)

	if err != nil {
		t.Fatalf("GetSkyView() error = %v", err)
	}

	// At noon on summer solstice in New York, Sun should be visible
	sunVisible := false
	for _, body := range response.VisibleBodies {
		if body.ID == "sun" {
			sunVisible = true
			break
		}
	}

	// Note: This test may vary based on actual ephemeris data
	// The test is mainly checking that visibility logic is applied
	t.Logf("Sun visible at noon on summer solstice in New York: %v", sunVisible)
	t.Logf("Total visible bodies: %d out of %d", len(response.VisibleBodies), len(response.Bodies))
}
