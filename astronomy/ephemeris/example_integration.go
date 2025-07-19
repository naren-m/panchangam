package ephemeris

import (
	"context"
	"fmt"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

// ExampleIntegration demonstrates how to use the ephemeris system
func ExampleIntegration() {
	// Initialize observability
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "ephemeris.ExampleIntegration")
	defer span.End()

	// Create providers
	jplProvider := NewJPLProvider()
	swissProvider := NewSwissProvider()
	
	// Create cache with 1000 entries and 1 hour TTL
	cache := NewMemoryCache(1000, 1*time.Hour)
	
	// Create manager with JPL as primary, Swiss as fallback
	manager := NewManager(jplProvider, swissProvider, cache)
	defer manager.Close()

	// Test date: July 18, 2024
	testDate := time.Date(2024, 7, 18, 12, 0, 0, 0, time.UTC)
	jd := TimeToJulianDay(testDate)
	
	span.SetAttributes(
		attribute.String("test_date", testDate.Format("2006-01-02")),
		attribute.Float64("julian_day", float64(jd)),
	)

	fmt.Printf("=== Ephemeris Integration Example ===\n")
	fmt.Printf("Date: %s\n", testDate.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("Julian Day: %.6f\n\n", jd)

	// Get Sun position
	fmt.Println("--- Sun Position ---")
	sunPos, err := manager.GetSunPosition(ctx, jd)
	if err != nil {
		fmt.Printf("Error getting sun position: %v\n", err)
		return
	}
	
	fmt.Printf("Longitude: %.6f°\n", sunPos.Longitude)
	fmt.Printf("Right Ascension: %.6f°\n", sunPos.RightAscension)
	fmt.Printf("Declination: %.6f°\n", sunPos.Declination)
	fmt.Printf("Distance: %.6f AU\n", sunPos.Distance)
	fmt.Printf("Equation of Time: %.2f minutes\n\n", sunPos.EquationOfTime)

	// Get Moon position
	fmt.Println("--- Moon Position ---")
	moonPos, err := manager.GetMoonPosition(ctx, jd)
	if err != nil {
		fmt.Printf("Error getting moon position: %v\n", err)
		return
	}
	
	fmt.Printf("Longitude: %.6f°\n", moonPos.Longitude)
	fmt.Printf("Latitude: %.6f°\n", moonPos.Latitude)
	fmt.Printf("Distance: %.2f km\n", moonPos.Distance)
	fmt.Printf("Phase: %.3f (%.1f%% illuminated)\n", moonPos.Phase, moonPos.Illumination)
	fmt.Printf("Angular Diameter: %.2f arcseconds\n\n", moonPos.AngularDiameter)

	// Get all planetary positions
	fmt.Println("--- Planetary Positions ---")
	positions, err := manager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		fmt.Printf("Error getting planetary positions: %v\n", err)
		return
	}

	planets := map[string]Position{
		"Mercury": positions.Mercury,
		"Venus":   positions.Venus,
		"Mars":    positions.Mars,
		"Jupiter": positions.Jupiter,
		"Saturn":  positions.Saturn,
		"Uranus":  positions.Uranus,
		"Neptune": positions.Neptune,
		"Pluto":   positions.Pluto,
	}

	for name, pos := range planets {
		fmt.Printf("%s: Long=%.2f°, Lat=%.2f°, Dist=%.3f AU, Speed=%.3f°/day\n",
			name, pos.Longitude, pos.Latitude, pos.Distance, pos.Speed)
	}

	// Demonstrate caching
	fmt.Println("\n--- Cache Performance ---")
	start := time.Now()
	_, _ = manager.GetSunPosition(ctx, jd) // Should hit cache
	cacheTime := time.Since(start)
	
	stats := cache.GetStats(ctx)
	fmt.Printf("Cache hits: %d\n", stats.Hits)
	fmt.Printf("Cache misses: %d\n", stats.Misses)
	fmt.Printf("Hit rate: %.2f%%\n", stats.HitRate*100)
	fmt.Printf("Cache lookup time: %v\n", cacheTime)

	// Health status
	fmt.Println("\n--- Provider Health ---")
	healthStatuses, err := manager.GetHealthStatus(ctx)
	if err != nil {
		fmt.Printf("Error getting health status: %v\n", err)
		return
	}

	for name, status := range healthStatuses {
		fmt.Printf("%s: Available=%t, Response=%v, Version=%s\n",
			name, status.Available, status.ResponseTime, status.Version)
	}

	span.AddEvent("Integration example completed successfully")
}

// ExamplePanchangamCalculation demonstrates using ephemeris for Panchangam calculations
func ExamplePanchangamCalculation() {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "ephemeris.ExamplePanchangamCalculation")
	defer span.End()

	// Create ephemeris system
	jplProvider := NewJPLProvider()
	swissProvider := NewSwissProvider()
	cache := NewMemoryCache(1000, 1*time.Hour)
	manager := NewManager(jplProvider, swissProvider, cache)
	defer manager.Close()

	// Test date for Panchangam calculation
	testDate := time.Date(2024, 7, 18, 6, 0, 0, 0, time.UTC) // Sunrise time
	jd := TimeToJulianDay(testDate)

	fmt.Printf("=== Panchangam Calculation Example ===\n")
	fmt.Printf("Date: %s\n", testDate.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("Julian Day: %.6f\n\n", jd)

	// Get Sun and Moon positions for Panchangam elements
	sunPos, err := manager.GetSunPosition(ctx, jd)
	if err != nil {
		fmt.Printf("Error getting sun position: %v\n", err)
		return
	}

	moonPos, err := manager.GetMoonPosition(ctx, jd)
	if err != nil {
		fmt.Printf("Error getting moon position: %v\n", err)
		return
	}

	// Calculate Panchangam elements
	fmt.Println("--- Panchangam Elements ---")
	
	// Tithi calculation (lunar day)
	sunLongitude := sunPos.Longitude
	moonLongitude := moonPos.Longitude
	
	// Tithi = (Moon longitude - Sun longitude) / 12
	tithiDegrees := moonLongitude - sunLongitude
	if tithiDegrees < 0 {
		tithiDegrees += 360
	}
	tithi := int(tithiDegrees/12) + 1
	
	fmt.Printf("Sun Longitude: %.6f°\n", sunLongitude)
	fmt.Printf("Moon Longitude: %.6f°\n", moonLongitude)
	fmt.Printf("Tithi: %d (%.2f°)\n", tithi, tithiDegrees)

	// Nakshatra calculation (lunar mansion)
	// Each nakshatra is 13°20' (13.333...)
	nakshatra := int(moonLongitude/13.333333) + 1
	fmt.Printf("Nakshatra: %d (Moon at %.6f°)\n", nakshatra, moonLongitude)

	// Yoga calculation
	// Yoga = (Sun longitude + Moon longitude) / 13.333...
	yogaDegrees := sunLongitude + moonLongitude
	if yogaDegrees >= 360 {
		yogaDegrees -= 360
	}
	yoga := int(yogaDegrees/13.333333) + 1
	fmt.Printf("Yoga: %d (%.2f°)\n", yoga, yogaDegrees)

	// Karana calculation
	// Karana = Tithi / 2 (simplified)
	karana := ((tithi - 1) * 2) + 1
	if karana > 60 {
		karana -= 60
	}
	fmt.Printf("Karana: %d\n", karana)

	// Vara (weekday) - already known from date
	weekday := testDate.Weekday()
	fmt.Printf("Vara: %s\n", weekday)

	// Additional astronomical data
	fmt.Println("\n--- Additional Astronomical Data ---")
	fmt.Printf("Moon Phase: %.3f (%.1f%% illuminated)\n", moonPos.Phase, moonPos.Illumination)
	fmt.Printf("Moon Distance: %.2f km\n", moonPos.Distance)
	fmt.Printf("Sun Distance: %.6f AU\n", sunPos.Distance)
	fmt.Printf("Equation of Time: %.2f minutes\n", sunPos.EquationOfTime)

	span.SetAttributes(
		attribute.Int("tithi", tithi),
		attribute.Int("nakshatra", nakshatra),
		attribute.Int("yoga", yoga),
		attribute.Int("karana", karana),
		attribute.String("vara", weekday.String()),
	)

	span.AddEvent("Panchangam calculation completed")
}

// ExamplePerformanceTest demonstrates performance characteristics
func ExamplePerformanceTest() {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "ephemeris.ExamplePerformanceTest")
	defer span.End()

	// Create ephemeris system
	jplProvider := NewJPLProvider()
	swissProvider := NewSwissProvider()
	cache := NewMemoryCache(1000, 1*time.Hour)
	manager := NewManager(jplProvider, swissProvider, cache)
	defer manager.Close()

	testJD := JulianDay(2451545.0) // J2000.0
	iterations := 1000

	fmt.Printf("=== Performance Test ===\n")
	fmt.Printf("Iterations: %d\n", iterations)
	fmt.Printf("Julian Day: %.6f\n\n", testJD)

	// Test without cache (first call)
	fmt.Println("--- Without Cache ---")
	start := time.Now()
	_, err := manager.GetSunPosition(ctx, testJD)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	firstCallTime := time.Since(start)
	fmt.Printf("First call (no cache): %v\n", firstCallTime)

	// Test with cache
	fmt.Println("\n--- With Cache ---")
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_, err := manager.GetSunPosition(ctx, testJD)
		if err != nil {
			fmt.Printf("Error on iteration %d: %v\n", i, err)
			return
		}
	}
	totalTime := time.Since(start)
	avgTime := totalTime / time.Duration(iterations)
	
	fmt.Printf("Total time for %d cached calls: %v\n", iterations, totalTime)
	fmt.Printf("Average time per cached call: %v\n", avgTime)
	fmt.Printf("Cache speedup: %.2fx\n", float64(firstCallTime)/float64(avgTime))

	// Cache statistics
	stats := cache.GetStats(ctx)
	fmt.Printf("\n--- Cache Statistics ---\n")
	fmt.Printf("Hits: %d\n", stats.Hits)
	fmt.Printf("Misses: %d\n", stats.Misses)
	fmt.Printf("Hit rate: %.2f%%\n", stats.HitRate*100)
	fmt.Printf("Average latency: %v\n", stats.AverageLatency)

	span.SetAttributes(
		attribute.Int("iterations", iterations),
		attribute.Int64("first_call_ms", firstCallTime.Milliseconds()),
		attribute.Int64("avg_cached_call_ns", avgTime.Nanoseconds()),
		attribute.Float64("cache_speedup", float64(firstCallTime)/float64(avgTime)),
	)

	span.AddEvent("Performance test completed")
}

// ExampleHealthMonitoring demonstrates health monitoring
func ExampleHealthMonitoring() {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "ephemeris.ExampleHealthMonitoring")
	defer span.End()

	// Create providers
	jplProvider := NewJPLProvider()
	swissProvider := NewSwissProvider()
	cache := NewMemoryCache(1000, 1*time.Hour)
	manager := NewManager(jplProvider, swissProvider, cache)
	defer manager.Close()

	fmt.Printf("=== Health Monitoring Example ===\n")

	// Get initial health status
	healthStatuses, err := manager.GetHealthStatus(ctx)
	if err != nil {
		fmt.Printf("Error getting health status: %v\n", err)
		return
	}

	fmt.Println("--- Initial Health Status ---")
	for name, status := range healthStatuses {
		fmt.Printf("%s:\n", name)
		fmt.Printf("  Available: %t\n", status.Available)
		fmt.Printf("  Version: %s\n", status.Version)
		fmt.Printf("  Response Time: %v\n", status.ResponseTime)
		fmt.Printf("  Data Range: %.1f to %.1f JD\n", status.DataStartJD, status.DataEndJD)
		fmt.Printf("  Last Check: %s\n", status.LastCheck.Format("2006-01-02 15:04:05"))
		if status.ErrorMessage != "" {
			fmt.Printf("  Error: %s\n", status.ErrorMessage)
		}
		fmt.Println()
	}

	// Start health checker
	if manager.healthChecker != nil {
		manager.healthChecker.Start()
		defer manager.healthChecker.Stop()

		// Wait for a health check cycle
		time.Sleep(2 * time.Second)

		fmt.Println("--- Health Checker Metrics ---")
		metrics := manager.healthChecker.GetMetrics()
		for key, value := range metrics {
			fmt.Printf("%s: %v\n", key, value)
		}

		fmt.Println("\n--- Individual Provider Status ---")
		allStatuses := manager.healthChecker.GetAllStatuses()
		for name, status := range allStatuses {
			fmt.Printf("%s: Available=%t, ResponseTime=%v\n", 
				name, status.Available, status.ResponseTime)
		}
	}

	span.AddEvent("Health monitoring example completed")
}

// ExampleErrorHandling demonstrates error handling and fallback
func ExampleErrorHandling() {
	observer := observability.Observer()
	ctx, span := observer.CreateSpan(context.Background(), "ephemeris.ExampleErrorHandling")
	defer span.End()

	fmt.Printf("=== Error Handling Example ===\n")

	// Create providers
	jplProvider := NewJPLProvider()
	swissProvider := NewSwissProvider()
	cache := NewMemoryCache(1000, 1*time.Hour)

	// Test with nil primary provider to trigger fallback
	fmt.Println("--- Testing Fallback Mechanism ---")
	managerWithNilPrimary := NewManager(nil, swissProvider, cache)
	
	testJD := JulianDay(2451545.0)
	position, err := managerWithNilPrimary.GetSunPosition(ctx, testJD)
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
	} else {
		fmt.Printf("Successfully got position from fallback: Long=%.2f°\n", position.Longitude)
	}

	// Test with invalid Julian day
	fmt.Println("\n--- Testing Invalid Julian Day ---")
	manager := NewManager(jplProvider, swissProvider, cache)
	defer manager.Close()

	invalidJD := JulianDay(1000000.0) // Way outside valid range
	_, err = manager.GetSunPosition(ctx, invalidJD)
	if err != nil {
		fmt.Printf("Error for invalid JD (expected): %v\n", err)
	} else {
		fmt.Printf("Unexpected success for invalid JD\n")
	}

	// Test with valid JD
	fmt.Println("\n--- Testing Valid Operation ---")
	validJD := JulianDay(2451545.0)
	position, err = manager.GetSunPosition(ctx, validJD)
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
	} else {
		fmt.Printf("Successfully got position: Long=%.2f°\n", position.Longitude)
	}

	span.AddEvent("Error handling example completed")
}