package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestVersionCommandDoesNotAdvertiseUnimplementedCommands(t *testing.T) {
	previousOutputFormat := outputFormat
	outputFormat = "json"
	defer func() {
		outputFormat = previousOutputFormat
	}()

	output := captureStdoutForTest(t, func() {
		if err := runVersionCommand(); err != nil {
			t.Fatalf("runVersionCommand returned error: %v", err)
		}
	})

	var version struct {
		SupportedFeatures []string `json:"supported_features"`
	}
	if err := json.Unmarshal([]byte(output), &version); err != nil {
		t.Fatalf("failed to parse version JSON: %v\noutput: %s", err, output)
	}

	unimplemented := map[string]bool{
		"nakshatra_calculation": false,
		"yoga_calculation":      false,
		"karana_calculation":    false,
		"ephemeris_data":        false,
		"event_calculations":    false,
		"muhurta_calculation":   false,
		"date_range":            false,
	}
	for _, feature := range version.SupportedFeatures {
		if _, exists := unimplemented[feature]; exists {
			unimplemented[feature] = true
		}
	}

	for feature, advertised := range unimplemented {
		if advertised {
			t.Fatalf("version output advertises unimplemented feature %q", feature)
		}
	}
}

func TestCliDoesNotShipUnimplementedLocalCommandScaffolding(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed to locate project root: %v", err)
	}

	staleSnippets := []string{
		"createNakshatraCommand",
		"createYogaCommand",
		"createKaranaCommand",
		"createEphemerisCommand",
		"createDateRangeCommand",
		"createEventsCommand",
		"createMuhurtaCommand",
		"runNakshatraCommand",
		"runYogaCommand",
		"runKaranaCommand",
		"runEphemerisCommand",
		"runDateRangeCommand",
		"runEventsCommand",
		"runMuhurtaCommand",
		"not implemented yet",
		"Framework Ready",
		"implementation pending",
		"future commands as placeholders",
	}

	paths := []string{
		"cmd/panchangam-cli/main.go",
		"cmd/panchangam-cli/commands.go",
		"cmd/panchangam-cli/local_commands.go",
		"cmd/panchangam-cli/README.md",
		"cmd/panchangam-cli/COMMANDS.md",
		"cmd/panchangam-cli/USER_GUIDE.md",
	}

	for _, relativePath := range paths {
		content, err := os.ReadFile(filepath.Join(projectRoot, relativePath))
		if err != nil {
			t.Fatalf("failed to read %s: %v", relativePath, err)
		}
		for _, staleSnippet := range staleSnippets {
			if strings.Contains(string(content), staleSnippet) {
				t.Fatalf("%s still contains stale unimplemented CLI snippet %q", relativePath, staleSnippet)
			}
		}
	}
}

func TestLocalCommandErrorsWrapUnderlyingCause(t *testing.T) {
	err := runTithiCommand("bad-date", 13.0827, 80.2707, "Asia/Kolkata", "", false)
	if err == nil {
		t.Fatal("expected invalid date to fail")
	}
	if errors.Unwrap(err) == nil {
		t.Fatalf("expected invalid date error to wrap the parse error, got %v", err)
	}

	err = runSunTimesCommand("2024-01-01", 13.0827, 80.2707, "Not/AZone", "", false)
	if err == nil {
		t.Fatal("expected invalid timezone to fail")
	}
	if errors.Unwrap(err) == nil {
		t.Fatalf("expected invalid timezone error to wrap the location error, got %v", err)
	}
}

func TestResolveLocalCommandLocation(t *testing.T) {
	lat, lon, tz, err := resolveLocalCommandLocation(1, 2, "Asia/Kolkata", "nyc", "Asia/Kolkata")
	if err != nil {
		t.Fatalf("resolveLocalCommandLocation returned error: %v", err)
	}
	if lat != 40.7128 || lon != -74.0060 || tz != "America/New_York" {
		t.Fatalf("expected nyc preset with preset timezone, got lat=%f lon=%f tz=%s", lat, lon, tz)
	}

	lat, lon, tz, err = resolveLocalCommandLocation(1, 2, "UTC", "nyc", "Asia/Kolkata")
	if err != nil {
		t.Fatalf("resolveLocalCommandLocation returned error: %v", err)
	}
	if lat != 40.7128 || lon != -74.0060 || tz != "UTC" {
		t.Fatalf("expected nyc preset with explicit timezone preserved, got lat=%f lon=%f tz=%s", lat, lon, tz)
	}

	_, _, _, err = resolveLocalCommandLocation(1, 2, "UTC", "unknown", "Asia/Kolkata")
	if err == nil {
		t.Fatal("expected unknown location to return an error")
	}
	if !strings.Contains(err.Error(), "unknown location: unknown") {
		t.Fatalf("expected unknown location error, got %v", err)
	}
}

func TestLocalCommandDateInTimezone(t *testing.T) {
	date, err := localCommandDateInTimezone(" 2024-01-15 ", " Asia/Kolkata ")
	if err != nil {
		t.Fatalf("localCommandDateInTimezone returned error: %v", err)
	}
	if date.Year() != 2024 || date.Month() != time.January || date.Day() != 15 {
		t.Fatalf("expected 2024-01-15, got %s", date.Format(time.RFC3339))
	}
	if date.Location().String() != "Asia/Kolkata" {
		t.Fatalf("expected Asia/Kolkata location, got %s", date.Location())
	}
	if date.Hour() != 0 || date.Minute() != 0 || date.Second() != 0 {
		t.Fatalf("expected midnight in timezone, got %s", date.Format(time.RFC3339))
	}
}

func TestTithiCommandDoesNotReturnHardCodedPurnima(t *testing.T) {
	previousOutputFormat := outputFormat
	outputFormat = "json"
	defer func() {
		outputFormat = previousOutputFormat
	}()

	output := captureStdoutForTest(t, func() {
		if err := runTithiCommand("2024-01-01", 13.0827, 80.2707, "Asia/Kolkata", "", false); err != nil {
			t.Fatalf("runTithiCommand returned error: %v", err)
		}
	})

	var tithi struct {
		Number int    `json:"number"`
		Name   string `json:"name"`
	}
	if err := json.Unmarshal([]byte(output), &tithi); err != nil {
		t.Fatalf("failed to parse tithi JSON: %v\noutput: %s", err, output)
	}
	if tithi.Number == 15 && tithi.Name == "Purnima" {
		t.Fatalf("tithi command returned the old hard-coded Purnima payload: %+v", tithi)
	}
}

func TestHealthCommandReportsRealEphemerisStatus(t *testing.T) {
	previousOutputFormat := outputFormat
	outputFormat = "json"
	defer func() {
		outputFormat = previousOutputFormat
	}()

	output := captureStdoutForTest(t, func() {
		if err := runHealthCommand(); err != nil {
			t.Fatalf("runHealthCommand returned error: %v", err)
		}
	})

	var health struct {
		EphemerisStatus string            `json:"ephemeris_status"`
		Providers       map[string]string `json:"providers"`
	}
	if err := json.Unmarshal([]byte(output), &health); err != nil {
		t.Fatalf("failed to parse health JSON: %v\noutput: %s", err, output)
	}
	if health.EphemerisStatus == "healthy (demo mode)" {
		t.Fatal("health command still reports demo-mode ephemeris status")
	}
	if health.Providers["primary"] == "" || health.Providers["fallback"] == "" {
		t.Fatalf("expected real primary and fallback provider health, got %+v", health.Providers)
	}
}

func captureStdoutForTest(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	os.Stdout = writePipe
	defer func() {
		os.Stdout = originalStdout
	}()

	fn()
	if err := writePipe.Close(); err != nil {
		t.Fatalf("failed to close stdout pipe: %v", err)
	}

	output, err := io.ReadAll(readPipe)
	if err != nil {
		t.Fatalf("failed to read stdout pipe: %v", err)
	}
	if err := readPipe.Close(); err != nil {
		t.Fatalf("failed to close stdout reader: %v", err)
	}
	return string(output)
}
