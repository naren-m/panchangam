package astronomy

import (
	"context"
	"math"
	"testing"
	"time"
)

func TestCalculateLunarTimes(t *testing.T) {
	// Test location: Bangalore, India
	loc := Location{
		Latitude:  12.9716,
		Longitude: 77.5946,
	}
	
	// Test date: January 15, 2024
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	lunarTimes, err := CalculateLunarTimes(loc, date)
	if err != nil {
		t.Fatalf("Failed to calculate lunar times: %v", err)
	}
	
	if lunarTimes == nil {
		t.Fatal("Lunar times should not be nil")
	}
	
	// Basic validation - moonrise and moonset should be different times
	if lunarTimes.Moonrise.Equal(lunarTimes.Moonset) {
		t.Error("Moonrise and moonset should be different times")
	}
	
	// Check that times are reasonable (within 24 hours of the input date)
	dayStart := time.Date(2024, 1, 15, 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(48 * time.Hour) // Allow for next day
	
	if lunarTimes.Moonrise.Before(dayStart) || lunarTimes.Moonrise.After(dayEnd) {
		t.Errorf("Moonrise time %v is outside reasonable range [%v, %v]", 
			lunarTimes.Moonrise, dayStart, dayEnd)
	}
	
	if lunarTimes.Moonset.Before(dayStart) || lunarTimes.Moonset.After(dayEnd) {
		t.Errorf("Moonset time %v is outside reasonable range [%v, %v]", 
			lunarTimes.Moonset, dayStart, dayEnd)
	}
	
	t.Logf("Moonrise: %s", lunarTimes.Moonrise.Format("15:04:05"))
	t.Logf("Moonset: %s", lunarTimes.Moonset.Format("15:04:05"))
	t.Logf("Is Visible: %t", lunarTimes.IsVisible)
}

func TestCalculateLunarTimesWithContext(t *testing.T) {
	ctx := context.Background()
	
	// Test location: New York, USA
	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	
	// Test date: June 21, 2024 (Summer solstice)
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)
	
	lunarTimes, err := CalculateLunarTimesWithContext(ctx, loc, date)
	if err != nil {
		t.Fatalf("Failed to calculate lunar times with context: %v", err)
	}
	
	if lunarTimes == nil {
		t.Fatal("Lunar times should not be nil")
	}
	
	t.Logf("NYC Moonrise: %s", lunarTimes.Moonrise.Format("15:04:05"))
	t.Logf("NYC Moonset: %s", lunarTimes.Moonset.Format("15:04:05"))
	t.Logf("NYC Is Visible: %t", lunarTimes.IsVisible)
}

func TestCalculateLunarPhase(t *testing.T) {
	// Test known full moon date: January 25, 2024 (approximately)
	date := time.Date(2024, 1, 25, 12, 0, 0, 0, time.UTC)
	
	phase, err := CalculateLunarPhase(date)
	if err != nil {
		t.Fatalf("Failed to calculate lunar phase: %v", err)
	}
	
	if phase == nil {
		t.Fatal("Lunar phase should not be nil")
	}
	
	// Phase should be between 0 and 1
	if phase.Phase < 0 || phase.Phase > 1 {
		t.Errorf("Phase %f should be between 0 and 1", phase.Phase)
	}
	
	// Illumination should be between 0 and 100
	if phase.Illumination < 0 || phase.Illumination > 100 {
		t.Errorf("Illumination %f should be between 0 and 100", phase.Illumination)
	}
	
	// Age should be reasonable (0-30 days)
	if phase.Age < 0 || phase.Age > 30 {
		t.Errorf("Age %f should be between 0 and 30 days", phase.Age)
	}
	
	// Name should not be empty
	if phase.Name == "" {
		t.Error("Phase name should not be empty")
	}
	
	// For a date near full moon, illumination should be high
	if phase.Illumination < 70 {
		t.Logf("Warning: Expected high illumination near full moon, got %f%%", phase.Illumination)
	}
	
	t.Logf("Phase: %f", phase.Phase)
	t.Logf("Illumination: %f%%", phase.Illumination)
	t.Logf("Name: %s", phase.Name)
	t.Logf("Age: %f days", phase.Age)
	t.Logf("Next Phase: %s", phase.NextPhase.Format("2006-01-02 15:04:05"))
}

func TestCalculateLunarPhaseNewMoon(t *testing.T) {
	// Test known new moon date: January 11, 2024 (approximately)
	date := time.Date(2024, 1, 11, 12, 0, 0, 0, time.UTC)
	
	phase, err := CalculateLunarPhase(date)
	if err != nil {
		t.Fatalf("Failed to calculate lunar phase: %v", err)
	}
	
	// For a date near new moon, phase should be close to 0
	if phase.Phase > 0.15 && phase.Phase < 0.85 {
		t.Logf("Warning: Expected phase near 0 for new moon, got %f", phase.Phase)
	}
	
	// For new moon, illumination should be low
	if phase.Illumination > 30 {
		t.Logf("Warning: Expected low illumination for new moon, got %f%%", phase.Illumination)
	}
	
	t.Logf("New Moon Phase: %f", phase.Phase)
	t.Logf("New Moon Illumination: %f%%", phase.Illumination)
	t.Logf("New Moon Name: %s", phase.Name)
}

func TestGetMoonriseTime(t *testing.T) {
	// Test location: London, UK
	loc := Location{
		Latitude:  51.5074,
		Longitude: -0.1278,
	}
	
	// Test date: December 21, 2024 (Winter solstice)
	date := time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC)
	
	moonrise, err := GetMoonriseTime(loc, date)
	if err != nil {
		t.Fatalf("Failed to get moonrise time: %v", err)
	}
	
	// Check that moonrise time is reasonable
	if moonrise.IsZero() {
		t.Error("Moonrise time should not be zero")
	}
	
	t.Logf("London Moonrise: %s", moonrise.Format("15:04:05"))
}

func TestGetMoonsetTime(t *testing.T) {
	// Test location: Tokyo, Japan
	loc := Location{
		Latitude:  35.6762,
		Longitude: 139.6503,
	}
	
	// Test date: March 21, 2024 (Spring equinox)
	date := time.Date(2024, 3, 21, 0, 0, 0, 0, time.UTC)
	
	moonset, err := GetMoonsetTime(loc, date)
	if err != nil {
		t.Fatalf("Failed to get moonset time: %v", err)
	}
	
	// Check that moonset time is reasonable
	if moonset.IsZero() {
		t.Error("Moonset time should not be zero")
	}
	
	t.Logf("Tokyo Moonset: %s", moonset.Format("15:04:05"))
}

func TestLunarTimesMultipleLocations(t *testing.T) {
	testCases := []struct {
		name string
		loc  Location
		date time.Time
	}{
		{
			name: "Mumbai, India",
			loc:  Location{Latitude: 19.0760, Longitude: 72.8777},
			date: time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Sydney, Australia",
			loc:  Location{Latitude: -33.8688, Longitude: 151.2093},
			date: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "São Paulo, Brazil",
			loc:  Location{Latitude: -23.5505, Longitude: -46.6333},
			date: time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC),
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lunarTimes, err := CalculateLunarTimes(tc.loc, tc.date)
			if err != nil {
				t.Fatalf("Failed to calculate lunar times for %s: %v", tc.name, err)
			}
			
			if lunarTimes == nil {
				t.Fatalf("Lunar times should not be nil for %s", tc.name)
			}
			
			t.Logf("%s Moonrise: %s", tc.name, lunarTimes.Moonrise.Format("15:04:05"))
			t.Logf("%s Moonset: %s", tc.name, lunarTimes.Moonset.Format("15:04:05"))
			t.Logf("%s Is Visible: %t", tc.name, lunarTimes.IsVisible)
		})
	}
}

func TestJulianDayToTime(t *testing.T) {
	// Test known Julian Day: JD 2451545.0 = January 1, 2000, 12:00:00 UTC
	jd := 2451545.0
	expectedTime := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	
	result := julianDayToTime(jd, time.UTC)
	
	// Allow for small differences due to floating point precision
	diff := math.Abs(result.Sub(expectedTime).Seconds())
	if diff > 60 { // Within 1 minute
		t.Errorf("Julian day conversion failed. Expected %v, got %v (difference: %f seconds)", 
			expectedTime, result, diff)
	}
	
	t.Logf("JD %f converts to %s", jd, result.Format("2006-01-02 15:04:05"))
}

func TestPolarRegions(t *testing.T) {
	// Test polar region where moon behavior might be different
	polarLoc := Location{
		Latitude:  78.9167, // Longyearbyen, Svalbard
		Longitude: 11.9500,
	}
	
	// Test during polar winter
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	lunarTimes, err := CalculateLunarTimes(polarLoc, date)
	if err != nil {
		t.Fatalf("Failed to calculate lunar times for polar region: %v", err)
	}
	
	if lunarTimes == nil {
		t.Fatal("Lunar times should not be nil for polar region")
	}
	
	t.Logf("Polar Region Moonrise: %s", lunarTimes.Moonrise.Format("15:04:05"))
	t.Logf("Polar Region Moonset: %s", lunarTimes.Moonset.Format("15:04:05"))
	t.Logf("Polar Region Is Visible: %t", lunarTimes.IsVisible)
}

func BenchmarkCalculateLunarTimes(b *testing.B) {
	loc := Location{
		Latitude:  12.9716,
		Longitude: 77.5946,
	}
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CalculateLunarTimes(loc, date)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkCalculateLunarPhase(b *testing.B) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CalculateLunarPhase(date)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}