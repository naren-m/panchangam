package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestWriteSunriseSimpleError(t *testing.T) {
	var out bytes.Buffer

	if err := writeSunriseSimpleError(&out, errors.New("invalid date")); err != nil {
		t.Fatalf("writeSunriseSimpleError returned error: %v", err)
	}

	want := "invalid date\n"
	if got := out.String(); got != want {
		t.Fatalf("writeSunriseSimpleError() output = %q, want %q", got, want)
	}
}

func TestResolveSunriseSimpleLocationUsesPreset(t *testing.T) {
	name, lat, lon, err := resolveSunriseSimpleLocation(" tokyo ", 1, 2)
	if err != nil {
		t.Fatalf("resolveSunriseSimpleLocation returned error: %v", err)
	}
	if name != "Tokyo" {
		t.Fatalf("expected Tokyo preset name, got %q", name)
	}
	if lat != 35.6762 || lon != 139.6503 {
		t.Fatalf("expected Tokyo coordinates, got lat=%f lon=%f", lat, lon)
	}
}

func TestResolveSunriseSimpleLocationPreservesCustomCoordinates(t *testing.T) {
	name, lat, lon, err := resolveSunriseSimpleLocation("", 12.3, 45.6)
	if err != nil {
		t.Fatalf("resolveSunriseSimpleLocation returned error: %v", err)
	}
	if name != "" {
		t.Fatalf("expected no preset name, got %q", name)
	}
	if lat != 12.3 || lon != 45.6 {
		t.Fatalf("expected custom coordinates, got lat=%f lon=%f", lat, lon)
	}
}

func TestResolveSunriseSimpleLocationReportsAvailablePresets(t *testing.T) {
	_, _, _, err := resolveSunriseSimpleLocation("unknown", 1, 2)
	if err == nil {
		t.Fatal("expected unknown location to return an error")
	}
	if !strings.Contains(err.Error(), "Available: nyc, london, tokyo, sydney, mumbai, capetown") {
		t.Fatalf("expected available preset list in error, got %v", err)
	}
}

func TestParseSunriseSimpleDateTrimsInput(t *testing.T) {
	date, err := parseSunriseSimpleDate(" 2024-06-21 ")
	if err != nil {
		t.Fatalf("parseSunriseSimpleDate returned error: %v", err)
	}
	if date.Year() != 2024 || date.Month() != time.June || date.Day() != 21 {
		t.Fatalf("expected 2024-06-21, got %s", date.Format(time.RFC3339))
	}
}
