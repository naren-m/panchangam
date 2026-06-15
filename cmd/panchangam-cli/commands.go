package main

import (
	"time"

	"github.com/spf13/cobra"
)

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
		Use:     "get",
		Short:   "Get panchangam data for a specific date and location",
		Long:    `Retrieve sunrise/sunset times and other panchangam data for a given date and location.`,
		Example: "  panchangam-cli get -l nyc\n  panchangam-cli get -l london -d 2024-06-21\n  panchangam-cli get --lat 37.7749 --lon -122.4194 --tz \"America/Los_Angeles\"\n  panchangam-cli get -l tokyo -o json",
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
	return &cobra.Command{
		Use:   "locations",
		Short: "List available predefined locations",
		Long:  `Display all available predefined locations with their coordinates and timezones.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLocationsCommand()
		},
	}
}

func createValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate server connectivity and basic functionality",
		Long:  `Test the connection to the gRPC server and validate basic functionality.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidateCommand()
		},
	}
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

func createTithiCommand() *cobra.Command {
	var (
		date      string
		latitude  float64
		longitude float64
		timezone  string
		location  string
		detailed  bool
	)

	cmd := &cobra.Command{
		Use:     "tithi",
		Short:   "Calculate Tithi (lunar day) for a specific date and location",
		Long:    `Calculate detailed Tithi information including timing, percentage, and characteristics.`,
		Example: "  panchangam-cli tithi -l mumbai\n  panchangam-cli tithi -d 2024-06-21 --lat 19.0760 --lon 72.8777 --detailed",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTithiCommand(date, latitude, longitude, timezone, location, detailed)
		},
	}

	today := time.Now().Format("2006-01-02")
	cmd.Flags().StringVarP(&date, "date", "d", today, "Date in YYYY-MM-DD format")
	cmd.Flags().Float64Var(&latitude, "lat", 19.0760, "Latitude (-90 to 90)")
	cmd.Flags().Float64Var(&longitude, "lon", 72.8777, "Longitude (-180 to 180)")
	cmd.Flags().StringVar(&timezone, "tz", "Asia/Kolkata", "Timezone")
	cmd.Flags().StringVarP(&location, "location", "l", "", "Predefined location")
	cmd.Flags().BoolVar(&detailed, "detailed", false, "Show detailed Tithi information")

	return cmd
}

func createSunTimesCommand() *cobra.Command {
	var (
		date      string
		latitude  float64
		longitude float64
		timezone  string
		location  string
		detailed  bool
	)

	cmd := &cobra.Command{
		Use:     "sun",
		Short:   "Calculate detailed sun timing information",
		Long:    `Calculate comprehensive sun timing information including sunrise, sunset, solar noon, and day length.`,
		Example: "  panchangam-cli sun -l tokyo\n  panchangam-cli sun -d 2024-06-21 --lat 35.6762 --lon 139.6503 --detailed",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSunTimesCommand(date, latitude, longitude, timezone, location, detailed)
		},
	}

	today := time.Now().Format("2006-01-02")
	cmd.Flags().StringVarP(&date, "date", "d", today, "Date in YYYY-MM-DD format")
	cmd.Flags().Float64Var(&latitude, "lat", 35.6762, "Latitude (-90 to 90)")
	cmd.Flags().Float64Var(&longitude, "lon", 139.6503, "Longitude (-180 to 180)")
	cmd.Flags().StringVar(&timezone, "tz", "Asia/Tokyo", "Timezone")
	cmd.Flags().StringVarP(&location, "location", "l", "", "Predefined location")
	cmd.Flags().BoolVar(&detailed, "detailed", false, "Show detailed sun information")

	return cmd
}

func createVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Display version information for the CLI and service.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersionCommand()
		},
	}
}

func createHealthCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check local CLI and ephemeris health",
		Long:  `Check local CLI status and local ephemeris provider health.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHealthCommand()
		},
	}
}
