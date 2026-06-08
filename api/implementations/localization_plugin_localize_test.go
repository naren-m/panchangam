package implementations

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/api"
)

func TestLocalizeTithi(t *testing.T) {
	plugin := NewMultiLanguageLocalizationPlugin()
	plugin.Initialize(context.Background(), nil)

	testCases := []struct {
		locale        string
		tithiName     string
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
		locale        string
		yogaName      string
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
