package implementations

import (
	"context"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

// CalculationMethodPlugin handles different Hindu astronomical calculation methods
// Primarily Drik Ganita (observational/modern) vs Vakya (traditional/tabular)
type CalculationMethodPlugin struct {
	enabled          bool
	config           map[string]interface{}
	tithiCalculator  *astronomy.TithiCalculator
	ephemerisManager *ephemeris.Manager
}

// NewCalculationMethodPlugin creates a new calculation method plugin
func NewCalculationMethodPlugin(ephemerisManager *ephemeris.Manager) *CalculationMethodPlugin {
	return &CalculationMethodPlugin{
		enabled:          false,
		config:           make(map[string]interface{}),
		tithiCalculator:  astronomy.NewTithiCalculator(ephemerisManager),
		ephemerisManager: ephemerisManager,
	}
}

// GetInfo returns plugin metadata
func (c *CalculationMethodPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:         "calculation_method_plugin",
		Version:      api.Version{Major: 1, Minor: 0, Patch: 0},
		Description:  "Handles different Hindu astronomical calculation methods: Drik Ganita vs Vakya",
		Author:       "Panchangam Team",
		Capabilities: []string{},
		Dependencies: []string{"astronomy", "ephemeris"},
		Metadata: map[string]interface{}{
			"calculation_methods": []string{"drik", "vakya", "auto"},
			"precision_levels":    []string{"high", "medium", "traditional"},
			"ayanamsa_systems":    []string{"lahiri", "krishnamurti", "raman", "traditional"},
		},
	}
}

// Initialize sets up the plugin with configuration
func (c *CalculationMethodPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	c.config = config
	c.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (c *CalculationMethodPlugin) IsEnabled() bool {
	return c.enabled
}

// Shutdown cleans up plugin resources
func (c *CalculationMethodPlugin) Shutdown(ctx context.Context) error {
	c.enabled = false
	return nil
}

// GetSupportedMethods returns calculation methods this plugin supports
func (c *CalculationMethodPlugin) GetSupportedMethods() []api.CalculationMethod {
	return []api.CalculationMethod{
		api.MethodDrik,
		api.MethodVakya,
		api.MethodAuto,
	}
}
