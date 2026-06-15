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
		{"", "en"},  // Default
		{"x", "en"}, // Single character - defaults to en
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
