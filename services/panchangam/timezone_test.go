package panchangam

import (
	"testing"
	"time"
)

func TestParseTimezone_IANATimezones(t *testing.T) {
	parser := NewTimezoneParser()

	tests := []struct {
		name     string
		timezone string
		wantErr  bool
	}{
		{
			name:     "Asia/Kolkata",
			timezone: "Asia/Kolkata",
			wantErr:  false,
		},
		{
			name:     "America/New_York",
			timezone: "America/New_York",
			wantErr:  false,
		},
		{
			name:     "America/Los_Angeles",
			timezone: "America/Los_Angeles",
			wantErr:  false,
		},
		{
			name:     "Europe/London",
			timezone: "Europe/London",
			wantErr:  false,
		},
		{
			name:     "Australia/Sydney",
			timezone: "Australia/Sydney",
			wantErr:  false,
		},
		{
			name:     "UTC",
			timezone: "UTC",
			wantErr:  false,
		},
		{
			name:     "GMT",
			timezone: "GMT",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := parser.ParseTimezone(tt.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimezone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && loc == nil {
				t.Errorf("ParseTimezone() returned nil location")
			}
			if !tt.wantErr {
				t.Logf("Successfully parsed %s: %s", tt.timezone, loc.String())
			}
		})
	}
}

func TestParseTimezone_TrimsWhitespace(t *testing.T) {
	parser := NewTimezoneParser()

	loc, err := parser.ParseTimezone(" Asia/Kolkata ")
	if err != nil {
		t.Fatalf("ParseTimezone() with IANA whitespace returned error: %v", err)
	}
	if loc.String() != "Asia/Kolkata" {
		t.Fatalf("expected normalized IANA timezone, got %s", loc.String())
	}

	loc, err = parser.ParseTimezone(" +05:30 ")
	if err != nil {
		t.Fatalf("ParseTimezone() with offset whitespace returned error: %v", err)
	}
	_, offset := time.Date(2025, 1, 1, 12, 0, 0, 0, loc).Zone()
	if offset != 5*3600+30*60 {
		t.Fatalf("expected +05:30 offset, got %d", offset)
	}
}

func TestParseTimezone_InvalidFormats(t *testing.T) {
	parser := NewTimezoneParser()

	tests := []struct {
		name     string
		timezone string
		wantErr  bool
	}{
		{
			name:     "Invalid IANA timezone",
			timezone: "Invalid/Timezone",
			wantErr:  true,
		},
		{
			name:     "Invalid offset format",
			timezone: "+99:99",
			wantErr:  true,
		},
		{
			name:     "Offset out of range",
			timezone: "+15:00",
			wantErr:  true,
		},
		{
			name:     "Invalid characters",
			timezone: "ABC123",
			wantErr:  true,
		},
		{
			name:     "Missing sign",
			timezone: "05:30",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := parser.ParseTimezone(tt.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimezone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && loc == nil {
				t.Errorf("ParseTimezone() returned nil location without error")
			}
			if tt.wantErr && loc != nil {
				t.Errorf("ParseTimezone() returned location %v, expected error", loc)
			}
		})
	}
}

func TestParseTimezone_EdgeCases(t *testing.T) {
	parser := NewTimezoneParser()

	tests := []struct {
		name     string
		timezone string
		wantErr  bool
	}{
		{
			name:     "Empty string defaults to UTC",
			timezone: "",
			wantErr:  false,
		},
		{
			name:     "UTC",
			timezone: "UTC",
			wantErr:  false,
		},
		{
			name:     "GMT",
			timezone: "GMT",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := parser.ParseTimezone(tt.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimezone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && loc == nil {
				t.Errorf("ParseTimezone() returned nil location")
			}
		})
	}
}
