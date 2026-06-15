package gateway

import (
	"testing"
	"time"

	"github.com/naren-m/panchangam/cache"
	"github.com/stretchr/testify/assert"
)

func TestWriteCacheHealthResponseReturnsWriteError(t *testing.T) {
	err := writeCacheHealthResponse(errorResponseWriter{}, time.Unix(0, 0).UTC())

	if err == nil {
		t.Fatal("expected write error")
	}
	assert.Contains(t, err.Error(), "write failed")
}

func TestConvertCacheToResponsePreservesEvents(t *testing.T) {
	response := convertCacheToResponse(&cache.PanchangamCacheData{
		Date:        "2024-01-15",
		Tithi:       "Purnima",
		Nakshatra:   "Hasta",
		Yoga:        "Siddhi",
		Karana:      "Vishti",
		SunriseTime: "06:45",
		SunsetTime:  "18:10",
		Events: []cache.Event{
			{Name: "Rahu Kalam", Time: "07:30", EventType: "RAHU_KALAM"},
		},
	})

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		t.Fatalf("expected response map, got %T", response)
	}
	if got := responseMap["date"]; got != "2024-01-15" {
		t.Fatalf("expected date to be preserved, got %v", got)
	}

	events, ok := responseMap["events"].([]map[string]interface{})
	if !ok {
		t.Fatalf("expected events list, got %T", responseMap["events"])
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if got := events[0]["name"]; got != "Rahu Kalam" {
		t.Fatalf("expected event name to be preserved, got %v", got)
	}
	if got := events[0]["time"]; got != "07:30" {
		t.Fatalf("expected event time to be preserved, got %v", got)
	}
	if got := events[0]["event_type"]; got != "RAHU_KALAM" {
		t.Fatalf("expected event type to be preserved, got %v", got)
	}
}
