package implementations

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/naren-m/panchangam/api"
)

func TestInauspiciousEventsPlugin(t *testing.T) {
	plugin := NewInauspiciousEventsPlugin()

	info := plugin.GetInfo()
	if info.Name != "inauspicious_events_plugin" {
		t.Errorf("Expected plugin name 'inauspicious_events_plugin', got %s", info.Name)
	}

	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	events, err := plugin.GetEvents(context.Background(), testDate, testLocation, api.RegionGlobal)
	if err != nil {
		t.Fatalf("Failed to get events: %v", err)
	}

	if len(events) != 3 {
		t.Errorf("Expected 3 events (Rahu Kalam, Yamagandam, Gulika Kalam), got %d", len(events))
	}

	expectedTypes := map[api.EventType]bool{
		api.EventTypeRahukalam:   false,
		api.EventTypeYamagandam:  false,
		api.EventTypeGulikakalam: false,
	}

	for _, event := range events {
		if _, exists := expectedTypes[event.Type]; exists {
			expectedTypes[event.Type] = true
		}
	}

	for eventType, found := range expectedTypes {
		if !found {
			t.Errorf("Expected event type %s not found", eventType)
		}
	}

	endDate := testDate.AddDate(0, 0, 2)
	rangeEvents, err := plugin.GetEventsInRange(context.Background(), testDate, endDate, testLocation, api.RegionGlobal)
	if err != nil {
		t.Fatalf("Failed to get events in range: %v", err)
	}

	expectedEventCount := 3 * 3
	if len(rangeEvents) != expectedEventCount {
		t.Errorf("Expected %d events for 3 days, got %d", expectedEventCount, len(rangeEvents))
	}
}

func TestInauspiciousEventsSummary(t *testing.T) {
	plugin := NewInauspiciousEventsPlugin()
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	summary, err := plugin.GetEventSummary(context.Background(), testDate, testLocation, api.RegionGlobal)
	if err != nil {
		t.Fatalf("Failed to get event summary: %v", err)
	}

	if summary["date"] != testDate.Format("2006-01-02") {
		t.Fatalf("Unexpected summary date: %v", summary["date"])
	}
	if summary["total_events"] != 3 {
		t.Fatalf("Expected 3 summary events, got %v", summary["total_events"])
	}

	events, ok := summary["events"].([]map[string]interface{})
	if !ok {
		t.Fatalf("Expected events summary list, got %T", summary["events"])
	}
	if len(events) != 3 {
		t.Fatalf("Expected 3 event summaries, got %d", len(events))
	}
	if events[0]["duration"] == "" {
		t.Fatal("Expected formatted duration in event summary")
	}

	if got := plugin.formatDuration(90 * time.Minute); got != "1 hour(s) 30 minute(s)" {
		t.Fatalf("Unexpected hour duration: %s", got)
	}
	if got := plugin.formatDuration(45 * time.Minute); got != "45 minute(s)" {
		t.Fatalf("Unexpected minute duration: %s", got)
	}
}

func TestInauspiciousEventsPluginRejectsEventsInRangeWhenDisabled(t *testing.T) {
	plugin := NewInauspiciousEventsPlugin()

	_, err := plugin.GetEventsInRange(context.Background(), testDate, testDate, testLocation, api.RegionGlobal)
	if err == nil {
		t.Fatal("Expected disabled plugin to return an error")
	}
	if err.Error() != "inauspicious events plugin is not enabled" {
		t.Fatalf("Expected disabled plugin error, got %v", err)
	}
}

func TestInauspiciousEventsPluginRejectsInvalidTimezone(t *testing.T) {
	plugin := NewInauspiciousEventsPlugin()
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	location := testLocation
	location.Timezone = "Not/AZone"

	_, err = plugin.GetEvents(context.Background(), testDate, location, api.RegionGlobal)
	if err == nil {
		t.Fatal("Expected invalid timezone to return an error")
	}
	if !strings.Contains(err.Error(), "invalid timezone") {
		t.Fatalf("Expected invalid timezone error, got %v", err)
	}
}

func TestInauspiciousEventsPluginRejectsEndBeforeStart(t *testing.T) {
	plugin := NewInauspiciousEventsPlugin()
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	start := testDate.AddDate(0, 0, 1)
	_, err = plugin.GetEventsInRange(context.Background(), start, testDate, testLocation, api.RegionGlobal)
	if err == nil {
		t.Fatal("Expected invalid date range to return an error")
	}
	if !strings.Contains(err.Error(), "end date") || !strings.Contains(err.Error(), "before start date") {
		t.Fatalf("Expected clear invalid date range error, got %v", err)
	}
}

func TestRegionalNames(t *testing.T) {
	plugin := NewInauspiciousEventsPlugin()
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	testRegions := []api.Region{
		api.RegionTamilNadu,
		api.RegionKerala,
		api.RegionBengal,
		api.RegionGujarat,
		api.RegionMaha,
		api.RegionNorthIndia,
		api.RegionGlobal,
	}

	for _, region := range testRegions {
		events, err := plugin.GetEvents(context.Background(), testDate, testLocation, region)
		if err != nil {
			t.Fatalf("Failed to get events for region %s: %v", region, err)
		}

		for _, event := range events {
			if event.NameLocal == "" {
				t.Errorf("Local name missing for event %s in region %s", event.Name, region)
			}

			if event.Region != region {
				t.Errorf("Event region mismatch: expected %s, got %s", region, event.Region)
			}

			t.Logf("Region %s: %s (%s)", region, event.Name, event.NameLocal)
		}
	}
}

func BenchmarkEventCalculation(b *testing.B) {
	plugin := NewInauspiciousEventsPlugin()
	plugin.Initialize(context.Background(), map[string]interface{}{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := plugin.GetEvents(context.Background(), testDate, testLocation, api.RegionGlobal)
		if err != nil {
			b.Fatalf("Failed to get events: %v", err)
		}
	}
}
