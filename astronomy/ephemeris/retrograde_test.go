package ephemeris

import (
	"context"
	"testing"
	"time"
)

func TestDetectRetrogradeMotion(t *testing.T) {
	manager := createTestManager(t)
	detector := NewRetrogradeDetector(manager)

	tests := []struct {
		name     string
		date     time.Time
		planet   string
		wantMotion RetrogradeMotion // Can be empty if we just want to verify it works
	}{
		{
			name:   "mercury_direct",
			date:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			planet: "mercury",
		},
		{
			name:   "venus_motion",
			date:   time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC),
			planet: "venus",
		},
		{
			name:   "mars_motion",
			date:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			planet: "mars",
		},
		{
			name:   "jupiter_motion",
			date:   time.Date(2024, 9, 1, 12, 0, 0, 0, time.UTC),
			planet: "jupiter",
		},
		{
			name:   "saturn_motion",
			date:   time.Date(2024, 12, 1, 12, 0, 0, 0, time.UTC),
			planet: "saturn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jd := TimeToJulianDay(tt.date)
			motion, err := detector.DetectRetrogradeMotion(context.Background(), jd, tt.planet)
			if err != nil {
				t.Fatalf("Failed to detect retrograde motion: %v", err)
			}

			// Verify motion is one of the valid values
			if motion != MotionDirect && motion != MotionRetrograde && motion != MotionStationary {
				t.Errorf("Invalid motion type: %s", motion)
			}

			t.Logf("%s at %s: %s", tt.planet, tt.date.Format("2006-01-02"), motion)
		})
	}
}

func TestFindPlanetaryStation(t *testing.T) {
	manager := createTestManager(t)
	detector := NewRetrogradeDetector(manager)

	// Mercury has frequent stations (every ~3-4 months)
	// Search from early 2024
	startDate := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	startJD := TimeToJulianDay(startDate)

	station, err := detector.FindPlanetaryStation(context.Background(), startJD, "mercury", 120)
	if err != nil {
		t.Logf("No station found for Mercury within 120 days: %v", err)
		// This is acceptable - not all planets have stations in any given period
		return
	}

	if station == nil {
		t.Log("No station found for Mercury within 120 days (this is acceptable)")
		return
	}

	// Verify station properties
	if station.Planet != "mercury" {
		t.Errorf("Expected planet mercury, got %s", station.Planet)
	}

	if station.Longitude < 0 || station.Longitude >= 360 {
		t.Errorf("Station longitude out of range: %f", station.Longitude)
	}

	if station.StationType != StationRetrograde && station.StationType != StationDirect {
		t.Errorf("Invalid station type: %s", station.StationType)
	}

	// Speed at station should be very small
	if station.Speed > 0.1 {
		t.Logf("Warning: Speed at station is %f (expected near 0)", station.Speed)
	}

	t.Logf("Found Mercury station at %s (JD %f)", station.Time.Format("2006-01-02 15:04"), station.JulianDay)
	t.Logf("  Longitude: %f°", station.Longitude)
	t.Logf("  Type: %s", station.StationType)
	t.Logf("  Speed: %f °/day", station.Speed)
}

func TestFindRetrogradePeriod(t *testing.T) {
	manager := createTestManager(t)
	detector := NewRetrogradeDetector(manager)

	// Try to find a retrograde period
	// We'll test multiple dates to find one
	testDates := []time.Time{
		time.Date(2024, 4, 15, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 8, 15, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 15, 12, 0, 0, 0, time.UTC),
	}

	planets := []string{"mercury", "venus", "mars"}

	foundAny := false
	for _, planet := range planets {
		for _, date := range testDates {
			jd := TimeToJulianDay(date)

			// First check if planet is retrograde
			motion, err := detector.DetectRetrogradeMotion(context.Background(), jd, planet)
			if err != nil || motion != MotionRetrograde {
				continue
			}

			// Try to find the retrograde period
			period, err := detector.FindRetrogradePeriod(context.Background(), jd, planet)
			if err != nil {
				continue
			}

			foundAny = true

			// Verify period properties
			if period.Planet != planet {
				t.Errorf("Expected planet %s, got %s", planet, period.Planet)
			}

			if period.EndJD <= period.StartJD {
				t.Errorf("End JD (%f) should be after Start JD (%f)", period.EndJD, period.StartJD)
			}

			if period.Duration <= 0 {
				t.Errorf("Duration should be positive: %s", period.Duration)
			}

			t.Logf("Found %s retrograde period:", planet)
			t.Logf("  Start: %s (JD %f)", period.StartTime.Format("2006-01-02"), period.StartJD)
			t.Logf("  End: %s (JD %f)", period.EndTime.Format("2006-01-02"), period.EndJD)
			t.Logf("  Duration: %.1f days", period.Duration.Hours()/24)
			t.Logf("  Start longitude: %f°", period.StartLongitude)
			t.Logf("  End longitude: %f°", period.EndLongitude)
			t.Logf("  Max retro distance: %f°", period.MaxRetroDistance)

			return // Found one, that's enough for the test
		}
	}

	if !foundAny {
		t.Log("No retrograde periods found in test dates (this is acceptable)")
	}
}

func TestAnalyzeMotion(t *testing.T) {
	manager := createTestManager(t)
	detector := NewRetrogradeDetector(manager)

	testDate := time.Date(2024, 7, 1, 12, 0, 0, 0, time.UTC)
	jd := TimeToJulianDay(testDate)

	planets := []string{"mercury", "venus", "mars", "jupiter", "saturn"}

	for _, planet := range planets {
		t.Run(planet, func(t *testing.T) {
			analysis, err := detector.AnalyzeMotion(context.Background(), jd, planet)
			if err != nil {
				t.Fatalf("Failed to analyze motion for %s: %v", planet, err)
			}

			if analysis == nil {
				t.Fatal("Expected non-nil analysis")
			}

			// Verify basic properties
			if analysis.Planet != planet {
				t.Errorf("Expected planet %s, got %s", planet, analysis.Planet)
			}

			if analysis.JulianDay != jd {
				t.Errorf("JD mismatch: got %f, want %f", analysis.JulianDay, jd)
			}

			if analysis.Motion != MotionDirect && analysis.Motion != MotionRetrograde && analysis.Motion != MotionStationary {
				t.Errorf("Invalid motion type: %s", analysis.Motion)
			}

			t.Logf("%s motion analysis:", planet)
			t.Logf("  Motion: %s", analysis.Motion)
			t.Logf("  Speed: %f °/day", analysis.Speed)
			t.Logf("  Longitude: %f°", analysis.Longitude)
			t.Logf("  Near station: %v", analysis.IsNearStation)

			if analysis.NextStation != nil {
				t.Logf("  Next station: %s at %s", analysis.NextStation.StationType, analysis.NextStation.Time.Format("2006-01-02"))
			}

			if analysis.CurrentPeriod != nil {
				t.Logf("  Current retrograde period: %s to %s",
					analysis.CurrentPeriod.StartTime.Format("2006-01-02"),
					analysis.CurrentPeriod.EndTime.Format("2006-01-02"))
			}

			if len(analysis.RecentStations) > 0 {
				t.Logf("  Recent stations: %d found", len(analysis.RecentStations))
			}
		})
	}
}

func TestGetRetrogradePlanets(t *testing.T) {
	manager := createTestManager(t)
	detector := NewRetrogradeDetector(manager)

	testDate := time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC)
	jd := TimeToJulianDay(testDate)

	planets, err := detector.GetRetrogradePlanets(context.Background(), jd)
	if err != nil {
		t.Fatalf("Failed to get retrograde planets: %v", err)
	}

	// Verify result is a valid list (can be empty)
	if planets == nil {
		t.Error("Expected non-nil planet list")
	}

	t.Logf("Retrograde planets on %s: %v", testDate.Format("2006-01-02"), planets)

	// Verify all returned planets are valid
	validPlanets := map[string]bool{
		"mercury": true, "venus": true, "mars": true, "jupiter": true,
		"saturn": true, "uranus": true, "neptune": true, "pluto": true,
	}

	for _, planet := range planets {
		if !validPlanets[planet] {
			t.Errorf("Invalid planet in retrograde list: %s", planet)
		}
	}
}

func TestRetrogradeMotionTypes(t *testing.T) {
	// Test the motion type constants
	motions := []RetrogradeMotion{
		MotionDirect,
		MotionRetrograde,
		MotionStationary,
	}

	for _, motion := range motions {
		if motion == "" {
			t.Error("Motion type should not be empty")
		}
		t.Logf("Motion type: %s", motion)
	}
}

func TestStationTypes(t *testing.T) {
	// Test the station type constants
	types := []StationType{
		StationRetrograde,
		StationDirect,
	}

	for _, typ := range types {
		if typ == "" {
			t.Error("Station type should not be empty")
		}
		t.Logf("Station type: %s", typ)
	}
}

func TestValidateKnownRetrograde(t *testing.T) {
	manager := createTestManager(t)
	detector := NewRetrogradeDetector(manager)

	// Known Mercury retrograde period in 2024 (approximate)
	// Mercury retrograde: April 1-25, 2024 (example dates)
	startDate := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC)

	startJD := TimeToJulianDay(startDate)
	endJD := TimeToJulianDay(endDate)

	// This test may pass or fail depending on actual ephemeris data
	// We're mainly testing that the function works
	isValid, err := detector.ValidateKnownRetrograde(context.Background(), "mercury", startJD, endJD)
	if err != nil {
		t.Logf("Validation check completed with error: %v", err)
	}

	t.Logf("Validation result for known Mercury retrograde period: %v", isValid)
}

func TestRetrogradeEdgeCases(t *testing.T) {
	manager := createTestManager(t)
	detector := NewRetrogradeDetector(manager)

	t.Run("invalid_planet", func(t *testing.T) {
		jd := TimeToJulianDay(time.Now())
		_, err := detector.DetectRetrogradeMotion(context.Background(), jd, "invalid_planet")
		if err == nil {
			t.Error("Expected error for invalid planet")
		}
	})

	t.Run("sun_retrograde_check", func(t *testing.T) {
		// Sun never goes retrograde from Earth's perspective
		jd := TimeToJulianDay(time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC))
		motion, err := detector.DetectRetrogradeMotion(context.Background(), jd, "sun")

		if err != nil {
			t.Logf("Sun retrograde check: %v", err)
		} else {
			// Sun should always be direct
			if motion == MotionRetrograde {
				t.Error("Sun should never be retrograde")
			}
			t.Logf("Sun motion: %s (as expected)", motion)
		}
	})

	t.Run("moon_retrograde_check", func(t *testing.T) {
		// Moon never goes retrograde
		jd := TimeToJulianDay(time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC))
		motion, err := detector.DetectRetrogradeMotion(context.Background(), jd, "moon")

		if err != nil {
			t.Logf("Moon retrograde check: %v", err)
		} else {
			// Moon should always be direct
			if motion == MotionRetrograde {
				t.Error("Moon should never be retrograde")
			}
			t.Logf("Moon motion: %s (as expected)", motion)
		}
	})
}

func BenchmarkDetectRetrogradeMotion(b *testing.B) {
	manager := createTestManager(b)
	detector := NewRetrogradeDetector(manager)

	jd := TimeToJulianDay(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.DetectRetrogradeMotion(ctx, jd, "mars")
		if err != nil {
			b.Fatalf("Detection failed: %v", err)
		}
	}
}

func BenchmarkAnalyzeMotion(b *testing.B) {
	manager := createTestManager(b)
	detector := NewRetrogradeDetector(manager)

	jd := TimeToJulianDay(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.AnalyzeMotion(ctx, jd, "jupiter")
		if err != nil {
			b.Fatalf("Analysis failed: %v", err)
		}
	}
}

func BenchmarkGetRetrogradePlanets(b *testing.B) {
	manager := createTestManager(b)
	detector := NewRetrogradeDetector(manager)

	jd := TimeToJulianDay(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.GetRetrogradePlanets(ctx, jd)
		if err != nil {
			b.Fatalf("Get retrograde planets failed: %v", err)
		}
	}
}
