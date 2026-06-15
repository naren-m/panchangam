package implementations

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy"
)

// getFestivalsByTithi returns festivals that occur on specific Tithis
func (a *AdvancedFestivalPlugin) getFestivalsByTithi(ctx context.Context, date time.Time, tithi *astronomy.TithiInfo, region api.Region) []api.Event {
	var events []api.Event
	month := date.Month()

	// Diwali - Amavasya in Kartik month (October/November)
	if tithi.Number == 30 && (month == time.October || month == time.November) {
		events = append(events, a.createDiwaliEvent(date, region))
	}

	// Holi - Purnima in Phalgun month (March)
	if tithi.Number == 15 && month == time.March {
		events = append(events, a.createHoliEvent(date, region))
	}

	// Janmashtami - Ashtami in Bhadrapada month (August/September)
	if tithi.Number == 23 && (month == time.August || month == time.September) {
		events = append(events, a.createJanmashtamiEvent(date, region))
	}

	// Ram Navami - Navami in Chaitra month (March/April)
	if tithi.Number == 9 && (month == time.March || month == time.April) {
		events = append(events, a.createRamNavamiEvent(date, region))
	}

	return events
}

// calculateEkadashi creates an Ekadashi event
func (a *AdvancedFestivalPlugin) calculateEkadashi(ctx context.Context, date time.Time, tithi *astronomy.TithiInfo, region api.Region) api.Event {
	// Determine Ekadashi name based on month and paksha
	ekadashiName := a.getEkadashiName(date, tithi.IsShukla)

	regionalNames := map[api.Region]string{
		api.RegionTamilNadu:  ekadashiName + " ஏகாதசி",
		api.RegionKerala:     ekadashiName + " ഏകാദശി",
		api.RegionBengal:     ekadashiName + " একাদশী",
		api.RegionGujarat:    ekadashiName + " એકાદશી",
		api.RegionMaha:       ekadashiName + " एकादशी",
		api.RegionNorthIndia: ekadashiName + " एकादशी",
		api.RegionSouthIndia: ekadashiName + " एकादशी",
		api.RegionGlobal:     ekadashiName + " एकादशी",
	}

	return api.Event{
		Name:         ekadashiName + " Ekadashi",
		NameLocal:    regionalNames[region],
		Type:         api.EventTypeEkadashi,
		StartTime:    tithi.StartTime,
		EndTime:      tithi.EndTime,
		Significance: "Fasting day dedicated to Lord Vishnu for spiritual purification",
		Region:       region,
		Metadata: map[string]interface{}{
			"deity":              "Vishnu",
			"fasting_type":       "nirjala_or_phalahar",
			"breaking_time":      "next_day_after_sunrise",
			"spiritual_benefits": []string{"purification", "devotion", "karma_cleansing"},
			"tithi_number":       tithi.Number,
			"paksha":             map[bool]string{true: "Shukla", false: "Krishna"}[tithi.IsShukla],
			"precise_timing":     true,
		},
	}
}

// calculateAmavasya creates an Amavasya event
func (a *AdvancedFestivalPlugin) calculateAmavasya(ctx context.Context, date time.Time, tithi *astronomy.TithiInfo, region api.Region) api.Event {
	regionalNames := map[api.Region]string{
		api.RegionTamilNadu:  "அமாவாசை",
		api.RegionKerala:     "അമാവാസ്യ",
		api.RegionBengal:     "অমাবস্যা",
		api.RegionGujarat:    "અમાવસ્યા",
		api.RegionMaha:       "अमावस्या",
		api.RegionNorthIndia: "अमावस्या",
		api.RegionSouthIndia: "अमावस्या",
		api.RegionGlobal:     "अमावस्या",
	}

	return api.Event{
		Name:         "Amavasya",
		NameLocal:    regionalNames[region],
		Type:         api.EventTypeAmavasya,
		StartTime:    tithi.StartTime,
		EndTime:      tithi.EndTime,
		Significance: "New moon day for ancestral worship and spiritual practices",
		Region:       region,
		Metadata: map[string]interface{}{
			"lunar_phase":            "new_moon",
			"spiritual_significance": "ancestral_worship",
			"recommended_activities": []string{"pitru_puja", "charity", "meditation", "fasting"},
			"avoid_activities":       []string{"new_ventures", "marriage", "travel"},
			"tithi_precision":        true,
		},
	}
}

// calculatePurnima creates a Purnima event
func (a *AdvancedFestivalPlugin) calculatePurnima(ctx context.Context, date time.Time, tithi *astronomy.TithiInfo, region api.Region) api.Event {
	regionalNames := map[api.Region]string{
		api.RegionTamilNadu:  "பௌர்ணமி",
		api.RegionKerala:     "പൂർണ്ണിമ",
		api.RegionBengal:     "পূর্ণিমা",
		api.RegionGujarat:    "પૂર્ણિમા",
		api.RegionMaha:       "पूर्णिमा",
		api.RegionNorthIndia: "पूर्णिमा",
		api.RegionSouthIndia: "पूर्णिमा",
		api.RegionGlobal:     "पूर्णिमा",
	}

	return api.Event{
		Name:         "Purnima",
		NameLocal:    regionalNames[region],
		Type:         api.EventTypePurnima,
		StartTime:    tithi.StartTime,
		EndTime:      tithi.EndTime,
		Significance: "Full moon day for worship and spiritual activities",
		Region:       region,
		Metadata: map[string]interface{}{
			"lunar_phase":            "full_moon",
			"spiritual_significance": "heightened_spiritual_energy",
			"recommended_activities": []string{"fasting", "temple_visit", "charity", "meditation"},
			"benefits":               []string{"mental_clarity", "spiritual_growth", "positive_energy"},
			"tithi_precision":        true,
		},
	}
}

// getSpecialAmavasyas returns special Amavasya festivals
func (a *AdvancedFestivalPlugin) getSpecialAmavasyas(ctx context.Context, date time.Time, tithi *astronomy.TithiInfo, region api.Region) []api.Event {
	var events []api.Event
	month := date.Month()

	// Diwali Amavasya
	if month == time.October || month == time.November {
		events = append(events, a.createDiwaliEvent(date, region))
	}

	// Mahalaya Amavasya (September/October)
	if month == time.September || month == time.October {
		events = append(events, a.createMahalayaAmavasya(date, region))
	}

	return events
}

// getSpecialPurnimas returns special Purnima festivals
func (a *AdvancedFestivalPlugin) getSpecialPurnimas(ctx context.Context, date time.Time, tithi *astronomy.TithiInfo, region api.Region) []api.Event {
	var events []api.Event
	month := date.Month()

	// Holi Purnima
	if month == time.March {
		events = append(events, a.createHoliEvent(date, region))
	}

	// Guru Purnima (July)
	if month == time.July {
		events = append(events, a.createGuruPurnima(date, region))
	}

	// Kartik Purnima (November)
	if month == time.November {
		events = append(events, a.createKartikPurnima(date, region))
	}

	return events
}

// getAshtamiFestivals returns festivals occurring on Ashtami
func (a *AdvancedFestivalPlugin) getAshtamiFestivals(ctx context.Context, date time.Time, tithi *astronomy.TithiInfo, region api.Region) []api.Event {
	var events []api.Event
	month := date.Month()

	// Janmashtami (Krishna Ashtami in Bhadrapada)
	if tithi.Number == 23 && (month == time.August || month == time.September) {
		events = append(events, a.createJanmashtamiEvent(date, region))
	}

	// Durga Ashtami (during Navaratri)
	if tithi.Number == 8 && (month == time.September || month == time.October) {
		events = append(events, a.createDurgaAshtami(date, region))
	}

	return events
}

// getNavamiFestivals returns festivals occurring on Navami
func (a *AdvancedFestivalPlugin) getNavamiFestivals(ctx context.Context, date time.Time, tithi *astronomy.TithiInfo, region api.Region) []api.Event {
	var events []api.Event
	month := date.Month()

	// Ram Navami (Chaitra Shukla Navami)
	if tithi.Number == 9 && (month == time.March || month == time.April) {
		events = append(events, a.createRamNavamiEvent(date, region))
	}

	// Maha Navami (during Navaratri)
	if tithi.Number == 9 && (month == time.September || month == time.October) {
		events = append(events, a.createMahaNavami(date, region))
	}

	return events
}

// calculateSankashtiChaturthi creates a Sankashti Chaturthi event
func (a *AdvancedFestivalPlugin) calculateSankashtiChaturthi(ctx context.Context, date time.Time, tithi *astronomy.TithiInfo, region api.Region) api.Event {
	regionalNames := map[api.Region]string{
		api.RegionTamilNadu:  "சங்கஷ்டி சதுர்த்தி",
		api.RegionKerala:     "സങ്കഷ്ടി ചതുർത്ഥി",
		api.RegionBengal:     "সংকষ্টি চতুর্থী",
		api.RegionGujarat:    "સંકષ્ટી ચતુર્થી",
		api.RegionMaha:       "संकष्टी चतुर्थी",
		api.RegionNorthIndia: "संकष्टी चतुर्थी",
		api.RegionSouthIndia: "संकष्टी चतुर्थी",
		api.RegionGlobal:     "संकष्टी चतुर्थी",
	}

	return api.Event{
		Name:         "Sankashti Chaturthi",
		NameLocal:    regionalNames[region],
		Type:         api.EventTypeSankashti,
		StartTime:    tithi.StartTime,
		EndTime:      tithi.EndTime,
		Significance: "Monthly fasting day dedicated to Lord Ganesha for removal of obstacles",
		Region:       region,
		Metadata: map[string]interface{}{
			"deity":              "Ganesha",
			"fasting_type":       "until_moonrise",
			"breaking_time":      "after_moonrise",
			"spiritual_benefits": []string{"obstacle_removal", "prosperity", "wisdom"},
			"monthly_occurrence": true,
			"tithi_number":       19,
			"paksha":             "Krishna",
		},
	}
}

// getEkadashiName returns the name of Ekadashi based on month and paksha
func (a *AdvancedFestivalPlugin) getEkadashiName(date time.Time, isShukla bool) string {
	month := date.Month()

	// Simplified mapping - in production this would be more complex
	ekadashiNames := map[time.Month]map[bool]string{
		time.January:   {true: "Saphala", false: "Putrada"},
		time.February:  {true: "Shattila", false: "Jaya"},
		time.March:     {true: "Vijaya", false: "Amalaki"},
		time.April:     {true: "Papamochani", false: "Kamada"},
		time.May:       {true: "Varuthini", false: "Mohini"},
		time.June:      {true: "Apara", false: "Nirjala"},
		time.July:      {true: "Yogini", false: "Devshayani"},
		time.August:    {true: "Kamika", false: "Shravana"},
		time.September: {true: "Aja", false: "Parsva"},
		time.October:   {true: "Indira", false: "Papankusha"},
		time.November:  {true: "Rama", false: "Haribodhini"},
		time.December:  {true: "Utpanna", false: "Mokshada"},
	}

	if names, exists := ekadashiNames[month]; exists {
		if name, exists := names[isShukla]; exists {
			return name
		}
	}

	return "Ekadashi"
}
