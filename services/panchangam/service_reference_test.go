package panchangam

import (
	"context"
	"strings"
	"testing"
	"time"

	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/require"
)

func TestSanJoseJune42026MatchesReferencePanchangam(t *testing.T) {
	server := NewPanchangamServer()

	resp, err := server.Get(context.Background(), &ppb.GetPanchangamRequest{
		Date:              "2026-06-04",
		Latitude:          37.3382,
		Longitude:         -121.8863,
		Timezone:          "America/Los_Angeles",
		Region:            "California",
		CalculationMethod: "Drik",
		Locale:            "en",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.PanchangamData)

	data := resp.PanchangamData

	require.Contains(t, data.Tithi, "Chathurthi")
	require.Contains(t, data.Nakshatra, "Uttara Ashadha")
	require.Contains(t, data.Yoga, "Brahma")
	require.Contains(t, data.Karana, "Balava")
	require.Equal(t, "Guruvara", panchangamEventValue(data.Events, "VARA", "Vara: "))
	require.Equal(t, "Makara", panchangamEventValue(data.Events, "RAASI", "Raasi: "))
	requirePanchangamClockBetween(t, data.Events, "TITHI_END", "10:45:00", "11:15:00")

	abhijitStart := panchangamEventTime(t, data.Events, "ABHIJIT_MUHURTA")
	require.False(t, abhijitStart.Before(panchangamClockTime(t, "12:00:00")), "Abhijit Muhurta should be near midday")
	require.False(t, abhijitStart.After(panchangamClockTime(t, "13:30:00")), "Abhijit Muhurta should be near midday")
}

func TestSanJoseJune42026CurrentTimeMatchesReferencePanchangam(t *testing.T) {
	server := NewPanchangamServer()

	resp, err := server.Get(context.Background(), &ppb.GetPanchangamRequest{
		Date:              "2026-06-04T15:28:04-07:00",
		Latitude:          37.3382,
		Longitude:         -121.8863,
		Timezone:          "America/Los_Angeles",
		Region:            "California",
		CalculationMethod: "Drik",
		Locale:            "en",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.PanchangamData)

	data := resp.PanchangamData
	require.Equal(t, "2026-06-04", data.Date)
	require.Contains(t, data.Tithi, "Panchami")
	require.Contains(t, data.Nakshatra, "Shravana")
	require.Contains(t, data.Yoga, "Brahma")
	require.Contains(t, data.Karana, "Kaulava")
	require.Equal(t, "Guruvara", panchangamEventValue(data.Events, "VARA", "Vara: "))
	require.Equal(t, "Makara", panchangamEventValue(data.Events, "RAASI", "Raasi: "))
	requirePanchangamClockBetween(t, data.Events, "TITHI_START", "10:45:00", "11:15:00")
	requirePanchangamClockBetween(t, data.Events, "TITHI_END", "12:35:00", "13:05:00")
}

func panchangamEventValue(events []*ppb.PanchangamEvent, eventType, prefix string) string {
	for _, event := range events {
		if event != nil && event.EventType == eventType {
			return strings.TrimPrefix(event.Name, prefix)
		}
	}
	return ""
}

func panchangamEventTime(t *testing.T, events []*ppb.PanchangamEvent, eventType string) time.Time {
	t.Helper()

	for _, event := range events {
		if event == nil || event.EventType != eventType {
			continue
		}
		return panchangamClockTime(t, event.Time)
	}

	t.Fatalf("missing %s event", eventType)
	return time.Time{}
}

func requirePanchangamClockBetween(t *testing.T, events []*ppb.PanchangamEvent, eventType, minClock, maxClock string) {
	t.Helper()

	actual := panchangamEventTime(t, events, eventType)
	require.False(t, actual.Before(panchangamClockTime(t, minClock)), "%s should be after %s, got %s", eventType, minClock, actual.Format("15:04:05"))
	require.False(t, actual.After(panchangamClockTime(t, maxClock)), "%s should be before %s, got %s", eventType, maxClock, actual.Format("15:04:05"))
}

func panchangamClockTime(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse("15:04:05", value)
	require.NoError(t, err)
	return parsed
}
