package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
)

// MaharashtraRegionalPlugin provides Maharashtra specific calculations
type MaharashtraRegionalPlugin struct {
	enabled bool
	config  map[string]interface{}
}

func NewMaharashtraRegionalPlugin() *MaharashtraRegionalPlugin {
	return &MaharashtraRegionalPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

func (m *MaharashtraRegionalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "maharashtra_regional_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Maharashtra regional calendar rules and Marathi names",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
		},
		Metadata: map[string]interface{}{
			"calendar_system": "purnimanta",
			"language":        "marathi",
		},
	}
}

func (m *MaharashtraRegionalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	m.config = config
	m.enabled = true
	return nil
}

func (m *MaharashtraRegionalPlugin) IsEnabled() bool {
	return m.enabled
}

func (m *MaharashtraRegionalPlugin) Shutdown(ctx context.Context) error {
	m.enabled = false
	return nil
}

func (m *MaharashtraRegionalPlugin) GetRegion() api.Region {
	return api.RegionMaha
}

func (m *MaharashtraRegionalPlugin) GetCalendarSystem() api.CalendarSystem {
	return api.CalendarPurnimanta
}

func (m *MaharashtraRegionalPlugin) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error {
	if !m.enabled {
		return fmt.Errorf("maharashtra_regional_plugin is not enabled")
	}

	data.CalendarSystem = api.CalendarPurnimanta
	return nil
}

func (m *MaharashtraRegionalPlugin) GetRegionalEvents(ctx context.Context, date time.Time, location api.Location) ([]api.Event, error) {
	if !m.enabled {
		return nil, fmt.Errorf("maharashtra_regional_plugin is not enabled")
	}

	return nil, fmt.Errorf("regional events are unavailable for %s; use event plugins", m.GetRegion())
}

func (m *MaharashtraRegionalPlugin) GetRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location) ([]api.Muhurta, error) {
	if !m.enabled {
		return nil, fmt.Errorf("maharashtra_regional_plugin is not enabled")
	}

	return nil, fmt.Errorf("regional muhurtas are unavailable for %s; use muhurta_plugin", m.GetRegion())
}

func (m *MaharashtraRegionalPlugin) GetRegionalNames(locale string) map[string]string {
	marathiNames := map[string]string{
		"Sunday":    "रविवार",
		"Monday":    "सोमवार",
		"Tuesday":   "मंगळवार",
		"Wednesday": "बुधवार",
		"Thursday":  "गुरुवार",
		"Friday":    "शुक्रवार",
		"Saturday":  "शनिवार",
	}
	return marathiNames
}
