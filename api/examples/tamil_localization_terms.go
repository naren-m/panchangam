package examples

import "strings"

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

	for english, tamil := range commonPhrases {
		if strings.Contains(strings.ToLower(significance), english) {
			significance = strings.ReplaceAll(significance, english, tamil)
		}
	}

	return significance
}
