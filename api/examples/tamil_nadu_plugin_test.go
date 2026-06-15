package examples

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/naren-m/panchangam/api"
)

func TestTamilNaduExtensionOnlyAdvertisesImplementedPluginContracts(t *testing.T) {
	plugin := NewTamilNaduExtension()
	info := plugin.GetInfo()

	if _, ok := interface{}(plugin).(api.RegionalExtension); !ok {
		t.Fatal("test expects TamilNaduExtension to provide regional extension behavior")
	}
	if _, ok := interface{}(plugin).(api.EventPlugin); ok {
		t.Fatal("test expects TamilNaduExtension regional events to stay behind the RegionalExtension contract")
	}
	if _, ok := interface{}(plugin).(api.MuhurtaPlugin); ok {
		t.Fatal("test expects TamilNaduExtension regional muhurtas to stay behind the RegionalExtension contract")
	}

	for _, capability := range info.Capabilities {
		switch capability {
		case string(api.CapabilityEvent), string(api.CapabilityMuhurta):
			t.Fatalf("TamilNaduExtension should not advertise %s without implementing its plugin interface", capability)
		}
	}
}

func TestTamilNaduExtensionRejectsRegionalEventsWhenDisabled(t *testing.T) {
	plugin := NewTamilNaduExtension()

	_, err := plugin.GetRegionalEvents(context.Background(), time.Date(2024, time.April, 14, 0, 0, 0, 0, time.UTC), api.Location{})
	if err == nil {
		t.Fatal("Expected disabled extension to return an error")
	}
	if !strings.Contains(err.Error(), "tamil nadu extension is not enabled") {
		t.Fatalf("Expected disabled extension error, got %v", err)
	}
}

func TestTamilNaduExtensionRejectsRegionalMuhurtasWhenDisabled(t *testing.T) {
	plugin := NewTamilNaduExtension()

	_, err := plugin.GetRegionalMuhurtas(context.Background(), time.Date(2024, time.April, 14, 0, 0, 0, 0, time.UTC), api.Location{})
	if err == nil {
		t.Fatal("Expected disabled extension to return an error")
	}
	if !strings.Contains(err.Error(), "tamil nadu extension is not enabled") {
		t.Fatalf("Expected disabled extension error, got %v", err)
	}
}

func TestTamilNaduExtensionRejectsRegionalRulesWhenDisabled(t *testing.T) {
	plugin := NewTamilNaduExtension()

	err := plugin.ApplyRegionalRules(context.Background(), &api.PanchangamData{})
	if err == nil {
		t.Fatal("Expected disabled extension to return an error")
	}
	if !strings.Contains(err.Error(), "tamil nadu extension is not enabled") {
		t.Fatalf("Expected disabled extension error, got %v", err)
	}
}
