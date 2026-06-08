package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ppb "github.com/naren-m/panchangam/proto"
)

const availableSunriseDemoLocations = "nyc, london, tokyo, sydney, mumbai, capetown"

type sunriseDemoLocationPreset struct {
	name      string
	latitude  float64
	longitude float64
	timezone  string
}

var sunriseDemoLocationPresets = map[string]sunriseDemoLocationPreset{
	"nyc":      {name: "New York", latitude: 40.7128, longitude: -74.0060, timezone: "America/New_York"},
	"newyork":  {name: "New York", latitude: 40.7128, longitude: -74.0060, timezone: "America/New_York"},
	"london":   {name: "London", latitude: 51.5074, longitude: -0.1278, timezone: "Europe/London"},
	"tokyo":    {name: "Tokyo", latitude: 35.6762, longitude: 139.6503, timezone: "Asia/Tokyo"},
	"sydney":   {name: "Sydney", latitude: -33.8688, longitude: 151.2093, timezone: "Australia/Sydney"},
	"mumbai":   {name: "Mumbai", latitude: 19.0760, longitude: 72.8777, timezone: "Asia/Kolkata"},
	"capetown": {name: "Cape Town", latitude: -33.9249, longitude: 18.4241, timezone: "Africa/Johannesburg"},
}

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

	presetName, resolvedLat, resolvedLon, resolvedTimezone, err := resolveSunriseDemoLocation(*location, *latitude, *longitude, *timezone)
	if err != nil {
		exitSunriseDemo(err)
	}
	if presetName != "" {
		fmt.Printf("Using %s coordinates\n", presetName)
	}

	// Connect to gRPC server
	conn, err := grpc.NewClient(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		exitSunriseDemo(fmt.Errorf("failed to connect to server: %w", err))
	}
	defer closeConnectionSafely(conn)

	client := ppb.NewPanchangamClient(conn)

	// Create request
	req := &ppb.GetPanchangamRequest{
		Date:      *date,
		Latitude:  resolvedLat,
		Longitude: resolvedLon,
		Timezone:  resolvedTimezone,
	}

	// Call the service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("\nSunrise/Sunset Calculator\n")
	fmt.Printf("-------------------------\n")
	fmt.Printf("Date: %s\n", *date)
	fmt.Printf("Location: lat %.4f, lon %.4f\n", resolvedLat, resolvedLon)
	fmt.Printf("Timezone: %s\n", resolvedTimezone)
	fmt.Printf("Server: %s\n", *address)
	fmt.Printf("-------------------------\n")

	resp, err := client.Get(ctx, req)
	if err != nil {
		exitSunriseDemo(fmt.Errorf("failed to get panchangam data: %w", err))
	}

	// Display results
	data := resp.PanchangamData
	fmt.Printf("\nResults:\n")
	fmt.Printf("  Sunrise: %s\n", data.SunriseTime)
	fmt.Printf("  Sunset:  %s\n", data.SunsetTime)

	// Calculate day length
	sunrise, err := time.Parse("15:04:05", data.SunriseTime)
	if err == nil {
		sunset, err := time.Parse("15:04:05", data.SunsetTime)
		if err == nil {
			dayLength := sunset.Sub(sunrise)
			if dayLength < 0 {
				dayLength += 24 * time.Hour
			}
			fmt.Printf("Day Length: %v\n", dayLength)
		}
	}

	// Display other panchangam data
	fmt.Printf("\nTraditional Panchangam Data:\n")
	fmt.Printf("  Tithi: %s\n", data.Tithi)
	fmt.Printf("  Nakshatra: %s\n", data.Nakshatra)
	fmt.Printf("  Yoga: %s\n", data.Yoga)
	fmt.Printf("  Karana: %s\n", data.Karana)

	if len(data.Events) > 0 {
		fmt.Printf("\nEvents:\n")
		for _, event := range data.Events {
			fmt.Printf("  %s at %s\n", event.Name, event.Time)
		}
	}

	fmt.Printf("\nCalculation completed successfully.\n")
}

func resolveSunriseDemoLocation(location string, latitude, longitude float64, timezone string) (string, float64, float64, string, error) {
	location = strings.TrimSpace(location)
	if location == "" {
		return "", latitude, longitude, timezone, nil
	}

	preset, exists := sunriseDemoLocationPresets[location]
	if !exists {
		return "", 0, 0, "", fmt.Errorf("unknown location: %s. Available: %s", location, availableSunriseDemoLocations)
	}
	return preset.name, preset.latitude, preset.longitude, preset.timezone, nil
}

func exitSunriseDemo(err error) {
	if writeErr := writeSunriseDemoError(os.Stderr, err); writeErr != nil {
		fmt.Fprintf(os.Stderr, "failed to write sunrise-demo error: %v\n", writeErr)
	}
	os.Exit(1)
}

func writeSunriseDemoError(w io.Writer, err error) error {
	if err == nil {
		return nil
	}
	_, writeErr := fmt.Fprintln(w, err)
	return writeErr
}

func closeConnectionSafely(conn interface{ Close() error }) {
	if conn == nil {
		return
	}
	if err := conn.Close(); err != nil {
		log.Printf("failed to close gRPC connection: %v", err)
	}
}
