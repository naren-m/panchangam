package skyview

import (
	"context"
	"math"
	"testing"
	"time"
)

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

	response, err := service.GetSkyView(
		context.Background(),
		observer,
		time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
	)

	if err != nil {
		t.Fatalf("GetSkyView() error = %v", err)
	}
	if response == nil {
		t.Fatal("GetSkyView() returned nil response")
	}
	if response.Observer.Latitude != observer.Latitude {
		t.Errorf("Observer.Latitude = %v, want %v", response.Observer.Latitude, observer.Latitude)
	}
	if len(response.Bodies) == 0 {
		t.Error("GetSkyView() returned no bodies")
	}

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

	if len(response.VisibleBodies) > len(response.Bodies) {
		t.Error("More visible bodies than total bodies")
	}

	expectedJD := 2451545.0
	if math.Abs(response.JulianDay-expectedJD) > 0.1 {
		t.Errorf("JulianDay = %v, want %v ± 0.1", response.JulianDay, expectedJD)
	}

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

	_, err := service.GetSkyView(context.Background(), observer, time.Now())

	if err == nil {
		t.Error("GetSkyView() with nil provider should return error")
	}
}
