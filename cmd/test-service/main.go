package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/naren-m/panchangam/services/panchangam"
)

func main() {
	fmt.Println("Testing Panchangam Service End-to-End")
	fmt.Println("=========================================")

	// Initialize observability for testing
	observability.NewLocalObserver()

	// Create service
	server := panchangam.NewPanchangamServer()

	// Test cases
	testCases := []struct {
		name    string
		request *ppb.GetPanchangamRequest
	}{
		{
			name: "Bangalore, India - January 15, 2024",
			request: &ppb.GetPanchangamRequest{
				Date:              "2024-01-15",
				Latitude:          12.9716,
				Longitude:         77.5946,
				Timezone:          "Asia/Kolkata",
				Region:            "India",
				CalculationMethod: "traditional",
				Locale:            "en",
			},
		},
		{
			name: "New York, USA - June 21, 2024 (Summer Solstice)",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  40.7128,
				Longitude: -74.0060,
				Timezone:  "America/New_York",
				Region:    "USA",
			},
		},
		{
			name: "London, UK - December 21, 2024 (Winter Solstice)",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-12-21",
				Latitude:  51.5074,
				Longitude: -0.1278,
				Timezone:  "Europe/London",
				Region:    "UK",
			},
		},
	}

	// Run tests
	fmt.Println("Running service tests...")
	successCount := 0
	totalTime := time.Duration(0)

	for i, tc := range testCases {
		fmt.Printf("\nTest %d: %s\n", i+1, tc.name)
		fmt.Printf("   Date: %s\n", tc.request.Date)
		fmt.Printf("   Location: %.4f, %.4f\n", tc.request.Latitude, tc.request.Longitude)
		fmt.Printf("   Timezone: %s\n", tc.request.Timezone)

		start := time.Now()
		ctx := context.Background()
		response, err := server.Get(ctx, tc.request)
		duration := time.Since(start)
		totalTime += duration

		if err != nil {
			fmt.Printf("   FAIL: %v\n", err)
			continue
		}

		data := response.PanchangamData
		if data == nil {
			fmt.Printf("   FAIL: no data in response\n")
			continue
		}

		// Validate response
		fmt.Printf("   PASS: completed in %v\n", duration)
		fmt.Printf("      Tithi: %s\n", data.Tithi)
		fmt.Printf("      Nakshatra: %s\n", data.Nakshatra)
		fmt.Printf("      Yoga: %s\n", data.Yoga)
		fmt.Printf("      Karana: %s\n", data.Karana)
		fmt.Printf("      Sunrise: %s\n", data.SunriseTime)
		fmt.Printf("      Sunset: %s\n", data.SunsetTime)
		fmt.Printf("      Events: %d\n", len(data.Events))

		successCount++
	}

	fmt.Println("\nTest Results Summary")
	fmt.Println("======================")
	fmt.Printf("Successful: %d/%d\n", successCount, len(testCases))
	fmt.Printf("Total Time: %v\n", totalTime)
	fmt.Printf("Average Time: %v\n", totalTime/time.Duration(len(testCases)))

	if successCount == len(testCases) {
		fmt.Println("All tests passed. Service is working correctly.")
	} else {
		if err := writeServiceTestFailure(os.Stderr, len(testCases)-successCount); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write service test failure: %v\n", err)
		}
		os.Exit(1)
	}

	// Test performance
	fmt.Println("\nPerformance Test")
	fmt.Println("===================")
	performanceRuns := 10
	performanceStart := time.Now()

	for i := 0; i < performanceRuns; i++ {
		ctx := context.Background()
		_, err := server.Get(ctx, testCases[0].request)
		if err != nil {
			fmt.Printf("Performance test run %d failed: %v\n", i+1, err)
		}
	}

	performanceDuration := time.Since(performanceStart)
	avgPerRequest := performanceDuration / time.Duration(performanceRuns)

	fmt.Printf("Performance Results:\n")
	fmt.Printf("   Requests: %d in %v\n", performanceRuns, performanceDuration)
	fmt.Printf("   Average: %v per request\n", avgPerRequest)
	fmt.Printf("   Rate: %.1f requests/second\n", float64(performanceRuns)/performanceDuration.Seconds())

	if avgPerRequest < 50*time.Millisecond {
		fmt.Println("PASS: performance target met (< 50ms average)")
	} else {
		fmt.Printf("WARN: performance target missed (average %v > 50ms)\n", avgPerRequest)
	}

	fmt.Println("\nService validation complete.")
}

func writeServiceTestFailure(w io.Writer, failedTests int) error {
	if _, err := fmt.Fprintf(w, "WARN: %d tests failed. Check error messages above.\n", failedTests); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w, "Service test failed")
	return err
}
