package main

import (
	"context"
	"fmt"
	"os"
	"time"

	ppb "github.com/naren-m/panchangam/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runGetCommand(date string, lat, lon float64, tz, location, region, method, locale string) error {
	if location != "" {
		preset, exists := locationPresets[location]
		if !exists {
			return fmt.Errorf("unknown location: %s. Use 'locations' command to see available locations", location)
		}
		lat = preset.Lat
		lon = preset.Lon
		if tz == "" {
			tz = preset.TZ
		}
		fmt.Printf("Using preset location: %s\n", preset.Name)
	}

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer closeClientConnectionSafely(conn)

	client := ppb.NewPanchangamClient(conn)
	req := &ppb.GetPanchangamRequest{
		Date:              date,
		Latitude:          lat,
		Longitude:         lon,
		Timezone:          tz,
		Region:            region,
		CalculationMethod: method,
		Locale:            locale,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := client.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get panchangam data: %w", err)
	}

	switch outputFormat {
	case "json":
		return outputJSON(resp)
	case "yaml":
		return outputYAML(resp)
	default:
		return outputTable(resp, req)
	}
}

func runLocationsCommand() error {
	fmt.Println("Available Predefined Locations:")
	fmt.Println("-------------------------------")
	fmt.Printf("%-12s %-25s %-15s %-20s\n", "CODE", "NAME", "COORDINATES", "TIMEZONE")
	fmt.Println("---------------------------------------------------------------------------")

	for code, preset := range locationPresets {
		coords := fmt.Sprintf("%.4f,%.4f", preset.Lat, preset.Lon)
		fmt.Printf("%-12s %-25s %-15s %-20s\n", code, preset.Name, coords, preset.TZ)
	}

	fmt.Println("\nUsage: panchangam-cli get -l <code>")
	fmt.Println("Example: panchangam-cli get -l london")
	return nil
}

func runValidateCommand() error {
	fmt.Printf("Validating connection to %s...\n", serverAddress)

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer closeClientConnectionSafely(conn)

	client := ppb.NewPanchangamClient(conn)
	req := &ppb.GetPanchangamRequest{
		Date:      time.Now().Format("2006-01-02"),
		Latitude:  40.7128,
		Longitude: -74.0060,
		Timezone:  "America/New_York",
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := client.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	fmt.Println("PASS: connection successful")
	fmt.Printf("PASS: server responded with data for %s\n", resp.PanchangamData.Date)
	fmt.Printf("PASS: sunrise time: %s\n", resp.PanchangamData.SunriseTime)
	fmt.Printf("PASS: sunset time: %s\n", resp.PanchangamData.SunsetTime)
	fmt.Println("PASS: basic validation passed")
	return nil
}

func runBenchmarkCommand(requests, workers int) error {
	fmt.Printf("Benchmarking %s with %d requests using %d workers...\n", serverAddress, requests, workers)

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer closeClientConnectionSafely(conn)

	client := ppb.NewPanchangamClient(conn)
	req := &ppb.GetPanchangamRequest{
		Date:      time.Now().Format("2006-01-02"),
		Latitude:  40.7128,
		Longitude: -74.0060,
		Timezone:  "America/New_York",
	}

	start := time.Now()
	errors := 0

	for i := 0; i < requests; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		_, err := client.Get(ctx, req)
		cancel()

		if err != nil {
			errors++
		}

		if i%10 == 0 {
			fmt.Print(".")
		}
	}

	duration := time.Since(start)
	fmt.Printf("\n\nBenchmark Results:\n")
	fmt.Println("------------------")
	fmt.Printf("Total requests: %d\n", requests)
	fmt.Printf("Successful: %d\n", requests-errors)
	fmt.Printf("Errors: %d\n", errors)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Requests/sec: %.2f\n", float64(requests)/duration.Seconds())
	fmt.Printf("Average latency: %v\n", duration/time.Duration(requests))
	return nil
}

func closeClientConnectionSafely(conn interface{ Close() error }) {
	if conn == nil {
		return
	}
	if err := conn.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close gRPC connection: %v\n", err)
	}
}
