package examples

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/api"
)

// TamilNaduExtension provides Tamil Nadu specific regional calculations
type TamilNaduExtension struct {
	enabled bool
	config  map[string]interface{}
}

// NewTamilNaduExtension creates a new Tamil Nadu regional extension plugin
func NewTamilNaduExtension() *TamilNaduExtension {
	return &TamilNaduExtension{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

// GetInfo returns plugin metadata
func (t *TamilNaduExtension) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "tamil_nadu_extension",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Tamil Nadu regional calculation extensions",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
			string(api.CapabilityEvent),
			string(api.CapabilityMuhurta),
		},
		Dependencies: []string{},
		Metadata: map[string]interface{}{
			"region":          "tamil_nadu",
			"calendar_system": "amanta",
			"language":        "tamil",
		},
	}
}

// Initialize sets up the plugin with configuration
func (t *TamilNaduExtension) Initialize(ctx context.Context, config map[string]interface{}) error {
	t.config = config
	t.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (t *TamilNaduExtension) IsEnabled() bool {
	return t.enabled
}

// Shutdown cleans up plugin resources
func (t *TamilNaduExtension) Shutdown(ctx context.Context) error {
	t.enabled = false
	return nil
}

// GetRegion returns the region this extension supports
func (t *TamilNaduExtension) GetRegion() api.Region {
	return api.RegionTamilNadu
}

// GetCalendarSystem returns the calendar system used
func (t *TamilNaduExtension) GetCalendarSystem() api.CalendarSystem {
	return api.CalendarAmanta
}

// ApplyRegionalRules applies Tamil Nadu specific rules to Panchangam data
func (t *TamilNaduExtension) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error {
	// Apply Tamil calendar adjustments
	if data.CalendarSystem == api.CalendarAmanta {
		// Adjust tithi calculations for Tamil traditions
		t.adjustTithiForTamilCalendar(&data.Tithi)

		// Apply Tamil nakshatra naming conventions
		t.applyTamilNakshatraNames(&data.Nakshatra)

		// Adjust yoga calculations for Tamil preferences
		t.adjustYogaForTamilTraditions(&data.Yoga)
	}

	return nil
}

// GetRegionalEvents returns Tamil Nadu specific events
func (t *TamilNaduExtension) GetRegionalEvents(ctx context.Context, date time.Time, location api.Location) ([]api.Event, error) {
	var events []api.Event

	// Tamil New Year (Puthandu) - April 13/14
	if t.isTamilNewYear(date) {
		events = append(events, api.Event{
			Name:         "Puthandu",
			NameLocal:    "புத்தாண்டு",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Tamil New Year celebration",
			Region:       api.RegionTamilNadu,
			Metadata: map[string]interface{}{
				"category":   "new_year",
				"importance": "high",
				"traditions": []string{"mango_leaves", "kolam", "special_prayers"},
			},
		})
	}

	// Chithirai Festival (April-May)
	if t.isChithiraiFestival(date) {
		events = append(events, api.Event{
			Name:         "Chithirai Festival",
			NameLocal:    "சித்திரை திருவிழா",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Celebrating the divine marriage of Meenakshi and Sundareshwar",
			Region:       api.RegionTamilNadu,
			Metadata: map[string]interface{}{
				"location": "Madurai",
				"duration": "10 days",
				"deity":    "Meenakshi Amman",
			},
		})
	}

	// Thai Pusam
	if t.isThaiPusam(date) {
		events = append(events, api.Event{
			Name:         "Thai Pusam",
			NameLocal:    "தைப்பூசம்",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Festival dedicated to Lord Murugan",
			Region:       api.RegionTamilNadu,
			Metadata: map[string]interface{}{
				"deity":   "Murugan",
				"star":    "Pusam",
				"month":   "Thai",
				"rituals": []string{"kavadi", "piercing", "fasting"},
			},
		})
	}

	return events, nil
}

// GetRegionalMuhurtas returns Tamil Nadu specific muhurtas
func (t *TamilNaduExtension) GetRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location) ([]api.Muhurta, error) {
	var muhurtas []api.Muhurta

	// Abhijit Muhurta (noon time - highly auspicious in Tamil tradition)
	abhijitStart := time.Date(date.Year(), date.Month(), date.Day(), 11, 44, 0, 0, date.Location())
	abhijitEnd := abhijitStart.Add(48 * time.Minute)

	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Abhijit Muhurta",
		NameLocal:    "அபிஜித் முகூர்த்தம்",
		StartTime:    abhijitStart,
		EndTime:      abhijitEnd,
		Quality:      api.QualityHighly,
		Purpose:      []string{"new_ventures", "important_decisions", "travel", "ceremonies"},
		Significance: "Most auspicious time of the day in Tamil tradition",
		Region:       api.RegionTamilNadu,
		Metadata: map[string]interface{}{
			"duration_minutes": 48,
			"daily_occurrence": true,
			"lord":             "Surya",
		},
	})

	// Brahma Muhurta (early morning - ideal for spiritual practices)
	brahmaStart := time.Date(date.Year(), date.Month(), date.Day(), 4, 30, 0, 0, date.Location())
	brahmaEnd := brahmaStart.Add(96 * time.Minute)

	muhurtas = append(muhurtas, api.Muhurta{
		Name:         "Brahma Muhurta",
		NameLocal:    "பிரம்ம முகூர்த்தம்",
		StartTime:    brahmaStart,
		EndTime:      brahmaEnd,
		Quality:      api.QualityAuspicious,
		Purpose:      []string{"meditation", "prayers", "study", "spiritual_practices"},
		Significance: "Divine time for spiritual activities",
		Region:       api.RegionTamilNadu,
		Metadata: map[string]interface{}{
			"duration_minutes": 96,
			"spiritual_value":  "highest",
			"recommended_for":  "daily_practice",
		},
	})

	return muhurtas, nil
}

// GetRegionalNames returns localized names for Panchangam elements
func (t *TamilNaduExtension) GetRegionalNames(locale string) map[string]string {
	if locale == "ta" || locale == "tamil" {
		return map[string]string{
			// Days of the week
			"Sunday":    "ஞாயிறு",
			"Monday":    "திங்கள்",
			"Tuesday":   "செவ்வாய்",
			"Wednesday": "புதன்",
			"Thursday":  "வியாழன்",
			"Friday":    "வெள்ளி",
			"Saturday":  "சனி",

			// Tithis
			"Pratipada": "பிரதமை",
			"Dwitiya":   "துவிதியை",
			"Tritiya":   "திருதியை",
			"Chaturthi": "சதுர்த்தி",
			"Panchami":  "பஞ்சமி",

			// Nakshatras
			"Ashwini":   "அசுவினி",
			"Bharani":   "பரணி",
			"Krittika":  "கார்த்திகை",
			"Rohini":    "ரோகிணி",
			"Mrigasira": "மிருகசீர்ஷம்",

			// Yogas
			"Vishkambha": "விஷ்கம்பா",
			"Preeti":     "ப்ரீதி",
			"Ayushman":   "ஆயுஷ்மான்",

			// Karanas
			"Bava":    "பவ",
			"Balava":  "பாலவ",
			"Kaulava": "கௌலவ",
		}
	}

	return make(map[string]string)
}

// Helper methods for date calculations

func (t *TamilNaduExtension) isTamilNewYear(date time.Time) bool {
	// Tamil New Year typically falls on April 13/14
	return (date.Month() == time.April && (date.Day() == 13 || date.Day() == 14))
}

func (t *TamilNaduExtension) isChithiraiFestival(date time.Time) bool {
	// Chithirai festival in Madurai - typically in April/May
	// This is a simplified check - actual calculation would involve lunar calendar
	return date.Month() == time.April && date.Day() >= 15 && date.Day() <= 25
}

func (t *TamilNaduExtension) isThaiPusam(date time.Time) bool {
	// Thai Pusam falls during Pusam nakshatra in Thai month (Jan/Feb)
	// This is a simplified check - actual calculation would involve nakshatra positions
	return date.Month() == time.January && date.Day() >= 25 || (date.Month() == time.February && date.Day() <= 5)
}

func (t *TamilNaduExtension) adjustTithiForTamilCalendar(tithi *api.Tithi) {
	// Apply Tamil-specific tithi adjustments
	// Tamil calendar follows Amanta system where month ends on new moon
	// This affects how tithis are counted and named

	if tithi.Number == 15 { // Amavasya in Amanta system
		tithi.Quality = "new_moon_day"
		tithi.Lord = "Shiva"
	}
}

func (t *TamilNaduExtension) applyTamilNakshatraNames(nakshatra *api.Nakshatra) {
	// Apply Tamil naming conventions for nakshatras
	tamilNames := t.GetRegionalNames("tamil")
	if localName, exists := tamilNames[nakshatra.Name]; exists {
		nakshatra.NameLocal = localName
	}
}

func (t *TamilNaduExtension) adjustYogaForTamilTraditions(yoga *api.Yoga) {
	// Apply Tamil-specific yoga interpretations
	// Different regions may have varying interpretations of yoga qualities

	// Tamil tradition places special emphasis on certain yogas
	specialYogas := map[string]string{
		"Siddha":    "highly_auspicious",
		"Sadhya":    "moderately_auspicious",
		"Subha":     "auspicious",
		"Sukla":     "pure",
		"Brahma":    "divine",
		"Mahendra":  "royal",
		"Vaidhriti": "avoid",
		"Vyaghata":  "avoid",
	}

	if quality, exists := specialYogas[yoga.Name]; exists {
		yoga.Quality = quality
	}
}
