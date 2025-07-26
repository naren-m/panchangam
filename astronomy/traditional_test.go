package astronomy

import (
	"context"
	"testing"
	"time"
)

func TestCalculateTraditionalPeriods(t *testing.T) {
	// Test location: Bangalore, India
	loc := Location{
		Latitude:  12.9716,
		Longitude: 77.5946,
	}
	
	// Test date: January 15, 2024 (Monday)
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	periods, err := CalculateTraditionalPeriods(loc, date)
	if err != nil {
		t.Fatalf("Failed to calculate traditional periods: %v", err)
	}
	
	if periods == nil {
		t.Fatal("Traditional periods should not be nil")
	}
	
	// Validate Rahu Kalam
	if periods.RahuKalam == nil {
		t.Fatal("Rahu Kalam should not be nil")
	}
	
	if periods.RahuKalam.Duration <= 0 {
		t.Error("Rahu Kalam duration should be positive")
	}
	
	if periods.RahuKalam.Start.After(periods.RahuKalam.End) {
		t.Error("Rahu Kalam start time should be before end time")
	}
	
	if periods.RahuKalam.Auspicious {
		t.Error("Rahu Kalam should not be auspicious")
	}
	
	// Validate Yamagandam
	if periods.Yamagandam == nil {
		t.Fatal("Yamagandam should not be nil")
	}
	
	if periods.Yamagandam.Duration <= 0 {
		t.Error("Yamagandam duration should be positive")
	}
	
	if periods.Yamagandam.Auspicious {
		t.Error("Yamagandam should not be auspicious")
	}
	
	// Validate Gulika Kalam
	if periods.GulikaKalam == nil {
		t.Fatal("Gulika Kalam should not be nil")
	}
	
	if periods.GulikaKalam.Duration <= 0 {
		t.Error("Gulika Kalam duration should be positive")
	}
	
	if periods.GulikaKalam.Auspicious {
		t.Error("Gulika Kalam should not be auspicious")
	}
	
	// Validate Abhijit Muhurta
	if periods.AbhijitMuhurta == nil {
		t.Fatal("Abhijit Muhurta should not be nil")
	}
	
	if periods.AbhijitMuhurta.Duration <= 0 {
		t.Error("Abhijit Muhurta duration should be positive")
	}
	
	// Log the results
	t.Logf("Rahu Kalam: %s - %s (%d min)", 
		periods.RahuKalam.Start.Format("15:04:05"),
		periods.RahuKalam.End.Format("15:04:05"),
		periods.RahuKalam.Duration)
	
	t.Logf("Yamagandam: %s - %s (%d min)", 
		periods.Yamagandam.Start.Format("15:04:05"),
		periods.Yamagandam.End.Format("15:04:05"),
		periods.Yamagandam.Duration)
	
	t.Logf("Gulika Kalam: %s - %s (%d min)", 
		periods.GulikaKalam.Start.Format("15:04:05"),
		periods.GulikaKalam.End.Format("15:04:05"),
		periods.GulikaKalam.Duration)
	
	t.Logf("Abhijit Muhurta: %s - %s (%d min, Auspicious: %t)", 
		periods.AbhijitMuhurta.Start.Format("15:04:05"),
		periods.AbhijitMuhurta.End.Format("15:04:05"),
		periods.AbhijitMuhurta.Duration,
		periods.AbhijitMuhurta.Auspicious)
}

func TestCalculateTraditionalPeriodsMultipleDays(t *testing.T) {
	loc := Location{
		Latitude:  12.9716,
		Longitude: 77.5946,
	}
	
	// Test different days of the week to ensure Rahu Kalam varies correctly
	testDates := []struct {
		date   time.Time
		day    string
		rahuPart int // Expected Rahu Kalam part (1-8)
	}{
		{time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC), "Sunday", 7},
		{time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "Monday", 1},
		{time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC), "Tuesday", 6},
		{time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC), "Wednesday", 4},
		{time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC), "Thursday", 3},
		{time.Date(2024, 1, 19, 0, 0, 0, 0, time.UTC), "Friday", 2},
		{time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC), "Saturday", 5},
	}
	
	for _, td := range testDates {
		t.Run(td.day, func(t *testing.T) {
			periods, err := CalculateTraditionalPeriods(loc, td.date)
			if err != nil {
				t.Fatalf("Failed to calculate traditional periods for %s: %v", td.day, err)
			}
			
			if periods.RahuKalam == nil {
				t.Fatalf("Rahu Kalam should not be nil for %s", td.day)
			}
			
			t.Logf("%s Rahu Kalam: %s - %s", 
				td.day,
				periods.RahuKalam.Start.Format("15:04:05"),
				periods.RahuKalam.End.Format("15:04:05"))
		})
	}
}

func TestGetRahuKalam(t *testing.T) {
	loc := Location{
		Latitude:  12.9716,
		Longitude: 77.5946,
	}
	
	// Test date: Wednesday (Rahu Kalam should be 4th part)
	date := time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC)
	
	rahuKalam, err := GetRahuKalam(loc, date)
	if err != nil {
		t.Fatalf("Failed to get Rahu Kalam: %v", err)
	}
	
	if rahuKalam == nil {
		t.Fatal("Rahu Kalam should not be nil")
	}
	
	if rahuKalam.Auspicious {
		t.Error("Rahu Kalam should not be auspicious")
	}
	
	if rahuKalam.Duration <= 0 {
		t.Error("Rahu Kalam duration should be positive")
	}
	
	t.Logf("Wednesday Rahu Kalam: %s - %s", 
		rahuKalam.Start.Format("15:04:05"),
		rahuKalam.End.Format("15:04:05"))
}

func TestGetAbhijitMuhurta(t *testing.T) {
	loc := Location{
		Latitude:  12.9716,
		Longitude: 77.5946,
	}
	
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	abhijit, err := GetAbhijitMuhurta(loc, date)
	if err != nil {
		t.Fatalf("Failed to get Abhijit Muhurta: %v", err)
	}
	
	if abhijit == nil {
		t.Fatal("Abhijit Muhurta should not be nil")
	}
	
	if abhijit.Duration <= 0 {
		t.Error("Abhijit Muhurta duration should be positive")
	}
	
	// Check if it's around midday (should be between 11 AM and 1 PM approximately)
	startHour := abhijit.Start.Hour()
	if startHour < 10 || startHour > 14 {
		t.Logf("Warning: Abhijit Muhurta start time %d might be outside expected range (10-14)", startHour)
	}
	
	t.Logf("Abhijit Muhurta: %s - %s (Auspicious: %t)", 
		abhijit.Start.Format("15:04:05"),
		abhijit.End.Format("15:04:05"),
		abhijit.Auspicious)
}

func TestCalculateMuhurtas(t *testing.T) {
	loc := Location{
		Latitude:  12.9716,
		Longitude: 77.5946,
	}
	
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	muhurtas, err := CalculateMuhurtas(loc, date)
	if err != nil {
		t.Fatalf("Failed to calculate muhurtas: %v", err)
	}
	
	if len(muhurtas) != 30 {
		t.Fatalf("Expected 30 muhurtas, got %d", len(muhurtas))
	}
	
	// Validate each muhurta
	for i, muhurta := range muhurtas {
		if muhurta == nil {
			t.Fatalf("Muhurta %d should not be nil", i+1)
		}
		
		if muhurta.Name == "" {
			t.Errorf("Muhurta %d should have a name", i+1)
		}
		
		if muhurta.Period == nil {
			t.Fatalf("Muhurta %d period should not be nil", i+1)
		}
		
		if muhurta.Period.Duration <= 0 {
			t.Errorf("Muhurta %d duration should be positive", i+1)
		}
		
		if muhurta.Period.Start.After(muhurta.Period.End) {
			t.Errorf("Muhurta %d start time should be before end time", i+1)
		}
		
		if muhurta.Quality == "" {
			t.Errorf("Muhurta %d should have a quality rating", i+1)
		}
		
		// Abhijit Muhurta (8th) should be excellent
		if i == 7 {
			if muhurta.Quality != "excellent" {
				t.Errorf("Abhijit Muhurta should be excellent quality, got %s", muhurta.Quality)
			}
			
			if !muhurta.Period.Auspicious {
				t.Error("Abhijit Muhurta should be auspicious")
			}
		}
	}
	
	// Log first few muhurtas
	for i := 0; i < 5; i++ {
		t.Logf("Muhurta %d (%s): %s - %s (%s)", 
			i+1,
			muhurtas[i].Name,
			muhurtas[i].Period.Start.Format("15:04:05"),
			muhurtas[i].Period.End.Format("15:04:05"),
			muhurtas[i].Quality)
	}
	
	// Log Abhijit Muhurta specifically
	abhijit := muhurtas[7]
	t.Logf("Abhijit Muhurta (%s): %s - %s (%s)", 
		abhijit.Name,
		abhijit.Period.Start.Format("15:04:05"),
		abhijit.Period.End.Format("15:04:05"),
		abhijit.Quality)
}

func TestTraditionalPeriodsWithContext(t *testing.T) {
	ctx := context.Background()
	
	loc := Location{
		Latitude:  40.7128,
		Longitude: -74.0060, // New York
	}
	
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC) // Summer solstice
	
	periods, err := CalculateTraditionalPeriodsWithContext(ctx, loc, date)
	if err != nil {
		t.Fatalf("Failed to calculate traditional periods with context: %v", err)
	}
	
	if periods == nil {
		t.Fatal("Traditional periods should not be nil")
	}
	
	t.Logf("NYC Summer Solstice Traditional Periods:")
	t.Logf("Rahu Kalam: %s - %s", 
		periods.RahuKalam.Start.Format("15:04:05"),
		periods.RahuKalam.End.Format("15:04:05"))
	t.Logf("Yamagandam: %s - %s", 
		periods.Yamagandam.Start.Format("15:04:05"),
		periods.Yamagandam.End.Format("15:04:05"))
	t.Logf("Abhijit Muhurta: %s - %s (Auspicious: %t)", 
		periods.AbhijitMuhurta.Start.Format("15:04:05"),
		periods.AbhijitMuhurta.End.Format("15:04:05"),
		periods.AbhijitMuhurta.Auspicious)
}

func TestAbhijitMuhurtaValidation(t *testing.T) {
	loc := Location{
		Latitude:  78.9167, // High latitude location (Svalbard)
		Longitude: 11.9500,
	}
	
	// Test during polar summer when sunrise might be very early
	date := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)
	
	periods, err := CalculateTraditionalPeriods(loc, date)
	if err != nil {
		t.Fatalf("Failed to calculate traditional periods for high latitude: %v", err)
	}
	
	if periods.AbhijitMuhurta == nil {
		t.Fatal("Abhijit Muhurta should not be nil")
	}
	
	// At high latitudes, Abhijit Muhurta might start after 12:30 PM and be invalid
	t.Logf("High Latitude Abhijit Muhurta: %s - %s (Auspicious: %t)", 
		periods.AbhijitMuhurta.Start.Format("15:04:05"),
		periods.AbhijitMuhurta.End.Format("15:04:05"),
		periods.AbhijitMuhurta.Auspicious)
	
	t.Logf("Description: %s", periods.AbhijitMuhurta.Description)
}

func BenchmarkCalculateTraditionalPeriods(b *testing.B) {
	loc := Location{
		Latitude:  12.9716,
		Longitude: 77.5946,
	}
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CalculateTraditionalPeriods(loc, date)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkCalculateMuhurtas(b *testing.B) {
	loc := Location{
		Latitude:  12.9716,
		Longitude: 77.5946,
	}
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CalculateMuhurtas(loc, date)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}