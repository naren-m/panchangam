package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
)

// BengalRegionalPlugin provides Bengal specific calculations
type BengalRegionalPlugin struct {
	enabled bool
	config  map[string]interface{}
}

func NewBengalRegionalPlugin() *BengalRegionalPlugin {
	return &BengalRegionalPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

func (b *BengalRegionalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "bengal_regional_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Bengal regional extensions with Amanta calendar and Bengali calendar",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
		},
		Metadata: map[string]interface{}{
			"calendar_system": "amanta",
			"language":        "bengali",
			"festivals":       []string{"Durga Puja", "Pohela Boishakh", "Kali Puja"},
		},
	}
}

func (b *BengalRegionalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	b.config = config
	b.enabled = true
	return nil
}

func (b *BengalRegionalPlugin) IsEnabled() bool {
	return b.enabled
}

func (b *BengalRegionalPlugin) Shutdown(ctx context.Context) error {
	b.enabled = false
	return nil
}

func (b *BengalRegionalPlugin) GetRegion() api.Region {
	return api.RegionBengal
}

func (b *BengalRegionalPlugin) GetCalendarSystem() api.CalendarSystem {
	return api.CalendarAmanta
}

func (b *BengalRegionalPlugin) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error {
	if !b.enabled {
		return fmt.Errorf("bengal_regional_plugin is not enabled")
	}

	data.CalendarSystem = api.CalendarAmanta
	return nil
}

func (b *BengalRegionalPlugin) GetRegionalEvents(ctx context.Context, date time.Time, location api.Location) ([]api.Event, error) {
	if !b.enabled {
		return nil, fmt.Errorf("bengal_regional_plugin is not enabled")
	}

	var events []api.Event

	// Pohela Boishakh (Bengali New Year) - April 14/15
	if date.Month() == time.April && (date.Day() == 14 || date.Day() == 15) {
		events = append(events, api.Event{
			Name:         "Pohela Boishakh",
			NameLocal:    "পহেলা বৈশাখ",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Bengali New Year celebration",
			Region:       api.RegionBengal,
			Metadata: map[string]interface{}{
				"importance": "highest",
				"new_year":   true,
			},
		})
	}

	return events, nil
}

func (b *BengalRegionalPlugin) GetRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location) ([]api.Muhurta, error) {
	if !b.enabled {
		return nil, fmt.Errorf("bengal_regional_plugin is not enabled")
	}

	return nil, fmt.Errorf("regional muhurtas are unavailable for %s; use muhurta_plugin", b.GetRegion())
}

func (b *BengalRegionalPlugin) GetRegionalNames(locale string) map[string]string {
	bengaliNames := map[string]string{
		"Sunday":    "রবিবার",
		"Monday":    "সোমবার",
		"Tuesday":   "মঙ্গলবার",
		"Wednesday": "বুধবার",
		"Thursday":  "বৃহস্পতিবার",
		"Friday":    "শুক্রবার",
		"Saturday":  "শনিবার",
	}
	return bengaliNames
}
