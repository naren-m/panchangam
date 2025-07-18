package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ppb "github.com/naren-m/panchangam/proto/panchangam"
)

func main() {
	// Command line flags
	var (
		address   = flag.String("address", "localhost:8080", "gRPC server address")
		date      = flag.String("date", time.Now().Format("2006-01-02"), "Date in YYYY-MM-DD format")
		latitude  = flag.Float64("lat", 40.7128, "Latitude (-90 to 90)")
		longitude = flag.Float64("lon", -74.0060, "Longitude (-180 to 180)")
		timezone  = flag.String("tz", "America/New_York", "Timezone (e.g., America/New_York, Asia/Kolkata)")
		location  = flag.String("location", "", "Predefined location (nyc, london, tokyo, sydney, mumbai, capetown)")
	)
	flag.Parse()

	// Handle predefined locations
	if *location != "" {
		switch *location {
		case "nyc", "newyork":
			*latitude = 40.7128
			*longitude = -74.0060
			*timezone = "America/New_York"
			fmt.Println("📍 Using New York coordinates")
		case "london":
			*latitude = 51.5074
			*longitude = -0.1278
			*timezone = "Europe/London"
			fmt.Println("📍 Using London coordinates")
		case "tokyo":
			*latitude = 35.6762
			*longitude = 139.6503
			*timezone = "Asia/Tokyo"
			fmt.Println("📍 Using Tokyo coordinates")
		case "sydney":
			*latitude = -33.8688
			*longitude = 151.2093
			*timezone = "Australia/Sydney"
			fmt.Println("📍 Using Sydney coordinates")
		case "mumbai":
			*latitude = 19.0760
			*longitude = 72.8777
			*timezone = "Asia/Kolkata"
			fmt.Println("📍 Using Mumbai coordinates")
		case "capetown":
			*latitude = -33.9249
			*longitude = 18.4241
			*timezone = "Africa/Johannesburg"
			fmt.Println("📍 Using Cape Town coordinates")
		default:
			log.Fatalf("Unknown location: %s. Available: nyc, london, tokyo, sydney, mumbai, capetown", *location)
		}
	}

	// Connect to gRPC server
	conn, err := grpc.NewClient(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := ppb.NewPanchangamClient(conn)

	// Create request
	req := &ppb.GetPanchangamRequest{
		Date:      *date,
		Latitude:  *latitude,
		Longitude: *longitude,
		Timezone:  *timezone,
	}

	// Call the service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("\n🌅 Sunrise/Sunset Calculator\n")
	fmt.Printf("═══════════════════════════════\n")
	fmt.Printf("📅 Date: %s\n", *date)
	fmt.Printf("📍 Location: %.4f°N, %.4f°E\n", *latitude, *longitude)
	fmt.Printf("🌐 Timezone: %s\n", *timezone)
	fmt.Printf("🔗 Server: %s\n", *address)
	fmt.Printf("═══════════════════════════════\n")

	resp, err := client.Get(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get panchangam data: %v", err)
	}

	// Display results
	data := resp.PanchangamData
	fmt.Printf("\n📊 Results:\n")
	fmt.Printf("┌─────────────────────────────────┐\n")
	fmt.Printf("│ 🌅 Sunrise: %-18s │\n", data.SunriseTime)
	fmt.Printf("│ 🌇 Sunset:  %-18s │\n", data.SunsetTime)
	fmt.Printf("└─────────────────────────────────┘\n")

	// Calculate day length
	sunrise, err := time.Parse("15:04:05", data.SunriseTime)
	if err == nil {
		sunset, err := time.Parse("15:04:05", data.SunsetTime)
		if err == nil {
			dayLength := sunset.Sub(sunrise)
			if dayLength < 0 {
				dayLength += 24 * time.Hour
			}
			fmt.Printf("☀️  Day Length: %v\n", dayLength)
		}
	}

	// Display other panchangam data
	fmt.Printf("\n📜 Traditional Panchangam Data:\n")
	fmt.Printf("• Tithi: %s\n", data.Tithi)
	fmt.Printf("• Nakshatra: %s\n", data.Nakshatra)
	fmt.Printf("• Yoga: %s\n", data.Yoga)
	fmt.Printf("• Karana: %s\n", data.Karana)

	if len(data.Events) > 0 {
		fmt.Printf("\n📅 Events:\n")
		for _, event := range data.Events {
			fmt.Printf("• %s at %s\n", event.Name, event.Time)
		}
	}

	fmt.Printf("\n✨ Calculation completed successfully!\n")
}