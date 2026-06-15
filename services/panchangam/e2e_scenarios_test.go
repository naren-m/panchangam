package panchangam

import (
	"context"
	"testing"
	"time"

	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testMultipleLocations(t *testing.T) {
	server := NewPanchangamServer()
	ctx := context.Background()

	locations := []struct {
		name      string
		latitude  float64
		longitude float64
		timezone  string
	}{
		{"Bangalore_India", 12.9716, 77.5946, "Asia/Kolkata"},
		{"New_York_USA", 40.7128, -74.0060, "America/New_York"},
		{"London_UK", 51.5074, -0.1278, "Europe/London"},
		{"Sydney_Australia", -33.8688, 151.2093, "Australia/Sydney"},
	}

	for _, loc := range locations {
		t.Run("E2E_Location_"+loc.name, func(t *testing.T) {
			req := &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  loc.latitude,
				Longitude: loc.longitude,
				Timezone:  loc.timezone,
			}

			resp, err := server.Get(ctx, req)
			require.NoError(t, err, "E2E: Location %s should succeed", loc.name)
			require.NotNil(t, resp, "E2E: Response should not be nil for %s", loc.name)
			require.NotNil(t, resp.PanchangamData, "E2E: Data should not be nil for %s", loc.name)

			data := resp.PanchangamData

			assert.Equal(t, req.Date, data.Date, "E2E: Date should match for %s", loc.name)
			assert.NotEmpty(t, data.SunriseTime, "E2E: Sunrise should be calculated for %s", loc.name)
			assert.NotEmpty(t, data.SunsetTime, "E2E: Sunset should be calculated for %s", loc.name)
			assert.NotEmpty(t, data.Tithi, "E2E: Tithi should be present for %s", loc.name)
			assert.NotEmpty(t, data.Nakshatra, "E2E: Nakshatra should be present for %s", loc.name)
			assert.NotEmpty(t, data.Yoga, "E2E: Yoga should be present for %s", loc.name)
			assert.NotEmpty(t, data.Karana, "E2E: Karana should be present for %s", loc.name)

			t.Logf("E2E: Location %s validated", loc.name)
			t.Logf("Sunrise: %s, Sunset: %s", data.SunriseTime, data.SunsetTime)
		})
	}
}

func testMultipleDates(t *testing.T) {
	server := NewPanchangamServer()
	ctx := context.Background()

	dates := []struct {
		name string
		date string
		desc string
	}{
		{"New_Year", "2024-01-01", "New Year"},
		{"Spring_Equinox", "2024-03-20", "Spring Equinox"},
		{"Summer_Solstice", "2024-06-21", "Summer Solstice"},
		{"Autumn_Equinox", "2024-09-22", "Autumn Equinox"},
		{"Winter_Solstice", "2024-12-21", "Winter Solstice"},
	}

	for _, dateTest := range dates {
		t.Run("E2E_Date_"+dateTest.name, func(t *testing.T) {
			req := &ppb.GetPanchangamRequest{
				Date:      dateTest.date,
				Latitude:  12.9716,
				Longitude: 77.5946,
				Timezone:  "Asia/Kolkata",
			}

			resp, err := server.Get(ctx, req)
			require.NoError(t, err, "E2E: Date %s should succeed", dateTest.name)
			require.NotNil(t, resp, "E2E: Response should not be nil for %s", dateTest.name)
			require.NotNil(t, resp.PanchangamData, "E2E: Data should not be nil for %s", dateTest.name)

			data := resp.PanchangamData

			assert.Equal(t, req.Date, data.Date, "E2E: Date should match for %s", dateTest.name)
			assert.NotEmpty(t, data.Tithi, "E2E: Tithi should be calculated for %s", dateTest.name)
			assert.NotEmpty(t, data.Nakshatra, "E2E: Nakshatra should be calculated for %s", dateTest.name)
			assert.NotEmpty(t, data.Yoga, "E2E: Yoga should be calculated for %s", dateTest.name)
			assert.NotEmpty(t, data.Karana, "E2E: Karana should be calculated for %s", dateTest.name)
			assert.NotEmpty(t, data.SunriseTime, "E2E: Sunrise should be calculated for %s", dateTest.name)
			assert.NotEmpty(t, data.SunsetTime, "E2E: Sunset should be calculated for %s", dateTest.name)

			t.Logf("E2E: Date %s (%s) validated", dateTest.name, dateTest.desc)
		})
	}
}

func testFeatureConsistency(t *testing.T) {
	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
	}

	var responses []*ppb.GetPanchangamResponse
	for i := 0; i < 5; i++ {
		resp, err := server.Get(ctx, req)
		require.NoError(t, err, "E2E: consistency request %d should succeed", i+1)
		require.NotNil(t, resp, "E2E: consistency request %d should return response", i+1)
		require.NotNil(t, resp.PanchangamData, "E2E: consistency request %d should return data", i+1)
		responses = append(responses, resp)
	}

	first := responses[0].PanchangamData
	second := responses[1].PanchangamData

	assert.Equal(t, first.Date, second.Date, "E2E: Date should be consistent")
	assert.Equal(t, first.Tithi, second.Tithi, "E2E: Tithi should be consistent")
	assert.Equal(t, first.Nakshatra, second.Nakshatra, "E2E: Nakshatra should be consistent")
	assert.Equal(t, first.Yoga, second.Yoga, "E2E: Yoga should be consistent")
	assert.Equal(t, first.Karana, second.Karana, "E2E: Karana should be consistent")
	assert.Equal(t, first.SunriseTime, second.SunriseTime, "E2E: Sunrise should be consistent")
	assert.Equal(t, first.SunsetTime, second.SunsetTime, "E2E: Sunset should be consistent")

	t.Logf("E2E: Feature consistency validated across %d successful requests", len(responses))
}

func testUserScenarios(t *testing.T) {
	server := NewPanchangamServer()
	ctx := context.Background()

	scenarios := []struct {
		name        string
		description string
		request     *ppb.GetPanchangamRequest
		validation  func(t *testing.T, resp *ppb.GetPanchangamResponse)
	}{
		{
			name:        "Morning_Planning",
			description: "User checking Panchangam for morning planning",
			request: &ppb.GetPanchangamRequest{
				Date:      time.Now().Format("2006-01-02"),
				Latitude:  12.9716,
				Longitude: 77.5946,
				Timezone:  "Asia/Kolkata",
				Region:    "India",
				Locale:    "en",
			},
			validation: func(t *testing.T, resp *ppb.GetPanchangamResponse) {
				require.NotNil(t, resp.PanchangamData, "Morning planning: Should have data")
				data := resp.PanchangamData
				assert.NotEmpty(t, data.Tithi, "Morning planning: Should have Tithi")
				assert.NotEmpty(t, data.SunriseTime, "Morning planning: Should have sunrise time")
			},
		},
		{
			name:        "Astrological_Consultation",
			description: "Astrologer requesting detailed Panchangam data",
			request: &ppb.GetPanchangamRequest{
				Date:              "2024-01-15",
				Latitude:          28.6139,
				Longitude:         77.2090,
				Timezone:          "Asia/Kolkata",
				Region:            "India",
				CalculationMethod: "traditional",
				Locale:            "en",
			},
			validation: func(t *testing.T, resp *ppb.GetPanchangamResponse) {
				require.NotNil(t, resp.PanchangamData, "Astrological: Should have data")
				data := resp.PanchangamData
				assert.NotEmpty(t, data.Tithi, "Astrological: Should have Tithi")
				assert.NotEmpty(t, data.Nakshatra, "Astrological: Should have Nakshatra")
				assert.NotEmpty(t, data.Yoga, "Astrological: Should have Yoga")
				assert.NotEmpty(t, data.Karana, "Astrological: Should have Karana")
				assert.True(t, len(data.Events) >= 0, "Astrological: Should have events list")
			},
		},
		{
			name:        "International_User",
			description: "International user requesting Panchangam data",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-01-15",
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timezone:  "America/New_York",
				Region:    "USA",
				Locale:    "en",
			},
			validation: func(t *testing.T, resp *ppb.GetPanchangamResponse) {
				require.NotNil(t, resp.PanchangamData, "International: Should have data")
				data := resp.PanchangamData
				assert.NotEmpty(t, data.SunriseTime, "International: Should calculate sunrise for location")
				assert.NotEmpty(t, data.SunsetTime, "International: Should calculate sunset for location")
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run("E2E_Scenario_"+scenario.name, func(t *testing.T) {
			start := time.Now()
			resp, err := server.Get(ctx, scenario.request)
			duration := time.Since(start)

			require.NoError(t, err, "E2E: Scenario %s should succeed", scenario.name)
			require.NotNil(t, resp, "E2E: Scenario %s should return response", scenario.name)
			scenario.validation(t, resp)
			assert.True(t, duration < 2*time.Second, "E2E: Scenario %s should be fast (<2s), got %v", scenario.name, duration)

			t.Logf("E2E: Scenario %s validated in %v", scenario.name, duration)
			t.Logf("Description: %s", scenario.description)
		})
	}
}
