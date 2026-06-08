package panchangam

import (
	"strings"
	"testing"
	"time"
)

func TestValidateTimezoneForLocation(t *testing.T) {
	parser := NewTimezoneParser()

	tests := []struct {
		name        string
		timezone    string
		latitude    float64
		longitude   float64
		date        time.Time
		expectWarn  bool
		wantWarning string
		description string
	}{
		{
			name:        "New York coordinates with New York timezone",
			timezone:    "America/New_York",
			latitude:    40.7128,
			longitude:   -74.0060,
			date:        time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
			expectWarn:  false,
			description: "Should validate correctly",
		},
		{
			name:        "India coordinates with IST",
			timezone:    "Asia/Kolkata",
			latitude:    13.0827,
			longitude:   80.2707,
			date:        time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
			expectWarn:  false,
			description: "Should validate correctly",
		},
		{
			name:        "Wrong timezone for location",
			timezone:    "Asia/Kolkata",
			latitude:    40.7128,
			longitude:   -74.0060,
			date:        time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
			expectWarn:  true,
			description: "Should warn about mismatch",
		},
		{
			name:        "UTC offset matching location",
			timezone:    "-08:00",
			latitude:    37.7749,
			longitude:   -122.4194,
			date:        time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
			expectWarn:  false,
			description: "Should validate correctly",
		},
		{
			name:        "warning uses request date offset",
			timezone:    "America/New_York",
			latitude:    13.0827,
			longitude:   80.2707,
			date:        time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
			expectWarn:  true,
			wantWarning: "offset -05:00",
			description: "Should report the timezone offset for the request date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := parser.ParseTimezone(tt.timezone)
			if err != nil {
				t.Fatalf("Failed to parse timezone: %v", err)
			}

			isValid, warning := parser.ValidateTimezoneForLocation(loc, tt.latitude, tt.longitude, tt.date)
			if tt.expectWarn && isValid {
				t.Errorf("ValidateTimezoneForLocation() expected warning but got valid")
			}
			if !tt.expectWarn && !isValid {
				t.Errorf("ValidateTimezoneForLocation() unexpected warning: %s", warning)
			}
			if tt.wantWarning != "" && !strings.Contains(warning, tt.wantWarning) {
				t.Errorf("ValidateTimezoneForLocation() warning = %q, want to contain %q", warning, tt.wantWarning)
			}
			t.Logf("Validation result for %s at (%.4f, %.4f): valid=%v, warning=%s",
				tt.timezone, tt.latitude, tt.longitude, isValid, warning)
		})
	}
}
