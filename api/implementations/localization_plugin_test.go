package implementations

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/api"
)

func TestMultiLanguageLocalizationPlugin(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()

	// Test plugin info
	info := plugin.GetInfo()
	if info.Name != "multi_language_localization_plugin" {
		t.Errorf("Expected plugin name 'multi_language_localization_plugin', got %s", info.Name)
	}

	// Test initialization
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled after initialization")
	}

	// Test supported locales
	locales := plugin.GetSupportedLocales()
	if len(locales) == 0 {
		t.Fatal("Expected supported locales, got none")
	}

	expectedLocales := []string{"ta", "ml", "bn", "gu", "mr"}
	for _, expected := range expectedLocales {
		found := false
		for _, locale := range locales {
			if locale == expected || locale == expected+"_IN" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected locale '%s' not found in supported locales", expected)
		}
	}

	// Test supported regions
	regions := plugin.GetSupportedRegions()
	if len(regions) == 0 {
		t.Fatal("Expected supported regions, got none")
	}

	expectedRegions := []api.Region{
		api.RegionTamilNadu,
		api.RegionKerala,
		api.RegionBengal,
		api.RegionGujarat,
		api.RegionMaha,
	}

	for _, expected := range expectedRegions {
		found := false
		for _, region := range regions {
			if region == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected region '%s' not found in supported regions", expected)
		}
	}

	// Test shutdown
	err = plugin.Shutdown(context.Background())
	if err != nil {
		t.Fatalf("Failed to shutdown plugin: %v", err)
	}

	if plugin.IsEnabled() {
		t.Error("Plugin should be disabled after shutdown")
	}
}

func TestLocalizeTithi(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()
	plugin.Initialize(context.Background(), nil)

	testCases := []struct {
		locale       string
		tithiName    string
		expectedLocal string
	}{
		{"ta", "Pratipada", "பிரதமை"},
		{"ta", "Purnima", "பௌர்ணமி"},
		{"ta", "Amavasya", "அமாவாசை"},
		{"ta", "Ekadashi", "ஏகாதசி"},
		{"ml", "Purnima", "പൂർണ്ണിമ"},
		{"ml", "Amavasya", "അമാവാസ്യ"},
		{"bn", "Purnima", "পূর্ণিমা"},
		{"bn", "Amavasya", "অমাবস্যা"},
		{"gu", "Purnima", "પૂર્ણિમા"},
		{"mr", "Purnima", "पूर्णिमा"},
	}

	for _, tc := range testCases {
		t.Run(tc.locale+"_"+tc.tithiName, func(t *testing.T) {
			tithi := &api.Tithi{
				Name:      tc.tithiName,
				Number:    1,
				StartTime: time.Now(),
				EndTime:   time.Now().Add(24 * time.Hour),
			}

			err := plugin.LocalizeTithi(tithi, tc.locale, api.RegionGlobal)
			if err != nil {
				t.Fatalf("Failed to localize tithi: %v", err)
			}

			if tithi.NameLocal != tc.expectedLocal {
				t.Errorf("Expected local name '%s', got '%s'", tc.expectedLocal, tithi.NameLocal)
			}
		})
	}
}

func TestLocalizeNakshatra(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()
	plugin.Initialize(context.Background(), nil)

	testCases := []struct {
		locale        string
		nakshatraName string
		expectedLocal string
	}{
		{"ta", "Ashwini", "அஸ்வினி"},
		{"ta", "Bharani", "பரணி"},
		{"ta", "Krittika", "கார்த்திகை"},
		{"ta", "Rohini", "ரோகிணி"},
		{"ta", "Ardra", "திருவாதிரை"},
		{"ta", "Chitra", "சித்திரை"},
		{"ml", "Ashwini", "അശ്വതി"},
		{"ml", "Bharani", "ഭരണി"},
		{"bn", "Ashwini", "অশ্বিনী"},
	}

	for _, tc := range testCases {
		t.Run(tc.locale+"_"+tc.nakshatraName, func(t *testing.T) {
			nakshatra := &api.Nakshatra{
				Name:      tc.nakshatraName,
				Number:    1,
				StartTime: time.Now(),
				EndTime:   time.Now().Add(24 * time.Hour),
			}

			err := plugin.LocalizeNakshatra(nakshatra, tc.locale, api.RegionGlobal)
			if err != nil {
				t.Fatalf("Failed to localize nakshatra: %v", err)
			}

			if nakshatra.NameLocal != tc.expectedLocal {
				t.Errorf("Expected local name '%s', got '%s'", tc.expectedLocal, nakshatra.NameLocal)
			}
		})
	}
}

func TestLocalizeYoga(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()
	plugin.Initialize(context.Background(), nil)

	testCases := []struct {
		locale       string
		yogaName     string
		expectedLocal string
	}{
		{"ta", "Vishkambha", "விஷ்கம்பம்"},
		{"ta", "Priti", "பிரீதி"},
		{"ta", "Ayushman", "ஆயுஷ்மான்"},
	}

	for _, tc := range testCases {
		t.Run(tc.locale+"_"+tc.yogaName, func(t *testing.T) {
			yoga := &api.Yoga{
				Name:      tc.yogaName,
				Number:    1,
				StartTime: time.Now(),
				EndTime:   time.Now().Add(24 * time.Hour),
			}

			err := plugin.LocalizeYoga(yoga, tc.locale, api.RegionGlobal)
			if err != nil {
				t.Fatalf("Failed to localize yoga: %v", err)
			}

			if yoga.NameLocal != tc.expectedLocal {
				t.Errorf("Expected local name '%s', got '%s'", tc.expectedLocal, yoga.NameLocal)
			}
		})
	}
}

func TestLocalizeKarana(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()
	plugin.Initialize(context.Background(), nil)

	testCases := []struct {
		locale        string
		karanaName    string
		expectedLocal string
	}{
		{"ta", "Bava", "பவ"},
		{"ta", "Balava", "பாலவ"},
		{"ta", "Vishti", "விஷ்டி"},
	}

	for _, tc := range testCases {
		t.Run(tc.locale+"_"+tc.karanaName, func(t *testing.T) {
			karana := &api.Karana{
				Name:      tc.karanaName,
				Number:    1,
				StartTime: time.Now(),
				EndTime:   time.Now().Add(12 * time.Hour),
			}

			err := plugin.LocalizeKarana(karana, tc.locale, api.RegionGlobal)
			if err != nil {
				t.Fatalf("Failed to localize karana: %v", err)
			}

			if karana.NameLocal != tc.expectedLocal {
				t.Errorf("Expected local name '%s', got '%s'", tc.expectedLocal, karana.NameLocal)
			}
		})
	}
}

func TestLocalizeEvent(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()
	plugin.Initialize(context.Background(), nil)

	testCases := []struct {
		locale        string
		eventName     string
		expectedLocal string
	}{
		{"ta", "Diwali", "தீபாவளி"},
		{"ta", "Rahu Kalam", "ராகு காலம்"},
		{"ta", "Yamagandam", "யமகண்டம்"},
		{"ml", "Diwali", "ദീപാവലി"},
		{"ml", "Vishu", "വിഷു"},
		{"bn", "Diwali", "দীপাবলি"},
		{"bn", "Durga Puja", "দুর্গা পূজা"},
		{"gu", "Diwali", "દિવાળી"},
		{"mr", "Diwali", "दिवाळी"},
	}

	for _, tc := range testCases {
		t.Run(tc.locale+"_"+tc.eventName, func(t *testing.T) {
			event := &api.Event{
				Name:      tc.eventName,
				Type:      api.EventTypeFestival,
				StartTime: time.Now(),
				EndTime:   time.Now().Add(24 * time.Hour),
			}

			err := plugin.LocalizeEvent(event, tc.locale, api.RegionGlobal)
			if err != nil {
				t.Fatalf("Failed to localize event: %v", err)
			}

			if event.NameLocal != tc.expectedLocal {
				t.Errorf("Expected local name '%s', got '%s'", tc.expectedLocal, event.NameLocal)
			}
		})
	}
}

func TestLocalizeMuhurta(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()
	plugin.Initialize(context.Background(), nil)

	testCases := []struct {
		locale        string
		muhurtaName   string
		expectedLocal string
	}{
		{"ta", "Brahma Muhurta", "பிரம்ம முகூர்த்தம்"},
		{"ta", "Abhijit Muhurta", "அபிஜித் முகூர்த்தம்"},
		{"ta", "Godhuli Muhurta", "கோதூளி முகூர்த்தம்"},
		{"ml", "Brahma Muhurta", "ബ്രഹ്മമുഹൂർത്തം"},
		{"bn", "Brahma Muhurta", "ব্রহ্ম মুহূর্ত"},
	}

	for _, tc := range testCases {
		t.Run(tc.locale+"_"+tc.muhurtaName, func(t *testing.T) {
			muhurta := &api.Muhurta{
				Name:      tc.muhurtaName,
				StartTime: time.Now(),
				EndTime:   time.Now().Add(96 * time.Minute),
				Quality:   api.QualityHighly,
			}

			err := plugin.LocalizeMuhurta(muhurta, tc.locale, api.RegionGlobal)
			if err != nil {
				t.Fatalf("Failed to localize muhurta: %v", err)
			}

			if muhurta.NameLocal != tc.expectedLocal {
				t.Errorf("Expected local name '%s', got '%s'", tc.expectedLocal, muhurta.NameLocal)
			}
		})
	}
}

func TestGetLanguageFromLocale(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()

	testCases := []struct {
		locale   string
		expected string
	}{
		{"ta", "ta"},
		{"ta_IN", "ta"},
		{"ml_IN", "ml"},
		{"en_US", "en"},
		{"hi", "hi"},
		{"", "en"}, // Default
		{"x", "en"},  // Single character - defaults to en
	}

	for _, tc := range testCases {
		result := plugin.getLanguageFromLocale(tc.locale)
		if result != tc.expected {
			t.Errorf("For locale '%s', expected language '%s', got '%s'", tc.locale, tc.expected, result)
		}
	}
}

func TestLocalizationWithUnknownLocale(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()
	plugin.Initialize(context.Background(), nil)

	// Test with unknown locale - should not error, just not localize
	tithi := &api.Tithi{
		Name:      "Pratipada",
		Number:    1,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}

	err := plugin.LocalizeTithi(tithi, "xx", api.RegionGlobal) // Unknown locale
	if err != nil {
		t.Fatalf("Should not error on unknown locale: %v", err)
	}

	// NameLocal might be empty or same as Name
	if tithi.Name != "Pratipada" {
		t.Error("Original name should not be changed")
	}
}

func TestLocalizationWithDisabledPlugin(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()
	// Don't initialize - plugin should be disabled

	tithi := &api.Tithi{
		Name:   "Pratipada",
		Number: 1,
	}

	err := plugin.LocalizeTithi(tithi, "ta", api.RegionGlobal)
	if err == nil {
		t.Error("Expected error when plugin is not enabled")
	}

	if err.Error() != "localization plugin is not enabled" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestAllLanguagesHaveTranslations(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()
	plugin.Initialize(context.Background(), nil)

	languages := []string{"ta", "ml", "bn", "gu", "mr"}

	for _, lang := range languages {
		t.Run("Language_"+lang, func(t *testing.T) {
			// Test that each language has at least some translations
			tithiTranslations := plugin.getTithiTranslations()
			if _, exists := tithiTranslations[lang]; !exists {
				t.Errorf("Language '%s' has no Tithi translations", lang)
			}

			// At least Purnima and Amavasya should be translated for all languages
			commonTithis := []string{"Purnima", "Amavasya"}
			for _, tithi := range commonTithis {
				if _, exists := tithiTranslations[lang][tithi]; !exists {
					t.Errorf("Language '%s' missing translation for '%s'", lang, tithi)
				}
			}
		})
	}
}
