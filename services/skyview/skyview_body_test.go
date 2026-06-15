package skyview

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

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

	body := service.createCelestialBody(
		"test",
		"Test Body",
		"परीक्षण",
		"परीक्षण",
		"planet",
		ephemeris.Position{
			Longitude: 90.0,
			Latitude:  0.0,
			Distance:  1.0,
			Speed:     1.0,
		},
		-2.0,
		"#ffffff",
		observer,
		180.0,
	)

	if body.ID != "test" {
		t.Errorf("ID = %v, want test", body.ID)
	}
	if body.Name != "Test Body" {
		t.Errorf("Name = %v, want Test Body", body.Name)
	}
	if body.EclipticCoords.Longitude != 90.0 {
		t.Errorf("Ecliptic longitude = %v, want 90.0", body.EclipticCoords.Longitude)
	}
	if body.EquatorialCoords == nil {
		t.Error("EquatorialCoords is nil")
	}
	if body.HorizontalCoords == nil {
		t.Error("HorizontalCoords is nil")
	}
	if body.HorizontalCoords.Azimuth < 0 || body.HorizontalCoords.Azimuth >= 360 {
		t.Errorf("Azimuth = %v, want value in [0, 360)", body.HorizontalCoords.Azimuth)
	}
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

	response, err := service.GetSkyView(
		context.Background(),
		observer,
		time.Date(2024, 6, 21, 12, 0, 0, 0, time.UTC),
	)

	if err != nil {
		t.Fatalf("GetSkyView() error = %v", err)
	}

	sunVisible := false
	for _, body := range response.VisibleBodies {
		if body.ID == "sun" {
			sunVisible = true
			break
		}
	}

	t.Logf("Sun visible at noon on summer solstice in New York: %v", sunVisible)
	t.Logf("Total visible bodies: %d out of %d", len(response.VisibleBodies), len(response.Bodies))
}
