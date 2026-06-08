package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
)

// GujaratRegionalPlugin provides Gujarat specific calculations
type GujaratRegionalPlugin struct {
	enabled bool
	config  map[string]interface{}
}

func NewGujaratRegionalPlugin() *GujaratRegionalPlugin {
	return &GujaratRegionalPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

func (g *GujaratRegionalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "gujarat_regional_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Gujarat regional extensions with Purnimanta calendar and Gujarati calendar",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
		},
		Metadata: map[string]interface{}{
			"calendar_system": "purnimanta",
			"language":        "gujarati",
			"festivals":       []string{"Uttarayan", "Navratri", "Bestu Varas"},
		},
	}
}

func (g *GujaratRegionalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	g.config = config
	g.enabled = true
	return nil
}

func (g *GujaratRegionalPlugin) IsEnabled() bool {
	return g.enabled
}

func (g *GujaratRegionalPlugin) Shutdown(ctx context.Context) error {
	g.enabled = false
	return nil
}

func (g *GujaratRegionalPlugin) GetRegion() api.Region {
	return api.RegionGujarat
}

func (g *GujaratRegionalPlugin) GetCalendarSystem() api.CalendarSystem {
	return api.CalendarPurnimanta
}

func (g *GujaratRegionalPlugin) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error {
	if !g.enabled {
		return fmt.Errorf("gujarat_regional_plugin is not enabled")
	}

	data.CalendarSystem = api.CalendarPurnimanta
	return nil
}

func (g *GujaratRegionalPlugin) GetRegionalEvents(ctx context.Context, date time.Time, location api.Location) ([]api.Event, error) {
	if !g.enabled {
		return nil, fmt.Errorf("gujarat_regional_plugin is not enabled")
	}

	var events []api.Event

	// Uttarayan - January 14
	if date.Month() == time.January && date.Day() == 14 {
		events = append(events, api.Event{
			Name:         "Uttarayan",
			NameLocal:    "ઉત્તરાયણ",
			Type:         api.EventTypeSolar,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Kite festival celebrating the sun's northward journey",
			Region:       api.RegionGujarat,
			Metadata: map[string]interface{}{
				"importance":    "highest",
				"solar_event":   true,
				"kite_festival": true,
			},
		})
	}

	return events, nil
}

func (g *GujaratRegionalPlugin) GetRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location) ([]api.Muhurta, error) {
	if !g.enabled {
		return nil, fmt.Errorf("gujarat_regional_plugin is not enabled")
	}

	return nil, fmt.Errorf("regional muhurtas are unavailable for %s; use muhurta_plugin", g.GetRegion())
}

func (g *GujaratRegionalPlugin) GetRegionalNames(locale string) map[string]string {
	gujaratiNames := map[string]string{
		"Sunday":    "રવિવાર",
		"Monday":    "સોમવાર",
		"Tuesday":   "મંગળવાર",
		"Wednesday": "બુધવાર",
		"Thursday":  "ગુરુવાર",
		"Friday":    "શુક્રવાર",
		"Saturday":  "શનિવાર",
	}
	return gujaratiNames
}
