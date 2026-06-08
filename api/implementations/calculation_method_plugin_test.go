package implementations

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

func TestCalculationMethodPlugin(t *testing.T) {
	ephemerisManager := &ephemeris.Manager{}
	plugin := NewCalculationMethodPlugin(ephemerisManager)

	info := plugin.GetInfo()
	if info.Name != "calculation_method_plugin" {
		t.Errorf("Expected plugin name 'calculation_method_plugin', got %s", info.Name)
	}

	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

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

func TestCalculationMethodPluginRejectsNakshatraCalculationWhenDisabled(t *testing.T) {
	plugin := NewCalculationMethodPlugin(&ephemeris.Manager{})

	_, err := plugin.CalculateNakshatra(
		context.Background(),
		time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC),
		api.Location{},
		api.CalculationMethod("bad"),
	)
	if err == nil {
		t.Fatal("Expected disabled plugin to return an error")
	}
	if err.Error() != "calculation method plugin is not enabled" {
		t.Fatalf("Expected disabled plugin error, got %v", err)
	}
}

func TestCalculationMethodPluginDoesNotAdvertiseRoutedCalculationCapability(t *testing.T) {
	plugin := NewCalculationMethodPlugin(&ephemeris.Manager{})
	info := plugin.GetInfo()

	if _, ok := interface{}(plugin).(api.CalculationPlugin); ok {
		t.Fatal("test expects CalculationMethodPlugin to stay behind its own partial calculation methods")
	}

	for _, capability := range info.Capabilities {
		if capability == string(api.CapabilityCalculation) {
			t.Fatal("CalculationMethodPlugin should not advertise calculation capability without implementing CalculationPlugin")
		}
	}
}
