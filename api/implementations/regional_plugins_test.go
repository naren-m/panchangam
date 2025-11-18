package implementations

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/api"
)

var (
	testDateRegional = time.Date(2024, 4, 14, 12, 0, 0, 0, time.UTC)
	testLocationTN   = api.Location{
		Latitude:  13.0827,
		Longitude: 80.2707,
		Timezone:  "Asia/Kolkata",
		Name:      "Chennai",
	}
	testLocationKerala = api.Location{
		Latitude:  10.8505,
		Longitude: 76.2711,
		Timezone:  "Asia/Kolkata",
		Name:      "Kochi",
	}
)

func TestTamilNaduRegionalPlugin(t *testing.T) {
	plugin := NewTamilNaduRegionalPlugin()

	// Test plugin info
	info := plugin.GetInfo()
	if info.Name != "tamil_nadu_regional_plugin" {
		t.Errorf("Expected plugin name 'tamil_nadu_regional_plugin', got %s", info.Name)
	}

	// Test initialization
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled after initialization")
	}

	// Test region and calendar system
	if plugin.GetRegion() != api.RegionTamilNadu {
		t.Errorf("Expected region Tamil Nadu, got %s", plugin.GetRegion())
	}

	if plugin.GetCalendarSystem() != api.CalendarAmanta {
		t.Errorf("Expected Amanta calendar system, got %s", plugin.GetCalendarSystem())
	}

	// Test applying regional rules
	testData := &api.PanchangamData{
		Date:           testDateRegional,
		Location:       testLocationTN,
		Region:         api.RegionTamilNadu,
		CalendarSystem: api.CalendarPurnimanta, // Will be changed to Amanta
		Events:         []api.Event{},
	}

	err = plugin.ApplyRegionalRules(context.Background(), testData)
	if err != nil {
		t.Fatalf("Failed to apply regional rules: %v", err)
	}

	if testData.CalendarSystem != api.CalendarAmanta {
		t.Errorf("Expected calendar system to be set to Amanta, got %s", testData.CalendarSystem)
	}

	// Test regional events - Tamil New Year on April 14
	events, err := plugin.GetRegionalEvents(context.Background(), testDateRegional, testLocationTN)
	if err != nil {
		t.Fatalf("Failed to get regional events: %v", err)
	}

	foundTamilNewYear := false
	for _, event := range events {
		if event.Name == "Tamil New Year" {
			foundTamilNewYear = true
			if event.NameLocal != "தமிழ் புத்தாண்டு" {
				t.Errorf("Expected Tamil name 'தமிழ் புத்தாண்டு', got %s", event.NameLocal)
			}
			if event.Region != api.RegionTamilNadu {
				t.Errorf("Expected region Tamil Nadu, got %s", event.Region)
			}
		}
	}

	if !foundTamilNewYear {
		t.Error("Tamil New Year event not found on April 14")
	}

	// Test Pongal - January 15
	pongalDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	pongalEvents, err := plugin.GetRegionalEvents(context.Background(), pongalDate, testLocationTN)
	if err != nil {
		t.Fatalf("Failed to get Pongal events: %v", err)
	}

	foundPongal := false
	for _, event := range pongalEvents {
		if event.Name == "Thai Pongal" {
			foundPongal = true
		}
	}

	if !foundPongal {
		t.Error("Thai Pongal event not found on January 15")
	}

	// Test regional names
	names := plugin.GetRegionalNames("ta")
	if names == nil || len(names) == 0 {
		t.Error("Expected regional names, got none")
	}

	if names["Sunday"] != "ஞாயிறு" {
		t.Errorf("Expected Sunday in Tamil to be 'ஞாயிறு', got %s", names["Sunday"])
	}

	// Test shutdown
	err = plugin.Shutdown(context.Background())
	if err != nil {
		t.Fatalf("Failed to shutdown plugin: %v", err)
	}

	if plugin.IsEnabled() {
		t.Error("Plugin should be disabled after shutdown")
	}
}

func TestKeralaRegionalPlugin(t *testing.T) {
	plugin := NewKeralaRegionalPlugin()

	// Test plugin info
	info := plugin.GetInfo()
	if info.Name != "kerala_regional_plugin" {
		t.Errorf("Expected plugin name 'kerala_regional_plugin', got %s", info.Name)
	}

	// Test initialization
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Test region and calendar system
	if plugin.GetRegion() != api.RegionKerala {
		t.Errorf("Expected region Kerala, got %s", plugin.GetRegion())
	}

	if plugin.GetCalendarSystem() != api.CalendarAmanta {
		t.Errorf("Expected Amanta calendar system, got %s", plugin.GetCalendarSystem())
	}

	// Test Vishu event on April 14/15
	events, err := plugin.GetRegionalEvents(context.Background(), testDateRegional, testLocationKerala)
	if err != nil {
		t.Fatalf("Failed to get regional events: %v", err)
	}

	foundVishu := false
	for _, event := range events {
		if event.Name == "Vishu" {
			foundVishu = true
			if event.NameLocal != "വിഷു" {
				t.Errorf("Expected Malayalam name 'വിഷു', got %s", event.NameLocal)
			}
		}
	}

	if !foundVishu {
		t.Error("Vishu event not found on April 14")
	}

	// Test regional names
	names := plugin.GetRegionalNames("ml")
	if names["Sunday"] != "ഞായർ" {
		t.Errorf("Expected Sunday in Malayalam to be 'ഞായർ', got %s", names["Sunday"])
	}
}

func TestBengalRegionalPlugin(t *testing.T) {
	plugin := NewBengalRegionalPlugin()

	// Test plugin info
	info := plugin.GetInfo()
	if info.Name != "bengal_regional_plugin" {
		t.Errorf("Expected plugin name 'bengal_regional_plugin', got %s", info.Name)
	}

	// Test initialization
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Test region and calendar system
	if plugin.GetRegion() != api.RegionBengal {
		t.Errorf("Expected region Bengal, got %s", plugin.GetRegion())
	}

	if plugin.GetCalendarSystem() != api.CalendarAmanta {
		t.Errorf("Expected Amanta calendar system, got %s", plugin.GetCalendarSystem())
	}

	// Test Pohela Boishakh
	events, err := plugin.GetRegionalEvents(context.Background(), testDateRegional, testLocationTN)
	if err != nil {
		t.Fatalf("Failed to get regional events: %v", err)
	}

	foundNewYear := false
	for _, event := range events {
		if event.Name == "Pohela Boishakh" {
			foundNewYear = true
			if event.NameLocal != "পহেলা বৈশাখ" {
				t.Errorf("Expected Bengali name 'পহেলা বৈশাখ', got %s", event.NameLocal)
			}
		}
	}

	if !foundNewYear {
		t.Error("Pohela Boishakh event not found on April 14")
	}

	// Test regional names
	names := plugin.GetRegionalNames("bn")
	if names["Sunday"] != "রবিবার" {
		t.Errorf("Expected Sunday in Bengali to be 'রবিবার', got %s", names["Sunday"])
	}
}

func TestGujaratRegionalPlugin(t *testing.T) {
	plugin := NewGujaratRegionalPlugin()

	// Test plugin info
	info := plugin.GetInfo()
	if info.Name != "gujarat_regional_plugin" {
		t.Errorf("Expected plugin name 'gujarat_regional_plugin', got %s", info.Name)
	}

	// Test initialization
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Test region and calendar system
	if plugin.GetRegion() != api.RegionGujarat {
		t.Errorf("Expected region Gujarat, got %s", plugin.GetRegion())
	}

	if plugin.GetCalendarSystem() != api.CalendarPurnimanta {
		t.Errorf("Expected Purnimanta calendar system, got %s", plugin.GetCalendarSystem())
	}

	// Test Uttarayan event on January 14
	uttarayanDate := time.Date(2024, 1, 14, 12, 0, 0, 0, time.UTC)
	events, err := plugin.GetRegionalEvents(context.Background(), uttarayanDate, testLocationTN)
	if err != nil {
		t.Fatalf("Failed to get regional events: %v", err)
	}

	foundUttarayan := false
	for _, event := range events {
		if event.Name == "Uttarayan" {
			foundUttarayan = true
			if event.NameLocal != "ઉત્તરાયણ" {
				t.Errorf("Expected Gujarati name 'ઉત્તરાયણ', got %s", event.NameLocal)
			}
		}
	}

	if !foundUttarayan {
		t.Error("Uttarayan event not found on January 14")
	}

	// Test regional names
	names := plugin.GetRegionalNames("gu")
	if names["Sunday"] != "રવિવાર" {
		t.Errorf("Expected Sunday in Gujarati to be 'રવિવાર', got %s", names["Sunday"])
	}
}

func TestMaharashtraRegionalPlugin(t *testing.T) {
	plugin := NewMaharashtraRegionalPlugin()

	// Test plugin info
	info := plugin.GetInfo()
	if info.Name != "maharashtra_regional_plugin" {
		t.Errorf("Expected plugin name 'maharashtra_regional_plugin', got %s", info.Name)
	}

	// Test initialization
	err := plugin.Initialize(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Test region and calendar system
	if plugin.GetRegion() != api.RegionMaha {
		t.Errorf("Expected region Maharashtra, got %s", plugin.GetRegion())
	}

	if plugin.GetCalendarSystem() != api.CalendarPurnimanta {
		t.Errorf("Expected Purnimanta calendar system, got %s", plugin.GetCalendarSystem())
	}

	// Test regional names
	names := plugin.GetRegionalNames("mr")
	if names["Sunday"] != "रविवार" {
		t.Errorf("Expected Sunday in Marathi to be 'रविवार', got %s", names["Sunday"])
	}
}

func TestAllRegionalPluginsCapabilities(t *testing.T) {
	plugins := []struct {
		name   string
		plugin api.RegionalExtension
		region api.Region
		calendar api.CalendarSystem
	}{
		{"Tamil Nadu", NewTamilNaduRegionalPlugin(), api.RegionTamilNadu, api.CalendarAmanta},
		{"Kerala", NewKeralaRegionalPlugin(), api.RegionKerala, api.CalendarAmanta},
		{"Bengal", NewBengalRegionalPlugin(), api.RegionBengal, api.CalendarAmanta},
		{"Gujarat", NewGujaratRegionalPlugin(), api.RegionGujarat, api.CalendarPurnimanta},
		{"Maharashtra", NewMaharashtraRegionalPlugin(), api.RegionMaha, api.CalendarPurnimanta},
	}

	for _, p := range plugins {
		t.Run(p.name, func(t *testing.T) {
			err := p.plugin.Initialize(context.Background(), nil)
			if err != nil {
				t.Fatalf("Failed to initialize %s plugin: %v", p.name, err)
			}

			if p.plugin.GetRegion() != p.region {
				t.Errorf("%s: Expected region %s, got %s", p.name, p.region, p.plugin.GetRegion())
			}

			if p.plugin.GetCalendarSystem() != p.calendar {
				t.Errorf("%s: Expected calendar %s, got %s", p.name, p.calendar, p.plugin.GetCalendarSystem())
			}

			names := p.plugin.GetRegionalNames("en")
			if names == nil {
				t.Errorf("%s: Expected regional names, got nil", p.name)
			}
		})
	}
}
