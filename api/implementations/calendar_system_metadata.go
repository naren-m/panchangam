package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
)

func (c *CalendarSystemPlugin) adjustMonthData(ctx context.Context, data *api.PanchangamData, monthInfo *MonthInfo) error {
	// PanchangamData has no top-level metadata field, so month details are
	// exposed as a lunar event with machine-readable metadata.
	monthEvent := api.Event{
		Name:         "Lunar Month",
		NameLocal:    monthInfo.NameLocal,
		Type:         api.EventTypeLunar,
		StartTime:    monthInfo.StartDate,
		EndTime:      monthInfo.EndDate,
		Significance: fmt.Sprintf("%s month in %s calendar system", monthInfo.Name, monthInfo.CalendarSystem),
		Region:       data.Region,
		Metadata: map[string]interface{}{
			"type":            "month_info",
			"name":            monthInfo.Name,
			"name_local":      monthInfo.NameLocal,
			"number":          monthInfo.Number,
			"calendar_system": string(monthInfo.CalendarSystem),
			"is_adhika_masa":  monthInfo.IsAdhikaMasa,
		},
	}

	data.Events = append(data.Events, monthEvent)
	return nil
}

func (c *CalendarSystemPlugin) addCalendarSystemMetadata(data *api.PanchangamData, monthInfo *MonthInfo) {
	// Calendar-system details use the same event-metadata contract as month details.
	calendarSystemEvent := api.Event{
		Name:         "Calendar System Info",
		NameLocal:    "",
		Type:         api.EventTypeLunar,
		StartTime:    data.Date,
		EndTime:      data.Date.Add(24 * time.Hour),
		Significance: c.getCalendarSystemDescription(data.CalendarSystem),
		Region:       data.Region,
		Metadata: map[string]interface{}{
			"type":                  "calendar_system_info",
			"system":                string(data.CalendarSystem),
			"description":           c.getCalendarSystemDescription(data.CalendarSystem),
			"month_boundary_rule":   c.getMonthBoundaryRule(data.CalendarSystem),
			"regional_preference":   c.getRegionalPreference(data.Region),
			"calculation_precision": "astronomical",
		},
	}

	data.Events = append(data.Events, calendarSystemEvent)
}

func (c *CalendarSystemPlugin) getCalendarSystemDescription(system api.CalendarSystem) string {
	descriptions := map[api.CalendarSystem]string{
		api.CalendarAmanta:     "South Indian system where lunar months end on Amavasya (New Moon)",
		api.CalendarPurnimanta: "North Indian system where lunar months end on Purnima (Full Moon)",
		api.CalendarLunar:      "Pure lunar calendar following regional preferences",
		api.CalendarSolar:      "Solar calendar based on sun's movement through zodiac signs",
	}
	return descriptions[system]
}

func (c *CalendarSystemPlugin) getMonthBoundaryRule(system api.CalendarSystem) string {
	rules := map[api.CalendarSystem]string{
		api.CalendarAmanta:     "Month begins day after Amavasya, ends on next Amavasya",
		api.CalendarPurnimanta: "Month begins day after Purnima, ends on next Purnima",
		api.CalendarLunar:      "Follows regional calendar system preference",
		api.CalendarSolar:      "Month begins when sun enters new zodiac sign",
	}
	return rules[system]
}

func (c *CalendarSystemPlugin) getRegionalPreference(region api.Region) string {
	preferences := map[api.Region]string{
		api.RegionNorthIndia: "Purnimanta",
		api.RegionSouthIndia: "Amanta",
		api.RegionTamilNadu:  "Amanta",
		api.RegionKerala:     "Amanta",
		api.RegionBengal:     "Amanta",
		api.RegionGujarat:    "Purnimanta",
		api.RegionMaha:       "Purnimanta",
		api.RegionGlobal:     "Mixed (region-dependent)",
	}
	return preferences[region]
}
