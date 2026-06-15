package examples

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
