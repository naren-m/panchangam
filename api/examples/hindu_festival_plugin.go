package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
)

// HinduFestivalPlugin provides comprehensive Hindu festival and event calculations
type HinduFestivalPlugin struct {
	enabled bool
	config  map[string]interface{}
}

// NewHinduFestivalPlugin creates a new Hindu festival plugin
func NewHinduFestivalPlugin() *HinduFestivalPlugin {
	return &HinduFestivalPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

// GetInfo returns plugin metadata
func (h *HinduFestivalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "hindu_festival_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Comprehensive Hindu festival and event calculations",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityEvent),
		},
		Dependencies: []string{},
		Metadata: map[string]interface{}{
			"festival_count":    "100+",
			"regional_support":  true,
			"lunar_calendar":    true,
			"solar_calendar":    true,
			"regional_variants": true,
		},
	}
}

// Initialize sets up the plugin with configuration
func (h *HinduFestivalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	h.config = config
	h.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (h *HinduFestivalPlugin) IsEnabled() bool {
	return h.enabled
}

// Shutdown cleans up plugin resources
func (h *HinduFestivalPlugin) Shutdown(ctx context.Context) error {
	h.enabled = false
	return nil
}

// GetSupportedRegions returns regions this plugin supports
func (h *HinduFestivalPlugin) GetSupportedRegions() []api.Region {
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
func (h *HinduFestivalPlugin) GetSupportedEventTypes() []api.EventType {
	return []api.EventType{
		api.EventTypeFestival,
		api.EventTypeEkadashi,
		api.EventTypeAmavasya,
		api.EventTypePurnima,
		api.EventTypeVrat,
		api.EventTypeSankashti,
		api.EventTypeAshtami,
		api.EventTypeNavami,
		api.EventTypeLunar,
		api.EventTypeSolar,
	}
}

// GetEvents returns events for a specific date and location
func (h *HinduFestivalPlugin) GetEvents(ctx context.Context, date time.Time, location api.Location, region api.Region) ([]api.Event, error) {
	if !h.enabled {
		return nil, fmt.Errorf("hindu festival plugin is not enabled")
	}

	var events []api.Event

	// Major festivals based on lunar calendar
	events = append(events, h.getLunarFestivals(date, region)...)

	// Solar festivals
	events = append(events, h.getSolarFestivals(date, region)...)

	// Ekadashi dates
	events = append(events, h.getEkadashiEvents(date, region)...)

	// Monthly observances
	events = append(events, h.getMonthlyObservances(date, region)...)

	// Regional specific festivals
	events = append(events, h.getRegionalFestivals(date, region)...)

	return events, nil
}

// GetEventsInRange returns events for a date range
func (h *HinduFestivalPlugin) GetEventsInRange(ctx context.Context, start, end time.Time, location api.Location, region api.Region) ([]api.Event, error) {
	if !h.enabled {
		return nil, fmt.Errorf("hindu festival plugin is not enabled")
	}

	if end.Before(start) {
		return nil, fmt.Errorf("end date %s is before start date %s", end.Format("2006-01-02"), start.Format("2006-01-02"))
	}

	var allEvents []api.Event

	current := start
	for current.Before(end) || current.Equal(end) {
		dayEvents, err := h.GetEvents(ctx, current, location, region)
		if err != nil {
			return nil, err
		}
		allEvents = append(allEvents, dayEvents...)
		current = current.AddDate(0, 0, 1)
	}

	return allEvents, nil
}
