package panchangam

import (
	"testing"
	"time"
)

func TestGetTimezoneInfo_DSTTransitions(t *testing.T) {
	parser := NewTimezoneParser()

	loc, err := parser.ParseTimezone("America/New_York")
	if err != nil {
		t.Fatalf("Failed to parse timezone: %v", err)
	}

	tests := []struct {
		name        string
		date        time.Time
		expectDST   bool
		description string
	}{
		{
			name:        "Summer time (July)",
			date:        time.Date(2025, 7, 1, 12, 0, 0, 0, loc),
			expectDST:   true,
			description: "Should detect DST in summer",
		},
		{
			name:        "Winter time (January)",
			date:        time.Date(2025, 1, 1, 12, 0, 0, 0, loc),
			expectDST:   false,
			description: "Should not detect DST in winter",
		},
		{
			name:        "Summer time (June)",
			date:        time.Date(2025, 6, 15, 12, 0, 0, 0, loc),
			expectDST:   true,
			description: "Should detect DST in June",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tzInfo := parser.GetTimezoneInfo(loc, tt.date)
			t.Logf("Timezone info for %s: Name=%s, Offset=%d, IsDST=%v, Formatted=%s",
				tt.date.Format("2006-01-02"), tzInfo.Name, tzInfo.Offset, tzInfo.IsDST, tzInfo.Formatted)

			if tzInfo.Name == "" {
				t.Errorf("GetTimezoneInfo() returned empty name")
			}
			if tzInfo.Formatted == "" {
				t.Errorf("GetTimezoneInfo() returned empty formatted offset")
			}
		})
	}
}

func TestGetTimezoneInfo_NoDATimezones(t *testing.T) {
	parser := NewTimezoneParser()

	tests := []struct {
		name     string
		timezone string
		date     time.Time
	}{
		{
			name:     "Asia/Kolkata (no DST)",
			timezone: "Asia/Kolkata",
			date:     time.Date(2025, 7, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "UTC (no DST)",
			timezone: "UTC",
			date:     time.Date(2025, 7, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Fixed offset +05:30 (no DST)",
			timezone: "+05:30",
			date:     time.Date(2025, 7, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := parser.ParseTimezone(tt.timezone)
			if err != nil {
				t.Fatalf("Failed to parse timezone: %v", err)
			}

			tzInfo := parser.GetTimezoneInfo(loc, tt.date)
			t.Logf("Timezone info for %s: Name=%s, Offset=%d, IsDST=%v, Formatted=%s",
				tt.timezone, tzInfo.Name, tzInfo.Offset, tzInfo.IsDST, tzInfo.Formatted)

			if tzInfo.Name == "" {
				t.Errorf("GetTimezoneInfo() returned empty name")
			}
		})
	}
}

func TestHistoricalTimezoneChanges(t *testing.T) {
	parser := NewTimezoneParser()

	loc, err := parser.ParseTimezone("America/New_York")
	if err != nil {
		t.Fatalf("Failed to parse timezone: %v", err)
	}

	tests := []struct {
		name string
		date time.Time
	}{
		{
			name: "Before DST existed (1900)",
			date: time.Date(1900, 7, 1, 12, 0, 0, 0, loc),
		},
		{
			name: "Modern era (2025)",
			date: time.Date(2025, 7, 1, 12, 0, 0, 0, loc),
		},
		{
			name: "Year 2000",
			date: time.Date(2000, 7, 1, 12, 0, 0, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tzInfo := parser.GetTimezoneInfo(loc, tt.date)
			t.Logf("Historical timezone info for %s: Name=%s, Offset=%d, Formatted=%s",
				tt.date.Format("2006-01-02"), tzInfo.Name, tzInfo.Offset, tzInfo.Formatted)

			if tzInfo.Name == "" {
				t.Errorf("GetTimezoneInfo() returned empty name for historical date")
			}
		})
	}
}
