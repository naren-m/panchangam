package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/naren-m/panchangam/observability"
)

func main() {
	// Initialize observability
	observability.NewLocalObserver()
	
	// Create ephemeris providers
	jplProvider := ephemeris.NewJPLProvider()
	swissProvider := ephemeris.NewSwissProvider()
	
	// Create cache and manager
	cache := ephemeris.NewMemoryCache(1000, 1*time.Hour)
	manager := ephemeris.NewManager(jplProvider, swissProvider, cache)
	defer manager.Close()
	
	// Test date
	testDate := time.Date(2024, 7, 18, 12, 0, 0, 0, time.UTC)
	jd := ephemeris.TimeToJulianDay(testDate)
	
	fmt.Printf("=== Ephemeris Integration Example ===\n")
	fmt.Printf("Date: %s\n", testDate.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("Julian Day: %.6f\n\n", jd)
	
	ctx := context.Background()
	
	// Get Sun position
	fmt.Println("--- Sun Position ---")
	sunPos, err := manager.GetSunPosition(ctx, jd)
	if err != nil {
		log.Fatalf("Error getting sun position: %v", err)
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
		log.Fatalf("Error getting moon position: %v", err)
	}
	
	fmt.Printf("Longitude: %.6f°\n", moonPos.Longitude)
	fmt.Printf("Latitude: %.6f°\n", moonPos.Latitude)
	fmt.Printf("Distance: %.2f km\n", moonPos.Distance)
	fmt.Printf("Phase: %.3f (%.1f%% illuminated)\n", moonPos.Phase, moonPos.Illumination)
	fmt.Printf("Angular Diameter: %.2f arcseconds\n\n", moonPos.AngularDiameter)
	
	// Calculate basic Panchangam elements
	fmt.Println("--- Basic Panchangam Elements ---")
	sunLongitude := sunPos.Longitude
	moonLongitude := moonPos.Longitude
	
	// Tithi calculation (lunar day)
	tithiDegrees := moonLongitude - sunLongitude
	if tithiDegrees < 0 {
		tithiDegrees += 360
	}
	tithi := int(tithiDegrees/12) + 1
	
	fmt.Printf("Sun Longitude: %.6f°\n", sunLongitude)
	fmt.Printf("Moon Longitude: %.6f°\n", moonLongitude)
	fmt.Printf("Tithi: %d (%.2f°)\n", tithi, tithiDegrees)
	
	// Nakshatra calculation (lunar mansion)
	nakshatra := int(moonLongitude/13.333333) + 1
	fmt.Printf("Nakshatra: %d (Moon at %.6f°)\n", nakshatra, moonLongitude)
	
	// Yoga calculation
	yogaDegrees := sunLongitude + moonLongitude
	if yogaDegrees >= 360 {
		yogaDegrees -= 360
	}
	yoga := int(yogaDegrees/13.333333) + 1
	fmt.Printf("Yoga: %d (%.2f°)\n", yoga, yogaDegrees)
	
	// Vara (weekday)
	weekday := testDate.Weekday()
	fmt.Printf("Vara: %s\n", weekday)
	
	// Health status
	fmt.Println("\n--- Provider Health ---")
	healthStatuses, err := manager.GetHealthStatus(ctx)
	if err != nil {
		log.Fatalf("Error getting health status: %v", err)
	}
	
	for name, status := range healthStatuses {
		fmt.Printf("%s: Available=%t, Response=%v\n", 
			name, status.Available, status.ResponseTime)
	}
	
	// Cache statistics
	fmt.Println("\n--- Cache Statistics ---")
	stats := cache.GetStats(ctx)
	fmt.Printf("Cache hits: %d\n", stats.Hits)
	fmt.Printf("Cache misses: %d\n", stats.Misses)
	fmt.Printf("Hit rate: %.2f%%\n", stats.HitRate*100)
	
	fmt.Println("\n=== Integration Complete ===")
}