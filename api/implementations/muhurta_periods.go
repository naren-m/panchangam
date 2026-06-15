package implementations

import (
	"time"

	"github.com/naren-m/panchangam/api"
)

// calculateRahuKalam calculates Rahu Kalam based on weekday and sunrise
func (m *MuhurtaPlugin) calculateRahuKalam(sunrise time.Time, dayLength time.Duration, weekday time.Weekday) api.Muhurta {
	// Rahu Kalam calculation based on traditional Vedic astronomy
	// Each day is divided into 8 parts, Rahu Kalam occupies one part based on weekday

	partDuration := dayLength / 8
	var rahuPart int

	switch weekday {
	case time.Sunday:
		rahuPart = 4 // 5th part (12:00-1:30 PM approximately)
	case time.Monday:
		rahuPart = 1 // 2nd part (7:30-9:00 AM approximately)
	case time.Tuesday:
		rahuPart = 6 // 7th part (3:00-4:30 PM approximately)
	case time.Wednesday:
		rahuPart = 3 // 4th part (10:30 AM-12:00 PM approximately)
	case time.Thursday:
		rahuPart = 5 // 6th part (1:30-3:00 PM approximately)
	case time.Friday:
		rahuPart = 2 // 3rd part (9:00-10:30 AM approximately)
	case time.Saturday:
		rahuPart = 7 // 8th part (4:30-6:00 PM approximately)
	}

	startTime := sunrise.Add(time.Duration(rahuPart) * partDuration)
	endTime := startTime.Add(partDuration)

	return api.Muhurta{
		Name:         "Rahu Kalam",
		NameLocal:    "राहु काल",
		StartTime:    startTime,
		EndTime:      endTime,
		Quality:      api.QualityInauspicious,
		Purpose:      []string{}, // Not suitable for any auspicious activities
		Avoid:        []string{"all_auspicious_activities", "new_ventures", "travel", "ceremonies"},
		Significance: "Inauspicious period ruled by Rahu, avoid starting new activities",
		Metadata: map[string]interface{}{
			"planetary_ruler": "Rahu",
			"weekday_part":    rahuPart + 1,
			"total_parts":     8,
			"calculation":     "vedic_astronomy",
		},
	}
}

// calculateYamagandam calculates Yamagandam based on weekday and sunrise
func (m *MuhurtaPlugin) calculateYamagandam(sunrise time.Time, dayLength time.Duration, weekday time.Weekday) api.Muhurta {
	// Yamagandam calculation - another inauspicious period
	partDuration := dayLength / 8
	var yamaPart int

	switch weekday {
	case time.Sunday:
		yamaPart = 2 // 3rd part
	case time.Monday:
		yamaPart = 5 // 6th part
	case time.Tuesday:
		yamaPart = 0 // 1st part
	case time.Wednesday:
		yamaPart = 4 // 5th part
	case time.Thursday:
		yamaPart = 6 // 7th part
	case time.Friday:
		yamaPart = 3 // 4th part
	case time.Saturday:
		yamaPart = 1 // 2nd part
	}

	startTime := sunrise.Add(time.Duration(yamaPart) * partDuration)
	endTime := startTime.Add(partDuration)

	return api.Muhurta{
		Name:         "Yamagandam",
		NameLocal:    "यमगण्डम्",
		StartTime:    startTime,
		EndTime:      endTime,
		Quality:      api.QualityInauspicious,
		Purpose:      []string{},
		Avoid:        []string{"important_decisions", "travel", "business_deals", "ceremonies"},
		Significance: "Inauspicious period ruled by Yama, lord of death",
		Metadata: map[string]interface{}{
			"planetary_ruler": "Yama",
			"weekday_part":    yamaPart + 1,
			"total_parts":     8,
		},
	}
}

// calculateGulikaKalam calculates Gulika Kalam (similar to Rahu Kalam but different timing)
func (m *MuhurtaPlugin) calculateGulikaKalam(sunrise time.Time, dayLength time.Duration, weekday time.Weekday) api.Muhurta {
	partDuration := dayLength / 8
	var gulikaPart int

	switch weekday {
	case time.Sunday:
		gulikaPart = 6 // 7th part
	case time.Monday:
		gulikaPart = 3 // 4th part
	case time.Tuesday:
		gulikaPart = 4 // 5th part
	case time.Wednesday:
		gulikaPart = 5 // 6th part
	case time.Thursday:
		gulikaPart = 2 // 3rd part
	case time.Friday:
		gulikaPart = 7 // 8th part
	case time.Saturday:
		gulikaPart = 0 // 1st part
	}

	startTime := sunrise.Add(time.Duration(gulikaPart) * partDuration)
	endTime := startTime.Add(partDuration)

	return api.Muhurta{
		Name:         "Gulika Kalam",
		NameLocal:    "गुलिक काल",
		StartTime:    startTime,
		EndTime:      endTime,
		Quality:      api.QualityInauspicious,
		Purpose:      []string{},
		Avoid:        []string{"financial_transactions", "investments", "property_deals"},
		Significance: "Inauspicious period ruled by Gulika (son of Shani)",
		Metadata: map[string]interface{}{
			"planetary_ruler": "Gulika",
			"weekday_part":    gulikaPart + 1,
			"total_parts":     8,
		},
	}
}

// calculateAbhijitMuhurta calculates the auspicious Abhijit period around noon
func (m *MuhurtaPlugin) calculateAbhijitMuhurta(sunrise, sunset time.Time) api.Muhurta {
	// Abhijit is the 8th nakshatra period, approximately 24 minutes around solar noon
	dayLength := sunset.Sub(sunrise)
	solarNoon := sunrise.Add(dayLength / 2)

	// Abhijit duration is approximately 24 minutes (1/60th of a day)
	abhijitDuration := 24 * time.Minute
	startTime := solarNoon.Add(-abhijitDuration / 2)
	endTime := solarNoon.Add(abhijitDuration / 2)

	return api.Muhurta{
		Name:         "Abhijit Muhurta",
		NameLocal:    "अभिजित मुहूर्त",
		StartTime:    startTime,
		EndTime:      endTime,
		Quality:      api.QualityHighly,
		Purpose:      []string{"all_auspicious_activities", "business_ventures", "education", "spiritual_practices"},
		Avoid:        []string{},
		Significance: "Highly auspicious period around solar noon, victory and success assured",
		Metadata: map[string]interface{}{
			"nakshatra":    "Abhijit",
			"duration":     "24_minutes",
			"calculation":  "solar_noon_based",
			"deity":        "Brahma",
			"significance": "victory_success",
		},
	}
}

// calculateBrahmaMuhurta calculates the pre-dawn auspicious period
func (m *MuhurtaPlugin) calculateBrahmaMuhurta(sunrise time.Time, date time.Time) api.Muhurta {
	// Brahma Muhurta is approximately 1.5 hours before sunrise
	// It's considered the most auspicious time for spiritual practices

	brahmaDuration := 96 * time.Minute // 1 hour 36 minutes (1/15th of day-night cycle)
	endTime := sunrise
	startTime := endTime.Add(-brahmaDuration)

	// Ensure it's on the same date (early morning)
	if startTime.Day() != date.Day() {
		startTime = time.Date(date.Year(), date.Month(), date.Day(), 4, 0, 0, 0, date.Location())
		endTime = startTime.Add(brahmaDuration)
	}

	return api.Muhurta{
		Name:         "Brahma Muhurta",
		NameLocal:    "ब्रह्म मुहूर्त",
		StartTime:    startTime,
		EndTime:      endTime,
		Quality:      api.QualityHighly,
		Purpose:      []string{"meditation", "prayer", "study", "yoga", "spiritual_practices"},
		Avoid:        []string{},
		Significance: "Most auspicious pre-dawn period for spiritual activities",
		Metadata: map[string]interface{}{
			"duration":               "96_minutes",
			"calculation":            "pre_sunrise",
			"deity":                  "Brahma",
			"time_of_day":            "pre_dawn",
			"spiritual_significance": "highest",
		},
	}
}
