package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/naren-m/panchangam/astronomy"
)

const availableSunriseSimpleLocations = "nyc, london, tokyo, sydney, mumbai, capetown"

type sunriseSimpleLocationPreset struct {
	name      string
	latitude  float64
	longitude float64
}

var sunriseSimpleLocationPresets = map[string]sunriseSimpleLocationPreset{
	"nyc":      {name: "New York", latitude: 40.7128, longitude: -74.0060},
	"newyork":  {name: "New York", latitude: 40.7128, longitude: -74.0060},
	"london":   {name: "London", latitude: 51.5074, longitude: -0.1278},
	"tokyo":    {name: "Tokyo", latitude: 35.6762, longitude: 139.6503},
	"sydney":   {name: "Sydney", latitude: -33.8688, longitude: 151.2093},
	"mumbai":   {name: "Mumbai", latitude: 19.0760, longitude: 72.8777},
	"capetown": {name: "Cape Town", latitude: -33.9249, longitude: 18.4241},
}

func main() {
	// Command line flags
	var (
		date      = flag.String("date", time.Now().Format("2006-01-02"), "Date in YYYY-MM-DD format")
		latitude  = flag.Float64("lat", 40.7128, "Latitude (-90 to 90)")
		longitude = flag.Float64("lon", -74.0060, "Longitude (-180 to 180)")
		location  = flag.String("location", "", "Predefined location (nyc, london, tokyo, sydney, mumbai, capetown)")
	)
	flag.Parse()

	presetName, resolvedLat, resolvedLon, err := resolveSunriseSimpleLocation(*location, *latitude, *longitude)
	if err != nil {
		exitSunriseSimple(err)
	}
	if presetName != "" {
		fmt.Printf("Using %s coordinates\n", presetName)
	}

	dateTime, err := parseSunriseSimpleDate(*date)
	if err != nil {
		exitSunriseSimple(fmt.Errorf("invalid date format: %w", err))
	}

	loc := astronomy.Location{
		Latitude:  resolvedLat,
		Longitude: resolvedLon,
	}

	sunTimes, err := astronomy.CalculateSunTimes(loc, dateTime)
	if err != nil {
		exitSunriseSimple(fmt.Errorf("failed to calculate sun times: %w", err))
	}

	// Display results
	fmt.Printf("\nSunrise/Sunset Calculator (Direct)\n")
	fmt.Printf("----------------------------------\n")
	fmt.Printf("Date: %s\n", dateTime.Format("2006-01-02"))
	fmt.Printf("Location: lat %.4f, lon %.4f\n", resolvedLat, resolvedLon)
	fmt.Printf("----------------------------------\n")

	fmt.Printf("\nResults (UTC):\n")
	fmt.Printf("  Sunrise: %s\n", sunTimes.Sunrise.Format("15:04:05"))
	fmt.Printf("  Sunset:  %s\n", sunTimes.Sunset.Format("15:04:05"))

	// Calculate day length
	dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)
	if dayLength < 0 {
		dayLength += 24 * time.Hour
	}
	fmt.Printf("Day Length: %v\n", dayLength)

	// Show full timestamps
	fmt.Printf("\nFull Timestamps:\n")
	fmt.Printf("  Sunrise: %s\n", sunTimes.Sunrise.Format(time.RFC3339))
	fmt.Printf("  Sunset:  %s\n", sunTimes.Sunset.Format(time.RFC3339))

	fmt.Printf("\nCalculation completed successfully.\n")
	fmt.Printf("Note: All times are in UTC. Convert to local timezone as needed.\n")
}

func resolveSunriseSimpleLocation(location string, latitude, longitude float64) (string, float64, float64, error) {
	location = strings.TrimSpace(location)
	if location == "" {
		return "", latitude, longitude, nil
	}

	preset, exists := sunriseSimpleLocationPresets[location]
	if !exists {
		return "", 0, 0, fmt.Errorf("unknown location: %s. Available: %s", location, availableSunriseSimpleLocations)
	}
	return preset.name, preset.latitude, preset.longitude, nil
}

func parseSunriseSimpleDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func exitSunriseSimple(err error) {
	if writeErr := writeSunriseSimpleError(os.Stderr, err); writeErr != nil {
		fmt.Fprintf(os.Stderr, "failed to write sunrise-simple error: %v\n", writeErr)
	}
	os.Exit(1)
}

func writeSunriseSimpleError(w io.Writer, err error) error {
	if err == nil {
		return nil
	}
	_, writeErr := fmt.Fprintln(w, err)
	return writeErr
}
