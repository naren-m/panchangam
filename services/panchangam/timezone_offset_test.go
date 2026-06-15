package panchangam

import (
	"testing"
	"time"
)

func TestParseTimezone_UTCOffsets(t *testing.T) {
	parser := NewTimezoneParser()

	tests := []struct {
		name           string
		timezone       string
		wantErr        bool
		expectedOffset int
	}{
		{
			name:           "Positive offset with colon",
			timezone:       "+05:30",
			wantErr:        false,
			expectedOffset: 5*3600 + 30*60,
		},
		{
			name:           "Negative offset with colon",
			timezone:       "-08:00",
			wantErr:        false,
			expectedOffset: -8 * 3600,
		},
		{
			name:           "UTC prefix positive",
			timezone:       "UTC+05:30",
			wantErr:        false,
			expectedOffset: 5*3600 + 30*60,
		},
		{
			name:           "GMT prefix negative",
			timezone:       "GMT-08:00",
			wantErr:        false,
			expectedOffset: -8 * 3600,
		},
		{
			name:           "Zero offset",
			timezone:       "+00:00",
			wantErr:        false,
			expectedOffset: 0,
		},
		{
			name:           "Offset without minutes",
			timezone:       "+05:00",
			wantErr:        false,
			expectedOffset: 5 * 3600,
		},
		{
			name:           "Maximum positive offset",
			timezone:       "+14:00",
			wantErr:        false,
			expectedOffset: 14 * 3600,
		},
		{
			name:           "Maximum negative offset",
			timezone:       "-14:00",
			wantErr:        false,
			expectedOffset: -14 * 3600,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := parser.ParseTimezone(tt.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimezone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if loc == nil {
					t.Errorf("ParseTimezone() returned nil location")
					return
				}
				testTime := time.Date(2025, 1, 1, 12, 0, 0, 0, loc)
				_, offset := testTime.Zone()
				if offset != tt.expectedOffset {
					t.Errorf("ParseTimezone() offset = %d, want %d", offset, tt.expectedOffset)
				}
				t.Logf("Successfully parsed %s with offset %d seconds", tt.timezone, offset)
			}
		})
	}
}

func TestFormatTimezoneOffset(t *testing.T) {
	tests := []struct {
		name          string
		offsetSeconds int
		expected      string
	}{
		{
			name:          "Positive offset with minutes",
			offsetSeconds: 5*3600 + 30*60,
			expected:      "+05:30",
		},
		{
			name:          "Negative offset",
			offsetSeconds: -8 * 3600,
			expected:      "-08:00",
		},
		{
			name:          "Zero offset",
			offsetSeconds: 0,
			expected:      "+00:00",
		},
		{
			name:          "Positive offset without minutes",
			offsetSeconds: 5 * 3600,
			expected:      "+05:00",
		},
		{
			name:          "Negative offset with minutes",
			offsetSeconds: -(8*3600 + 30*60),
			expected:      "-08:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimezoneOffset(tt.offsetSeconds)
			if result != tt.expected {
				t.Errorf("formatTimezoneOffset(%d) = %s, want %s", tt.offsetSeconds, result, tt.expected)
			}
		})
	}
}
