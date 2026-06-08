package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockRealServiceFlow(t *testing.T, ctx context.Context, req *ppb.GetPanchangamRequest) {
	t.Helper()

	date, err := time.Parse("2006-01-02", req.Date)
	require.NoError(t, err, "Date should parse correctly")

	location := astronomy.Location{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, date)
	require.NoError(t, err, "Sunrise/sunset calculation should work")

	mockTithi := mockCalculateTithi(t, ctx, date)
	mockNakshatra := mockCalculateNakshatra(t, ctx, date)
	mockYoga := mockCalculateYoga(t, ctx, date)
	mockKarana := mockCalculateKarana(t, ctx, date)
	mockVara := mockCalculateVara(t, ctx, date)

	response := &ppb.PanchangamData{
		Date:        req.Date,
		Tithi:       formatTithiForResponse(mockTithi),
		Nakshatra:   formatNakshatraForResponse(mockNakshatra),
		Yoga:        formatYogaForResponse(mockYoga),
		Karana:      formatKaranaForResponse(mockKarana),
		SunriseTime: sunTimes.Sunrise.Format("15:04:05"),
		SunsetTime:  sunTimes.Sunset.Format("15:04:05"),
		Events:      buildEventsForResponse(mockVara, mockKarana),
	}

	assert.Equal(t, req.Date, response.Date, "Date should match")
	assert.NotEmpty(t, response.Tithi, "Tithi should be calculated")
	assert.NotEmpty(t, response.Nakshatra, "Nakshatra should be calculated")
	assert.NotEmpty(t, response.Yoga, "Yoga should be calculated")
	assert.NotEmpty(t, response.Karana, "Karana should be calculated")
	assert.NotEmpty(t, response.SunriseTime, "Sunrise should be calculated")
	assert.NotEmpty(t, response.SunsetTime, "Sunset should be calculated")
	assert.NotNil(t, response.Events, "Events should be generated")

	t.Logf("Mocked real service flow successfully")
	t.Logf("Date: %s", response.Date)
	t.Logf("Tithi: %s", response.Tithi)
	t.Logf("Nakshatra: %s", response.Nakshatra)
	t.Logf("Yoga: %s", response.Yoga)
	t.Logf("Karana: %s", response.Karana)
	t.Logf("Sunrise: %s", response.SunriseTime)
	t.Logf("Sunset: %s", response.SunsetTime)
	t.Logf("Events: %d", len(response.Events))
}

func mockCalculateTithi(t *testing.T, ctx context.Context, date time.Time) *astronomy.TithiInfo {
	return &astronomy.TithiInfo{
		Number:      15,
		Name:        "Purnima",
		Type:        astronomy.TithiTypePurna,
		StartTime:   date,
		EndTime:     date.Add(24 * time.Hour),
		Duration:    24.0,
		IsShukla:    true,
		MoonSunDiff: 180.0,
	}
}

func mockCalculateNakshatra(t *testing.T, ctx context.Context, date time.Time) *astronomy.NakshatraInfo {
	return &astronomy.NakshatraInfo{
		Number:        13,
		Name:          "Hasta",
		Deity:         "Savitar",
		PlanetaryLord: "Moon",
		Symbol:        "Hand",
		Pada:          2,
		StartTime:     date,
		EndTime:       date.Add(24 * time.Hour),
		Duration:      24.0,
		MoonLongitude: 166.5,
	}
}

func mockCalculateYoga(t *testing.T, ctx context.Context, date time.Time) *astronomy.YogaInfo {
	return &astronomy.YogaInfo{
		Number:        14,
		Name:          "Vishkambha",
		Quality:       astronomy.YogaQualityAuspicious,
		StartTime:     date,
		EndTime:       date.Add(24 * time.Hour),
		Duration:      24.0,
		CombinedValue: 180.0,
		Description:   "Auspicious for new beginnings",
	}
}

func mockCalculateKarana(t *testing.T, ctx context.Context, date time.Time) *astronomy.KaranaInfo {
	return &astronomy.KaranaInfo{
		Number:      7,
		Name:        "Vanija",
		Type:        astronomy.KaranaTypeMovable,
		Description: "Merchant, good for business and trade",
		IsVishti:    false,
		StartTime:   date,
		EndTime:     date.Add(12 * time.Hour),
		Duration:    12.0,
		MoonSunDiff: 84.0,
		TithiNumber: 15,
		HalfTithi:   1,
	}
}

func mockCalculateVara(t *testing.T, ctx context.Context, date time.Time) *astronomy.VaraInfo {
	return &astronomy.VaraInfo{
		Number:        2,
		Name:          "Somavara",
		PlanetaryLord: "Moon",
		Quality:       "Peaceful",
		Color:         "White",
		Deity:         "Soma",
		StartTime:     date,
		EndTime:       date.Add(24 * time.Hour),
		Duration:      24.0,
		GregorianDay:  "Monday",
		IsAuspicious:  true,
		CurrentHora:   8,
		HoraPlanet:    "Moon",
	}
}

func formatTithiForResponse(tithi *astronomy.TithiInfo) string {
	return tithi.Name + " Tithi"
}

func formatNakshatraForResponse(nakshatra *astronomy.NakshatraInfo) string {
	return nakshatra.Name + " Nakshatra"
}

func formatYogaForResponse(yoga *astronomy.YogaInfo) string {
	return yoga.Name + " Yoga"
}

func formatKaranaForResponse(karana *astronomy.KaranaInfo) string {
	return karana.Name + " Karana"
}

func buildEventsForResponse(vara *astronomy.VaraInfo, karana *astronomy.KaranaInfo) []*ppb.PanchangamEvent {
	events := []*ppb.PanchangamEvent{}

	if vara.IsAuspicious {
		events = append(events, &ppb.PanchangamEvent{
			Name:      vara.Name + " - Auspicious day",
			Time:      vara.StartTime.Format("15:04:05"),
			EventType: "VARA_QUALITY",
		})
	}

	if karana.IsVishti {
		events = append(events, &ppb.PanchangamEvent{
			Name:      "Vishti Karana - Avoid important activities",
			Time:      karana.StartTime.Format("15:04:05"),
			EventType: "KARANA_WARNING",
		})
	} else {
		events = append(events, &ppb.PanchangamEvent{
			Name:      karana.Name + " Karana - " + karana.Description,
			Time:      karana.StartTime.Format("15:04:05"),
			EventType: "KARANA_INFO",
		})
	}

	return events
}
