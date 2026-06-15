package implementations

import (
	"time"

	"github.com/naren-m/panchangam/api"
)

// Month name and numbering logic

func (c *CalendarSystemPlugin) getAmantaMonthInfo(date time.Time, region api.Region) (int, string) {
	// Amanta month names based on when the month's Amavasya falls
	// This is a simplified implementation
	monthNames := []string{
		"Chaitra", "Vaisakha", "Jyeshtha", "Ashadha",
		"Shravana", "Bhadrapada", "Ashwin", "Kartik",
		"Margashirsha", "Pausha", "Magha", "Phalgun",
	}

	// Simplified mapping based on Gregorian calendar
	// In actual implementation, this would be based on precise lunar calculations
	monthIndex := (int(date.Month()) + 10) % 12 // Rough approximation
	return monthIndex + 1, monthNames[monthIndex]
}

func (c *CalendarSystemPlugin) getPurnimantaMonthInfo(date time.Time, region api.Region) (int, string) {
	// Purnimanta month names based on when the month's Purnima falls
	// The key difference is that Purnimanta months start ~15 days later than Amanta
	monthNames := []string{
		"Chaitra", "Vaisakha", "Jyeshtha", "Ashadha",
		"Shravana", "Bhadrapada", "Ashwin", "Kartik",
		"Margashirsha", "Pausha", "Magha", "Phalgun",
	}

	// For Purnimanta, adjust by approximately 15 days
	adjustedDate := date.AddDate(0, 0, -15)
	monthIndex := (int(adjustedDate.Month()) + 10) % 12
	return monthIndex + 1, monthNames[monthIndex]
}

func (c *CalendarSystemPlugin) getSolarMonthInfo(sunLongitude float64, region api.Region) (int, string) {
	// Solar months based on sun's zodiac position
	solarMonths := []string{
		"Mesha", "Vrishabha", "Mithuna", "Karkataka",
		"Simha", "Kanya", "Tula", "Vrishchika",
		"Dhanus", "Makara", "Kumbha", "Meena",
	}

	// Each zodiac sign is 30 degrees
	monthIndex := int(sunLongitude / 30.0)
	if monthIndex >= 12 {
		monthIndex = 11
	}

	return monthIndex + 1, solarMonths[monthIndex]
}

func (c *CalendarSystemPlugin) getLocalMonthName(sanskritName string, region api.Region) string {
	// Regional month name mappings
	monthMappings := map[api.Region]map[string]string{
		api.RegionTamilNadu: {
			"Chaitra":      "சித்திரை",
			"Vaisakha":     "வைகாசி",
			"Jyeshtha":     "ஜெயிஷ்டா",
			"Ashadha":      "ஆஷாட",
			"Shravana":     "ஸ்ராவண",
			"Bhadrapada":   "பத்ரபத",
			"Ashwin":       "ஆஸ்வின",
			"Kartik":       "கார்த்திக",
			"Margashirsha": "மார்கழி",
			"Pausha":       "பௌஷ",
			"Magha":        "மாக",
			"Phalgun":      "பால்குன",
		},
		api.RegionBengal: {
			"Chaitra":      "চৈত্র",
			"Vaisakha":     "বৈশাখ",
			"Jyeshtha":     "জ্যৈষ্ঠ",
			"Ashadha":      "আষাঢ়",
			"Shravana":     "শ্রাবণ",
			"Bhadrapada":   "ভাদ্র",
			"Ashwin":       "আশ্বিন",
			"Kartik":       "কার্তিক",
			"Margashirsha": "অগ্রহায়ণ",
			"Pausha":       "পৌষ",
			"Magha":        "মাঘ",
			"Phalgun":      "ফাল্গুন",
		},
		// Add more regional mappings as needed
	}

	if regionalMap, exists := monthMappings[region]; exists {
		if localName, exists := regionalMap[sanskritName]; exists {
			return localName
		}
	}

	return sanskritName // Fallback to Sanskrit name
}

// Helper methods

func (c *CalendarSystemPlugin) isNorthIndianRegion(region api.Region) bool {
	northIndianRegions := map[api.Region]bool{
		api.RegionNorthIndia: true,
		api.RegionGujarat:    true,
		api.RegionMaha:       true,
	}
	return northIndianRegions[region]
}
