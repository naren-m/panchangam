package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

// AdvancedFestivalPlugin provides precise Hindu festival calculations using lunar astronomy
type AdvancedFestivalPlugin struct {
	enabled          bool
	config           map[string]interface{}
	tithiCalculator  *astronomy.TithiCalculator
	ephemerisManager *ephemeris.Manager
}

// NewAdvancedFestivalPlugin creates a new advanced festival calculation plugin
func NewAdvancedFestivalPlugin(ephemerisManager *ephemeris.Manager) *AdvancedFestivalPlugin {
	return &AdvancedFestivalPlugin{
		enabled:          false,
		config:           make(map[string]interface{}),
		tithiCalculator:  astronomy.NewTithiCalculator(ephemerisManager),
		ephemerisManager: ephemerisManager,
	}
}

// GetInfo returns plugin metadata
func (a *AdvancedFestivalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "advanced_festival_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Precise Hindu festival calculations using astronomical data and lunar calendar",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityEvent),
		},
		Dependencies: []string{"astronomy", "ephemeris"},
		Metadata: map[string]interface{}{
			"calculation_precision": "astronomical",
			"festival_count":        "50+",
			"lunar_calculations":    true,
			"regional_variations":   true,
			"ayanamsa_aware":        true,
		},
	}
}

// Initialize sets up the plugin with configuration
func (a *AdvancedFestivalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	a.config = config
	a.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (a *AdvancedFestivalPlugin) IsEnabled() bool {
	return a.enabled
}

// Shutdown cleans up plugin resources
func (a *AdvancedFestivalPlugin) Shutdown(ctx context.Context) error {
	a.enabled = false
	return nil
}

// GetSupportedRegions returns regions this plugin supports
func (a *AdvancedFestivalPlugin) GetSupportedRegions() []api.Region {
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
func (a *AdvancedFestivalPlugin) GetSupportedEventTypes() []api.EventType {
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
	}
}

// GetEvents returns festival events for a specific date and location
func (a *AdvancedFestivalPlugin) GetEvents(ctx context.Context, date time.Time, location api.Location, region api.Region) ([]api.Event, error) {
	if !a.enabled {
		return nil, fmt.Errorf("advanced festival plugin is not enabled")
	}

	var events []api.Event

	// Get Tithi information for precise lunar calculations
	tithi, err := a.tithiCalculator.GetTithiForDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate tithi: %w", err)
	}

	// Check for festivals based on Tithi
	tithiEvents := a.getFestivalsByTithi(ctx, date, tithi, region)
	events = append(events, tithiEvents...)

	// Check for Ekadashi
	if tithi.Number == 11 || tithi.Number == 26 { // 11th day of both fortnights
		ekadashiEvent := a.calculateEkadashi(ctx, date, tithi, region)
		events = append(events, ekadashiEvent)
	}

	// Check for Amavasya (New Moon)
	if tithi.Number == 30 {
		amavasya := a.calculateAmavasya(ctx, date, tithi, region)
		events = append(events, amavasya)

		// Check for special Amavasya festivals
		specialAmavasya := a.getSpecialAmavasyas(ctx, date, tithi, region)
		events = append(events, specialAmavasya...)
	}

	// Check for Purnima (Full Moon)
	if tithi.Number == 15 {
		purnima := a.calculatePurnima(ctx, date, tithi, region)
		events = append(events, purnima)

		// Check for special Purnima festivals
		specialPurnima := a.getSpecialPurnimas(ctx, date, tithi, region)
		events = append(events, specialPurnima...)
	}

	// Check for Ashtami festivals (8th day)
	if tithi.Number == 8 || tithi.Number == 23 {
		ashtamiEvents := a.getAshtamiFestivals(ctx, date, tithi, region)
		events = append(events, ashtamiEvents...)
	}

	// Check for Navami festivals (9th day)
	if tithi.Number == 9 || tithi.Number == 24 {
		navamiEvents := a.getNavamiFestivals(ctx, date, tithi, region)
		events = append(events, navamiEvents...)
	}

	// Check for Sankashti Chaturthi (4th day of Krishna Paksha)
	if tithi.Number == 19 {
		sankashti := a.calculateSankashtiChaturthi(ctx, date, tithi, region)
		events = append(events, sankashti)
	}

	// Regional specific festivals
	regionalEvents := a.getRegionalFestivals(ctx, date, tithi, region)
	events = append(events, regionalEvents...)

	return events, nil
}

// GetEventsInRange returns festival events for a date range
func (a *AdvancedFestivalPlugin) GetEventsInRange(ctx context.Context, start, end time.Time, location api.Location, region api.Region) ([]api.Event, error) {
	if !a.enabled {
		return nil, fmt.Errorf("advanced festival plugin is not enabled")
	}

	if err := validateEventDateRange(start, end); err != nil {
		return nil, err
	}

	var allEvents []api.Event

	current := start
	for current.Before(end) || current.Equal(end) {
		dayEvents, err := a.GetEvents(ctx, current, location, region)
		if err != nil {
			return nil, fmt.Errorf("failed to get events for %s: %w", current.Format("2006-01-02"), err)
		}
		allEvents = append(allEvents, dayEvents...)
		current = current.AddDate(0, 0, 1)
	}

	return allEvents, nil
}
