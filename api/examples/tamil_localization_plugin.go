package examples

import (
	"context"
	"fmt"

	"github.com/naren-m/panchangam/api"
)

// TamilLocalizationPlugin provides Tamil language localization for Panchangam elements
type TamilLocalizationPlugin struct {
	enabled bool
	config  map[string]interface{}

	// Translation dictionaries
	tithiNames     map[string]string
	nakshatraNames map[string]string
	yogaNames      map[string]string
	karanaNames    map[string]string
	varaNames      map[string]string
	eventNames     map[string]string
	muhurtaNames   map[string]string
}

// NewTamilLocalizationPlugin creates a new Tamil localization plugin
func NewTamilLocalizationPlugin() *TamilLocalizationPlugin {
	plugin := &TamilLocalizationPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}

	plugin.initializeTranslations()
	return plugin
}

// GetInfo returns plugin metadata
func (t *TamilLocalizationPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "tamil_localization",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Tamil language localization for Panchangam elements",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityLocalization),
		},
		Dependencies: []string{},
		Metadata: map[string]interface{}{
			"language":       "tamil",
			"script":         "tamil",
			"locale_codes":   []string{"ta", "ta-IN", "tamil"},
			"encoding":       "UTF-8",
			"region_support": []string{"tamil_nadu", "south_india", "global"},
		},
	}
}

// Initialize sets up the plugin with configuration
func (t *TamilLocalizationPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	t.config = config
	t.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (t *TamilLocalizationPlugin) IsEnabled() bool {
	return t.enabled
}

// Shutdown cleans up plugin resources
func (t *TamilLocalizationPlugin) Shutdown(ctx context.Context) error {
	t.enabled = false
	return nil
}

// GetSupportedLocales returns supported locale codes
func (t *TamilLocalizationPlugin) GetSupportedLocales() []string {
	return []string{"ta", "ta-IN", "tamil"}
}

// GetSupportedRegions returns regions this plugin supports
func (t *TamilLocalizationPlugin) GetSupportedRegions() []api.Region {
	return []api.Region{
		api.RegionTamilNadu,
		api.RegionSouthIndia,
		api.RegionGlobal,
	}
}

func (t *TamilLocalizationPlugin) ensureEnabled() error {
	if t.enabled {
		return nil
	}
	return fmt.Errorf("tamil localization plugin is not enabled")
}
