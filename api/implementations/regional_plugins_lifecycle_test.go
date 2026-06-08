package implementations

import (
	"context"
	"strings"
	"testing"

	"github.com/naren-m/panchangam/api"
)

func TestRegionalPluginsRejectRegionalCallsWhenDisabled(t *testing.T) {
	plugins := []api.RegionalExtension{
		NewTamilNaduRegionalPlugin(),
		NewKeralaRegionalPlugin(),
		NewBengalRegionalPlugin(),
		NewGujaratRegionalPlugin(),
		NewMaharashtraRegionalPlugin(),
	}

	for _, plugin := range plugins {
		t.Run(plugin.GetInfo().Name, func(t *testing.T) {
			assertDisabledRegionalPlugin(t, plugin)

			if err := plugin.Initialize(context.Background(), map[string]interface{}{}); err != nil {
				t.Fatalf("failed to initialize plugin: %v", err)
			}
			if err := plugin.Shutdown(context.Background()); err != nil {
				t.Fatalf("failed to shutdown plugin: %v", err)
			}

			assertDisabledRegionalPlugin(t, plugin)
		})
	}
}

func assertDisabledRegionalPlugin(t *testing.T, plugin api.RegionalExtension) {
	t.Helper()

	want := plugin.GetInfo().Name + " is not enabled"
	data := &api.PanchangamData{
		Date:     testDateRegional,
		Location: testLocationTN,
	}

	if err := plugin.ApplyRegionalRules(context.Background(), data); err == nil {
		t.Fatal("expected disabled regional rules call to return an error")
	} else if !strings.Contains(err.Error(), want) {
		t.Fatalf("expected disabled regional rules error %q, got %v", want, err)
	}

	events, err := plugin.GetRegionalEvents(context.Background(), testDateRegional, testLocationTN)
	if err == nil {
		t.Fatal("expected disabled regional event call to return an error")
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("expected disabled regional event error %q, got %v", want, err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no events from disabled regional call, got %d", len(events))
	}

	muhurtas, err := plugin.GetRegionalMuhurtas(context.Background(), testDateRegional, testLocationTN)
	if err == nil {
		t.Fatal("expected disabled regional muhurta call to return an error")
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("expected disabled regional muhurta error %q, got %v", want, err)
	}
	if len(muhurtas) != 0 {
		t.Fatalf("expected no muhurtas from disabled regional call, got %d", len(muhurtas))
	}
}
