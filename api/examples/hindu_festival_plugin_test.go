package examples

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/naren-m/panchangam/api"
)

func TestHinduFestivalPluginDoesNotAdvertiseUnsupportedMuhurtas(t *testing.T) {
	plugin := NewHinduFestivalPlugin()
	info := plugin.GetInfo()

	if _, ok := interface{}(plugin).(api.MuhurtaPlugin); ok {
		t.Fatal("test expects HinduFestivalPlugin to only provide event calculations")
	}

	for _, capability := range info.Capabilities {
		if capability == string(api.CapabilityMuhurta) {
			t.Fatal("HinduFestivalPlugin should not advertise muhurta capability without implementing MuhurtaPlugin")
		}
	}
}

func TestHinduFestivalPluginRejectsGetEventsWhenDisabled(t *testing.T) {
	plugin := NewHinduFestivalPlugin()

	_, err := plugin.GetEvents(context.Background(), time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC), api.Location{}, api.RegionGlobal)
	if err == nil {
		t.Fatal("Expected disabled plugin to return an error")
	}
	if !strings.Contains(err.Error(), "hindu festival plugin is not enabled") {
		t.Fatalf("Expected disabled plugin error, got %v", err)
	}
}

func TestHinduFestivalPluginRejectsGetEventsInRangeWhenDisabled(t *testing.T) {
	plugin := NewHinduFestivalPlugin()
	start := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, time.June, 2, 0, 0, 0, 0, time.UTC)

	_, err := plugin.GetEventsInRange(context.Background(), start, end, api.Location{}, api.RegionGlobal)
	if err == nil {
		t.Fatal("Expected disabled plugin to return an error")
	}
	if !strings.Contains(err.Error(), "hindu festival plugin is not enabled") {
		t.Fatalf("Expected disabled plugin error, got %v", err)
	}
}

func TestHinduFestivalPluginRejectsGetEventsInRangeWhenDisabledBeforeValidatingDates(t *testing.T) {
	plugin := NewHinduFestivalPlugin()
	start := time.Date(2024, time.June, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)

	_, err := plugin.GetEventsInRange(context.Background(), start, end, api.Location{}, api.RegionGlobal)
	if err == nil {
		t.Fatal("Expected disabled plugin to return an error")
	}
	if err.Error() != "hindu festival plugin is not enabled" {
		t.Fatalf("Expected disabled plugin error, got %v", err)
	}
}

func TestHinduFestivalPluginRejectsEndBeforeStart(t *testing.T) {
	plugin := NewHinduFestivalPlugin()
	if err := plugin.Initialize(context.Background(), map[string]interface{}{}); err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	start := time.Date(2024, time.June, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)

	_, err := plugin.GetEventsInRange(context.Background(), start, end, api.Location{}, api.RegionGlobal)
	if err == nil {
		t.Fatal("Expected invalid date range to return an error")
	}
	if !strings.Contains(err.Error(), "end date") || !strings.Contains(err.Error(), "before start date") {
		t.Fatalf("Expected clear invalid date range error, got %v", err)
	}
}
