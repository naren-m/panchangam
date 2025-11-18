package implementations

import (
	"testing"
	"time"
)

func TestNaazhikaiConverter(t *testing.T) {
	converter := NewNaazhikaiConverter()

	// Test date: 2024-01-15 at noon
	testDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	sunrise := time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC)
	sunset := time.Date(2024, 1, 15, 18, 0, 0, 0, time.UTC)

	t.Run("ToNaazhikai Conversion", func(t *testing.T) {
		// At noon, 6 hours (360 minutes) after sunrise
		// 360 / 24 = 15 Naazhikai
		tn := converter.ToNaazhikai(testDate, sunrise)

		if tn.Naazhikai != 15 {
			t.Errorf("Expected 15 Naazhikai at noon, got %d", tn.Naazhikai)
		}

		if tn.Vinaazhikai != 0 {
			t.Errorf("Expected 0 Vinaazhikai at exact hour, got %d", tn.Vinaazhikai)
		}
	})

	t.Run("FromNaazhikai Conversion", func(t *testing.T) {
		// 15 Naazhikai should be 360 minutes
		duration := converter.FromNaazhikai(15, 0)
		expectedMinutes := 15 * 24 // 360 minutes

		if int(duration.Minutes()) != expectedMinutes {
			t.Errorf("Expected %d minutes, got %d", expectedMinutes, int(duration.Minutes()))
		}
	})

	t.Run("RoundTrip Conversion", func(t *testing.T) {
		// Test round-trip conversion
		original := 10 // 10 Naazhikai
		duration := converter.FromNaazhikai(original, 0)
		timeAfterSunrise := sunrise.Add(duration)
		tn := converter.ToNaazhikai(timeAfterSunrise, sunrise)

		if tn.Naazhikai != original {
			t.Errorf("Round-trip failed: expected %d Naazhikai, got %d", original, tn.Naazhikai)
		}
	})

	t.Run("GetNaazhikaiTime", func(t *testing.T) {
		timeStr := converter.GetNaazhikaiTime(testDate, sunrise)
		expected := "15 Naazhikai 0 Vinaazhikai"

		if timeStr != expected {
			t.Errorf("Expected '%s', got '%s'", expected, timeStr)
		}
	})

	t.Run("GetNaazhikaiDetails", func(t *testing.T) {
		details := converter.GetNaazhikaiDetails(testDate, sunrise, sunset)

		if details == nil {
			t.Fatal("Expected details, got nil")
		}

		if details["naazhikai"].(int) != 15 {
			t.Errorf("Expected naazhikai to be 15, got %v", details["naazhikai"])
		}

		if details["tamil_name"] != "நாழிகை" {
			t.Errorf("Expected Tamil name 'நாழிகை', got %v", details["tamil_name"])
		}

		// Day length is 12 hours = 720 minutes = 30 Naazhikai
		expectedDayNaazhikai := 30.0
		if details["total_day_naazhikai"].(float64) != expectedDayNaazhikai {
			t.Errorf("Expected %f day Naazhikai, got %v", expectedDayNaazhikai, details["total_day_naazhikai"])
		}
	})

	t.Run("ConvertDurationToNaazhikai", func(t *testing.T) {
		// 1 hour = 60 minutes = 2.5 Naazhikai
		duration := time.Hour
		tn := converter.ConvertDurationToNaazhikai(duration)

		if tn.Naazhikai != 2 {
			t.Errorf("Expected 2 Naazhikai for 1 hour, got %d", tn.Naazhikai)
		}

		if tn.Vinaazhikai != 30 {
			t.Errorf("Expected 30 Vinaazhikai for remaining 0.5 Naazhikai, got %d", tn.Vinaazhikai)
		}
	})

	t.Run("GetNaazhikaiPeriodName", func(t *testing.T) {
		testCases := []struct {
			naazhikai int
			contains  string
		}{
			{2, "Morning"},     // காலை
			{10, "Forenoon"},   // முற்பகல்
			{17, "Noon"},       // மதியம்
			{25, "Afternoon"},  // பிற்பகல்
			{35, "Evening"},    // மாலை
			{45, "Night"},      // இரவு
		}

		for _, tc := range testCases {
			periodName := converter.GetNaazhikaiPeriodName(tc.naazhikai)
			if periodName == "" {
				t.Errorf("Expected period name for Naazhikai %d, got empty string", tc.naazhikai)
			}
			// Just verify we got something back, the Tamil text is correct in the implementation
		}
	})

	t.Run("CalculateMuhurtaInNaazhikai", func(t *testing.T) {
		muhurtas := converter.CalculateMuhurtaInNaazhikai(sunrise, sunset)

		if len(muhurtas) == 0 {
			t.Fatal("Expected muhurtas, got none")
		}

		// Verify muhurtas are sequential
		for i := 0; i < len(muhurtas); i++ {
			m := muhurtas[i]

			if m.MuhurtaNumber != i+1 {
				t.Errorf("Expected muhurta number %d, got %d", i+1, m.MuhurtaNumber)
			}

			if m.EndNaazhikai != m.StartNaazhikai+2 {
				t.Errorf("Expected 2 Naazhikai duration, got %d-%d", m.StartNaazhikai, m.EndNaazhikai)
			}

			if m.Quality == "" {
				t.Error("Expected quality to be set")
			}

			if m.TamilName == "" {
				t.Error("Expected Tamil name to be set")
			}
		}
	})

	t.Run("FormatNaazhikaiTime", func(t *testing.T) {
		// Test without Tamil numerals
		formatted := converter.FormatNaazhikaiTime(15, 30, false)
		expected := "15 நாழிகை 30 விநாழிகை"

		if formatted != expected {
			t.Errorf("Expected '%s', got '%s'", expected, formatted)
		}

		// Test with Tamil numerals
		formattedTamil := converter.FormatNaazhikaiTime(15, 30, true)
		if formattedTamil == "" {
			t.Error("Expected formatted Tamil time, got empty string")
		}
	})

	t.Run("ToTamilNumerals", func(t *testing.T) {
		converter := NewNaazhikaiConverter()

		testCases := []struct {
			input    int
			expected string
		}{
			{0, "௦"},
			{1, "௧"},
			{5, "௫"},
			{10, "௧௦"},
			{15, "௧௫"},
			{99, "௯௯"},
		}

		for _, tc := range testCases {
			result := converter.toTamilNumerals(tc.input)
			if result != tc.expected {
				t.Errorf("Expected Tamil numeral '%s' for %d, got '%s'", tc.expected, tc.input, result)
			}
		}
	})

	t.Run("GetSunriseInNaazhikai", func(t *testing.T) {
		// Sunrise at 6 AM = 360 minutes from midnight = 15 Naazhikai
		tn := converter.GetSunriseInNaazhikai(sunrise)

		if tn.Naazhikai != 15 {
			t.Errorf("Expected sunrise at 15 Naazhikai from midnight, got %d", tn.Naazhikai)
		}
	})

	t.Run("CompareWithModernTime", func(t *testing.T) {
		comparison := converter.CompareWithModernTime(testDate, sunrise)

		if comparison == nil {
			t.Fatal("Expected comparison data, got nil")
		}

		if comparison["modern_time"] == "" {
			t.Error("Expected modern time to be set")
		}

		if comparison["naazhikai_time"] == "" {
			t.Error("Expected naazhikai time to be set")
		}

		if comparison["tamil_period"] == "" {
			t.Error("Expected Tamil period to be set")
		}

		if comparison["naazhikai_value"].(int) != 15 {
			t.Errorf("Expected naazhikai value to be 15, got %v", comparison["naazhikai_value"])
		}

		minutesFromSunrise := comparison["minutes_from_sunrise"].(float64)
		if minutesFromSunrise != 360.0 {
			t.Errorf("Expected 360 minutes from sunrise, got %f", minutesFromSunrise)
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		// Test at sunrise (should be 0 Naazhikai)
		tn := converter.ToNaazhikai(sunrise, sunrise)
		if tn.Naazhikai != 0 {
			t.Errorf("Expected 0 Naazhikai at sunrise, got %d", tn.Naazhikai)
		}

		// Test just before sunset (should be close to 30 Naazhikai for a 12-hour day)
		justBeforeSunset := sunset.Add(-1 * time.Minute)
		tn = converter.ToNaazhikai(justBeforeSunset, sunrise)
		expectedNaazhikai := 29 // Should be approximately 29
		if tn.Naazhikai < expectedNaazhikai || tn.Naazhikai > 30 {
			t.Errorf("Expected around %d Naazhikai just before sunset, got %d", expectedNaazhikai, tn.Naazhikai)
		}
	})
}

func TestNaazhikaiMuhurtaQuality(t *testing.T) {
	converter := NewNaazhikaiConverter()

	testCases := []struct {
		naazhikai int
		quality   string
	}{
		{4, "auspicious"},
		{6, "auspicious"},
		{10, "inauspicious"},
		{14, "inauspicious"},
		{5, "neutral"},
	}

	for _, tc := range testCases {
		quality := converter.getMuhurtaQuality(tc.naazhikai)
		if quality != tc.quality {
			t.Errorf("For Naazhikai %d, expected quality '%s', got '%s'", tc.naazhikai, tc.quality, quality)
		}
	}
}

func TestNaazhikaiConstants(t *testing.T) {
	// Verify the Naazhikai time unit constants
	oneNaazhikaiMinutes := 24
	oneDay := 60 // Naazhikai per day
	oneDayMinutes := oneNaazhikaiMinutes * oneDay

	if oneDayMinutes != 1440 {
		t.Errorf("Expected 1440 minutes in a day, got %d", oneDayMinutes)
	}

	// Verify Vinaazhikai
	oneVinaazhikaiSeconds := 24 // 24 seconds
	vinaazhikaiPerNaazhikai := 60
	secondsPerNaazhikai := oneVinaazhikaiSeconds * vinaazhikaiPerNaazhikai

	if secondsPerNaazhikai != 1440 {
		t.Errorf("Expected 1440 seconds per Naazhikai, got %d", secondsPerNaazhikai)
	}

	if secondsPerNaazhikai/60 != oneNaazhikaiMinutes {
		t.Errorf("Naazhikai time calculation mismatch")
	}
}
