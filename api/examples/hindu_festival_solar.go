package examples

import (
	"time"

	"github.com/naren-m/panchangam/api"
)

func (h *HinduFestivalPlugin) getSolarFestivals(date time.Time, region api.Region) []api.Event {
	var events []api.Event

	if h.isMakarSankranti(date) {
		events = append(events, api.Event{
			Name:         "Makar Sankranti",
			NameLocal:    "मकर संक्रांति",
			Type:         api.EventTypeSolar,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Sun's transition into Capricorn, marking the end of winter solstice",
			Region:       region,
			Metadata: map[string]interface{}{
				"importance":      "high",
				"solar_event":     "sun_capricorn_entry",
				"seasonal_change": "winter_to_spring",
				"rituals":         []string{"kite_flying", "til_gud", "holy_bath", "charity"},
				"regional_names": map[string]string{
					"tamil_nadu": "Pongal",
					"punjab":     "Lohri",
					"assam":      "Magh Bihu",
					"kerala":     "Makara Vilakku",
					"karnataka":  "Makara Sankramana",
					"gujarat":    "Uttarayan",
				},
				"special_foods": []string{"til_gud_laddu", "khichdi", "jaggery_sweets"},
			},
		})
	}

	if h.isRamNavami(date) {
		events = append(events, api.Event{
			Name:         "Ram Navami",
			NameLocal:    "राम नवमी",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Celebration of Lord Rama's birth",
			Region:       region,
			Metadata: map[string]interface{}{
				"importance":          "high",
				"lunar_month":         "Chaitra",
				"lunar_day":           "Navami",
				"deity":               "Rama",
				"rituals":             []string{"rama_bhajan", "temple_visit", "fasting", "procession"},
				"sacred_places":       []string{"Ayodhya", "Rameswaram", "Bhadrachalam"},
				"special_recitations": []string{"Ramayana", "Rama_Chalisa", "Rama_Stotram"},
			},
		})
	}

	return events
}
