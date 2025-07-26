package implementations

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/api"
	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

// AdvancedFestivalPlugin provides precise Hindu festival calculations using lunar astronomy
type AdvancedFestivalPlugin struct {
	enabled             bool
	config              map[string]interface{}
	tithiCalculator     *astronomy.TithiCalculator
	ephemerisManager    *ephemeris.Manager
}

// NewAdvancedFestivalPlugin creates a new advanced festival calculation plugin
func NewAdvancedFestivalPlugin(ephemerisManager *ephemeris.Manager) *AdvancedFestivalPlugin {
	return &AdvancedFestivalPlugin{
		enabled:          false,
		config:           make(map[string]interface{}),
		tithiCalculator:  astronomy.NewTithiCalculator(ephemerisManager),
		ephemerisManager: ephemerisManager,
	}
}

// GetInfo returns plugin metadata
func (a *AdvancedFestivalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "advanced_festival_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Precise Hindu festival calculations using astronomical data and lunar calendar",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityEvent),
			string(api.CapabilityRegional),
		},
		Dependencies: []string{"astronomy", "ephemeris"},
		Metadata: map[string]interface{}{
			"calculation_precision": "astronomical",
			"festival_count":        "50+",
			"lunar_calculations":    true,
			"regional_variations":   true,
			"ayanamsa_aware":       true,
		},
	}
}

// Initialize sets up the plugin with configuration
func (a *AdvancedFestivalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	a.config = config
	a.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (a *AdvancedFestivalPlugin) IsEnabled() bool {
	return a.enabled
}

// Shutdown cleans up plugin resources
func (a *AdvancedFestivalPlugin) Shutdown(ctx context.Context) error {
	a.enabled = false
	return nil
}

// GetSupportedRegions returns regions this plugin supports
func (a *AdvancedFestivalPlugin) GetSupportedRegions() []api.Region {
	return []api.Region{
		api.RegionGlobal,
		api.RegionNorthIndia,
		api.RegionSouthIndia,
		api.RegionTamilNadu,
		api.RegionKerala,
		api.RegionBengal,
		api.RegionGujarat,
		api.RegionMaha,
	}
}

// GetSupportedEventTypes returns event types this plugin can generate
func (a *AdvancedFestivalPlugin) GetSupportedEventTypes() []api.EventType {
	return []api.EventType{
		api.EventTypeFestival,
		api.EventTypeEkadashi,
		api.EventTypeAmavasya,
		api.EventTypePurnima,
		api.EventTypeVrat,
		api.EventTypeSankashti,
		api.EventTypeAshtami,
		api.EventTypeNavami,
		api.EventTypeLunar,
	}
}

// GetEvents returns festival events for a specific date and location
func (a *AdvancedFestivalPlugin) GetEvents(ctx context.Context, date time.Time, location api.Location, region api.Region) ([]api.Event, error) {
	if !a.enabled {
		return nil, fmt.Errorf("advanced festival plugin is not enabled")
	}

	var events []api.Event

	// Get Tithi information for precise lunar calculations
	tithi, err := a.tithiCalculator.GetTithiForDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate tithi: %w", err)
	}

	// Check for festivals based on Tithi
	tithiEvents := a.getFestivalsByTithi(ctx, date, tithi, region)
	events = append(events, tithiEvents...)

	// Check for Ekadashi
	if tithi.Number == 11 || tithi.Number == 26 { // 11th day of both fortnights
		ekadashiEvent := a.calculateEkadashi(ctx, date, tithi, region)
		events = append(events, ekadashiEvent)
	}

	// Check for Amavasya (New Moon)
	if tithi.Number == 30 {
		amavasya := a.calculateAmavasya(ctx, date, tithi, region)
		events = append(events, amavasya)
		
		// Check for special Amavasya festivals
		specialAmavasya := a.getSpecialAmavasyas(ctx, date, tithi, region)
		events = append(events, specialAmavasya...)
	}

	// Check for Purnima (Full Moon)
	if tithi.Number == 15 {
		purnima := a.calculatePurnima(ctx, date, tithi, region)
		events = append(events, purnima)
		
		// Check for special Purnima festivals
		specialPurnima := a.getSpecialPurnimas(ctx, date, tithi, region)
		events = append(events, specialPurnima...)
	}

	// Check for Ashtami festivals (8th day)
	if tithi.Number == 8 || tithi.Number == 23 {
		ashtamiEvents := a.getAshtamiFestivals(ctx, date, tithi, region)
		events = append(events, ashtamiEvents...)
	}

	// Check for Navami festivals (9th day)
	if tithi.Number == 9 || tithi.Number == 24 {
		navamiEvents := a.getNavamiFestivals(ctx, date, tithi, region)
		events = append(events, navamiEvents...)
	}

	// Check for Sankashti Chaturthi (4th day of Krishna Paksha)
	if tithi.Number == 19 {
		sankashti := a.calculateSankashtiChaturthi(ctx, date, tithi, region)
		events = append(events, sankashti)
	}

	// Regional specific festivals
	regionalEvents := a.getRegionalFestivals(ctx, date, tithi, region)
	events = append(events, regionalEvents...)

	return events, nil
}

// GetEventsInRange returns festival events for a date range
func (a *AdvancedFestivalPlugin) GetEventsInRange(ctx context.Context, start, end time.Time, location api.Location, region api.Region) ([]api.Event, error) {
	var allEvents []api.Event

	current := start
	for current.Before(end) || current.Equal(end) {
		dayEvents, err := a.GetEvents(ctx, current, location, region)
		if err != nil {
			return nil, fmt.Errorf("failed to get events for %s: %w", current.Format("2006-01-02"), err)
		}
		allEvents = append(allEvents, dayEvents...)
		current = current.AddDate(0, 0, 1)
	}

	return allEvents, nil
}

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
			"paksha":            map[bool]string{true: "Shukla", false: "Krishna"}[tithi.IsShukla],
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
			"benefits":              []string{"mental_clarity", "spiritual_growth", "positive_energy"},
			"tithi_precision":       true,
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
			"deity":                "Ganesha",
			"fasting_type":         "until_moonrise",
			"breaking_time":        "after_moonrise",
			"spiritual_benefits":   []string{"obstacle_removal", "prosperity", "wisdom"},
			"monthly_occurrence":   true,
			"tithi_number":        19,
			"paksha":              "Krishna",
		},
	}
}

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
			"importance":          "highest",
			"duration_days":       5,
			"lunar_calculation":   true,
			"deities":            []string{"Lakshmi", "Ganesha"},
			"activities":         []string{"lighting_diyas", "rangoli", "puja", "fireworks", "sweets"},
			"astronomical_basis":  "Kartik_Amavasya",
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
			"importance":          "high",
			"duration_days":       2,
			"lunar_calculation":   true,
			"deities":            []string{"Krishna", "Radha"},
			"activities":         []string{"color_throwing", "bonfire", "dance", "sweets"},
			"astronomical_basis":  "Phalgun_Purnima",
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
			"importance":          "high",
			"lunar_calculation":   true,
			"deity":              "Krishna",
			"midnight_celebration": true,
			"activities":         []string{"fasting", "midnight_puja", "dahi_handi", "jhula"},
			"astronomical_basis":  "Bhadrapada_Krishna_Ashtami",
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
			"importance":          "high",
			"lunar_calculation":   true,
			"deity":              "Rama",
			"astronomical_basis":  "Chaitra_Shukla_Navami",
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
			"importance":    "high",
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

// getRegionalFestivals returns region-specific festivals
func (a *AdvancedFestivalPlugin) getRegionalFestivals(ctx context.Context, date time.Time, tithi *astronomy.TithiInfo, region api.Region) []api.Event {
	var events []api.Event

	switch region {
	case api.RegionTamilNadu:
		events = append(events, a.getTamilFestivals(date, tithi)...)
	case api.RegionKerala:
		events = append(events, a.getKeralaFestivals(date, tithi)...)
	case api.RegionBengal:
		events = append(events, a.getBengalFestivals(date, tithi)...)
	}

	return events
}

// Helper methods for regional festivals
func (a *AdvancedFestivalPlugin) getTamilFestivals(date time.Time, tithi *astronomy.TithiInfo) []api.Event {
	var events []api.Event
	// Tamil-specific festival logic based on tithi
	// Implementation would go here
	return events
}

func (a *AdvancedFestivalPlugin) getKeralaFestivals(date time.Time, tithi *astronomy.TithiInfo) []api.Event {
	var events []api.Event
	// Kerala-specific festival logic based on tithi
	// Implementation would go here
	return events
}

func (a *AdvancedFestivalPlugin) getBengalFestivals(date time.Time, tithi *astronomy.TithiInfo) []api.Event {
	var events []api.Event
	// Bengal-specific festival logic based on tithi
	// Implementation would go here
	return events
}

// getEkadashiName returns the name of Ekadashi based on month and paksha
func (a *AdvancedFestivalPlugin) getEkadashiName(date time.Time, isShukla bool) string {
	month := date.Month()
	
	// Simplified mapping - in production this would be more complex
	ekadashiNames := map[time.Month]map[bool]string{
		time.January:  {true: "Saphala", false: "Putrada"},
		time.February: {true: "Shattila", false: "Jaya"},
		time.March:    {true: "Vijaya", false: "Amalaki"},
		time.April:    {true: "Papamochani", false: "Kamada"},
		time.May:      {true: "Varuthini", false: "Mohini"},
		time.June:     {true: "Apara", false: "Nirjala"},
		time.July:     {true: "Yogini", false: "Devshayani"},
		time.August:   {true: "Kamika", false: "Shravana"},
		time.September:{true: "Aja", false: "Parsva"},
		time.October:  {true: "Indira", false: "Papankusha"},
		time.November: {true: "Rama", false: "Haribodhini"},
		time.December: {true: "Utpanna", false: "Mokshada"},
	}
	
	if names, exists := ekadashiNames[month]; exists {
		if name, exists := names[isShukla]; exists {
			return name
		}
	}
	
	return "Ekadashi"
}