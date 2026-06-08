package implementations

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/naren-m/panchangam/api"
)

func TestMuhurtaPlugin(t *testing.T) {
	plugin := NewMuhurtaPlugin()

	info := plugin.GetInfo()
	if info.Name != "muhurta_plugin" {
		t.Errorf("Expected plugin name 'muhurta_plugin', got %s", info.Name)
	}

	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled after initialization")
	}

	muhurtas, err := plugin.GetMuhurtas(context.Background(), testDate, testLocation, api.RegionGlobal)
	if err != nil {
		t.Fatalf("Failed to get muhurtas: %v", err)
	}

	if len(muhurtas) == 0 {
		t.Error("Expected at least some muhurtas, got none")
	}

	expectedMuhurtas := []string{"Rahu Kalam", "Yamagandam", "Gulika Kalam", "Abhijit Muhurta", "Brahma Muhurta"}
	muhurtaNames := make(map[string]bool)
	for _, m := range muhurtas {
		muhurtaNames[m.Name] = true
	}

	for _, expected := range expectedMuhurtas {
		if !muhurtaNames[expected] {
			t.Errorf("Expected muhurta '%s' not found", expected)
		}
	}

	activities := []string{"business", "education"}
	auspiciousTimes, err := plugin.FindAuspiciousTimes(context.Background(), testDate, testLocation, activities)
	if err != nil {
		t.Fatalf("Failed to find auspicious times: %v", err)
	}

	if len(auspiciousTimes) == 0 {
		t.Error("Expected to find at least one auspicious time")
	}

	testTime := time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC)
	isAuspicious, message, err := plugin.IsTimeAuspicious(context.Background(), testTime, testLocation, activities)
	if err != nil {
		t.Fatalf("Failed to check if time is auspicious: %v", err)
	}

	t.Logf("Time auspiciousness check: %v - %s", isAuspicious, message)
}

func TestMuhurtaPluginDoesNotAdvertiseUnsupportedEvents(t *testing.T) {
	plugin := NewMuhurtaPlugin()
	info := plugin.GetInfo()

	if _, ok := interface{}(plugin).(api.EventPlugin); ok {
		t.Fatal("test expects MuhurtaPlugin to only provide muhurta calculations")
	}

	for _, capability := range info.Capabilities {
		if capability == string(api.CapabilityEvent) {
			t.Fatal("MuhurtaPlugin should not advertise event capability without implementing EventPlugin")
		}
	}
}

func TestMuhurtaPluginRejectsInvalidTimezone(t *testing.T) {
	plugin := NewMuhurtaPlugin()
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	location := testLocation
	location.Timezone = "Not/AZone"

	_, err = plugin.GetMuhurtas(context.Background(), testDate, location, api.RegionGlobal)
	if err == nil {
		t.Fatal("Expected invalid timezone to return an error")
	}
	if !strings.Contains(err.Error(), "invalid timezone") {
		t.Fatalf("Expected invalid timezone error, got %v", err)
	}
}

func TestWeekdayCalculations(t *testing.T) {
	plugin := NewMuhurtaPlugin()
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	testDates := []struct {
		date    time.Time
		weekday string
	}{
		{time.Date(2024, 1, 14, 12, 0, 0, 0, time.UTC), "Sunday"},
		{time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), "Monday"},
		{time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC), "Tuesday"},
		{time.Date(2024, 1, 17, 12, 0, 0, 0, time.UTC), "Wednesday"},
		{time.Date(2024, 1, 18, 12, 0, 0, 0, time.UTC), "Thursday"},
		{time.Date(2024, 1, 19, 12, 0, 0, 0, time.UTC), "Friday"},
		{time.Date(2024, 1, 20, 12, 0, 0, 0, time.UTC), "Saturday"},
	}

	for _, testCase := range testDates {
		muhurtas, err := plugin.GetMuhurtas(context.Background(), testCase.date, testLocation, api.RegionGlobal)
		if err != nil {
			t.Fatalf("Failed to get muhurtas for %s: %v", testCase.weekday, err)
		}

		var rahuKalam *api.Muhurta
		for _, m := range muhurtas {
			if m.Name == "Rahu Kalam" {
				rahuKalam = &m
				break
			}
		}

		if rahuKalam == nil {
			t.Errorf("Rahu Kalam not found for %s", testCase.weekday)
			continue
		}

		if rahuKalam.StartTime.After(rahuKalam.EndTime) {
			t.Errorf("Rahu Kalam start time is after end time for %s", testCase.weekday)
		}

		duration := rahuKalam.EndTime.Sub(rahuKalam.StartTime)
		if duration <= 0 || duration > 3*time.Hour {
			t.Errorf("Rahu Kalam duration seems invalid for %s: %v", testCase.weekday, duration)
		}

		t.Logf("%s Rahu Kalam: %s to %s (%v)", testCase.weekday,
			rahuKalam.StartTime.Format("15:04"),
			rahuKalam.EndTime.Format("15:04"),
			duration)
	}
}

func TestMuhurtaActivitySuitability(t *testing.T) {
	plugin := NewMuhurtaPlugin()

	tests := []struct {
		name       string
		muhurta    api.Muhurta
		activities []string
		expected   bool
	}{
		{
			name: "matching purpose is suitable",
			muhurta: api.Muhurta{
				Purpose: []string{"business", "education"},
			},
			activities: []string{"business"},
			expected:   true,
		},
		{
			name: "all auspicious purpose matches any activity",
			muhurta: api.Muhurta{
				Purpose: []string{"all_auspicious_activities"},
			},
			activities: []string{"travel"},
			expected:   true,
		},
		{
			name: "empty purpose is not suitable",
			muhurta: api.Muhurta{
				Avoid: []string{"travel"},
			},
			activities: []string{"travel"},
			expected:   false,
		},
		{
			name: "avoid activity rejects when purpose does not match",
			muhurta: api.Muhurta{
				Purpose: []string{"education"},
				Avoid:   []string{"travel"},
			},
			activities: []string{"travel"},
			expected:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := plugin.isActivitySuitable(test.muhurta, test.activities)
			if got != test.expected {
				t.Fatalf("expected %t, got %t", test.expected, got)
			}
		})
	}
}

func BenchmarkMuhurtaCalculation(b *testing.B) {
	plugin := NewMuhurtaPlugin()
	plugin.Initialize(context.Background(), map[string]interface{}{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := plugin.GetMuhurtas(context.Background(), testDate, testLocation, api.RegionGlobal)
		if err != nil {
			b.Fatalf("Failed to get muhurtas: %v", err)
		}
	}
}
