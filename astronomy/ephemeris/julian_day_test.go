package ephemeris

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJulianDayConversion(t *testing.T) {
	tests := []struct {
		name      string
		time      time.Time
		expected  JulianDay
		tolerance float64
	}{
		{
			name:      "J2000.0 epoch",
			time:      time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
			expected:  JulianDay(2451545.0),
			tolerance: 0.001,
		},
		{
			name:      "Unix epoch",
			time:      time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:  JulianDay(2440587.5),
			tolerance: 0.001,
		},
		{
			name:      "Current test date",
			time:      time.Date(2024, 7, 18, 0, 0, 0, 0, time.UTC),
			expected:  JulianDay(2460509.5),
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jd := TimeToJulianDay(tt.time)
			assert.InDelta(t, float64(tt.expected), float64(jd), tt.tolerance)

			converted := JulianDayToTime(jd)
			assert.WithinDuration(t, tt.time, converted, time.Minute)
		})
	}
}
