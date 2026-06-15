package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
)

// KeralaRegionalPlugin provides Kerala specific calculations
type KeralaRegionalPlugin struct {
	enabled bool
	config  map[string]interface{}
}

// NewKeralaRegionalPlugin creates a new Kerala regional plugin
func NewKeralaRegionalPlugin() *KeralaRegionalPlugin {
	return &KeralaRegionalPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

// GetInfo returns plugin metadata
func (k *KeralaRegionalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "kerala_regional_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Kerala regional extensions with Amanta calendar and Malayalam calendar",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
		},
		Dependencies: []string{"astronomy"},
		Metadata: map[string]interface{}{
			"calendar_system": "amanta",
			"language":        "malayalam",
			"festivals":       []string{"Onam", "Vishu", "Thiruvathira"},
		},
	}
}

func (k *KeralaRegionalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	k.config = config
	k.enabled = true
	return nil
}

func (k *KeralaRegionalPlugin) IsEnabled() bool {
	return k.enabled
}

func (k *KeralaRegionalPlugin) Shutdown(ctx context.Context) error {
	k.enabled = false
	return nil
}

func (k *KeralaRegionalPlugin) GetRegion() api.Region {
	return api.RegionKerala
}

func (k *KeralaRegionalPlugin) GetCalendarSystem() api.CalendarSystem {
	return api.CalendarAmanta
}

func (k *KeralaRegionalPlugin) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error {
	if !k.enabled {
		return fmt.Errorf("kerala_regional_plugin is not enabled")
	}

	data.CalendarSystem = api.CalendarAmanta
	return nil
}

func (k *KeralaRegionalPlugin) GetRegionalEvents(ctx context.Context, date time.Time, location api.Location) ([]api.Event, error) {
	if !k.enabled {
		return nil, fmt.Errorf("kerala_regional_plugin is not enabled")
	}

	var events []api.Event

	// Vishu (Kerala New Year) - April 14/15
	if date.Month() == time.April && (date.Day() == 14 || date.Day() == 15) {
		events = append(events, api.Event{
			Name:         "Vishu",
			NameLocal:    "വിഷു",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Malayalam New Year celebrated with Vishukkani",
			Region:       api.RegionKerala,
			Metadata: map[string]interface{}{
				"importance":     "highest",
				"new_year":       true,
				"solar_festival": true,
			},
		})
	}

	return events, nil
}

func (k *KeralaRegionalPlugin) GetRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location) ([]api.Muhurta, error) {
	if !k.enabled {
		return nil, fmt.Errorf("kerala_regional_plugin is not enabled")
	}

	return nil, fmt.Errorf("regional muhurtas are unavailable for %s; use muhurta_plugin", k.GetRegion())
}

func (k *KeralaRegionalPlugin) GetRegionalNames(locale string) map[string]string {
	malayalamNames := map[string]string{
		"Sunday":    "ഞായർ",
		"Monday":    "തിങ്കൾ",
		"Tuesday":   "ചൊവ്വ",
		"Wednesday": "ബുധൻ",
		"Thursday":  "വ്യാഴം",
		"Friday":    "വെള്ളി",
		"Saturday":  "ശനി",
	}
	return malayalamNames
}
