package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
)

// TamilNaduRegionalPlugin provides Tamil Nadu specific calculations
type TamilNaduRegionalPlugin struct {
	enabled bool
	config  map[string]interface{}
}

// NewTamilNaduRegionalPlugin creates a new Tamil Nadu regional plugin
func NewTamilNaduRegionalPlugin() *TamilNaduRegionalPlugin {
	return &TamilNaduRegionalPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

// GetInfo returns plugin metadata
func (t *TamilNaduRegionalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "tamil_nadu_regional_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Tamil Nadu regional extensions with Amanta calendar and Naazhikai time system",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
		},
		Dependencies: []string{"astronomy"},
		Metadata: map[string]interface{}{
			"calendar_system": "amanta",
			"time_system":     "naazhikai",
			"language":        "tamil",
			"festivals":       []string{"Pongal", "Tamil New Year", "Aadi Perukku"},
		},
	}
}

// Initialize sets up the plugin
func (t *TamilNaduRegionalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	t.config = config
	t.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is enabled
func (t *TamilNaduRegionalPlugin) IsEnabled() bool {
	return t.enabled
}

// Shutdown cleans up resources
func (t *TamilNaduRegionalPlugin) Shutdown(ctx context.Context) error {
	t.enabled = false
	return nil
}

// GetRegion returns the region this plugin supports
func (t *TamilNaduRegionalPlugin) GetRegion() api.Region {
	return api.RegionTamilNadu
}

// GetCalendarSystem returns the calendar system used
func (t *TamilNaduRegionalPlugin) GetCalendarSystem() api.CalendarSystem {
	return api.CalendarAmanta
}

// ApplyRegionalRules applies Tamil Nadu specific rules
func (t *TamilNaduRegionalPlugin) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error {
	if !t.enabled {
		return fmt.Errorf("tamil_nadu_regional_plugin is not enabled")
	}

	// Set calendar system to Amanta
	data.CalendarSystem = api.CalendarAmanta

	// Add regional metadata via events
	regionalInfo := api.Event{
		Name:         "Tamil Nadu Regional Info",
		NameLocal:    "தமிழ்நாடு பிராந்திய தகவல்",
		Type:         api.EventTypeLunar,
		StartTime:    data.Date,
		EndTime:      data.Date.Add(24 * time.Hour),
		Significance: "Tamil Nadu follows Amanta calendar system and uses Naazhikai time units",
		Region:       api.RegionTamilNadu,
		Metadata: map[string]interface{}{
			"type":            "regional_info",
			"calendar_system": "amanta",
			"time_system":     "naazhikai",
			"language":        "tamil",
		},
	}

	data.Events = append(data.Events, regionalInfo)
	return nil
}

// GetRegionalEvents returns Tamil Nadu specific events
func (t *TamilNaduRegionalPlugin) GetRegionalEvents(ctx context.Context, date time.Time, location api.Location) ([]api.Event, error) {
	if !t.enabled {
		return nil, fmt.Errorf("tamil_nadu_regional_plugin is not enabled")
	}

	var events []api.Event

	// Tamil New Year (Puthandu) - April 14th
	if date.Month() == time.April && date.Day() == 14 {
		events = append(events, api.Event{
			Name:         "Tamil New Year",
			NameLocal:    "தமிழ் புத்தாண்டு",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Tamil New Year celebrated on the first day of Chithirai month",
			Region:       api.RegionTamilNadu,
			Metadata: map[string]interface{}{
				"importance":         "highest",
				"month":              "Chithirai",
				"solar_event":        true,
				"pan_tamil_festival": true,
			},
		})
	}

	// Pongal - January 14-17
	if date.Month() == time.January && date.Day() >= 14 && date.Day() <= 17 {
		pongalDay := ""
		switch date.Day() {
		case 14:
			pongalDay = "Bhogi Pongal"
		case 15:
			pongalDay = "Thai Pongal"
		case 16:
			pongalDay = "Maattu Pongal"
		case 17:
			pongalDay = "Kaanum Pongal"
		}

		events = append(events, api.Event{
			Name:         pongalDay,
			NameLocal:    "பொங்கல்",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Tamil harvest festival celebrating the Sun God",
			Region:       api.RegionTamilNadu,
			Metadata: map[string]interface{}{
				"importance":       "highest",
				"pongal_day":       pongalDay,
				"harvest_festival": true,
				"duration_days":    4,
			},
		})
	}

	return events, nil
}

// GetRegionalMuhurtas returns Tamil Nadu specific muhurtas
func (t *TamilNaduRegionalPlugin) GetRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location) ([]api.Muhurta, error) {
	if !t.enabled {
		return nil, fmt.Errorf("tamil_nadu_regional_plugin is not enabled")
	}

	return nil, fmt.Errorf("regional muhurtas are unavailable for %s; use muhurta_plugin", t.GetRegion())
}

// GetRegionalNames returns localized Tamil names
func (t *TamilNaduRegionalPlugin) GetRegionalNames(locale string) map[string]string {
	tamilNames := map[string]string{
		"Sunday":    "ஞாயிறு",
		"Monday":    "திங்கள்",
		"Tuesday":   "செவ்வாய்",
		"Wednesday": "புதன்",
		"Thursday":  "வியாழன்",
		"Friday":    "வெள்ளி",
		"Saturday":  "சனி",
	}
	return tamilNames
}
