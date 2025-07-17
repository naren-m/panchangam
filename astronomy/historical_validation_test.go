package astronomy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHistoricalValidation tests sunrise/sunset calculations against known historical data
// from TimeAndDate.com for January 15, 2020 across all continents
func TestHistoricalValidation(t *testing.T) {
	// Test date: January 15, 2020 (historical date)
	testDate := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	
	tests := []struct {
		name              string
		continent         string
		location          Location
		expectedSunrise   time.Time
		expectedSunset    time.Time
		tolerance         time.Duration
		source            string
		notes             string
	}{
		{
			name:      "New York - North America",
			continent: "North America",
			location: Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			expectedSunrise: time.Date(2020, 1, 15, 12, 18, 0, 0, time.UTC), // 7:18 EST = 12:18 UTC
			expectedSunset:  time.Date(2020, 1, 15, 21, 52, 0, 0, time.UTC), // 4:52 EST = 21:52 UTC
			tolerance:       15 * time.Minute,
			source:          "TimeAndDate.com",
			notes:           "EST timezone (UTC-5), winter solstice period",
		},
		{
			name:      "London - Europe",
			continent: "Europe",
			location: Location{
				Latitude:  51.5074,
				Longitude: -0.1278,
			},
			expectedSunrise: time.Date(2020, 1, 15, 7, 59, 0, 0, time.UTC),
			expectedSunset:  time.Date(2020, 1, 15, 16, 19, 0, 0, time.UTC),
			tolerance:       15 * time.Minute,
			source:          "TimeAndDate.com",
			notes:           "GMT timezone, shortest days of year",
		},
		{
			name:      "Tokyo - Asia",
			continent: "Asia",
			location: Location{
				Latitude:  35.6762,
				Longitude: 139.6503,
			},
			expectedSunrise: time.Date(2020, 1, 15, 21, 50, 0, 0, time.UTC), // 6:50 JST = 21:50 UTC (previous day)
			expectedSunset:  time.Date(2020, 1, 15, 7, 50, 0, 0, time.UTC),  // 4:50 JST = 7:50 UTC (next day)
			tolerance:       15 * time.Minute,
			source:          "TimeAndDate.com",
			notes:           "JST timezone (UTC+9), spans UTC days",
		},
		{
			name:      "Sydney - Australia/Oceania",
			continent: "Australia/Oceania",
			location: Location{
				Latitude:  -33.8688,
				Longitude: 151.2093,
			},
			expectedSunrise: time.Date(2020, 1, 15, 18, 59, 0, 0, time.UTC), // 5:59 AEDT = 18:59 UTC (previous day)
			expectedSunset:  time.Date(2020, 1, 15, 9, 9, 0, 0, time.UTC),   // 8:09 AEDT = 9:09 UTC (next day)
			tolerance:       15 * time.Minute,
			source:          "TimeAndDate.com",
			notes:           "AEDT timezone (UTC+11), summer in southern hemisphere",
		},
		{
			name:      "Mumbai - Asia",
			continent: "Asia",
			location: Location{
				Latitude:  19.0760,
				Longitude: 72.8777,
			},
			expectedSunrise: time.Date(2020, 1, 15, 1, 44, 0, 0, time.UTC), // 7:14 IST = 1:44 UTC
			expectedSunset:  time.Date(2020, 1, 15, 12, 51, 0, 0, time.UTC), // 6:21 IST = 12:51 UTC
			tolerance:       15 * time.Minute,
			source:          "TimeAndDate.com",
			notes:           "IST timezone (UTC+5:30), tropical location",
		},
		{
			name:      "Cape Town - Africa",
			continent: "Africa",
			location: Location{
				Latitude:  -33.9249,
				Longitude: 18.4241,
			},
			expectedSunrise: time.Date(2020, 1, 15, 3, 50, 0, 0, time.UTC), // 5:50 SAST = 3:50 UTC
			expectedSunset:  time.Date(2020, 1, 15, 18, 0, 0, 0, time.UTC),  // 8:00 SAST = 18:00 UTC
			tolerance:       15 * time.Minute,
			source:          "TimeAndDate.com",
			notes:           "SAST timezone (UTC+2), summer in southern hemisphere",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate sunrise and sunset using our algorithm
			sunTimes, err := CalculateSunTimes(tt.location, testDate)
			require.NoError(t, err, "Failed to calculate sun times for %s", tt.name)
			require.NotNil(t, sunTimes, "Sun times should not be nil for %s", tt.name)

			// Convert our calculated times to UTC for comparison
			// (Our algorithm returns times in the input date's timezone, which is UTC)
			
			// Check sunrise accuracy
			sunriseDiff := sunTimes.Sunrise.Sub(tt.expectedSunrise)
			if sunriseDiff < 0 {
				sunriseDiff = -sunriseDiff
			}
			// Handle day boundary crossing
			if sunriseDiff > 12*time.Hour {
				sunriseDiff = 24*time.Hour - sunriseDiff
			}
			
			// Check sunset accuracy
			sunsetDiff := sunTimes.Sunset.Sub(tt.expectedSunset)
			if sunsetDiff < 0 {
				sunsetDiff = -sunsetDiff
			}
			// Handle day boundary crossing
			if sunsetDiff > 12*time.Hour {
				sunsetDiff = 24*time.Hour - sunsetDiff
			}

			// Log detailed comparison
			t.Logf("=== %s (%s) ===", tt.name, tt.continent)
			t.Logf("Source: %s", tt.source)
			t.Logf("Notes: %s", tt.notes)
			t.Logf("Expected Sunrise: %s", tt.expectedSunrise.Format("15:04:05"))
			t.Logf("Calculated Sunrise: %s", sunTimes.Sunrise.Format("15:04:05"))
			t.Logf("Sunrise Difference: %v", sunriseDiff)
			t.Logf("Expected Sunset: %s", tt.expectedSunset.Format("15:04:05"))
			t.Logf("Calculated Sunset: %s", sunTimes.Sunset.Format("15:04:05"))
			t.Logf("Sunset Difference: %v", sunsetDiff)
			t.Logf("Day Length (Expected): %v", tt.expectedSunset.Sub(tt.expectedSunrise))
			t.Logf("Day Length (Calculated): %v", sunTimes.Sunset.Sub(sunTimes.Sunrise))
			t.Logf("Tolerance: %v", tt.tolerance)
			t.Logf("")

			// Validate sunrise time
			assert.True(t, sunriseDiff <= tt.tolerance,
				"Sunrise time difference (%v) exceeds tolerance (%v) for %s.\nExpected: %s\nGot: %s",
				sunriseDiff, tt.tolerance, tt.name,
				tt.expectedSunrise.Format("15:04:05"),
				sunTimes.Sunrise.Format("15:04:05"))

			// Validate sunset time
			assert.True(t, sunsetDiff <= tt.tolerance,
				"Sunset time difference (%v) exceeds tolerance (%v) for %s.\nExpected: %s\nGot: %s",
				sunsetDiff, tt.tolerance, tt.name,
				tt.expectedSunset.Format("15:04:05"),
				sunTimes.Sunset.Format("15:04:05"))

			// Validate day length is reasonable (between 6-18 hours for non-polar regions)
			dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)
			if dayLength < 0 {
				dayLength = dayLength + 24*time.Hour
			}
			assert.True(t, dayLength >= 6*time.Hour && dayLength <= 18*time.Hour,
				"Day length (%v) is outside reasonable range for %s", dayLength, tt.name)
		})
	}
}

// TestHistoricalValidationDifferentDates tests multiple historical dates
func TestHistoricalValidationDifferentDates(t *testing.T) {
	// Test multiple historical dates with New York as reference
	location := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	
	tests := []struct {
		name            string
		date            time.Time
		expectedSunrise time.Time
		expectedSunset  time.Time
		tolerance       time.Duration
		source          string
	}{
		{
			name:            "New York - Summer Solstice 2019",
			date:            time.Date(2019, 6, 21, 0, 0, 0, 0, time.UTC),
			expectedSunrise: time.Date(2019, 6, 21, 9, 25, 0, 0, time.UTC), // 5:25 EDT = 9:25 UTC
			expectedSunset:  time.Date(2019, 6, 21, 0, 30, 0, 0, time.UTC),  // 8:30 EDT = 0:30 UTC (next day)
			tolerance:       20 * time.Minute,
			source:          "TimeAndDate.com approximation",
		},
		{
			name:            "New York - Winter Solstice 2019",
			date:            time.Date(2019, 12, 21, 0, 0, 0, 0, time.UTC),
			expectedSunrise: time.Date(2019, 12, 21, 12, 20, 0, 0, time.UTC), // 7:20 EST = 12:20 UTC
			expectedSunset:  time.Date(2019, 12, 21, 21, 38, 0, 0, time.UTC), // 4:38 EST = 21:38 UTC
			tolerance:       20 * time.Minute,
			source:          "TimeAndDate.com approximation",
		},
		{
			name:            "New York - Spring Equinox 2020",
			date:            time.Date(2020, 3, 20, 0, 0, 0, 0, time.UTC),
			expectedSunrise: time.Date(2020, 3, 20, 11, 0, 0, 0, time.UTC), // 7:00 EDT = 11:00 UTC
			expectedSunset:  time.Date(2020, 3, 20, 23, 0, 0, 0, time.UTC), // 7:00 EDT = 23:00 UTC
			tolerance:       20 * time.Minute,
			source:          "TimeAndDate.com approximation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sunTimes, err := CalculateSunTimes(location, tt.date)
			require.NoError(t, err)
			require.NotNil(t, sunTimes)

			sunriseDiff := sunTimes.Sunrise.Sub(tt.expectedSunrise)
			if sunriseDiff < 0 {
				sunriseDiff = -sunriseDiff
			}
			
			sunsetDiff := sunTimes.Sunset.Sub(tt.expectedSunset)
			if sunsetDiff < 0 {
				sunsetDiff = -sunsetDiff
			}

			t.Logf("=== %s ===", tt.name)
			t.Logf("Date: %s", tt.date.Format("2006-01-02"))
			t.Logf("Expected Sunrise: %s, Got: %s, Diff: %v", 
				tt.expectedSunrise.Format("15:04:05"), 
				sunTimes.Sunrise.Format("15:04:05"), 
				sunriseDiff)
			t.Logf("Expected Sunset: %s, Got: %s, Diff: %v", 
				tt.expectedSunset.Format("15:04:05"), 
				sunTimes.Sunset.Format("15:04:05"), 
				sunsetDiff)

			assert.True(t, sunriseDiff <= tt.tolerance,
				"Sunrise difference (%v) exceeds tolerance (%v)", sunriseDiff, tt.tolerance)
			assert.True(t, sunsetDiff <= tt.tolerance,
				"Sunset difference (%v) exceeds tolerance (%v)", sunsetDiff, tt.tolerance)
		})
	}
}

// TestExtremeLatitudes tests locations near polar regions
func TestExtremeLatitudes(t *testing.T) {
	testDate := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	
	tests := []struct {
		name     string
		location Location
		desc     string
	}{
		{
			name: "Reykjavik, Iceland - High Latitude",
			location: Location{
				Latitude:  64.1466,
				Longitude: -21.9426,
			},
			desc: "High latitude, short winter days",
		},
		{
			name: "Anchorage, Alaska - Very High Latitude",
			location: Location{
				Latitude:  61.2181,
				Longitude: -149.9003,
			},
			desc: "Very high latitude, extreme seasonal variation",
		},
		{
			name: "Ushuaia, Argentina - Southern High Latitude",
			location: Location{
				Latitude:  -54.8019,
				Longitude: -68.3030,
			},
			desc: "Southern hemisphere high latitude, summer season",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sunTimes, err := CalculateSunTimes(tt.location, testDate)
			require.NoError(t, err)
			require.NotNil(t, sunTimes)

			t.Logf("=== %s ===", tt.name)
			t.Logf("Description: %s", tt.desc)
			t.Logf("Location: %.4f°N, %.4f°E", tt.location.Latitude, tt.location.Longitude)
			t.Logf("Sunrise: %s", sunTimes.Sunrise.Format("15:04:05"))
			t.Logf("Sunset: %s", sunTimes.Sunset.Format("15:04:05"))
			
			dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)
			if dayLength < 0 {
				dayLength = dayLength + 24*time.Hour
			}
			t.Logf("Day Length: %v", dayLength)

			// Basic validation - times should be reasonable for the date
			assert.NotZero(t, sunTimes.Sunrise)
			assert.NotZero(t, sunTimes.Sunset)
			assert.True(t, dayLength > 0 && dayLength < 24*time.Hour)
		})
	}
}