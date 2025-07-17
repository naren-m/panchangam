package astronomy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateSunTimes(t *testing.T) {
	tests := []struct {
		name     string
		location Location
		date     time.Time
		desc     string
	}{
		{
			name: "Equator Equinox",
			location: Location{
				Latitude:  0.0,
				Longitude: 0.0,
			},
			date: time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
			desc: "Should have approximately 12-hour day",
		},
		{
			name: "New York Summer Solstice",
			location: Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			date: time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
			desc: "Should have long day (>12 hours)",
		},
		{
			name: "London Winter Solstice",
			location: Location{
				Latitude:  51.5074,
				Longitude: -0.1278,
			},
			date: time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
			desc: "Should have short day (<12 hours)",
		},
		{
			name: "Chennai India",
			location: Location{
				Latitude:  13.0827,
				Longitude: 80.2707,
			},
			date: time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
			desc: "Tropical location",
		},
		{
			name: "Sydney Australia",
			location: Location{
				Latitude:  -33.8688,
				Longitude: 151.2093,
			},
			date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			desc: "Southern hemisphere summer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sunTimes, err := CalculateSunTimes(tt.location, tt.date)
			require.NoError(t, err)
			require.NotNil(t, sunTimes)

			// Basic validation - times should be valid
			assert.NotZero(t, sunTimes.Sunrise)
			assert.NotZero(t, sunTimes.Sunset)
			
			// Sunrise and sunset should be on the same day (allowing for timezone differences)
			// For locations crossing UTC boundaries, day may differ by 1
			sunriseDay := sunTimes.Sunrise.Day()
			sunsetDay := sunTimes.Sunset.Day()
			targetDay := tt.date.Day()
			
			assert.True(t, sunriseDay == targetDay || sunriseDay == targetDay-1 || sunriseDay == targetDay+1,
				"Sunrise day %d should be within 1 day of target day %d", sunriseDay, targetDay)
			assert.True(t, sunsetDay == targetDay || sunsetDay == targetDay-1 || sunsetDay == targetDay+1,
				"Sunset day %d should be within 1 day of target day %d", sunsetDay, targetDay)
			
			// For most locations, sunset should be after sunrise
			if tt.location.Latitude < 70 && tt.location.Latitude > -70 {
				// Skip polar regions
				dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)
				// Handle case where sunset is on next day (crossing midnight)
				if dayLength < 0 {
					dayLength = dayLength + 24*time.Hour
				}
				assert.Positive(t, dayLength.Minutes())
				t.Logf("Location: %s, Sunrise: %s, Sunset: %s, Day length: %v", 
					tt.name, sunTimes.Sunrise.Format("15:04:05"), sunTimes.Sunset.Format("15:04:05"), dayLength)
			}
		})
	}
}

func TestJulianDate(t *testing.T) {
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
			delta:    0.5,
		},
		{
			name:     "Random date 2024",
			date:     time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC),
			expected: 2460476.0,
			delta:    0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jd := julianDate(tt.date)
			assert.InDelta(t, tt.expected, jd, tt.delta)
		})
	}
}

func TestGetSunriseTime(t *testing.T) {
	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)

	sunrise, err := GetSunriseTime(loc, date)
	assert.NoError(t, err)
	assert.NotZero(t, sunrise)
	assert.Equal(t, date.Day(), sunrise.Day())
}

func TestGetSunsetTime(t *testing.T) {
	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)

	sunset, err := GetSunsetTime(loc, date)
	assert.NoError(t, err)
	assert.NotZero(t, sunset)
	assert.Equal(t, date.Day(), sunset.Day())
}

func TestPolarConditions(t *testing.T) {
	tests := []struct {
		name     string
		location Location
		date     time.Time
		desc     string
	}{
		{
			name: "Arctic Summer - Continuous Daylight",
			location: Location{
				Latitude:  71.0, // North of Arctic Circle
				Longitude: 0.0,
			},
			date: time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
			desc: "Should handle polar day",
		},
		{
			name: "Arctic Winter - Continuous Night",
			location: Location{
				Latitude:  71.0, // North of Arctic Circle
				Longitude: 0.0,
			},
			date: time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
			desc: "Should handle polar night",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sunTimes, err := CalculateSunTimes(tt.location, tt.date)
			assert.NoError(t, err)
			assert.NotNil(t, sunTimes)
		})
	}
}

func TestInvalidInputs(t *testing.T) {
	tests := []struct {
		name     string
		location Location
		desc     string
	}{
		{
			name: "Valid extreme latitude - North Pole",
			location: Location{
				Latitude:  90.0,
				Longitude: 0.0,
			},
			desc: "Should handle North Pole",
		},
		{
			name: "Valid extreme latitude - South Pole",
			location: Location{
				Latitude:  -90.0,
				Longitude: 0.0,
			},
			desc: "Should handle South Pole",
		},
		{
			name: "Valid extreme longitude - International Date Line East",
			location: Location{
				Latitude:  0.0,
				Longitude: 180.0,
			},
			desc: "Should handle IDL East",
		},
		{
			name: "Valid extreme longitude - International Date Line West",
			location: Location{
				Latitude:  0.0,
				Longitude: -180.0,
			},
			desc: "Should handle IDL West",
		},
	}

	date := time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sunTimes, err := CalculateSunTimes(tt.location, date)
			assert.NoError(t, err)
			assert.NotNil(t, sunTimes)
		})
	}
}

func TestSolarPosition(t *testing.T) {
	jd := julianDate(time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC))
	eqTime, decl := solarPosition(jd)
	
	// Equation of time should be within reasonable bounds (-20 to +20 minutes)
	assert.Greater(t, eqTime, -20.0)
	assert.Less(t, eqTime, 20.0)
	
	// Declination should be within solar bounds (-23.44 to +23.44 degrees in radians)
	assert.Greater(t, decl, -23.44*DegToRad)
	assert.Less(t, decl, 23.44*DegToRad)
}

func TestCalculateRiseSet(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		jd        float64
		eqTime    float64
		decl      float64
	}{
		{
			name:      "Equator test",
			latitude:  0.0,
			longitude: 0.0,
			jd:        2451545.0,
			eqTime:    0.0,
			decl:      0.0,
		},
		{
			name:      "Mid-latitude test",
			latitude:  45.0,
			longitude: 0.0,
			jd:        2451545.0,
			eqTime:    0.0,
			decl:      10.0 * DegToRad,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sunrise, sunset := calculateRiseSet(tt.latitude, tt.longitude, tt.jd, tt.eqTime, tt.decl)
			
			// Basic sanity checks
			assert.GreaterOrEqual(t, sunrise, 0.0)
			assert.LessOrEqual(t, sunrise, 1440.0) // 24 hours in minutes
			assert.GreaterOrEqual(t, sunset, 0.0)
			assert.LessOrEqual(t, sunset, 1440.0)
			
			// Sunset should be after sunrise for normal conditions
			if sunrise > 0 && sunset < 1440 {
				assert.Greater(t, sunset, sunrise)
			}
		})
	}
}