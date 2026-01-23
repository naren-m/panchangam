package implementations

import (
	"context"
	"fmt"

	"github.com/naren-m/panchangam/api"
)

// MultiLanguageLocalizationPlugin provides localization for multiple Indian languages
type MultiLanguageLocalizationPlugin struct {
	enabled bool
	config  map[string]interface{}
}

// NewMultiLanguageLocalizationPlugin creates a new localization plugin
func NewMultiLanguageLocalizationPlugin() *MultiLanguageLocalizationPlugin {
	return &MultiLanguageLocalizationPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

// GetInfo returns plugin metadata
func (m *MultiLanguageLocalizationPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "multi_language_localization_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Localization plugin supporting Tamil, Malayalam, Bengali, Gujarati, and Marathi",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityLocalization),
		},
		Metadata: map[string]interface{}{
			"supported_languages": []string{"ta", "ml", "bn", "gu", "mr", "hi", "en"},
			"locales":             []string{"ta_IN", "ml_IN", "bn_IN", "gu_IN", "mr_IN", "hi_IN", "en_IN"},
		},
	}
}

// Initialize sets up the plugin
func (m *MultiLanguageLocalizationPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	m.config = config
	m.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is enabled
func (m *MultiLanguageLocalizationPlugin) IsEnabled() bool {
	return m.enabled
}

// Shutdown cleans up resources
func (m *MultiLanguageLocalizationPlugin) Shutdown(ctx context.Context) error {
	m.enabled = false
	return nil
}

// GetSupportedLocales returns supported locale codes
func (m *MultiLanguageLocalizationPlugin) GetSupportedLocales() []string {
	return []string{"ta", "ta_IN", "ml", "ml_IN", "bn", "bn_IN", "gu", "gu_IN", "mr", "mr_IN", "hi", "hi_IN", "en", "en_IN"}
}

// GetSupportedRegions returns supported regions
func (m *MultiLanguageLocalizationPlugin) GetSupportedRegions() []api.Region {
	return []api.Region{
		api.RegionTamilNadu,
		api.RegionKerala,
		api.RegionBengal,
		api.RegionGujarat,
		api.RegionMaha,
		api.RegionNorthIndia,
		api.RegionSouthIndia,
		api.RegionGlobal,
	}
}

// LocalizeTithi localizes tithi information
func (m *MultiLanguageLocalizationPlugin) LocalizeTithi(tithi *api.Tithi, locale string, region api.Region) error {
	if !m.enabled {
		return fmt.Errorf("localization plugin is not enabled")
	}

	tithiNames := m.getTithiTranslations()
	language := m.getLanguageFromLocale(locale)

	if translations, exists := tithiNames[language]; exists {
		if localName, exists := translations[tithi.Name]; exists {
			tithi.NameLocal = localName
		}
	}

	return nil
}

// LocalizeNakshatra localizes nakshatra information
func (m *MultiLanguageLocalizationPlugin) LocalizeNakshatra(nakshatra *api.Nakshatra, locale string, region api.Region) error {
	if !m.enabled {
		return fmt.Errorf("localization plugin is not enabled")
	}

	nakshatraNames := m.getNakshatraTranslations()
	language := m.getLanguageFromLocale(locale)

	if translations, exists := nakshatraNames[language]; exists {
		if localName, exists := translations[nakshatra.Name]; exists {
			nakshatra.NameLocal = localName
		}
	}

	return nil
}

// LocalizeYoga localizes yoga information
func (m *MultiLanguageLocalizationPlugin) LocalizeYoga(yoga *api.Yoga, locale string, region api.Region) error {
	if !m.enabled {
		return fmt.Errorf("localization plugin is not enabled")
	}

	yogaNames := m.getYogaTranslations()
	language := m.getLanguageFromLocale(locale)

	if translations, exists := yogaNames[language]; exists {
		if localName, exists := translations[yoga.Name]; exists {
			yoga.NameLocal = localName
		}
	}

	return nil
}

// LocalizeKarana localizes karana information
func (m *MultiLanguageLocalizationPlugin) LocalizeKarana(karana *api.Karana, locale string, region api.Region) error {
	if !m.enabled {
		return fmt.Errorf("localization plugin is not enabled")
	}

	karanaNames := m.getKaranaTranslations()
	language := m.getLanguageFromLocale(locale)

	if translations, exists := karanaNames[language]; exists {
		if localName, exists := translations[karana.Name]; exists {
			karana.NameLocal = localName
		}
	}

	return nil
}

// LocalizeEvent localizes event information
func (m *MultiLanguageLocalizationPlugin) LocalizeEvent(event *api.Event, locale string, region api.Region) error {
	if !m.enabled {
		return fmt.Errorf("localization plugin is not enabled")
	}

	eventNames := m.getEventTranslations()
	language := m.getLanguageFromLocale(locale)

	if translations, exists := eventNames[language]; exists {
		if localName, exists := translations[event.Name]; exists {
			event.NameLocal = localName
		}
	}

	return nil
}

// LocalizeMuhurta localizes muhurta information
func (m *MultiLanguageLocalizationPlugin) LocalizeMuhurta(muhurta *api.Muhurta, locale string, region api.Region) error {
	if !m.enabled {
		return fmt.Errorf("localization plugin is not enabled")
	}

	muhurtaNames := m.getMuhurtaTranslations()
	language := m.getLanguageFromLocale(locale)

	if translations, exists := muhurtaNames[language]; exists {
		if localName, exists := translations[muhurta.Name]; exists {
			muhurta.NameLocal = localName
		}
	}

	return nil
}

// Helper methods for translations

func (m *MultiLanguageLocalizationPlugin) getLanguageFromLocale(locale string) string {
	// Extract language code from locale (e.g., "ta_IN" -> "ta")
	if len(locale) >= 2 {
		return locale[:2]
	}
	return "en"
}

func (m *MultiLanguageLocalizationPlugin) getTithiTranslations() map[string]map[string]string {
	return map[string]map[string]string{
		"ta": { // Tamil
			"Pratipada":  "பிரதமை",
			"Dwitiya":    "துவிதியை",
			"Tritiya":    "திருதியை",
			"Chaturthi":  "சதுர்த்தி",
			"Panchami":   "பஞ்சமி",
			"Shashthi":   "சஷ்டி",
			"Saptami":    "சப்தமி",
			"Ashtami":    "அஷ்டமி",
			"Navami":     "நவமி",
			"Dashami":    "தசமி",
			"Ekadashi":   "ஏகாதசி",
			"Dwadashi":   "துவாதசி",
			"Trayodashi": "த்ரயோதசி",
			"Chaturdashi": "சதுர்த்தசி",
			"Purnima":    "பௌர்ணமி",
			"Amavasya":   "அமாவாசை",
		},
		"ml": { // Malayalam
			"Pratipada":  "പ്രതിപദം",
			"Purnima":    "പൂർണ്ണിമ",
			"Amavasya":   "അമാവാസ്യ",
			"Ekadashi":   "ഏകാദശി",
		},
		"bn": { // Bengali
			"Pratipada":  "প্রতিপদ",
			"Purnima":    "পূর্ণিমা",
			"Amavasya":   "অমাবস্যা",
			"Ekadashi":   "একাদশী",
		},
		"gu": { // Gujarati
			"Pratipada":  "પ્રતિપદા",
			"Purnima":    "પૂર્ણિમા",
			"Amavasya":   "અમાવસ્યા",
			"Ekadashi":   "એકાદશી",
		},
		"mr": { // Marathi
			"Pratipada":  "प्रतिपदा",
			"Purnima":    "पूर्णिमा",
			"Amavasya":   "अमावस्या",
			"Ekadashi":   "एकादशी",
		},
	}
}

func (m *MultiLanguageLocalizationPlugin) getNakshatraTranslations() map[string]map[string]string {
	return map[string]map[string]string{
		"ta": { // Tamil
			"Ashwini":   "அஸ்வினி",
			"Bharani":   "பரணி",
			"Krittika":  "கார்த்திகை",
			"Rohini":    "ரோகிணி",
			"Mrigashira": "மிருகசீரிஷம்",
			"Ardra":     "திருவாதிரை",
			"Punarvasu": "புனர்பூசம்",
			"Pushya":    "பூசம்",
			"Ashlesha":  "ஆயில்யம்",
			"Magha":     "மகம்",
			"Purva Phalguni": "பூரம்",
			"Uttara Phalguni": "உத்திரம்",
			"Hasta":     "அஸ்தம்",
			"Chitra":    "சித்திரை",
			"Swati":     "சுவாதி",
			"Vishakha":  "விசாகம்",
			"Anuradha":  "அனுஷம்",
			"Jyeshtha":  "கேட்டை",
			"Mula":      "மூலம்",
			"Purva Ashadha": "பூராடம்",
			"Uttara Ashadha": "உத்திராடம்",
			"Shravana":  "திருவோணம்",
			"Dhanishta": "அவிட்டம்",
			"Shatabhisha": "சதயம்",
			"Purva Bhadrapada": "பூரட்டாதி",
			"Uttara Bhadrapada": "உத்திரட்டாதி",
			"Revati":    "ரேவதி",
		},
		"ml": { // Malayalam
			"Ashwini":   "അശ്വതി",
			"Bharani":   "ഭരണി",
			"Krittika":  "കാർത്തിക",
			"Rohini":    "രോഹിണി",
		},
		"bn": { // Bengali
			"Ashwini":   "অশ্বিনী",
			"Bharani":   "ভরণী",
			"Krittika":  "কৃত্তিকা",
		},
	}
}

func (m *MultiLanguageLocalizationPlugin) getYogaTranslations() map[string]map[string]string {
	return map[string]map[string]string{
		"ta": {
			"Vishkambha": "விஷ்கம்பம்",
			"Priti":      "பிரீதி",
			"Ayushman":   "ஆயுஷ்மான்",
			"Saubhagya":  "சௌபாக்யம்",
		},
		"ml": {
			"Vishkambha": "വിഷ്കംഭം",
		},
	}
}

func (m *MultiLanguageLocalizationPlugin) getKaranaTranslations() map[string]map[string]string {
	return map[string]map[string]string{
		"ta": {
			"Bava":    "பவ",
			"Balava":  "பாலவ",
			"Kaulava": "கௌலவ",
			"Taitila": "தைதில",
			"Garaja":  "கராஜ",
			"Vanija":  "வணிஜ",
			"Vishti":  "விஷ்டி",
		},
	}
}

func (m *MultiLanguageLocalizationPlugin) getEventTranslations() map[string]map[string]string {
	return map[string]map[string]string{
		"ta": {
			"Diwali":       "தீபாவளி",
			"Holi":         "ஹோலி",
			"Janmashtami":  "ஜென்மாஷ்டமி",
			"Ram Navami":   "ராம நவமி",
			"Rahu Kalam":   "ராகு காலம்",
			"Yamagandam":   "யமகண்டம்",
			"Gulika Kalam": "குளிக காலம்",
		},
		"ml": {
			"Diwali":       "ദീപാവലി",
			"Vishu":        "വിഷു",
			"Onam":         "ഓണം",
			"Rahu Kalam":   "രാഹുകാലം",
			"Yamagandam":   "യമഗണ്ഡം",
		},
		"bn": {
			"Diwali":       "দীপাবলি",
			"Durga Puja":   "দুর্গা পূজা",
			"Kali Puja":    "কালী পূজা",
		},
		"gu": {
			"Diwali":       "દિવાળી",
			"Uttarayan":    "ઉત્તરાયણ",
		},
		"mr": {
			"Diwali":       "दिवाळी",
			"Gudi Padwa":   "गुढी पाडवा",
		},
	}
}

func (m *MultiLanguageLocalizationPlugin) getMuhurtaTranslations() map[string]map[string]string {
	return map[string]map[string]string{
		"ta": {
			"Brahma Muhurta":  "பிரம்ம முகூர்த்தம்",
			"Abhijit Muhurta": "அபிஜித் முகூர்த்தம்",
			"Godhuli Muhurta": "கோதூளி முகூர்த்தம்",
		},
		"ml": {
			"Brahma Muhurta":  "ബ്രഹ്മമുഹൂർത്തം",
			"Abhijit Muhurta": "അഭിജിത് മുഹൂർത്തം",
		},
		"bn": {
			"Brahma Muhurta":  "ব্রহ্ম মুহূর্ত",
			"Abhijit Muhurta": "অভিজিৎ মুহূর্ত",
		},
	}
}
