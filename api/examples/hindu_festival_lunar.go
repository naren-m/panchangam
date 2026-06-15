package examples

import (
	"time"

	"github.com/naren-m/panchangam/api"
)

func (h *HinduFestivalPlugin) getLunarFestivals(date time.Time, region api.Region) []api.Event {
	var events []api.Event

	// This is a simplified implementation. A production version would use
	// precise lunar calendar calculations.
	if h.isDiwali(date) {
		events = append(events, api.Event{
			Name:         "Diwali",
			NameLocal:    "दीपावली",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Festival of lights celebrating the victory of light over darkness",
			Region:       region,
			Metadata: map[string]interface{}{
				"importance":    "highest",
				"duration_days": 5,
				"lunar_month":   "Kartik",
				"lunar_day":     "Amavasya",
				"deities":       []string{"Lakshmi", "Ganesha"},
				"rituals":       []string{"lighting_diyas", "rangoli", "puja", "fireworks"},
				"regional_names": map[string]string{
					"tamil":     "தீபாவளி",
					"telugu":    "దీపావళి",
					"kannada":   "ದೀಪಾವಳಿ",
					"malayalam": "ദീപാവലി",
					"gujarati":  "દિવાળી",
					"marathi":   "दिवाळी",
					"bengali":   "কালীপূজা",
				},
			},
		})
	}

	if h.isHoli(date) {
		events = append(events, api.Event{
			Name:         "Holi",
			NameLocal:    "होली",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Festival of colors celebrating the arrival of spring",
			Region:       region,
			Metadata: map[string]interface{}{
				"importance":    "high",
				"duration_days": 2,
				"lunar_month":   "Phalgun",
				"lunar_day":     "Purnima",
				"deities":       []string{"Krishna", "Radha"},
				"rituals":       []string{"color_throwing", "bonfire", "dance", "sweets"},
				"regional_names": map[string]string{
					"gujarati": "હોળી",
					"marathi":  "होळी",
					"punjabi":  "ਹੋਲੀ",
					"tamil":    "ஹோலி",
					"bengali":  "দোল",
				},
			},
		})
	}

	if h.isNavaratri(date) {
		navaratriDay := h.getNavaratriDay(date)
		events = append(events, api.Event{
			Name:         "Navaratri",
			NameLocal:    "नवरात्रि",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Nine nights dedicated to the worship of Goddess Durga",
			Region:       region,
			Metadata: map[string]interface{}{
				"importance":    "high",
				"duration_days": 9,
				"current_day":   navaratriDay,
				"lunar_month":   "Ashwin",
				"deity":         "Durga",
				"daily_goddess": h.getNavaratriGoddess(navaratriDay),
				"rituals":       []string{"fasting", "dancing", "puja", "garba"},
				"regional_variants": map[string]string{
					"gujarat":    "Garba celebrations",
					"bengal":     "Durga Puja",
					"karnataka":  "Dasara",
					"tamil_nadu": "Navarathri Golu",
				},
			},
		})
	}

	if h.isJanmashtami(date) {
		events = append(events, api.Event{
			Name:         "Janmashtami",
			NameLocal:    "जन्माष्टमी",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Celebration of Lord Krishna's birth",
			Region:       region,
			Metadata: map[string]interface{}{
				"importance":    "high",
				"lunar_month":   "Bhadrapada",
				"lunar_day":     "Ashtami",
				"deity":         "Krishna",
				"rituals":       []string{"midnight_celebration", "fasting", "dahi_handi", "jhula"},
				"special_foods": []string{"makhan", "mishri", "panchamrit"},
				"regional_celebrations": map[string]string{
					"mathura":   "Krishna Janmabhoomi",
					"vrindavan": "Banke Bihari Temple",
					"mumbai":    "Dahi Handi",
					"gujarat":   "Dwarkadheesh Temple",
				},
			},
		})
	}

	return events
}

func (h *HinduFestivalPlugin) getEkadashiEvents(date time.Time, region api.Region) []api.Event {
	var events []api.Event

	if h.isEkadashi(date) {
		ekadashiName := h.getEkadashiName(date)
		events = append(events, api.Event{
			Name:         ekadashiName,
			NameLocal:    ekadashiName + " एकादशी",
			Type:         api.EventTypeEkadashi,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Fasting day dedicated to Lord Vishnu",
			Region:       region,
			Metadata: map[string]interface{}{
				"importance":         "medium",
				"lunar_day":          "Ekadashi",
				"deity":              "Vishnu",
				"fasting_type":       "nirjala_or_phalahar",
				"breaking_time":      "next_day_after_sunrise",
				"rituals":            []string{"fasting", "vishnu_puja", "tulsi_worship", "charity"},
				"spiritual_benefits": "purification, devotion, karma_cleansing",
			},
		})
	}

	return events
}

func (h *HinduFestivalPlugin) getMonthlyObservances(date time.Time, region api.Region) []api.Event {
	var events []api.Event

	if h.isAmavasya(date) {
		events = append(events, api.Event{
			Name:         "Amavasya",
			NameLocal:    "अमावस्या",
			Type:         api.EventTypeAmavasya,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "New moon day for ancestral worship and spiritual practices",
			Region:       region,
			Metadata: map[string]interface{}{
				"lunar_phase":            "new_moon",
				"spiritual_significance": "ancestral_worship",
				"rituals":                []string{"pitru_puja", "charity", "meditation"},
				"recommendations":        []string{"avoid_travel", "spiritual_practices", "charity"},
			},
		})
	}

	if h.isPurnima(date) {
		events = append(events, api.Event{
			Name:         "Purnima",
			NameLocal:    "पूर्णिमा",
			Type:         api.EventTypePurnima,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Full moon day for worship and spiritual activities",
			Region:       region,
			Metadata: map[string]interface{}{
				"lunar_phase":            "full_moon",
				"spiritual_significance": "heightened_spiritual_energy",
				"rituals":                []string{"fasting", "temple_visit", "charity", "meditation"},
				"benefits":               []string{"mental_clarity", "spiritual_growth", "positive_energy"},
			},
		})
	}

	return events
}
