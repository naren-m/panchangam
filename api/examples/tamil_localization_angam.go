package examples

import "github.com/naren-m/panchangam/api"

// LocalizeTithi returns localized tithi information
func (t *TamilLocalizationPlugin) LocalizeTithi(tithi *api.Tithi, locale string, region api.Region) error {
	if err := t.ensureEnabled(); err != nil {
		return err
	}
	if !t.isSupported(locale) {
		return nil
	}

	if localName, exists := t.tithiNames[tithi.Name]; exists {
		tithi.NameLocal = localName
	}

	if tithi.Lord != "" {
		tithi.Lord = t.localizeDivineNames(tithi.Lord)
	}

	if tithi.Quality != "" {
		tithi.Quality = t.localizeQuality(tithi.Quality)
	}

	return nil
}

// LocalizeNakshatra returns localized nakshatra information
func (t *TamilLocalizationPlugin) LocalizeNakshatra(nakshatra *api.Nakshatra, locale string, region api.Region) error {
	if err := t.ensureEnabled(); err != nil {
		return err
	}
	if !t.isSupported(locale) {
		return nil
	}

	if localName, exists := t.nakshatraNames[nakshatra.Name]; exists {
		nakshatra.NameLocal = localName
	}

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
	if err := t.ensureEnabled(); err != nil {
		return err
	}
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
	if err := t.ensureEnabled(); err != nil {
		return err
	}
	if !t.isSupported(locale) {
		return nil
	}

	if localName, exists := t.karanaNames[karana.Name]; exists {
		karana.NameLocal = localName
	}

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
