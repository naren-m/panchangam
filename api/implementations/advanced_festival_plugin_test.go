package implementations

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

func TestAdvancedFestivalPlugin(t *testing.T) {
	ephemerisManager := &ephemeris.Manager{}
	plugin := NewAdvancedFestivalPlugin(ephemerisManager)

	info := plugin.GetInfo()
	if info.Name != "advanced_festival_plugin" {
		t.Errorf("Expected plugin name 'advanced_festival_plugin', got %s", info.Name)
	}

	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	regions := plugin.GetSupportedRegions()
	if len(regions) == 0 {
		t.Error("Expected supported regions, got none")
	}

	eventTypes := plugin.GetSupportedEventTypes()
	if len(eventTypes) == 0 {
		t.Error("Expected supported event types, got none")
	}
}

func TestAdvancedFestivalPluginDoesNotAdvertiseRegionalExtensionCapability(t *testing.T) {
	plugin := NewAdvancedFestivalPlugin(&ephemeris.Manager{})
	info := plugin.GetInfo()

	if _, ok := interface{}(plugin).(api.RegionalExtension); ok {
		t.Fatal("test expects AdvancedFestivalPlugin regional behavior to stay behind the EventPlugin contract")
	}

	for _, capability := range info.Capabilities {
		if capability == string(api.CapabilityRegional) {
			t.Fatal("AdvancedFestivalPlugin should not advertise regional capability without implementing RegionalExtension")
		}
	}
}

func TestAdvancedFestivalPluginRejectsEventsInRangeWhenDisabled(t *testing.T) {
	plugin := NewAdvancedFestivalPlugin(&ephemeris.Manager{})
	start := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)

	_, err := plugin.GetEventsInRange(context.Background(), start, start, api.Location{}, api.RegionGlobal)
	if err == nil {
		t.Fatal("Expected disabled plugin to return an error")
	}
	if err.Error() != "advanced festival plugin is not enabled" {
		t.Fatalf("Expected disabled plugin error, got %v", err)
	}
}

func TestAdvancedFestivalPluginRejectsEndBeforeStart(t *testing.T) {
	plugin := NewAdvancedFestivalPlugin(&ephemeris.Manager{})
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

func TestAdvancedFestivalPluginHelpers(t *testing.T) {
	plugin := NewAdvancedFestivalPlugin(&ephemeris.Manager{})

	tests := []struct {
		name     string
		date     time.Time
		isShukla bool
		want     string
	}{
		{
			name:     "june shukla",
			date:     time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC),
			isShukla: true,
			want:     "Apara",
		},
		{
			name:     "june krishna",
			date:     time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC),
			isShukla: false,
			want:     "Nirjala",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plugin.getEkadashiName(tt.date, tt.isShukla)
			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}

	date := time.Date(2024, time.November, 1, 0, 0, 0, 0, time.UTC)
	event := plugin.createDiwaliEvent(date, api.RegionGlobal)

	if event.Name != "Diwali" {
		t.Fatalf("expected Diwali event, got %s", event.Name)
	}
	if event.Type != api.EventTypeFestival {
		t.Fatalf("expected festival type, got %s", event.Type)
	}
	if event.Region != api.RegionGlobal {
		t.Fatalf("expected global region, got %s", event.Region)
	}
	if !event.StartTime.Equal(date) {
		t.Fatalf("expected start time %s, got %s", date, event.StartTime)
	}
	if !event.EndTime.Equal(date.Add(24 * time.Hour)) {
		t.Fatalf("expected one day event, got %s to %s", event.StartTime, event.EndTime)
	}
	if event.Metadata["astronomical_basis"] != "Kartik_Amavasya" {
		t.Fatalf("unexpected astronomical basis: %v", event.Metadata["astronomical_basis"])
	}
}
