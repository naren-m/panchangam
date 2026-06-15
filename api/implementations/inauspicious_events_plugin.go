package implementations

import (
	"context"

	"github.com/naren-m/panchangam/api"
)

// InauspiciousEventsPlugin provides calculations for inauspicious periods
// Specifically handles Rahu Kalam, Yamagandam, and Gulika Kalam as events
type InauspiciousEventsPlugin struct {
	enabled bool
	config  map[string]interface{}
}

// NewInauspiciousEventsPlugin creates a new inauspicious events plugin
func NewInauspiciousEventsPlugin() *InauspiciousEventsPlugin {
	return &InauspiciousEventsPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

// GetInfo returns plugin metadata
func (i *InauspiciousEventsPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "inauspicious_events_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Calculations for inauspicious periods: Rahu Kalam, Yamagandam, and Gulika Kalam",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityEvent),
		},
		Dependencies: []string{"astronomy"},
		Metadata: map[string]interface{}{
			"event_types":      []string{"rahukalam", "yamagandam", "gulikakalam"},
			"calculation_base": "vedic_astronomy",
			"precision":        "minute_level",
		},
	}
}

// Initialize sets up the plugin with configuration
func (i *InauspiciousEventsPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	i.config = config
	i.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (i *InauspiciousEventsPlugin) IsEnabled() bool {
	return i.enabled
}

// Shutdown cleans up plugin resources
func (i *InauspiciousEventsPlugin) Shutdown(ctx context.Context) error {
	i.enabled = false
	return nil
}

// GetSupportedRegions returns regions this plugin supports
func (i *InauspiciousEventsPlugin) GetSupportedRegions() []api.Region {
	return []api.Region{
		api.RegionGlobal,
		api.RegionNorthIndia,
		api.RegionSouthIndia,
		api.RegionTamilNadu,
		api.RegionKerala,
		api.RegionBengal,
		api.RegionGujarat,
		api.RegionMaha,
	}
}

// GetSupportedEventTypes returns event types this plugin can generate
func (i *InauspiciousEventsPlugin) GetSupportedEventTypes() []api.EventType {
	return []api.EventType{
		api.EventTypeRahukalam,
		api.EventTypeYamagandam,
		api.EventTypeGulikakalam,
	}
}
