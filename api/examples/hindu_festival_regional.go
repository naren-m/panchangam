package examples

import (
	"time"

	"github.com/naren-m/panchangam/api"
)

func (h *HinduFestivalPlugin) getRegionalFestivals(date time.Time, region api.Region) []api.Event {
	var events []api.Event

	switch region {
	case api.RegionTamilNadu:
		events = append(events, h.getTamilFestivals(date)...)
	case api.RegionKerala:
		events = append(events, h.getKeralaFestivals(date)...)
	case api.RegionBengal:
		events = append(events, h.getBengalFestivals(date)...)
	}

	return events
}

func (h *HinduFestivalPlugin) getTamilFestivals(date time.Time) []api.Event {
	var events []api.Event

	if h.isPongal(date) {
		pongalDay := h.getPongalDay(date)
		events = append(events, api.Event{
			Name:         pongalDay,
			NameLocal:    h.getPongalTamilName(pongalDay),
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Tamil harvest festival celebrating nature and prosperity",
			Region:       api.RegionTamilNadu,
			Metadata: map[string]interface{}{
				"festival_type":    "harvest",
				"duration_days":    4,
				"current_day":      pongalDay,
				"traditional_dish": "pongal_rice",
				"rituals":          h.getPongalRituals(pongalDay),
			},
		})
	}

	return events
}

func (h *HinduFestivalPlugin) getKeralaFestivals(date time.Time) []api.Event {
	var events []api.Event

	if h.isOnam(date) {
		onamDay := h.getOnamDay(date)
		events = append(events, api.Event{
			Name:         "Onam",
			NameLocal:    "ഓണം",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Kerala harvest festival celebrating King Mahabali's return",
			Region:       api.RegionKerala,
			Metadata: map[string]interface{}{
				"festival_type":    "harvest",
				"duration_days":    10,
				"current_day":      onamDay,
				"main_deity":       "Mahabali",
				"traditional_meal": "Onam_Sadhya",
				"cultural_events":  []string{"Kathakali", "Theyyam", "boat_race"},
			},
		})
	}

	return events
}

func (h *HinduFestivalPlugin) getBengalFestivals(date time.Time) []api.Event {
	var events []api.Event

	if h.isDurgaPuja(date) {
		events = append(events, api.Event{
			Name:         "Durga Puja",
			NameLocal:    "দুর্গা পূজা",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Bengal's grandest festival celebrating Goddess Durga",
			Region:       api.RegionBengal,
			Metadata: map[string]interface{}{
				"duration_days":         5,
				"main_deity":            "Durga",
				"cultural_significance": "homecoming_of_daughter",
				"rituals":               []string{"pandal_hopping", "dhunuchi_dance", "sindoor_khela"},
				"traditional_food":      []string{"bhog", "khichuri", "payesh"},
			},
		})
	}

	return events
}
