package api

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
)

type testCalculationPlugin struct {
	enabled bool
	methods []CalculationMethod
	regions []Region
}

func (p testCalculationPlugin) GetInfo() PluginInfo {
	return PluginInfo{
		Name:         "test_calculation_plugin",
		Version:      Version{Major: 1},
		Capabilities: []string{string(CapabilityCalculation)},
	}
}

func (p testCalculationPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	return nil
}

func (p testCalculationPlugin) IsEnabled() bool {
	return p.enabled
}

func (p testCalculationPlugin) Shutdown(ctx context.Context) error {
	return nil
}

func (p testCalculationPlugin) GetSupportedMethods() []CalculationMethod {
	return p.methods
}

func (p testCalculationPlugin) GetSupportedRegions() []Region {
	return p.regions
}

func (p testCalculationPlugin) CalculateTithi(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Tithi, error) {
	return nil, errors.New("test calculation plugin does not support tithi calculation")
}

func (p testCalculationPlugin) CalculateNakshatra(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Nakshatra, error) {
	return nil, errors.New("test calculation plugin does not support nakshatra calculation")
}

func (p testCalculationPlugin) CalculateYoga(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Yoga, error) {
	return nil, errors.New("test calculation plugin does not support yoga calculation")
}

func (p testCalculationPlugin) CalculateKarana(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Karana, error) {
	return nil, errors.New("test calculation plugin does not support karana calculation")
}

func (p testCalculationPlugin) CalculateSunMoonTimes(ctx context.Context, date time.Time, location Location) (*SunMoonTimes, error) {
	return nil, errors.New("test calculation plugin does not support sun and moon times")
}

type methodCheckingCalculationPlugin struct {
	testCalculationPlugin
}

func (p methodCheckingCalculationPlugin) CalculateTithi(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Tithi, error) {
	if method != MethodDrik {
		return nil, errors.New("expected default drik method")
	}
	return &Tithi{Number: 15, Name: "Purnima", IsRunning: true}, nil
}

func (p methodCheckingCalculationPlugin) CalculateNakshatra(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Nakshatra, error) {
	if method != MethodDrik {
		return nil, errors.New("expected default drik method")
	}
	return &Nakshatra{Number: 13, Name: "Hasta", Pada: 1, IsRunning: true}, nil
}

func (p methodCheckingCalculationPlugin) CalculateYoga(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Yoga, error) {
	if method != MethodDrik {
		return nil, errors.New("expected default drik method")
	}
	return &Yoga{Number: 14, Name: "Siddhi", IsRunning: true}, nil
}

func (p methodCheckingCalculationPlugin) CalculateKarana(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Karana, error) {
	if method != MethodDrik {
		return nil, errors.New("expected default drik method")
	}
	return &Karana{Number: 7, Name: "Vishti", IsRunning: true}, nil
}

func TestCorePanchangamAPISupportedValues(t *testing.T) {
	core := NewCorePanchangamAPI(observability.NewLocalObserver())

	if got := core.GetVersion().String(); got != "1.0.0-alpha" {
		t.Fatalf("unexpected version: %s", got)
	}
	if len(core.GetSupportedRegions()) != 8 {
		t.Fatalf("expected 8 supported regions, got %d", len(core.GetSupportedRegions()))
	}
	if len(core.GetSupportedMethods()) != 3 {
		t.Fatalf("expected 3 supported methods, got %d", len(core.GetSupportedMethods()))
	}
	if len(core.GetSupportedCalendars()) != 4 {
		t.Fatalf("expected 4 supported calendars, got %d", len(core.GetSupportedCalendars()))
	}
	if core.GetPluginManager() == nil {
		t.Fatal("expected plugin manager")
	}
}

func TestCorePanchangamAPIValidateRequest(t *testing.T) {
	core := NewCorePanchangamAPI(observability.NewLocalObserver())
	validRequest := PanchangamRequest{
		Date: time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
		Location: Location{
			Latitude:  12.9716,
			Longitude: 77.5946,
		},
	}

	if err := core.validateRequest(context.Background(), validRequest); err != nil {
		t.Fatalf("valid request failed validation: %v", err)
	}

	tests := []struct {
		name    string
		update  func(*PanchangamRequest)
		wantErr string
	}{
		{
			name: "latitude below range",
			update: func(req *PanchangamRequest) {
				req.Location.Latitude = -91
			},
			wantErr: "invalid latitude",
		},
		{
			name: "longitude above range",
			update: func(req *PanchangamRequest) {
				req.Location.Longitude = 181
			},
			wantErr: "invalid longitude",
		},
		{
			name: "date is required",
			update: func(req *PanchangamRequest) {
				req.Date = time.Time{}
			},
			wantErr: "date is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validRequest
			tt.update(&req)

			err := core.validateRequest(context.Background(), req)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestCorePanchangamAPIRequiresCalculationPlugin(t *testing.T) {
	core := NewCorePanchangamAPI(observability.NewLocalObserver())
	req := PanchangamRequest{
		Date: time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
		Location: Location{
			Latitude:  12.9716,
			Longitude: 77.5946,
		},
	}

	data, err := core.GetPanchangam(context.Background(), req)

	if err == nil {
		t.Fatalf("expected missing calculation plugin error, got data: %+v", data)
	}
	if data != nil {
		t.Fatalf("expected no data when calculation plugin is missing, got: %+v", data)
	}
	if !strings.Contains(err.Error(), "no enabled calculation plugin") {
		t.Fatalf("expected missing calculation plugin error, got %q", err.Error())
	}
}

func TestCorePanchangamAPIPassesDefaultMethodToCalculationPlugin(t *testing.T) {
	core := NewCorePanchangamAPI(observability.NewLocalObserver())
	if err := core.RegisterPlugin(methodCheckingCalculationPlugin{
		testCalculationPlugin: testCalculationPlugin{
			enabled: true,
			methods: []CalculationMethod{MethodDrik},
			regions: []Region{RegionGlobal},
		},
	}); err != nil {
		t.Fatalf("failed to register plugin: %v", err)
	}

	data, err := core.GetPanchangam(context.Background(), PanchangamRequest{
		Date: time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
		Location: Location{
			Latitude:  12.9716,
			Longitude: 77.5946,
		},
	})

	if err != nil {
		t.Fatalf("expected plugin calculation to use default method, got error: %v", err)
	}
	if data.CalculationMethod != MethodDrik {
		t.Fatalf("expected result to use default method %q, got %q", MethodDrik, data.CalculationMethod)
	}
}

func TestCorePanchangamAPIPluginSupportsMethodAndRegion(t *testing.T) {
	core := NewCorePanchangamAPI(observability.NewLocalObserver())

	plugin := testCalculationPlugin{
		enabled: true,
		methods: []CalculationMethod{MethodDrik},
		regions: []Region{RegionTamilNadu, RegionGlobal},
	}

	if !core.pluginSupportsMethodAndRegion(plugin, MethodDrik, RegionTamilNadu) {
		t.Fatal("expected explicit region support")
	}
	if !core.pluginSupportsMethodAndRegion(plugin, MethodDrik, RegionKerala) {
		t.Fatal("expected global region support")
	}
	if core.pluginSupportsMethodAndRegion(plugin, MethodVakya, RegionTamilNadu) {
		t.Fatal("did not expect unsupported method to match")
	}
}

func TestCorePanchangamAPIGetDateRangeRejectsEndBeforeStart(t *testing.T) {
	core := NewCorePanchangamAPI(observability.NewLocalObserver())
	start := time.Date(2024, time.January, 16, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)

	results, err := core.GetDateRange(context.Background(), start, end, Location{
		Latitude:  12.9716,
		Longitude: 77.5946,
	})

	if err == nil {
		t.Fatalf("expected invalid date range error, got results: %+v", results)
	}
	if results != nil {
		t.Fatalf("expected no results for invalid date range, got: %+v", results)
	}
	if !strings.Contains(err.Error(), "end date") || !strings.Contains(err.Error(), "before start date") {
		t.Fatalf("expected clear invalid date range error, got %q", err.Error())
	}
}
