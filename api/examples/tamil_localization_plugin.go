package examples

import (
	"context"
	"strings"

	"github.com/naren-m/panchangam/api"
)

// TamilLocalizationPlugin provides Tamil language localization for Panchangam elements
type TamilLocalizationPlugin struct {
	enabled bool
	config  map[string]interface{}

	// Translation dictionaries
	tithiNames     map[string]string
	nakshatraNames map[string]string
	yogaNames      map[string]string
	karanaNames    map[string]string
	varaNames      map[string]string
	eventNames     map[string]string
	muhurtaNames   map[string]string
}

// NewTamilLocalizationPlugin creates a new Tamil localization plugin
func NewTamilLocalizationPlugin() *TamilLocalizationPlugin {
	plugin := &TamilLocalizationPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}

	plugin.initializeTranslations()
	return plugin
}

// GetInfo returns plugin metadata
func (t *TamilLocalizationPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "tamil_localization",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Tamil language localization for Panchangam elements",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityLocalization),
		},
		Dependencies: []string{},
		Metadata: map[string]interface{}{
			"language":       "tamil",
			"script":         "tamil",
			"locale_codes":   []string{"ta", "ta-IN", "tamil"},
			"encoding":       "UTF-8",
			"region_support": []string{"tamil_nadu", "south_india", "global"},
		},
	}
}

// Initialize sets up the plugin with configuration
func (t *TamilLocalizationPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	t.config = config
	t.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (t *TamilLocalizationPlugin) IsEnabled() bool {
	return t.enabled
}

// Shutdown cleans up plugin resources
func (t *TamilLocalizationPlugin) Shutdown(ctx context.Context) error {
	t.enabled = false
	return nil
}

// GetSupportedLocales returns supported locale codes
func (t *TamilLocalizationPlugin) GetSupportedLocales() []string {
	return []string{"ta", "ta-IN", "tamil"}
}

// GetSupportedRegions returns regions this plugin supports
func (t *TamilLocalizationPlugin) GetSupportedRegions() []api.Region {
	return []api.Region{
		api.RegionTamilNadu,
		api.RegionSouthIndia,
		api.RegionGlobal,
	}
}

// LocalizeTithi returns localized tithi information
func (t *TamilLocalizationPlugin) LocalizeTithi(tithi *api.Tithi, locale string, region api.Region) error {
	if !t.isSupported(locale) {
		return nil // Not our responsibility
	}

	if localName, exists := t.tithiNames[tithi.Name]; exists {
		tithi.NameLocal = localName
	}

	// Localize tithi lord names
	if tithi.Lord != "" {
		tithi.Lord = t.localizeDivineNames(tithi.Lord)
	}

	// Localize quality descriptions
	if tithi.Quality != "" {
		tithi.Quality = t.localizeQuality(tithi.Quality)
	}

	return nil
}

// LocalizeNakshatra returns localized nakshatra information
func (t *TamilLocalizationPlugin) LocalizeNakshatra(nakshatra *api.Nakshatra, locale string, region api.Region) error {
	if !t.isSupported(locale) {
		return nil
	}

	if localName, exists := t.nakshatraNames[nakshatra.Name]; exists {
		nakshatra.NameLocal = localName
	}

	// Localize nakshatra attributes
	if nakshatra.Lord != "" {
		nakshatra.Lord = t.localizeDivineNames(nakshatra.Lord)
	}

	if nakshatra.Deity != "" {
		nakshatra.Deity = t.localizeDivineNames(nakshatra.Deity)
	}

	if nakshatra.Symbol != "" {
		nakshatra.Symbol = t.localizeSymbols(nakshatra.Symbol)
	}

	if nakshatra.Quality != "" {
		nakshatra.Quality = t.localizeQuality(nakshatra.Quality)
	}

	return nil
}

// LocalizeYoga returns localized yoga information
func (t *TamilLocalizationPlugin) LocalizeYoga(yoga *api.Yoga, locale string, region api.Region) error {
	if !t.isSupported(locale) {
		return nil
	}

	if localName, exists := t.yogaNames[yoga.Name]; exists {
		yoga.NameLocal = localName
	}

	if yoga.Quality != "" {
		yoga.Quality = t.localizeQuality(yoga.Quality)
	}

	return nil
}

// LocalizeKarana returns localized karana information
func (t *TamilLocalizationPlugin) LocalizeKarana(karana *api.Karana, locale string, region api.Region) error {
	if !t.isSupported(locale) {
		return nil
	}

	if localName, exists := t.karanaNames[karana.Name]; exists {
		karana.NameLocal = localName
	}

	// Localize karana type
	switch karana.Type {
	case "movable":
		karana.Type = "நகரக்கூடிய"
	case "fixed":
		karana.Type = "நிலையான"
	}

	if karana.Quality != "" {
		karana.Quality = t.localizeQuality(karana.Quality)
	}

	return nil
}

// LocalizeEvent returns localized event information
func (t *TamilLocalizationPlugin) LocalizeEvent(event *api.Event, locale string, region api.Region) error {
	if !t.isSupported(locale) {
		return nil
	}

	if localName, exists := t.eventNames[event.Name]; exists {
		event.NameLocal = localName
	}

	// Localize event significance
	if event.Significance != "" {
		event.Significance = t.localizeSignificance(event.Significance)
	}

	return nil
}

// LocalizeMuhurta returns localized muhurta information
func (t *TamilLocalizationPlugin) LocalizeMuhurta(muhurta *api.Muhurta, locale string, region api.Region) error {
	if !t.isSupported(locale) {
		return nil
	}

	if localName, exists := t.muhurtaNames[muhurta.Name]; exists {
		muhurta.NameLocal = localName
	}

	// Localize quality
	switch muhurta.Quality {
	case api.QualityAuspicious:
		muhurta.Quality = "शुभ"
	case api.QualityInauspicious:
		muhurta.Quality = "अशुभ"
	case api.QualityHighly:
		muhurta.Quality = "अत्यंत शुभ"
	case api.QualityAvoid:
		muhurta.Quality = "टाळावे"
	}

	// Localize purpose and significance
	if muhurta.Significance != "" {
		muhurta.Significance = t.localizeSignificance(muhurta.Significance)
	}

	return nil
}

// Helper methods

func (t *TamilLocalizationPlugin) isSupported(locale string) bool {
	locale = strings.ToLower(locale)
	supportedLocales := []string{"ta", "ta-in", "tamil"}

	for _, supported := range supportedLocales {
		if locale == supported {
			return true
		}
	}
	return false
}

func (t *TamilLocalizationPlugin) localizeDivineNames(name string) string {
	divineNames := map[string]string{
		"Shiva":     "சிவன்",
		"Vishnu":    "விஷ்ணு",
		"Brahma":    "பிரம்மா",
		"Indra":     "இந்திரன்",
		"Agni":      "அக்னி",
		"Vayu":      "வாயு",
		"Surya":     "சூர்யன்",
		"Chandra":   "சந்திரன்",
		"Mangal":    "செவ்வாய்",
		"Budha":     "புதன்",
		"Guru":      "குரு",
		"Shukra":    "சுக்ரன்",
		"Shani":     "சனி",
		"Rahu":      "ராகு",
		"Ketu":      "கேது",
		"Ganesha":   "கணேசன்",
		"Murugan":   "முருகன்",
		"Devi":      "தேவி",
		"Lakshmi":   "லக்ஷ்மி",
		"Saraswati": "சரஸ்வதி",
	}

	if tamilName, exists := divineNames[name]; exists {
		return tamilName
	}
	return name
}

func (t *TamilLocalizationPlugin) localizeSymbols(symbol string) string {
	symbols := map[string]string{
		"Horse's head":    "குதிரை தலை",
		"Elephant":        "யானை",
		"Knife":           "கத்தி",
		"Cart":            "வண்டி",
		"Deer's head":     "மான் தலை",
		"Jewel":           "ரத்னம்",
		"Flower":          "பூ",
		"Serpent":         "பாம்பு",
		"Drum":            "மிருதங்கம்",
		"Water pot":       "கமண்டலு",
		"Bed":             "படுக்கை",
		"Crown":           "கிரீடம்",
		"Hand":            "கை",
		"Pearl":           "முத்து",
		"Flute":           "புல்லாங்குழல்",
		"Fan":             "விசிறி",
		"Pot":             "பானை",
		"Tusk":            "தந்தம்",
		"Earring":         "காதணி",
		"Fish":            "மீன்",
		"Bamboo":          "மூங்கில்",
		"Two front teeth": "இரு முன்பற்கள்",
		"Tail of lion":    "சிங்க வால்",
		"Sword":           "வாள்",
		"Couch":           "மஞ்சம்",
		"Thunderbolt":     "வச்சிராயுதம்",
		"Arch":            "வில்",
	}

	if tamilSymbol, exists := symbols[symbol]; exists {
		return tamilSymbol
	}
	return symbol
}

func (t *TamilLocalizationPlugin) localizeQuality(quality string) string {
	qualities := map[string]string{
		"auspicious":        "நன்மை",
		"inauspicious":      "தீமை",
		"neutral":           "நடுநிலை",
		"highly_auspicious": "மிக நன்மை",
		"mildly_auspicious": "சிறிது நன்மை",
		"avoid":             "தவிர்க்க",
		"good":              "நல்ல",
		"bad":               "கெட்ட",
		"excellent":         "சிறந்த",
		"poor":              "மோசமான",
		"moderate":          "மிதமான",
	}

	if tamilQuality, exists := qualities[quality]; exists {
		return tamilQuality
	}
	return quality
}

func (t *TamilLocalizationPlugin) localizeSignificance(significance string) string {
	// This would typically involve more sophisticated translation
	// For now, we provide common patterns

	commonPhrases := map[string]string{
		"auspicious time":     "நன்மையான நேரம்",
		"divine time":         "தெய்வீக நேரம்",
		"spiritual practices": "ஆன்மீக நடைமுறைகள்",
		"new ventures":        "புதிய முயற்சிகள்",
		"important decisions": "முக்கியமான முடிவுகள்",
		"celebration":         "கொண்டாட்டம்",
		"festival":            "திருவிழா",
		"worship":             "வழிபாடு",
		"meditation":          "தியானம்",
		"prayers":             "பிரார்த்தனைகள்",
		"fasting":             "உபவாசம்",
		"avoid travel":        "பயணத்தை தவிர்க்க",
		"avoid new work":      "புதிய வேலையை தவிர்க்க",
		"good for marriage":   "திருமணத்திற்கு நல்லது",
		"good for business":   "வணிகத்திற்கு நல்லது",
	}

	// Simple phrase replacement
	for english, tamil := range commonPhrases {
		if strings.Contains(strings.ToLower(significance), english) {
			significance = strings.ReplaceAll(significance, english, tamil)
		}
	}

	return significance
}

// initializeTranslations sets up all translation dictionaries
func (t *TamilLocalizationPlugin) initializeTranslations() {
	t.initializeTithiNames()
	t.initializeNakshatraNames()
	t.initializeYogaNames()
	t.initializeKaranaNames()
	t.initializeVaraNames()
	t.initializeEventNames()
	t.initializeMuhurtaNames()
}

func (t *TamilLocalizationPlugin) initializeTithiNames() {
	t.tithiNames = map[string]string{
		"Pratipada":   "பிரதமை",
		"Dwitiya":     "துவிதியை",
		"Tritiya":     "திருதியை",
		"Chaturthi":   "சதுர்த்தி",
		"Panchami":    "பஞ்சமி",
		"Shashthi":    "சஷ்டி",
		"Saptami":     "சப்தமி",
		"Ashtami":     "அஷ்டமி",
		"Navami":      "நவமி",
		"Dashami":     "தசமி",
		"Ekadashi":    "ஏகாதசி",
		"Dwadashi":    "துவாதசி",
		"Trayodashi":  "திரயோதசி",
		"Chaturdashi": "சதுர்தசி",
		"Purnima":     "பூர்ணிமா",
		"Amavasya":    "அமாவாஸ்யா",
	}
}

func (t *TamilLocalizationPlugin) initializeNakshatraNames() {
	t.nakshatraNames = map[string]string{
		"Ashwini":           "அசுவினி",
		"Bharani":           "பரணி",
		"Krittika":          "கார்த்திகை",
		"Rohini":            "ரோகிணி",
		"Mrigasira":         "மிருகசீர்ஷம்",
		"Ardra":             "ஆத்ரா",
		"Punarvasu":         "புனர்வசு",
		"Pushya":            "பூசம்",
		"Ashlesha":          "ஆயில்யம்",
		"Magha":             "மகம்",
		"Purva Phalguni":    "பூரம்",
		"Uttara Phalguni":   "உத்தரம்",
		"Hasta":             "ஹஸ்தம்",
		"Chitra":            "சித்திரை",
		"Swati":             "சுவாதி",
		"Vishakha":          "விசாகம்",
		"Anuradha":          "அனுஷம்",
		"Jyeshtha":          "கேட்டை",
		"Mula":              "மூலம்",
		"Purva Ashadha":     "பூராடம்",
		"Uttara Ashadha":    "உத்திராடம்",
		"Shravana":          "திருவோணம்",
		"Dhanishta":         "அவிட்டம்",
		"Shatabhisha":       "சதயம்",
		"Purva Bhadrapada":  "பூரட்டாதி",
		"Uttara Bhadrapada": "உத்திரட்டாதி",
		"Revati":            "ரேவதி",
	}
}

func (t *TamilLocalizationPlugin) initializeYogaNames() {
	t.yogaNames = map[string]string{
		"Vishkambha": "விஷ்கம்பா",
		"Preeti":     "ப்ரீதி",
		"Ayushman":   "ஆயுஷ்மான்",
		"Saubhagya":  "சௌபாக்ய",
		"Shobhana":   "சோபன",
		"Atiganda":   "அதிகண்ட",
		"Sukarma":    "சுகர்மா",
		"Dhriti":     "திருதி",
		"Shula":      "சூல",
		"Ganda":      "கண்ட",
		"Vriddhi":    "வ்ருத்தி",
		"Dhruva":     "துருவ",
		"Vyaghata":   "வ்யாகாத",
		"Harshana":   "ஹர்ஷண",
		"Vajra":      "வஜ்ரா",
		"Siddhi":     "சித்தி",
		"Vyatipata":  "வ்யதிபாத",
		"Variyana":   "வரியான்",
		"Parigha":    "பரிக",
		"Shiva":      "சிவ",
		"Siddha":     "சித்த",
		"Sadhya":     "சாத்ய",
		"Subha":      "சுப",
		"Sukla":      "சுக்ல",
		"Brahma":     "பிரம்மா",
		"Mahendra":   "மகேந்திர",
		"Vaidhriti":  "வைதிருதி",
	}
}

func (t *TamilLocalizationPlugin) initializeKaranaNames() {
	t.karanaNames = map[string]string{
		"Bava":        "பவ",
		"Balava":      "பாலவ",
		"Kaulava":     "கௌலவ",
		"Taitila":     "தைதில",
		"Garija":      "காரிஜ",
		"Vanija":      "வணிஜ",
		"Vishti":      "விஷ்டி",
		"Shakuni":     "சகுனி",
		"Chatushpada": "சதுஷ்பாத",
		"Naga":        "நாக",
		"Kimstughna":  "கிம்ஸ்துக்ன",
	}
}

func (t *TamilLocalizationPlugin) initializeVaraNames() {
	t.varaNames = map[string]string{
		"Sunday":    "ஞாயிறு",
		"Monday":    "திங்கள்",
		"Tuesday":   "செவ்வாய்",
		"Wednesday": "புதன்",
		"Thursday":  "வியாழன்",
		"Friday":    "வெள்ளி",
		"Saturday":  "சனி",
	}
}

func (t *TamilLocalizationPlugin) initializeEventNames() {
	t.eventNames = map[string]string{
		"Diwali":           "தீபாவளி",
		"Holi":             "ஹோலி",
		"Dussehra":         "விஜயதசமி",
		"Navaratri":        "நவராத்திரி",
		"Karva Chauth":     "கர்வா சௌத்",
		"Raksha Bandhan":   "ராக்ஷா பந்தன்",
		"Janmashtami":      "ஜன்மாஷ்டமி",
		"Maha Shivaratri":  "மகா சிவராத்திரி",
		"Ram Navami":       "ராம நவமி",
		"Hanuman Jayanti":  "அனுமன் ஜயந்தி",
		"Ganesh Chaturthi": "விநாயகர் சதுர்த்தி",
		"Pongal":           "பொங்கல்",
		"Onam":             "ஓணம்",
		"Vishu":            "விஷு",
		"Baisakhi":         "வைசாகி",
		"Ekadashi":         "ஏகாதசி",
		"Purnima":          "பூர்ணிமா",
		"Amavasya":         "அமாவாஸ்யா",
	}
}

func (t *TamilLocalizationPlugin) initializeMuhurtaNames() {
	t.muhurtaNames = map[string]string{
		"Abhijit Muhurta": "அபிஜித் முகூர்த்தம்",
		"Brahma Muhurta":  "பிரம்ம முகூர்த்தம்",
		"Godhuli Muhurta": "கோதூளி முகூர்த்தம்",
		"Rahu Kalam":      "ராகு காலம்",
		"Yamagandam":      "யமகண்டம்",
		"Gulika Kalam":    "குளிக காலம்",
		"Dur Muhurta":     "துர் முகூர்த்தம்",
		"Amrit Kalam":     "அம்ருத காலம்",
		"Shubh Muhurta":   "சுப முகூர்த்தம்",
		"Vivah Muhurta":   "விவாஹ முகூர்த்தம்",
		"Griha Pravesh":   "க்ருஹ ப்ரவேசம்",
		"Upanayana":       "உபநயனம்",
		"Mundan":          "முண்டன்",
		"Anna Prashan":    "அன்ன பிராசன்",
		"Namkaran":        "நாமகரண்",
	}
}
