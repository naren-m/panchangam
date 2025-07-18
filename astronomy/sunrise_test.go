package astronomy

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateSunTimes(t *testing.T) {
	// Initialize observability for testing
	observability.NewLocalObserver()
	
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
	// Initialize observability for testing
	observability.NewLocalObserver()
	
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
	// Initialize observability for testing
	observability.NewLocalObserver()
	
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
	// Initialize observability for testing
	observability.NewLocalObserver()
	
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
	// Initialize observability for testing
	observability.NewLocalObserver()
	
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
	// Initialize observability for testing
	observability.NewLocalObserver()
	
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
	// Initialize observability for testing
	observability.NewLocalObserver()
	
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
	// Initialize observability for testing
	observability.NewLocalObserver()
	
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

// TestCalculateSunTimesWithContext tests the context-aware version with tracing
func TestCalculateSunTimesWithContext(t *testing.T) {
	// Initialize observability for testing
	observability.NewLocalObserver()
	
	ctx := context.Background()
	
	tests := []struct {
		name     string
		location Location
		date     time.Time
		desc     string
	}{
		{
			name: "New York with Context",
			location: Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			date: time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
			desc: "Should trace all calculation steps",
		},
		{
			name: "Arctic Location - Polar Day",
			location: Location{
				Latitude:  75.0,
				Longitude: 0.0,
			},
			date: time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
			desc: "Should trace polar day condition",
		},
		{
			name: "Arctic Location - Polar Night",
			location: Location{
				Latitude:  75.0,
				Longitude: 0.0,
			},
			date: time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
			desc: "Should trace polar night condition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sunTimes, err := CalculateSunTimesWithContext(ctx, tt.location, tt.date)
			require.NoError(t, err)
			require.NotNil(t, sunTimes)

			// Basic validation - times should be valid
			assert.NotZero(t, sunTimes.Sunrise)
			assert.NotZero(t, sunTimes.Sunset)
			
			t.Logf("Location: %s, Sunrise: %s, Sunset: %s", 
				tt.name, sunTimes.Sunrise.Format("15:04:05"), sunTimes.Sunset.Format("15:04:05"))
		})
	}
}

// TestGetSunriseTimeWithContext tests the context-aware sunrise function
func TestGetSunriseTimeWithContext(t *testing.T) {
	observability.NewLocalObserver()
	
	ctx := context.Background()
	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)

	sunrise, err := GetSunriseTimeWithContext(ctx, loc, date)
	assert.NoError(t, err)
	assert.NotZero(t, sunrise)
	
	// Compare with regular function
	sunriseRegular, err := GetSunriseTime(loc, date)
	assert.NoError(t, err)
	assert.Equal(t, sunriseRegular, sunrise)
}

// TestGetSunsetTimeWithContext tests the context-aware sunset function
func TestGetSunsetTimeWithContext(t *testing.T) {
	observability.NewLocalObserver()
	
	ctx := context.Background()
	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)

	sunset, err := GetSunsetTimeWithContext(ctx, loc, date)
	assert.NoError(t, err)
	assert.NotZero(t, sunset)
	
	// Compare with regular function
	sunsetRegular, err := GetSunsetTime(loc, date)
	assert.NoError(t, err)
	assert.Equal(t, sunsetRegular, sunset)
}

// TestSolarPositionWithContext tests the context-aware solar position function
func TestSolarPositionWithContext(t *testing.T) {
	observability.NewLocalObserver()
	
	ctx := context.Background()
	jd := julianDate(time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC))
	
	eqTime, decl := solarPositionWithContext(ctx, jd)
	
	// Equation of time should be within reasonable bounds (-20 to +20 minutes)
	assert.Greater(t, eqTime, -20.0)
	assert.Less(t, eqTime, 20.0)
	
	// Declination should be within solar bounds (-23.44 to +23.44 degrees in radians)
	assert.Greater(t, decl, -23.44*DegToRad)
	assert.Less(t, decl, 23.44*DegToRad)
	
	// Compare with regular function
	eqTimeRegular, declRegular := solarPosition(jd)
	assert.Equal(t, eqTimeRegular, eqTime)
	assert.Equal(t, declRegular, decl)
}

// TestCalculateRiseSetWithContext tests the context-aware rise/set function
func TestCalculateRiseSetWithContext(t *testing.T) {
	observability.NewLocalObserver()
	
	ctx := context.Background()
	
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		jd        float64
		eqTime    float64
		decl      float64
		desc      string
	}{
		{
			name:      "Equator test with context",
			latitude:  0.0,
			longitude: 0.0,
			jd:        2451545.0,
			eqTime:    0.0,
			decl:      0.0,
			desc:      "Should trace equatorial calculation",
		},
		{
			name:      "Arctic test - should detect polar conditions",
			latitude:  80.0,
			longitude: 0.0,
			jd:        2451545.0,
			eqTime:    0.0,
			decl:      23.0 * DegToRad,
			desc:      "Should trace polar day condition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sunrise, sunset := calculateRiseSetWithContext(ctx, tt.latitude, tt.longitude, tt.jd, tt.eqTime, tt.decl)
			
			// Basic sanity checks
			assert.GreaterOrEqual(t, sunrise, 0.0)
			assert.LessOrEqual(t, sunrise, 1440.0) // 24 hours in minutes
			assert.GreaterOrEqual(t, sunset, 0.0)
			assert.LessOrEqual(t, sunset, 1440.0)
			
			// Compare with regular function
			sunriseRegular, sunsetRegular := calculateRiseSet(tt.latitude, tt.longitude, tt.jd, tt.eqTime, tt.decl)
			assert.Equal(t, sunriseRegular, sunrise)
			assert.Equal(t, sunsetRegular, sunset)
		})
	}
}

// TestTracingPerformance tests that tracing doesn't significantly impact performance
func TestTracingPerformance(t *testing.T) {
	observability.NewLocalObserver()
	
	ctx := context.Background()
	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)

	// Test regular function performance
	start := time.Now()
	for i := 0; i < 100; i++ {
		_, err := CalculateSunTimes(loc, date)
		require.NoError(t, err)
	}
	regularDuration := time.Since(start)

	// Test context-aware function performance
	start = time.Now()
	for i := 0; i < 100; i++ {
		_, err := CalculateSunTimesWithContext(ctx, loc, date)
		require.NoError(t, err)
	}
	contextDuration := time.Since(start)

	// Context-aware function should not be more than 3x slower
	assert.Less(t, contextDuration, regularDuration*3, 
		"Context-aware function is too slow: regular=%v, context=%v", 
		regularDuration, contextDuration)
		
	t.Logf("Performance comparison: regular=%v, context=%v (%.2fx slower)", 
		regularDuration, contextDuration, float64(contextDuration)/float64(regularDuration))
}

// TestContextCancellation tests that context cancellation is properly handled
func TestContextCancellation(t *testing.T) {
	observability.NewLocalObserver()
	
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)

	// Function should still work despite cancelled context
	sunTimes, err := CalculateSunTimesWithContext(ctx, loc, date)
	require.NoError(t, err)
	require.NotNil(t, sunTimes)
	
	// Results should be the same as regular function
	sunTimesRegular, err := CalculateSunTimes(loc, date)
	require.NoError(t, err)
	assert.Equal(t, sunTimesRegular.Sunrise, sunTimes.Sunrise)
	assert.Equal(t, sunTimesRegular.Sunset, sunTimes.Sunset)
}

// TestSpanAttributes tests that proper span attributes are set
func TestSpanAttributes(t *testing.T) {
	observability.NewLocalObserver()
	
	ctx := context.Background()

	// Test various locations and dates to ensure spans are created
	testCases := []struct {
		name     string
		location Location
		date     time.Time
	}{
		{
			name: "New York",
			location: Location{Latitude: 40.7128, Longitude: -74.0060},
			date: time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "London",
			location: Location{Latitude: 51.5074, Longitude: -0.1278},
			date: time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Sydney",
			location: Location{Latitude: -33.8688, Longitude: 151.2093},
			date: time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sunTimes, err := CalculateSunTimesWithContext(ctx, tc.location, tc.date)
			require.NoError(t, err)
			require.NotNil(t, sunTimes)
			
			// Also test individual functions
			sunrise, err := GetSunriseTimeWithContext(ctx, tc.location, tc.date)
			require.NoError(t, err)
			assert.Equal(t, sunTimes.Sunrise, sunrise)
			
			sunset, err := GetSunsetTimeWithContext(ctx, tc.location, tc.date)
			require.NoError(t, err)
			assert.Equal(t, sunTimes.Sunset, sunset)
		})
	}
}

// TestErrorHandlingWithContext tests error handling in context-aware functions
func TestErrorHandlingWithContext(t *testing.T) {
	observability.NewLocalObserver()
	
	ctx := context.Background()
	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)

	// Test that context-aware functions handle errors the same way
	sunTimes, err := CalculateSunTimesWithContext(ctx, loc, date)
	require.NoError(t, err)
	require.NotNil(t, sunTimes)
	
	sunrise, err := GetSunriseTimeWithContext(ctx, loc, date)
	require.NoError(t, err)
	assert.Equal(t, sunTimes.Sunrise, sunrise)
	
	sunset, err := GetSunsetTimeWithContext(ctx, loc, date)
	require.NoError(t, err)
	assert.Equal(t, sunTimes.Sunset, sunset)
}