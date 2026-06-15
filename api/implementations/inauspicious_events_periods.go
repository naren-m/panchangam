package implementations

import (
	"time"

	"github.com/naren-m/panchangam/api"
)

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
			"weekday":            weekdayName,
			"weekday_part":       rahuPart + 1,
			"total_parts":        8,
			"duration_minutes":   int(partDuration.Minutes()),
			"calculation_method": "vedic_eight_parts",
			"warnings": []string{
				"Avoid starting new business",
				"Avoid important meetings",
				"Avoid travel",
				"Avoid ceremonies",
				"Avoid financial transactions",
			},
			"traditional_belief": "Period when Rahu's malefic influence is strongest",
			"modern_usage":       "Time for reflection, rest, or routine maintenance work",
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
			"weekday":            weekdayName,
			"weekday_part":       yamaPart + 1,
			"total_parts":        8,
			"duration_minutes":   int(partDuration.Minutes()),
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
			"parent_planet":      "Saturn",
			"weekday":            weekdayName,
			"weekday_part":       gulikaPart + 1,
			"total_parts":        8,
			"duration_minutes":   int(partDuration.Minutes()),
			"calculation_method": "vedic_eight_parts",
			"warnings": []string{
				"Avoid financial investments",
				"Avoid property transactions",
				"Avoid borrowing or lending money",
				"Avoid starting business partnerships",
				"Avoid gambling or speculation",
			},
			"traditional_belief": "Time when material losses and financial obstacles are likely",
			"severity":           "moderate_to_high",
			"regional_variations": map[string]string{
				"tamil_nadu": "Known as 'Kuli Kalam' - time for caution in monetary matters",
				"kerala":     "Considered particularly important for business decisions",
				"bengal":     "Associated with Saturn's restricting influence",
			},
		},
	}
}
