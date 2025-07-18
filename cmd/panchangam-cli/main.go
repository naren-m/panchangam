package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ppb "github.com/naren-m/panchangam/proto/panchangam"
)

var (
	serverAddress string
	outputFormat  string
	timeout       time.Duration
)

// Location presets
var locationPresets = map[string]struct {
	Lat  float64
	Lon  float64
	TZ   string
	Name string
}{
	"nyc": {40.7128, -74.0060, "America/New_York", "New York, USA"},
	"london": {51.5074, -0.1278, "Europe/London", "London, UK"},
	"tokyo": {35.6762, 139.6503, "Asia/Tokyo", "Tokyo, Japan"},
	"sydney": {-33.8688, 151.2093, "Australia/Sydney", "Sydney, Australia"},
	"mumbai": {19.0760, 72.8777, "Asia/Kolkata", "Mumbai, India"},
	"capetown": {-33.9249, 18.4241, "Africa/Johannesburg", "Cape Town, South Africa"},
	"paris": {48.8566, 2.3522, "Europe/Paris", "Paris, France"},
	"moscow": {55.7558, 37.6176, "Europe/Moscow", "Moscow, Russia"},
	"beijing": {39.9042, 116.4074, "Asia/Shanghai", "Beijing, China"},
	"cairo": {30.0444, 31.2357, "Africa/Cairo", "Cairo, Egypt"},
	"rio": {-22.9068, -43.1729, "America/Sao_Paulo", "Rio de Janeiro, Brazil"},
	"losangeles": {34.0522, -118.2437, "America/Los_Angeles", "Los Angeles, USA"},
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "panchangam-cli",
		Short: "CLI client for Panchangam gRPC service",
		Long: `A comprehensive CLI client for testing the Panchangam gRPC service.
Supports sunrise/sunset calculations for any location and date.`,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&serverAddress, "server", "s", "localhost:8080", "gRPC server address")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "Request timeout")

	// Add subcommands
	rootCmd.AddCommand(createGetCommand())
	rootCmd.AddCommand(createLocationsCommand())
	rootCmd.AddCommand(createValidateCommand())
	rootCmd.AddCommand(createBenchmarkCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createGetCommand() *cobra.Command {
	var (
		date      string
		latitude  float64
		longitude float64
		timezone  string
		location  string
		region    string
		method    string
		locale    string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get panchangam data for a specific date and location",
		Long:  `Retrieve sunrise/sunset times and other panchangam data for a given date and location.`,
		Example: `  # Get data for today in New York
  panchangam-cli get -l nyc

  # Get data for specific date in London
  panchangam-cli get -l london -d 2024-06-21

  # Get data for custom coordinates
  panchangam-cli get --lat 37.7749 --lon -122.4194 --tz "America/Los_Angeles"

  # Get data with JSON output
  panchangam-cli get -l tokyo -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCommand(date, latitude, longitude, timezone, location, region, method, locale)
		},
	}

	today := time.Now().Format("2006-01-02")
	cmd.Flags().StringVarP(&date, "date", "d", today, "Date in YYYY-MM-DD format")
	cmd.Flags().Float64Var(&latitude, "lat", 40.7128, "Latitude (-90 to 90)")
	cmd.Flags().Float64Var(&longitude, "lon", -74.0060, "Longitude (-180 to 180)")
	cmd.Flags().StringVar(&timezone, "tz", "", "Timezone (e.g., America/New_York)")
	cmd.Flags().StringVarP(&location, "location", "l", "", "Predefined location (use 'locations' command to see available)")
	cmd.Flags().StringVar(&region, "region", "", "Regional system (e.g., Tamil Nadu, Kerala)")
	cmd.Flags().StringVar(&method, "method", "", "Calculation method (e.g., Drik, Vakya)")
	cmd.Flags().StringVar(&locale, "locale", "", "Language/locale (e.g., en, ta)")

	return cmd
}

func createLocationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locations",
		Short: "List available predefined locations",
		Long:  `Display all available predefined locations with their coordinates and timezones.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLocationsCommand()
		},
	}

	return cmd
}

func createValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate server connectivity and basic functionality",
		Long:  `Test the connection to the gRPC server and validate basic functionality.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidateCommand()
		},
	}

	return cmd
}

func createBenchmarkCommand() *cobra.Command {
	var (
		requests int
		workers  int
	)

	cmd := &cobra.Command{
		Use:   "benchmark",
		Short: "Benchmark server performance",
		Long:  `Run performance benchmarks against the gRPC server.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBenchmarkCommand(requests, workers)
		},
	}

	cmd.Flags().IntVarP(&requests, "requests", "n", 100, "Number of requests to make")
	cmd.Flags().IntVarP(&workers, "workers", "w", 10, "Number of concurrent workers")

	return cmd
}

func runGetCommand(date string, lat, lon float64, tz, location, region, method, locale string) error {
	// Handle predefined locations
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
		fmt.Printf("📍 Using preset location: %s\n", preset.Name)
	}

	// Connect to server
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := ppb.NewPanchangamClient(conn)

	// Create request
	req := &ppb.GetPanchangamRequest{
		Date:               date,
		Latitude:           lat,
		Longitude:          lon,
		Timezone:           tz,
		Region:             region,
		CalculationMethod:  method,
		Locale:             locale,
	}

	// Make request
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := client.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get panchangam data: %v", err)
	}

	// Output results
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
	fmt.Println("📍 Available Predefined Locations:")
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("%-12s %-25s %-15s %-20s\n", "CODE", "NAME", "COORDINATES", "TIMEZONE")
	fmt.Println("───────────────────────────────────────────────────────────────────────────")

	for code, preset := range locationPresets {
		coords := fmt.Sprintf("%.4f,%.4f", preset.Lat, preset.Lon)
		fmt.Printf("%-12s %-25s %-15s %-20s\n", code, preset.Name, coords, preset.TZ)
	}

	fmt.Println("\n💡 Usage: panchangam-cli get -l <code>")
	fmt.Println("   Example: panchangam-cli get -l london")
	return nil
}

func runValidateCommand() error {
	fmt.Printf("🔍 Validating connection to %s...\n", serverAddress)

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("❌ Failed to connect: %v", err)
	}
	defer conn.Close()

	client := ppb.NewPanchangamClient(conn)

	// Test basic functionality
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
		return fmt.Errorf("❌ Request failed: %v", err)
	}

	fmt.Println("✅ Connection successful!")
	fmt.Printf("✅ Server responded with data for %s\n", resp.PanchangamData.Date)
	fmt.Printf("✅ Sunrise time: %s\n", resp.PanchangamData.SunriseTime)
	fmt.Printf("✅ Sunset time: %s\n", resp.PanchangamData.SunsetTime)
	fmt.Println("✅ Basic validation passed!")

	return nil
}

func runBenchmarkCommand(requests, workers int) error {
	fmt.Printf("🚀 Benchmarking %s with %d requests using %d workers...\n", serverAddress, requests, workers)

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := ppb.NewPanchangamClient(conn)

	// Create request
	req := &ppb.GetPanchangamRequest{
		Date:      time.Now().Format("2006-01-02"),
		Latitude:  40.7128,
		Longitude: -74.0060,
		Timezone:  "America/New_York",
	}

	// Run benchmark
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
			fmt.Printf(".")
		}
	}
	
	duration := time.Since(start)
	fmt.Printf("\n\n📊 Benchmark Results:\n")
	fmt.Printf("═══════════════════════\n")
	fmt.Printf("Total requests: %d\n", requests)
	fmt.Printf("Successful: %d\n", requests-errors)
	fmt.Printf("Errors: %d\n", errors)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Requests/sec: %.2f\n", float64(requests)/duration.Seconds())
	fmt.Printf("Average latency: %v\n", duration/time.Duration(requests))

	return nil
}

func outputTable(resp *ppb.GetPanchangamResponse, req *ppb.GetPanchangamRequest) error {
	data := resp.PanchangamData
	
	fmt.Printf("\n🌅 Panchangam Data\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("📅 Date: %s\n", data.Date)
	fmt.Printf("📍 Location: %.4f°N, %.4f°E\n", req.Latitude, req.Longitude)
	if req.Timezone != "" {
		fmt.Printf("🌐 Timezone: %s\n", req.Timezone)
	}
	fmt.Printf("═══════════════════════════════════════\n")

	fmt.Printf("\n📊 Sunrise/Sunset Times (UTC):\n")
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

	fmt.Printf("\n📜 Traditional Panchangam:\n")
	fmt.Printf("• Tithi: %s\n", data.Tithi)
	fmt.Printf("• Nakshatra: %s\n", data.Nakshatra)
	fmt.Printf("• Yoga: %s\n", data.Yoga)
	fmt.Printf("• Karana: %s\n", data.Karana)

	if len(data.Events) > 0 {
		fmt.Printf("\n📅 Events:\n")
		for _, event := range data.Events {
			fmt.Printf("• %s at %s", event.Name, event.Time)
			if event.EventType != "" {
				fmt.Printf(" (%s)", event.EventType)
			}
			fmt.Println()
		}
	}

	return nil
}

func outputJSON(resp *ppb.GetPanchangamResponse) error {
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func outputYAML(resp *ppb.GetPanchangamResponse) error {
	// Simple YAML-like output
	data := resp.PanchangamData
	fmt.Printf("panchangam_data:\n")
	fmt.Printf("  date: %s\n", data.Date)
	fmt.Printf("  sunrise_time: %s\n", data.SunriseTime)
	fmt.Printf("  sunset_time: %s\n", data.SunsetTime)
	fmt.Printf("  tithi: %s\n", data.Tithi)
	fmt.Printf("  nakshatra: %s\n", data.Nakshatra)
	fmt.Printf("  yoga: %s\n", data.Yoga)
	fmt.Printf("  karana: %s\n", data.Karana)
	
	if len(data.Events) > 0 {
		fmt.Printf("  events:\n")
		for _, event := range data.Events {
			fmt.Printf("    - name: %s\n", event.Name)
			fmt.Printf("      time: %s\n", event.Time)
			if event.EventType != "" {
				fmt.Printf("      type: %s\n", event.EventType)
			}
		}
	}
	
	return nil
}