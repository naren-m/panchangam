package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/naren-m/panchangam/astronomy"
)

func main() {
	// Command line flags
	var (
		date      = flag.String("date", time.Now().Format("2006-01-02"), "Date in YYYY-MM-DD format")
		latitude  = flag.Float64("lat", 40.7128, "Latitude (-90 to 90)")
		longitude = flag.Float64("lon", -74.0060, "Longitude (-180 to 180)")
		location  = flag.String("location", "", "Predefined location (nyc, london, tokyo, sydney, mumbai, capetown)")
	)
	flag.Parse()

	// Handle predefined locations
	if *location != "" {
		switch *location {
		case "nyc", "newyork":
			*latitude = 40.7128
			*longitude = -74.0060
			fmt.Println("📍 Using New York coordinates")
		case "london":
			*latitude = 51.5074
			*longitude = -0.1278
			fmt.Println("📍 Using London coordinates")
		case "tokyo":
			*latitude = 35.6762
			*longitude = 139.6503
			fmt.Println("📍 Using Tokyo coordinates")
		case "sydney":
			*latitude = -33.8688
			*longitude = 151.2093
			fmt.Println("📍 Using Sydney coordinates")
		case "mumbai":
			*latitude = 19.0760
			*longitude = 72.8777
			fmt.Println("📍 Using Mumbai coordinates")
		case "capetown":
			*latitude = -33.9249
			*longitude = 18.4241
			fmt.Println("📍 Using Cape Town coordinates")
		default:
			log.Fatalf("Unknown location: %s. Available: nyc, london, tokyo, sydney, mumbai, capetown", *location)
		}
	}

	// Parse date
	dateTime, err := time.Parse("2006-01-02", *date)
	if err != nil {
		log.Fatalf("Invalid date format: %v", err)
	}

	// Create location
	loc := astronomy.Location{
		Latitude:  *latitude,
		Longitude: *longitude,
	}

	// Calculate sunrise and sunset
	sunTimes, err := astronomy.CalculateSunTimes(loc, dateTime)
	if err != nil {
		log.Fatalf("Failed to calculate sun times: %v", err)
	}

	// Display results
	fmt.Printf("\n🌅 Sunrise/Sunset Calculator (Direct)\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("📅 Date: %s\n", *date)
	fmt.Printf("📍 Location: %.4f°N, %.4f°E\n", *latitude, *longitude)
	fmt.Printf("═══════════════════════════════════════\n")

	fmt.Printf("\n📊 Results (UTC):\n")
	fmt.Printf("┌─────────────────────────────────┐\n")
	fmt.Printf("│ 🌅 Sunrise: %-18s │\n", sunTimes.Sunrise.Format("15:04:05"))
	fmt.Printf("│ 🌇 Sunset:  %-18s │\n", sunTimes.Sunset.Format("15:04:05"))
	fmt.Printf("└─────────────────────────────────┘\n")

	// Calculate day length
	dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)
	if dayLength < 0 {
		dayLength += 24 * time.Hour
	}
	fmt.Printf("☀️  Day Length: %v\n", dayLength)

	// Show full timestamps
	fmt.Printf("\n📅 Full Timestamps:\n")
	fmt.Printf("• Sunrise: %s\n", sunTimes.Sunrise.Format(time.RFC3339))
	fmt.Printf("• Sunset:  %s\n", sunTimes.Sunset.Format(time.RFC3339))

	fmt.Printf("\n✨ Calculation completed successfully!\n")
	fmt.Printf("💡 Note: All times are in UTC. Convert to local timezone as needed.\n")
}