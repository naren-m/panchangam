package implementations

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

// Test location for Mumbai, India
var testLocation = api.Location{
	Latitude:  19.0760,
	Longitude: 72.8777,
	Timezone:  "Asia/Kolkata",
	Name:      "Mumbai",
}

// Test date
var testDate = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestMuhurtaPlugin(t *testing.T) {
	plugin := NewMuhurtaPlugin()

	// Test plugin info
	info := plugin.GetInfo()
	if info.Name != "muhurta_plugin" {
		t.Errorf("Expected plugin name 'muhurta_plugin', got %s", info.Name)
	}

	// Test initialization
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled after initialization")
	}

	// Test muhurta calculation
	muhurtas, err := plugin.GetMuhurtas(context.Background(), testDate, testLocation, api.RegionGlobal)
	if err != nil {
		t.Fatalf("Failed to get muhurtas: %v", err)
	}

	if len(muhurtas) == 0 {
		t.Error("Expected at least some muhurtas, got none")
	}

	// Verify we have the basic muhurtas
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

	// Test auspicious time finding
	activities := []string{"business", "education"}
	auspiciousTimes, err := plugin.FindAuspiciousTimes(context.Background(), testDate, testLocation, activities)
	if err != nil {
		t.Fatalf("Failed to find auspicious times: %v", err)
	}

	// Should find at least Abhijit Muhurta which is good for all activities
	if len(auspiciousTimes) == 0 {
		t.Error("Expected to find at least one auspicious time")
	}

	// Test time checking
	testTime := time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC) // During Abhijit Muhurta
	isAuspicious, message, err := plugin.IsTimeAuspicious(context.Background(), testTime, testLocation, activities)
	if err != nil {
		t.Fatalf("Failed to check if time is auspicious: %v", err)
	}

	t.Logf("Time auspiciousness check: %v - %s", isAuspicious, message)
}

func TestInauspiciousEventsPlugin(t *testing.T) {
	plugin := NewInauspiciousEventsPlugin()

	// Test plugin info
	info := plugin.GetInfo()
	if info.Name != "inauspicious_events_plugin" {
		t.Errorf("Expected plugin name 'inauspicious_events_plugin', got %s", info.Name)
	}

	// Test initialization
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Test event calculation
	events, err := plugin.GetEvents(context.Background(), testDate, testLocation, api.RegionGlobal)
	if err != nil {
		t.Fatalf("Failed to get events: %v", err)
	}

	if len(events) != 3 {
		t.Errorf("Expected 3 events (Rahu Kalam, Yamagandam, Gulika Kalam), got %d", len(events))
	}

	// Verify event types
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

	// Test event range calculation
	endDate := testDate.AddDate(0, 0, 2) // 3 days
	rangeEvents, err := plugin.GetEventsInRange(context.Background(), testDate, endDate, testLocation, api.RegionGlobal)
	if err != nil {
		t.Fatalf("Failed to get events in range: %v", err)
	}

	expectedEventCount := 3 * 3 // 3 events per day for 3 days
	if len(rangeEvents) != expectedEventCount {
		t.Errorf("Expected %d events for 3 days, got %d", expectedEventCount, len(rangeEvents))
	}
}

func TestAdvancedFestivalPlugin(t *testing.T) {
	// Skip if no ephemeris manager available
	ephemerisManager := &ephemeris.Manager{} // Mock manager for testing
	plugin := NewAdvancedFestivalPlugin(ephemerisManager)

	// Test plugin info
	info := plugin.GetInfo()
	if info.Name != "advanced_festival_plugin" {
		t.Errorf("Expected plugin name 'advanced_festival_plugin', got %s", info.Name)
	}

	// Test initialization
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Test supported regions
	regions := plugin.GetSupportedRegions()
	if len(regions) == 0 {
		t.Error("Expected supported regions, got none")
	}

	// Test supported event types
	eventTypes := plugin.GetSupportedEventTypes()
	if len(eventTypes) == 0 {
		t.Error("Expected supported event types, got none")
	}

	// Note: Actual festival calculation tests would require a working ephemeris manager
	// For now, we just test the basic plugin structure
}

func TestCalendarSystemPlugin(t *testing.T) {
	// Skip if no ephemeris manager available
	ephemerisManager := &ephemeris.Manager{} // Mock manager for testing
	plugin := NewCalendarSystemPlugin(ephemerisManager)

	// Test plugin info
	info := plugin.GetInfo()
	if info.Name != "calendar_system_plugin" {
		t.Errorf("Expected plugin name 'calendar_system_plugin', got %s", info.Name)
	}

	// Test initialization
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Test supported regions
	regions := plugin.GetSupportedRegions()
	if len(regions) == 0 {
		t.Error("Expected supported regions, got none")
	}

	// Test supported methods
	methods := plugin.GetSupportedMethods()
	expectedMethods := []api.CalculationMethod{api.MethodDrik, api.MethodVakya, api.MethodAuto}
	if len(methods) != len(expectedMethods) {
		t.Errorf("Expected %d supported methods, got %d", len(expectedMethods), len(methods))
	}
}

func TestCalculationMethodPlugin(t *testing.T) {
	// Skip if no ephemeris manager available
	ephemerisManager := &ephemeris.Manager{} // Mock manager for testing
	plugin := NewCalculationMethodPlugin(ephemerisManager)

	// Test plugin info
	info := plugin.GetInfo()
	if info.Name != "calculation_method_plugin" {
		t.Errorf("Expected plugin name 'calculation_method_plugin', got %s", info.Name)
	}

	// Test initialization
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Test supported methods
	methods := plugin.GetSupportedMethods()
	expectedMethods := []api.CalculationMethod{api.MethodDrik, api.MethodVakya, api.MethodAuto}
	if len(methods) != len(expectedMethods) {
		t.Errorf("Expected %d supported methods, got %d", len(expectedMethods), len(methods))
	}

	for _, expected := range expectedMethods {
		found := false
		for _, method := range methods {
			if method == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected method %s not found in supported methods", expected)
		}
	}
}

// Test weekday calculation for inauspicious periods
func TestWeekdayCalculations(t *testing.T) {
	plugin := NewMuhurtaPlugin()
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Test different weekdays to ensure Rahu Kalam timing is correct
	testDates := []struct {
		date    time.Time
		weekday string
	}{
		{time.Date(2024, 1, 14, 12, 0, 0, 0, time.UTC), "Sunday"},    // Sunday
		{time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), "Monday"},    // Monday
		{time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC), "Tuesday"},   // Tuesday
		{time.Date(2024, 1, 17, 12, 0, 0, 0, time.UTC), "Wednesday"}, // Wednesday
		{time.Date(2024, 1, 18, 12, 0, 0, 0, time.UTC), "Thursday"},  // Thursday
		{time.Date(2024, 1, 19, 12, 0, 0, 0, time.UTC), "Friday"},    // Friday
		{time.Date(2024, 1, 20, 12, 0, 0, 0, time.UTC), "Saturday"},  // Saturday
	}

	for _, testCase := range testDates {
		muhurtas, err := plugin.GetMuhurtas(context.Background(), testCase.date, testLocation, api.RegionGlobal)
		if err != nil {
			t.Fatalf("Failed to get muhurtas for %s: %v", testCase.weekday, err)
		}

		// Find Rahu Kalam
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

		// Verify Rahu Kalam has valid timing
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

// Test regional name variations
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
			
			// Verify region is set correctly
			if event.Region != region {
				t.Errorf("Event region mismatch: expected %s, got %s", region, event.Region)
			}
			
			t.Logf("Region %s: %s (%s)", region, event.Name, event.NameLocal)
		}
	}
}

// Benchmark tests
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