package implementations

import (
	"time"

	"github.com/naren-m/panchangam/api"
)

// Helper methods to create specific festival events

func (a *AdvancedFestivalPlugin) createDiwaliEvent(date time.Time, region api.Region) api.Event {
	regionalNames := map[api.Region]string{
		api.RegionTamilNadu:  "தீபாவளி",
		api.RegionKerala:     "ദീപാവലി",
		api.RegionBengal:     "কালীপূজা/দীপাবলি",
		api.RegionGujarat:    "દિવાળી",
		api.RegionMaha:       "दिवाळी",
		api.RegionNorthIndia: "दीपावली",
		api.RegionSouthIndia: "दीपावली",
		api.RegionGlobal:     "दीपावली",
	}

	return api.Event{
		Name:         "Diwali",
		NameLocal:    regionalNames[region],
		Type:         api.EventTypeFestival,
		StartTime:    date,
		EndTime:      date.Add(24 * time.Hour),
		Significance: "Festival of lights celebrating the victory of light over darkness",
		Region:       region,
		Metadata: map[string]interface{}{
			"importance":         "highest",
			"duration_days":      5,
			"lunar_calculation":  true,
			"deities":            []string{"Lakshmi", "Ganesha"},
			"activities":         []string{"lighting_diyas", "rangoli", "puja", "fireworks", "sweets"},
			"astronomical_basis": "Kartik_Amavasya",
		},
	}
}

func (a *AdvancedFestivalPlugin) createHoliEvent(date time.Time, region api.Region) api.Event {
	regionalNames := map[api.Region]string{
		api.RegionTamilNadu:  "ஹோலி",
		api.RegionKerala:     "ഹോളി",
		api.RegionBengal:     "দোল/হোলি",
		api.RegionGujarat:    "હોળી",
		api.RegionMaha:       "होळी",
		api.RegionNorthIndia: "होली",
		api.RegionSouthIndia: "होली",
		api.RegionGlobal:     "होली",
	}

	return api.Event{
		Name:         "Holi",
		NameLocal:    regionalNames[region],
		Type:         api.EventTypeFestival,
		StartTime:    date,
		EndTime:      date.Add(24 * time.Hour),
		Significance: "Festival of colors celebrating the arrival of spring and victory of good over evil",
		Region:       region,
		Metadata: map[string]interface{}{
			"importance":         "high",
			"duration_days":      2,
			"lunar_calculation":  true,
			"deities":            []string{"Krishna", "Radha"},
			"activities":         []string{"color_throwing", "bonfire", "dance", "sweets"},
			"astronomical_basis": "Phalgun_Purnima",
		},
	}
}

func (a *AdvancedFestivalPlugin) createJanmashtamiEvent(date time.Time, region api.Region) api.Event {
	regionalNames := map[api.Region]string{
		api.RegionTamilNadu:  "ஜென்மாஷ்டமி",
		api.RegionKerala:     "ജന്മാഷ്ടമി",
		api.RegionBengal:     "জন্মাষ্টমী",
		api.RegionGujarat:    "જન્માષ્ટમી",
		api.RegionMaha:       "जन्माष्टमी",
		api.RegionNorthIndia: "जन्माष्टमी",
		api.RegionSouthIndia: "जन्माष्टमी",
		api.RegionGlobal:     "जन्माष्टमी",
	}

	return api.Event{
		Name:         "Janmashtami",
		NameLocal:    regionalNames[region],
		Type:         api.EventTypeFestival,
		StartTime:    date,
		EndTime:      date.Add(24 * time.Hour),
		Significance: "Celebration of Lord Krishna's birth on Bhadrapada Krishna Ashtami",
		Region:       region,
		Metadata: map[string]interface{}{
			"importance":           "high",
			"lunar_calculation":    true,
			"deity":                "Krishna",
			"midnight_celebration": true,
			"activities":           []string{"fasting", "midnight_puja", "dahi_handi", "jhula"},
			"astronomical_basis":   "Bhadrapada_Krishna_Ashtami",
		},
	}
}

func (a *AdvancedFestivalPlugin) createRamNavamiEvent(date time.Time, region api.Region) api.Event {
	return api.Event{
		Name:         "Ram Navami",
		NameLocal:    "राम नवमी",
		Type:         api.EventTypeFestival,
		StartTime:    date,
		EndTime:      date.Add(24 * time.Hour),
		Significance: "Celebration of Lord Rama's birth on Chaitra Shukla Navami",
		Region:       region,
		Metadata: map[string]interface{}{
			"importance":         "high",
			"lunar_calculation":  true,
			"deity":              "Rama",
			"astronomical_basis": "Chaitra_Shukla_Navami",
		},
	}
}

// Additional festival creation methods
func (a *AdvancedFestivalPlugin) createMahalayaAmavasya(date time.Time, region api.Region) api.Event {
	return api.Event{
		Name:         "Mahalaya Amavasya",
		NameLocal:    "महालया अमावस्या",
		Type:         api.EventTypeAmavasya,
		StartTime:    date,
		EndTime:      date.Add(24 * time.Hour),
		Significance: "Pitru Paksha ends, beginning of Devi Paksha",
		Region:       region,
		Metadata: map[string]interface{}{
			"importance":   "high",
			"pitru_paksha": "end",
			"devi_paksha":  "beginning",
		},
	}
}

func (a *AdvancedFestivalPlugin) createGuruPurnima(date time.Time, region api.Region) api.Event {
	return api.Event{
		Name:         "Guru Purnima",
		NameLocal:    "गुरु पूर्णिमा",
		Type:         api.EventTypePurnima,
		StartTime:    date,
		EndTime:      date.Add(24 * time.Hour),
		Significance: "Day dedicated to honoring spiritual teachers and gurus",
		Region:       region,
		Metadata: map[string]interface{}{
			"importance": "high",
			"deity":      "Vyasa",
			"purpose":    "guru_worship",
		},
	}
}

func (a *AdvancedFestivalPlugin) createKartikPurnima(date time.Time, region api.Region) api.Event {
	return api.Event{
		Name:         "Kartik Purnima",
		NameLocal:    "कार्तिक पूर्णिमा",
		Type:         api.EventTypePurnima,
		StartTime:    date,
		EndTime:      date.Add(24 * time.Hour),
		Significance: "Sacred full moon in Kartik month, festival of lights on water",
		Region:       region,
		Metadata: map[string]interface{}{
			"importance": "medium",
			"activities": []string{"ganga_aarti", "deep_daan", "holy_bath"},
		},
	}
}

func (a *AdvancedFestivalPlugin) createDurgaAshtami(date time.Time, region api.Region) api.Event {
	return api.Event{
		Name:         "Durga Ashtami",
		NameLocal:    "दुर्गा अष्टमी",
		Type:         api.EventTypeAshtami,
		StartTime:    date,
		EndTime:      date.Add(24 * time.Hour),
		Significance: "Eighth day of Navaratri dedicated to Goddess Durga",
		Region:       region,
		Metadata: map[string]interface{}{
			"navaratri_day": 8,
			"deity":         "Durga",
			"importance":    "high",
		},
	}
}

func (a *AdvancedFestivalPlugin) createMahaNavami(date time.Time, region api.Region) api.Event {
	return api.Event{
		Name:         "Maha Navami",
		NameLocal:    "महा नवमी",
		Type:         api.EventTypeNavami,
		StartTime:    date,
		EndTime:      date.Add(24 * time.Hour),
		Significance: "Ninth day of Navaratri, grand celebration of Divine Mother",
		Region:       region,
		Metadata: map[string]interface{}{
			"navaratri_day": 9,
			"deity":         "Durga",
			"importance":    "highest",
		},
	}
}
