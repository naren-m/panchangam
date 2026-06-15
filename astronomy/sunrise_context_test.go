package astronomy

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateSunTimesWithContext(t *testing.T) {
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

			assert.NotZero(t, sunTimes.Sunrise)
			assert.NotZero(t, sunTimes.Sunset)

			t.Logf("Location: %s, Sunrise: %s, Sunset: %s",
				tt.name, sunTimes.Sunrise.Format("15:04:05"), sunTimes.Sunset.Format("15:04:05"))
		})
	}
}

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

	sunriseRegular, err := GetSunriseTime(loc, date)
	assert.NoError(t, err)
	assert.Equal(t, sunriseRegular, sunrise)
}

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

	sunsetRegular, err := GetSunsetTime(loc, date)
	assert.NoError(t, err)
	assert.Equal(t, sunsetRegular, sunset)
}

func TestSolarPositionWithContext(t *testing.T) {
	observability.NewLocalObserver()

	ctx := context.Background()
	jd := julianDate(time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC))

	eqTime, decl := solarPositionWithContext(ctx, jd)

	assert.Greater(t, eqTime, -20.0)
	assert.Less(t, eqTime, 20.0)
	assert.Greater(t, decl, -23.44*DegToRad)
	assert.Less(t, decl, 23.44*DegToRad)

	eqTimeRegular, declRegular := solarPosition(jd)
	assert.Equal(t, eqTimeRegular, eqTime)
	assert.Equal(t, declRegular, decl)
}

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

			assert.GreaterOrEqual(t, sunrise, 0.0)
			assert.LessOrEqual(t, sunrise, 1440.0)
			assert.GreaterOrEqual(t, sunset, 0.0)
			assert.LessOrEqual(t, sunset, 1440.0)

			sunriseRegular, sunsetRegular := calculateRiseSet(tt.latitude, tt.longitude, tt.jd, tt.eqTime, tt.decl)
			assert.Equal(t, sunriseRegular, sunrise)
			assert.Equal(t, sunsetRegular, sunset)
		})
	}
}
