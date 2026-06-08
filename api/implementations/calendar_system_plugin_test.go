package implementations

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

func TestCalendarSystemPlugin(t *testing.T) {
	ephemerisManager := &ephemeris.Manager{}
	plugin := NewCalendarSystemPlugin(ephemerisManager)

	info := plugin.GetInfo()
	if info.Name != "calendar_system_plugin" {
		t.Errorf("Expected plugin name 'calendar_system_plugin', got %s", info.Name)
	}

	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	regions := plugin.GetSupportedRegions()
	if len(regions) == 0 {
		t.Error("Expected supported regions, got none")
	}

	methods := plugin.GetSupportedMethods()
	expectedMethods := []api.CalculationMethod{api.MethodDrik, api.MethodVakya, api.MethodAuto}
	if len(methods) != len(expectedMethods) {
		t.Errorf("Expected %d supported methods, got %d", len(expectedMethods), len(methods))
	}
}

func TestCalendarSystemPluginDoesNotAdvertiseRoutedCapabilities(t *testing.T) {
	plugin := NewCalendarSystemPlugin(&ephemeris.Manager{})
	info := plugin.GetInfo()

	if _, ok := interface{}(plugin).(api.CalculationPlugin); ok {
		t.Fatal("test expects CalendarSystemPlugin to use its own calendar-system methods")
	}
	if _, ok := interface{}(plugin).(api.RegionalExtension); ok {
		t.Fatal("test expects CalendarSystemPlugin to use its own calendar-system methods")
	}

	for _, capability := range info.Capabilities {
		switch capability {
		case string(api.CapabilityCalculation), string(api.CapabilityRegional):
			t.Fatalf("CalendarSystemPlugin should not advertise %s without implementing its routed plugin interface", capability)
		}
	}
}

func TestCalendarSystemPluginRejectsCurrentMonthWhenDisabled(t *testing.T) {
	plugin := NewCalendarSystemPlugin(&ephemeris.Manager{})

	_, err := plugin.GetCurrentMonth(
		context.Background(),
		time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC),
		api.Location{},
		api.CalendarSystem("bad"),
		api.RegionGlobal,
	)
	if err == nil {
		t.Fatal("Expected disabled plugin to return an error")
	}
	if err.Error() != "calendar system plugin is not enabled" {
		t.Fatalf("Expected disabled plugin error, got %v", err)
	}
}

func TestCalendarSystemPluginHelpers(t *testing.T) {
	plugin := NewCalendarSystemPlugin(&ephemeris.Manager{})
	date := time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)

	monthNumber, monthName := plugin.getAmantaMonthInfo(date, api.RegionGlobal)
	if monthNumber != 12 || monthName != "Phalgun" {
		t.Fatalf("unexpected Amanta month: %d %s", monthNumber, monthName)
	}

	monthNumber, monthName = plugin.getPurnimantaMonthInfo(date, api.RegionGlobal)
	if monthNumber != 11 || monthName != "Magha" {
		t.Fatalf("unexpected Purnimanta month: %d %s", monthNumber, monthName)
	}

	monthNumber, monthName = plugin.getSolarMonthInfo(359, api.RegionGlobal)
	if monthNumber != 12 || monthName != "Meena" {
		t.Fatalf("unexpected solar month: %d %s", monthNumber, monthName)
	}

	if !plugin.isNorthIndianRegion(api.RegionMaha) {
		t.Fatal("expected Maharashtra to use a north Indian calendar preference")
	}
	if plugin.isNorthIndianRegion(api.RegionKerala) {
		t.Fatal("expected Kerala to use a south Indian calendar preference")
	}

	if got := plugin.getCalendarSystemDescription(api.CalendarAmanta); got == "" {
		t.Fatal("expected Amanta calendar description")
	}
	if got := plugin.getMonthBoundaryRule(api.CalendarSolar); got == "" {
		t.Fatal("expected solar month boundary rule")
	}
	if got := plugin.getRegionalPreference(api.RegionGujarat); got != "Purnimanta" {
		t.Fatalf("expected Gujarat preference Purnimanta, got %s", got)
	}
}
