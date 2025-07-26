package astronomy

import (
	"context"
	"testing"
	"time"
)

func TestNewFestivalCalendar(t *testing.T) {
	fc := NewFestivalCalendar()
	
	if fc == nil {
		t.Fatal("Festival calendar should not be nil")
	}
	
	if len(fc.fixedFestivals) == 0 {
		t.Error("Fixed festivals should be initialized")
	}
	
	if len(fc.lunarFestivals) == 0 {
		t.Error("Lunar festivals should be initialized")
	}
}

func TestGetFestivalsForDate(t *testing.T) {
	fc := NewFestivalCalendar()
	ctx := context.Background()
	
	testCases := []struct {
		name         string
		date         time.Time
		tithiNumber  int
		expectCount  int
		expectNames  []string
	}{
		{
			name:        "Republic Day",
			date:        time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC),
			tithiNumber: 15,
			expectCount: 2, // Republic Day + Purnima
			expectNames: []string{"Republic Day", "Pausha Purnima"},
		},
		{
			name:        "Independence Day",
			date:        time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
			tithiNumber: 10,
			expectCount: 1,
			expectNames: []string{"Independence Day"},
		},
		{
			name:        "Ekadashi",
			date:        time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			tithiNumber: 11,
			expectCount: 1,
			expectNames: []string{"Vaishakha Mohini Ekadashi"},
		},
		{
			name:        "Purnima (Full Moon)",
			date:        time.Date(2024, 8, 19, 0, 0, 0, 0, time.UTC),
			tithiNumber: 15,
			expectCount: 1,
			expectNames: []string{"Raksha Bandhan"},
		},
		{
			name:        "Amavasya (New Moon)",
			date:        time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC),
			tithiNumber: 30,
			expectCount: 1,
			expectNames: []string{"Diwali Amavasya"},
		},
		{
			name:        "Regular day with no festivals",
			date:        time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC),
			tithiNumber: 5,
			expectCount: 0,
			expectNames: []string{},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			festivals, err := fc.GetFestivalsForDate(ctx, tc.date, tc.tithiNumber)
			if err != nil {
				t.Fatalf("Failed to get festivals for date: %v", err)
			}
			
			if len(festivals) != tc.expectCount {
				t.Errorf("Expected %d festivals, got %d", tc.expectCount, len(festivals))
				for i, f := range festivals {
					t.Logf("Festival %d: %s", i, f.Name)
				}
			}
			
			// Check if expected festival names are present
			for _, expectedName := range tc.expectNames {
				found := false
				for _, festival := range festivals {
					if festival.Name == expectedName {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected festival '%s' not found", expectedName)
				}
			}
		})
	}
}

func TestGetSeasonalFestivals(t *testing.T) {
	fc := NewFestivalCalendar()
	
	testCases := []struct {
		name         string
		date         time.Time
		expectName   string
		shouldFind   bool
	}{
		{
			name:       "Spring Equinox",
			date:       time.Date(2024, 3, 21, 0, 0, 0, 0, time.UTC),
			expectName: "Spring Equinox",
			shouldFind: true,
		},
		{
			name:       "Summer Solstice",
			date:       time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
			expectName: "Summer Solstice",
			shouldFind: true,
		},
		{
			name:       "Autumn Equinox",
			date:       time.Date(2024, 9, 23, 0, 0, 0, 0, time.UTC),
			expectName: "Autumn Equinox",
			shouldFind: true,
		},
		{
			name:       "Winter Solstice",
			date:       time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
			expectName: "Winter Solstice",
			shouldFind: true,
		},
		{
			name:       "Regular day",
			date:       time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			expectName: "",
			shouldFind: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			festivals := fc.getSeasonalFestivals(tc.date)
			
			if tc.shouldFind {
				if len(festivals) == 0 {
					t.Errorf("Expected to find seasonal festival, but got none")
					return
				}
				
				found := false
				for _, festival := range festivals {
					if festival.Name == tc.expectName {
						found = true
						
						// Validate festival details
						if festival.Type != "seasonal" {
							t.Errorf("Expected seasonal festival type, got %s", festival.Type)
						}
						
						if festival.Significance == "" {
							t.Error("Festival should have significance")
						}
						
						if len(festival.Observances) == 0 {
							t.Error("Festival should have observances")
						}
						
						break
					}
				}
				
				if !found {
					t.Errorf("Expected to find festival '%s'", tc.expectName)
				}
			} else {
				if len(festivals) > 0 {
					t.Errorf("Expected no seasonal festivals, but found %d", len(festivals))
				}
			}
		})
	}
}

func TestGetMonthSpecificName(t *testing.T) {
	fc := NewFestivalCalendar()
	
	testCases := []struct {
		baseName    string
		month       time.Month
		expected    string
	}{
		{"Ekadashi", time.May, "Vaishakha Mohini Ekadashi"},
		{"Ekadashi", time.July, "Ashadha Yogini Ekadashi"},
		{"Purnima", time.August, "Raksha Bandhan"},
		{"Purnima", time.May, "Buddha Purnima"},
		{"Amavasya", time.October, "Diwali Amavasya"},
		{"Amavasya", time.January, "Amavasya"}, // No special name
		{"Unknown", time.January, "Unknown"}, // Unknown festival
	}
	
	for _, tc := range testCases {
		date := time.Date(2024, tc.month, 15, 0, 0, 0, 0, time.UTC)
		result := fc.getMonthSpecificName(tc.baseName, date, 15)
		
		if result != tc.expected {
			t.Errorf("Expected '%s' for %s in %s, got '%s'", 
				tc.expected, tc.baseName, tc.month, result)
		}
	}
}

func TestGetUpcomingFestivals(t *testing.T) {
	fc := NewFestivalCalendar()
	ctx := context.Background()
	
	// Test upcoming festivals for a week
	startDate := time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC) // Day before Republic Day
	festivals, err := fc.GetUpcomingFestivals(ctx, startDate, 7)
	
	if err != nil {
		t.Fatalf("Failed to get upcoming festivals: %v", err)
	}
	
	// Should find at least Republic Day
	found := false
	for _, festival := range festivals {
		if festival.Name == "Republic Day" {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Should find Republic Day in upcoming festivals")
	}
	
	t.Logf("Found %d upcoming festivals in 7 days", len(festivals))
	for _, festival := range festivals {
		t.Logf("  %s on %s (%s)", 
			festival.Name, 
			festival.Date.Format("2006-01-02"), 
			festival.Type)
	}
}

func TestGetFestivalNamesForDate(t *testing.T) {
	ctx := context.Background()
	
	// Test Republic Day
	date := time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC)
	names, err := GetFestivalNamesForDate(ctx, date, 15)
	
	if err != nil {
		t.Fatalf("Failed to get festival names: %v", err)
	}
	
	if len(names) == 0 {
		t.Error("Should find festivals on Republic Day")
	}
	
	found := false
	for _, name := range names {
		if name == "Republic Day" {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Should find Republic Day")
	}
	
	t.Logf("Festivals on %s: %v", date.Format("2006-01-02"), names)
}

func TestFestivalTypes(t *testing.T) {
	fc := NewFestivalCalendar()
	ctx := context.Background()
	
	// Test that all festivals have valid types
	validTypes := map[string]bool{
		"major":    true,
		"minor":    true,
		"national": true,
		"regional": true,
		"seasonal": true,
	}
	
	// Test a few dates
	dates := []time.Time{
		time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), // Republic Day
		time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC), // Ekadashi
		time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC), // Summer Solstice
	}
	
	for _, date := range dates {
		festivals, err := fc.GetFestivalsForDate(ctx, date, 11)
		if err != nil {
			t.Fatalf("Failed to get festivals: %v", err)
		}
		
		for _, festival := range festivals {
			if !validTypes[festival.Type] {
				t.Errorf("Invalid festival type '%s' for festival '%s'", 
					festival.Type, festival.Name)
			}
			
			if festival.Name == "" {
				t.Error("Festival name should not be empty")
			}
			
			if festival.Significance == "" {
				t.Errorf("Festival '%s' should have significance", festival.Name)
			}
		}
	}
}

func BenchmarkGetFestivalsForDate(b *testing.B) {
	fc := NewFestivalCalendar()
	ctx := context.Background()
	date := time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := fc.GetFestivalsForDate(ctx, date, 15)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkGetUpcomingFestivals(b *testing.B) {
	fc := NewFestivalCalendar()
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := fc.GetUpcomingFestivals(ctx, startDate, 30)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}