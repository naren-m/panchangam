package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"gopkg.in/yaml.v3"
)

func runTithiCommand(date string, lat, lon float64, tz, location string, detailed bool) error {
	lat, lon, tz, err := resolveLocalCommandLocation(lat, lon, tz, location, "Asia/Kolkata")
	if err != nil {
		return err
	}

	dateInTZ, err := localCommandDateInTimezone(date, tz)
	if err != nil {
		return err
	}

	manager := ephemeris.NewManager(
		ephemeris.NewJPLProvider(),
		ephemeris.NewSwissProvider(),
		ephemeris.NewMemoryCache(100, time.Hour),
	)
	tithiInfo, err := astronomy.NewTithiCalculator(manager).GetTithiForDate(context.Background(), dateInTZ)
	if err != nil {
		return fmt.Errorf("failed to calculate tithi: %w", err)
	}
	if err := manager.Close(); err != nil {
		return fmt.Errorf("failed to close ephemeris manager: %w", err)
	}

	switch outputFormat {
	case "json":
		return outputTithiJSON(tithiInfo)
	case "yaml":
		return outputTithiYAML(tithiInfo)
	default:
		return outputTithiTable(tithiInfo, detailed)
	}
}

func runVersionCommand() error {
	version := map[string]interface{}{
		"cli_version": "1.0.0",
		"api_version": "1.0.0",
		"build_date":  time.Now().Format("2006-01-02"),
		"go_version":  "1.23+",
		"supported_features": []string{
			"tithi_calculation",
			"sunrise_sunset",
			"regional_variations",
		},
	}

	switch outputFormat {
	case "json":
		data, err := json.MarshalIndent(version, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(version)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
	default:
		fmt.Println("Panchangam CLI")
		fmt.Println("---------------")
		fmt.Printf("CLI Version: %s\n", version["cli_version"])
		fmt.Printf("API Version: %s\n", version["api_version"])
		fmt.Printf("Build Date: %s\n", version["build_date"])
		fmt.Printf("Go Version: %s\n", version["go_version"])
		fmt.Println("\nSupported Features:")
		for _, feature := range version["supported_features"].([]string) {
			fmt.Printf("  - %s\n", strings.ReplaceAll(feature, "_", " "))
		}
	}
	return nil
}

func runHealthCommand() error {
	manager := ephemeris.NewManager(
		ephemeris.NewJPLProvider(),
		ephemeris.NewSwissProvider(),
		ephemeris.NewMemoryCache(100, time.Hour),
	)
	statuses, err := manager.GetHealthStatus(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get ephemeris health: %w", err)
	}
	if err := manager.Close(); err != nil {
		return fmt.Errorf("failed to close ephemeris manager: %w", err)
	}

	ephemerisStatus := "healthy"
	providers := make(map[string]string)
	if len(statuses) == 0 {
		ephemerisStatus = "unavailable"
	}
	for name, status := range statuses {
		if status.Available {
			providers[name] = "available"
		} else {
			providers[name] = "unavailable"
			ephemerisStatus = "degraded"
		}
	}

	health := map[string]interface{}{
		"timestamp":        time.Now().Format(time.RFC3339),
		"cli_status":       "healthy",
		"ephemeris_status": ephemerisStatus,
		"providers":        providers,
	}

	switch outputFormat {
	case "json":
		data, err := json.MarshalIndent(health, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(health)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
	default:
		fmt.Println("Local Health Check")
		fmt.Println("------------------")
		fmt.Printf("Timestamp: %s\n", health["timestamp"])
		fmt.Printf("CLI Status: %s\n", health["cli_status"])
		fmt.Printf("Ephemeris Status: %s\n", health["ephemeris_status"])
		fmt.Println("\nProviders:")
		providers := health["providers"].(map[string]string)
		for name, status := range providers {
			fmt.Printf("  - %s: %s\n", strings.ReplaceAll(name, "_", " "), status)
		}
	}
	return nil
}

func runSunTimesCommand(date string, lat, lon float64, tz, location string, detailed bool) error {
	lat, lon, tz, err := resolveLocalCommandLocation(lat, lon, tz, location, "Asia/Tokyo")
	if err != nil {
		return err
	}

	dateInTZ, err := localCommandDateInTimezone(date, tz)
	if err != nil {
		return err
	}

	astronomyLocation := astronomy.Location{Latitude: lat, Longitude: lon}

	sunTimes, err := astronomy.CalculateSunTimes(astronomyLocation, dateInTZ)
	if err != nil {
		return fmt.Errorf("failed to calculate sun times: %w", err)
	}

	switch outputFormat {
	case "json":
		return outputSunTimesJSON(sunTimes, astronomyLocation, dateInTZ)
	case "yaml":
		return outputSunTimesYAML(sunTimes, astronomyLocation, dateInTZ)
	default:
		return outputSunTimesTable(sunTimes, astronomyLocation, dateInTZ, detailed)
	}
}

func resolveLocalCommandLocation(lat, lon float64, tz, location, defaultTZ string) (float64, float64, string, error) {
	if location == "" {
		return lat, lon, tz, nil
	}

	preset, exists := locationPresets[location]
	if !exists {
		return 0, 0, "", fmt.Errorf("unknown location: %s", location)
	}

	if tz == defaultTZ {
		tz = preset.TZ
	}
	return preset.Lat, preset.Lon, tz, nil
}

func localCommandDateInTimezone(date, tz string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(date))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	loc, err := time.LoadLocation(strings.TrimSpace(tz))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone: %w", err)
	}

	return time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, loc), nil
}
