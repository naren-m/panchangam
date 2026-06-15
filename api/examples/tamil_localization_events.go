package examples

import "github.com/naren-m/panchangam/api"

// LocalizeEvent returns localized event information
func (t *TamilLocalizationPlugin) LocalizeEvent(event *api.Event, locale string, region api.Region) error {
	if err := t.ensureEnabled(); err != nil {
		return err
	}
	if !t.isSupported(locale) {
		return nil
	}

	if localName, exists := t.eventNames[event.Name]; exists {
		event.NameLocal = localName
	}

	if event.Significance != "" {
		event.Significance = t.localizeSignificance(event.Significance)
	}

	return nil
}

// LocalizeMuhurta returns localized muhurta information
func (t *TamilLocalizationPlugin) LocalizeMuhurta(muhurta *api.Muhurta, locale string, region api.Region) error {
	if err := t.ensureEnabled(); err != nil {
		return err
	}
	if !t.isSupported(locale) {
		return nil
	}

	if localName, exists := t.muhurtaNames[muhurta.Name]; exists {
		muhurta.NameLocal = localName
	}

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

	if muhurta.Significance != "" {
		muhurta.Significance = t.localizeSignificance(muhurta.Significance)
	}

	return nil
}
