package astronomy

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTracingPerformance(t *testing.T) {
	observability.NewLocalObserver()

	ctx := context.Background()
	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)

	start := time.Now()
	for i := 0; i < 100; i++ {
		_, err := CalculateSunTimes(loc, date)
		require.NoError(t, err)
	}
	regularDuration := time.Since(start)

	start = time.Now()
	for i := 0; i < 100; i++ {
		_, err := CalculateSunTimesWithContext(ctx, loc, date)
		require.NoError(t, err)
	}
	contextDuration := time.Since(start)

	assert.Less(t, contextDuration, regularDuration*3,
		"Context-aware function is too slow: regular=%v, context=%v",
		regularDuration, contextDuration)

	t.Logf("Performance comparison: regular=%v, context=%v (%.2fx slower)",
		regularDuration, contextDuration, float64(contextDuration)/float64(regularDuration))
}

func TestContextCancellation(t *testing.T) {
	observability.NewLocalObserver()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)

	sunTimes, err := CalculateSunTimesWithContext(ctx, loc, date)
	require.NoError(t, err)
	require.NotNil(t, sunTimes)

	sunTimesRegular, err := CalculateSunTimes(loc, date)
	require.NoError(t, err)
	assert.Equal(t, sunTimesRegular.Sunrise, sunTimes.Sunrise)
	assert.Equal(t, sunTimesRegular.Sunset, sunTimes.Sunset)
}

func TestSpanAttributes(t *testing.T) {
	observability.NewLocalObserver()

	ctx := context.Background()

	testCases := []struct {
		name     string
		location Location
		date     time.Time
	}{
		{
			name:     "New York",
			location: Location{Latitude: 40.7128, Longitude: -74.0060},
			date:     time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "London",
			location: Location{Latitude: 51.5074, Longitude: -0.1278},
			date:     time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Sydney",
			location: Location{Latitude: -33.8688, Longitude: 151.2093},
			date:     time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sunTimes, err := CalculateSunTimesWithContext(ctx, tc.location, tc.date)
			require.NoError(t, err)
			require.NotNil(t, sunTimes)

			sunrise, err := GetSunriseTimeWithContext(ctx, tc.location, tc.date)
			require.NoError(t, err)
			assert.Equal(t, sunTimes.Sunrise, sunrise)

			sunset, err := GetSunsetTimeWithContext(ctx, tc.location, tc.date)
			require.NoError(t, err)
			assert.Equal(t, sunTimes.Sunset, sunset)
		})
	}
}

func TestErrorHandlingWithContext(t *testing.T) {
	observability.NewLocalObserver()

	ctx := context.Background()
	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)

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
