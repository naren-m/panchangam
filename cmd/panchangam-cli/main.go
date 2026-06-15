package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	serverAddress string
	outputFormat  string
	timeout       time.Duration
	verbose       bool
	debug         bool
)

var locationPresets = map[string]struct {
	Lat  float64
	Lon  float64
	TZ   string
	Name string
}{
	"nyc":        {40.7128, -74.0060, "America/New_York", "New York, USA"},
	"london":     {51.5074, -0.1278, "Europe/London", "London, UK"},
	"tokyo":      {35.6762, 139.6503, "Asia/Tokyo", "Tokyo, Japan"},
	"sydney":     {-33.8688, 151.2093, "Australia/Sydney", "Sydney, Australia"},
	"mumbai":     {19.0760, 72.8777, "Asia/Kolkata", "Mumbai, India"},
	"capetown":   {-33.9249, 18.4241, "Africa/Johannesburg", "Cape Town, South Africa"},
	"paris":      {48.8566, 2.3522, "Europe/Paris", "Paris, France"},
	"moscow":     {55.7558, 37.6176, "Europe/Moscow", "Moscow, Russia"},
	"beijing":    {39.9042, 116.4074, "Asia/Shanghai", "Beijing, China"},
	"cairo":      {30.0444, 31.2357, "Africa/Cairo", "Cairo, Egypt"},
	"rio":        {-22.9068, -43.1729, "America/Sao_Paulo", "Rio de Janeiro, Brazil"},
	"losangeles": {34.0522, -118.2437, "America/Los_Angeles", "Los Angeles, USA"},
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "panchangam-cli",
		Short: "Comprehensive CLI for Panchangam astronomical calculations",
		Long: `Panchangam CLI - A comprehensive astronomical calculation tool for Hindu calendar systems.

This CLI provides access to various astronomical calculations including:
  - Tithi (lunar day) calculations with timing and characteristics
  - Detailed sunrise/sunset times with solar noon and day length

Features:
  - Support for multiple output formats (table, json, yaml, csv)
  - Predefined location presets for major cities worldwide
  - Custom coordinate input with timezone support
  - Detailed mode for comprehensive information display
  - Local ephemeris health checks

Examples:
  # Get today's Tithi for Mumbai
  panchangam-cli tithi -l mumbai

  # Get detailed sun times for Tokyo
  panchangam-cli sun -l tokyo --detailed

  # Get Panchangam data in JSON format
  panchangam-cli get -l london -o json

  # Use custom coordinates
  panchangam-cli tithi --lat 19.0760 --lon 72.8777 --tz "Asia/Kolkata"

  # Check local health
  panchangam-cli health

For more information on a specific command, use:
  panchangam-cli [command] --help`,
	}

	rootCmd.PersistentFlags().StringVarP(&serverAddress, "server", "s", "localhost:8080", "gRPC server address")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml, csv)")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "Request timeout")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")

	rootCmd.AddCommand(createGetCommand())
	rootCmd.AddCommand(createTithiCommand())
	rootCmd.AddCommand(createSunTimesCommand())
	rootCmd.AddCommand(createLocationsCommand())
	rootCmd.AddCommand(createValidateCommand())
	rootCmd.AddCommand(createBenchmarkCommand())
	rootCmd.AddCommand(createVersionCommand())
	rootCmd.AddCommand(createHealthCommand())

	if err := rootCmd.Execute(); err != nil {
		if writeErr := writeCommandError(os.Stderr, err); writeErr != nil {
			fmt.Fprintf(os.Stderr, "failed to write command error: %v\n", writeErr)
		}
		os.Exit(1)
	}
}

func writeCommandError(w io.Writer, err error) error {
	if err == nil {
		return nil
	}
	_, writeErr := fmt.Fprintln(w, err)
	return writeErr
}
