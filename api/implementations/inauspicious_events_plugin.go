package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy"
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

// GetEvents returns inauspicious events for a specific date and location
func (i *InauspiciousEventsPlugin) GetEvents(ctx context.Context, date time.Time, location api.Location, region api.Region) ([]api.Event, error) {
	if !i.enabled {
		return nil, fmt.Errorf("inauspicious events plugin is not enabled")
	}

	var events []api.Event

	// Calculate sunrise and sunset for the location
	astroLocation := astronomy.Location{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
	
	sunTimes, err := astronomy.CalculateSunTimes(astroLocation, date)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate sun times: %w", err)
	}

	// Convert to local timezone if provided
	if location.Timezone != "" {
		if tz, err := time.LoadLocation(location.Timezone); err == nil {
			sunTimes.Sunrise = sunTimes.Sunrise.In(tz)
			sunTimes.Sunset = sunTimes.Sunset.In(tz)
		}
	}

	// Calculate day length
	dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)
	weekday := date.Weekday()

	// Generate Rahu Kalam event
	rahuKalam := i.calculateRahuKalamEvent(sunTimes.Sunrise, dayLength, weekday, region)
	events = append(events, rahuKalam)

	// Generate Yamagandam event
	yamagandam := i.calculateYamgandamEvent(sunTimes.Sunrise, dayLength, weekday, region)
	events = append(events, yamagandam)

	// Generate Gulika Kalam event
	gulikaKalam := i.calculateGulikaKalamEvent(sunTimes.Sunrise, dayLength, weekday, region)
	events = append(events, gulikaKalam)

	return events, nil
}

// GetEventsInRange returns inauspicious events for a date range
func (i *InauspiciousEventsPlugin) GetEventsInRange(ctx context.Context, start, end time.Time, location api.Location, region api.Region) ([]api.Event, error) {
	var allEvents []api.Event

	current := start
	for current.Before(end) || current.Equal(end) {
		dayEvents, err := i.GetEvents(ctx, current, location, region)
		if err != nil {
			return nil, fmt.Errorf("failed to get events for %s: %w", current.Format("2006-01-02"), err)
		}
		allEvents = append(allEvents, dayEvents...)
		current = current.AddDate(0, 0, 1)
	}

	return allEvents, nil
}

// calculateRahuKalamEvent calculates Rahu Kalam as an event
func (i *InauspiciousEventsPlugin) calculateRahuKalamEvent(sunrise time.Time, dayLength time.Duration, weekday time.Weekday, region api.Region) api.Event {
	// Divide day into 8 equal parts
	partDuration := dayLength / 8
	
	// Determine which part belongs to Rahu based on weekday
	var rahuPart int
	var weekdayName string
	
	switch weekday {
	case time.Sunday:
		rahuPart = 4 // 5th part
		weekdayName = "Sunday"
	case time.Monday:
		rahuPart = 1 // 2nd part  
		weekdayName = "Monday"
	case time.Tuesday:
		rahuPart = 6 // 7th part
		weekdayName = "Tuesday"
	case time.Wednesday:
		rahuPart = 3 // 4th part
		weekdayName = "Wednesday"
	case time.Thursday:
		rahuPart = 5 // 6th part
		weekdayName = "Thursday"
	case time.Friday:
		rahuPart = 2 // 3rd part
		weekdayName = "Friday"
	case time.Saturday:
		rahuPart = 7 // 8th part
		weekdayName = "Saturday"
	}
	
	startTime := sunrise.Add(time.Duration(rahuPart) * partDuration)
	endTime := startTime.Add(partDuration)
	
	// Regional names
	localNames := map[api.Region]string{
		api.RegionTamilNadu:  "ராகு காலம்",
		api.RegionKerala:     "രാഹു കാലം",
		api.RegionBengal:     "রাহু কাল",
		api.RegionGujarat:    "રાહુ કાળ",
		api.RegionMaha:       "राहू काळ",
		api.RegionNorthIndia: "राहु काल",
		api.RegionSouthIndia: "राहु काल",
		api.RegionGlobal:     "राहु काल",
	}
	
	return api.Event{
		Name:         "Rahu Kalam",
		NameLocal:    localNames[region],
		Type:         api.EventTypeRahukalam,
		StartTime:    startTime,
		EndTime:      endTime,
		Significance: "Inauspicious period ruled by Rahu. Avoid starting new ventures, travel, and important activities.",
		Region:       region,
		Metadata: map[string]interface{}{
			"planetary_ruler":    "Rahu",
			"weekday":           weekdayName,
			"weekday_part":      rahuPart + 1,
			"total_parts":       8,
			"duration_minutes":  int(partDuration.Minutes()),
			"calculation_method": "vedic_eight_parts",
			"warnings": []string{
				"Avoid starting new business",
				"Avoid important meetings",
				"Avoid travel",
				"Avoid ceremonies",
				"Avoid financial transactions",
			},
			"traditional_belief": "Period when Rahu's malefic influence is strongest",
			"modern_usage":      "Time for reflection, rest, or routine maintenance work",
		},
	}
}

// calculateYamgandamEvent calculates Yamagandam as an event
func (i *InauspiciousEventsPlugin) calculateYamgandamEvent(sunrise time.Time, dayLength time.Duration, weekday time.Weekday, region api.Region) api.Event {
	partDuration := dayLength / 8
	
	var yamaPart int
	var weekdayName string
	
	switch weekday {
	case time.Sunday:
		yamaPart = 2 // 3rd part
		weekdayName = "Sunday"
	case time.Monday:
		yamaPart = 5 // 6th part
		weekdayName = "Monday"
	case time.Tuesday:
		yamaPart = 0 // 1st part
		weekdayName = "Tuesday"
	case time.Wednesday:
		yamaPart = 4 // 5th part
		weekdayName = "Wednesday"
	case time.Thursday:
		yamaPart = 6 // 7th part
		weekdayName = "Thursday"
	case time.Friday:
		yamaPart = 3 // 4th part
		weekdayName = "Friday"
	case time.Saturday:
		yamaPart = 1 // 2nd part
		weekdayName = "Saturday"
	}
	
	startTime := sunrise.Add(time.Duration(yamaPart) * partDuration)
	endTime := startTime.Add(partDuration)
	
	// Regional names
	localNames := map[api.Region]string{
		api.RegionTamilNadu:  "யமகண்டம்",
		api.RegionKerala:     "യമഗണ്ഡം",
		api.RegionBengal:     "যমগণ্ডম",
		api.RegionGujarat:    "યમગણ્ડમ",
		api.RegionMaha:       "यमगंडम",
		api.RegionNorthIndia: "यमगण्डम्",
		api.RegionSouthIndia: "यमगण्डम्",
		api.RegionGlobal:     "यमगण्डम्",
	}
	
	return api.Event{
		Name:         "Yamagandam",
		NameLocal:    localNames[region],
		Type:         api.EventTypeYamagandam,
		StartTime:    startTime,
		EndTime:      endTime,
		Significance: "Inauspicious period ruled by Yama (Lord of Death). Avoid major decisions and important activities.",
		Region:       region,
		Metadata: map[string]interface{}{
			"planetary_ruler":    "Yama",
			"weekday":           weekdayName,
			"weekday_part":      yamaPart + 1,
			"total_parts":       8,
			"duration_minutes":  int(partDuration.Minutes()),
			"calculation_method": "vedic_eight_parts",
			"warnings": []string{
				"Avoid major decisions",
				"Avoid signing contracts",
				"Avoid court proceedings",
				"Avoid surgery or medical procedures",
				"Avoid arguments or confrontations",
			},
			"traditional_belief": "Time when Yama's influence brings obstacles and delays",
			"alternative_activities": []string{
				"Complete pending tasks",
				"Administrative work",
				"Cleaning and organizing",
				"Planning future activities",
			},
		},
	}
}

// calculateGulikaKalamEvent calculates Gulika Kalam as an event
func (i *InauspiciousEventsPlugin) calculateGulikaKalamEvent(sunrise time.Time, dayLength time.Duration, weekday time.Weekday, region api.Region) api.Event {
	partDuration := dayLength / 8
	
	var gulikaPart int
	var weekdayName string
	
	switch weekday {
	case time.Sunday:
		gulikaPart = 6 // 7th part
		weekdayName = "Sunday"
	case time.Monday:
		gulikaPart = 3 // 4th part
		weekdayName = "Monday"
	case time.Tuesday:
		gulikaPart = 4 // 5th part
		weekdayName = "Tuesday"
	case time.Wednesday:
		gulikaPart = 5 // 6th part
		weekdayName = "Wednesday"
	case time.Thursday:
		gulikaPart = 2 // 3rd part
		weekdayName = "Thursday"
	case time.Friday:
		gulikaPart = 7 // 8th part
		weekdayName = "Friday"
	case time.Saturday:
		gulikaPart = 0 // 1st part
		weekdayName = "Saturday"
	}
	
	startTime := sunrise.Add(time.Duration(gulikaPart) * partDuration)
	endTime := startTime.Add(partDuration)
	
	// Regional names
	localNames := map[api.Region]string{
		api.RegionTamilNadu:  "குளிக காலம்",
		api.RegionKerala:     "ഗുളിക കാലം",
		api.RegionBengal:     "গুলিক কাল",
		api.RegionGujarat:    "ગુલિક કાળ",
		api.RegionMaha:       "गुलिक काळ",
		api.RegionNorthIndia: "गुलिक काल",
		api.RegionSouthIndia: "गुलिक काल",
		api.RegionGlobal:     "गुलिक काल",
	}
	
	return api.Event{
		Name:         "Gulika Kalam",
		NameLocal:    localNames[region],
		Type:         api.EventTypeGulikakalam,
		StartTime:    startTime,
		EndTime:      endTime,
		Significance: "Inauspicious period ruled by Gulika (son of Saturn). Particularly unfavorable for financial activities.",
		Region:       region,
		Metadata: map[string]interface{}{
			"planetary_ruler":    "Gulika",
			"parent_planet":     "Saturn",
			"weekday":           weekdayName,
			"weekday_part":      gulikaPart + 1,
			"total_parts":       8,
			"duration_minutes":  int(partDuration.Minutes()),
			"calculation_method": "vedic_eight_parts",
			"warnings": []string{
				"Avoid financial investments",
				"Avoid property transactions",
				"Avoid borrowing or lending money",
				"Avoid starting business partnerships",
				"Avoid gambling or speculation",
			},
			"traditional_belief": "Time when material losses and financial obstacles are likely",
			"severity":          "moderate_to_high",
			"regional_variations": map[string]string{
				"tamil_nadu": "Known as 'Kuli Kalam' - time for caution in monetary matters",
				"kerala":     "Considered particularly important for business decisions",
				"bengal":     "Associated with Saturn's restricting influence",
			},
		},
	}
}

// Helper function to format duration for display
func (i *InauspiciousEventsPlugin) formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%d hour(s) %d minute(s)", hours, minutes)
	}
	return fmt.Sprintf("%d minute(s)", minutes)
}

// GetEventSummary provides a summary of all inauspicious periods for a day
func (i *InauspiciousEventsPlugin) GetEventSummary(ctx context.Context, date time.Time, location api.Location, region api.Region) (map[string]interface{}, error) {
	events, err := i.GetEvents(ctx, date, location, region)
	if err != nil {
		return nil, err
	}
	
	summary := map[string]interface{}{
		"date":         date.Format("2006-01-02"),
		"weekday":      date.Weekday().String(),
		"location":     location.Name,
		"region":       string(region),
		"total_events": len(events),
		"events":       make([]map[string]interface{}, 0, len(events)),
	}
	
	for _, event := range events {
		eventSummary := map[string]interface{}{
			"name":       event.Name,
			"name_local": event.NameLocal,
			"type":       string(event.Type),
			"start_time": event.StartTime.Format("15:04"),
			"end_time":   event.EndTime.Format("15:04"),
			"duration":   i.formatDuration(event.EndTime.Sub(event.StartTime)),
		}
		summary["events"] = append(summary["events"].([]map[string]interface{}), eventSummary)
	}
	
	return summary, nil
}