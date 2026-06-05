package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetTithiSummaryReturnsPreciseTransitionTimes(t *testing.T) {
	baseTime := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	provider := newLinearTithiProvider(baseTime, 18, 12)
	manager := ephemeris.NewManager(provider, nil, ephemeris.NewNoOpCache())
	server := NewPanchangamServerWithDependencies(manager, DefaultConfig())

	resp, err := server.GetTithiSummary(context.Background(), &ppb.GetTithiSummaryRequest{
		At:                "2026-06-02T12:00:00Z",
		Latitude:          37.3382,
		Longitude:         -121.8863,
		Timezone:          "UTC",
		Region:            "California",
		CalculationMethod: "Drik",
		Locale:            "en",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Tithi)
	require.NotNil(t, resp.PanchaAnga)
	require.NotNil(t, resp.Calculation)
	require.NotNil(t, resp.Day)
	assert.Equal(t, "2026-06-02", resp.Date)
	assert.Equal(t, "Dwitiya", resp.Tithi.Name)
	assert.Equal(t, "Dvithiya", resp.Tithi.TraditionalName)
	assert.Equal(t, int32(2), resp.Tithi.Number)
	assert.Equal(t, "Shukla", resp.Tithi.Paksha)
	assert.Equal(t, int32(2), resp.Tithi.PakshaDay)
	assert.Equal(t, "Bhadra", resp.Tithi.Type)
	assert.Equal(t, "2026-06-02T00:00:00Z", resp.Tithi.StartTime)
	assert.Equal(t, "2026-06-03T00:00:00Z", resp.Tithi.EndTime)
	assert.Equal(t, "2026-06-02T12:00:00Z", resp.GeneratedAt)
	assert.Equal(t, "2026-06-03T00:00:00Z", resp.NextRefreshAt)
	assert.NotEmpty(t, resp.PanchaAnga.Nakshatra)
	assert.NotEmpty(t, resp.PanchaAnga.Yoga)
	assert.NotEmpty(t, resp.PanchaAnga.Karana)
	assert.NotEmpty(t, resp.PanchaAnga.Vara)
	assertHMTime(t, resp.Day.SunriseTime)
	assertHMTime(t, resp.Day.SunsetTime)
	require.NotNil(t, resp.Day.AbhijitMuhurta)
	assert.Equal(t, "Abhijit", resp.Day.AbhijitMuhurta.Name)
	assertHMTime(t, resp.Day.AbhijitMuhurta.StartTime)
	assertHMTime(t, resp.Day.AbhijitMuhurta.EndTime)
	assert.True(t, resp.Day.AbhijitMuhurta.Auspicious)
	assert.Equal(t, "UTC", resp.Calculation.Timezone)
	assert.Equal(t, "California", resp.Calculation.Region)
	assert.Equal(t, "Purnimanta", resp.Calculation.CalendarSystem)
	assert.Equal(t, "Drik", resp.Calculation.Method)
	assert.Equal(t, "en", resp.Calculation.Locale)
}

func TestGetTithiSummaryRejectsInvalidTimezone(t *testing.T) {
	server := NewPanchangamServerWithDependencies(ephemeris.NewManager(newLinearTithiProvider(time.Now(), 18, 12), nil, ephemeris.NewNoOpCache()), DefaultConfig())

	resp, err := server.GetTithiSummary(context.Background(), &ppb.GetTithiSummaryRequest{
		At:        "2026-06-02T12:00:00Z",
		Latitude:  37.3382,
		Longitude: -121.8863,
		Timezone:  "Invalid/Zone",
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), "invalid timezone")
}

func assertHMTime(t *testing.T, value string) {
	t.Helper()

	_, err := time.Parse("15:04", value)
	assert.NoError(t, err, "expected HH:MM time, got %q", value)
}

type linearTithiProvider struct {
	baseJD   ephemeris.JulianDay
	baseDiff float64
	rate     float64
}

func newLinearTithiProvider(baseTime time.Time, baseDiff, rate float64) *linearTithiProvider {
	return &linearTithiProvider{
		baseJD:   ephemeris.TimeToJulianDay(baseTime),
		baseDiff: baseDiff,
		rate:     rate,
	}
}

func (p *linearTithiProvider) GetPlanetaryPositions(ctx context.Context, jd ephemeris.JulianDay) (*ephemeris.PlanetaryPositions, error) {
	diff := p.diffAt(jd)
	position := ephemeris.Position{Longitude: diff, Speed: p.rate}
	return &ephemeris.PlanetaryPositions{
		JulianDay: jd,
		Sun:       ephemeris.Position{Longitude: 0},
		Moon:      position,
		Mercury:   position,
		Venus:     position,
		Mars:      position,
		Jupiter:   position,
		Saturn:    position,
		Uranus:    position,
		Neptune:   position,
		Pluto:     position,
	}, nil
}

func (p *linearTithiProvider) GetSunPosition(ctx context.Context, jd ephemeris.JulianDay) (*ephemeris.SolarPosition, error) {
	return &ephemeris.SolarPosition{JulianDay: jd, Longitude: 0}, nil
}

func (p *linearTithiProvider) GetMoonPosition(ctx context.Context, jd ephemeris.JulianDay) (*ephemeris.LunarPosition, error) {
	return &ephemeris.LunarPosition{JulianDay: jd, Longitude: p.diffAt(jd)}, nil
}

func (p *linearTithiProvider) IsAvailable(ctx context.Context) bool {
	return true
}

func (p *linearTithiProvider) GetDataRange() (ephemeris.JulianDay, ephemeris.JulianDay) {
	return p.baseJD - 3650, p.baseJD + 3650
}

func (p *linearTithiProvider) GetHealthStatus(ctx context.Context) (*ephemeris.HealthStatus, error) {
	return &ephemeris.HealthStatus{Available: true, Source: p.GetProviderName(), Version: p.GetVersion()}, nil
}

func (p *linearTithiProvider) GetProviderName() string {
	return "linear-tithi"
}

func (p *linearTithiProvider) GetVersion() string {
	return "test"
}

func (p *linearTithiProvider) Close() error {
	return nil
}

func (p *linearTithiProvider) diffAt(jd ephemeris.JulianDay) float64 {
	diff := p.baseDiff + (float64(jd-p.baseJD) * p.rate)
	for diff < 0 {
		diff += 360
	}
	for diff >= 360 {
		diff -= 360
	}
	return diff
}
