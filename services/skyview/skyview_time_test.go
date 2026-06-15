package skyview

import (
	"math"
	"testing"
	"time"
)

func TestDateToJulianDay(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected float64
		delta    float64
	}{
		{
			name:     "J2000 epoch",
			date:     time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: 2451545.0,
			delta:    0.01,
		},
		{
			name:     "Unix epoch",
			date:     time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 2440587.5,
			delta:    0.01,
		},
		{
			name:     "Year 2024",
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 2460310.5,
			delta:    0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jd := DateToJulianDay(tt.date)
			if math.Abs(jd-tt.expected) > tt.delta {
				t.Errorf("DateToJulianDay() = %v, want %v ± %v", jd, tt.expected, tt.delta)
			}
		})
	}
}

func TestGetLocalSiderealTime(t *testing.T) {
	tests := []struct {
		name      string
		longitude float64
		time      time.Time
	}{
		{
			name:      "Greenwich at J2000",
			longitude: 0.0,
			time:      time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "New York",
			longitude: -74.006,
			time:      time.Date(2024, 6, 21, 18, 0, 0, 0, time.UTC),
		},
		{
			name:      "Tokyo",
			longitude: 139.6917,
			time:      time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lst := GetLocalSiderealTime(tt.longitude, tt.time)

			if lst < 0 || lst >= 360 {
				t.Errorf("GetLocalSiderealTime() = %v, want value in [0, 360)", lst)
			}
		})
	}
}
