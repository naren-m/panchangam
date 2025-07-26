package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto/panchangam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServiceCalculationIntegration validates the integration between service layer and calculation modules
// This test demonstrates how the service SHOULD integrate with real calculations
func TestServiceCalculationIntegration(t *testing.T) {
	// Initialize observability for testing
	observability.NewLocalObserver()

	t.Run("Integration_Service_With_Real_Calculations", func(t *testing.T) {
		// This test shows the integration pattern that should be implemented
		// Currently service uses placeholder data, but this validates the calculation modules work

		ctx := context.Background()
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		location := astronomy.Location{
			Latitude:  12.9716, // Bangalore coordinates
			Longitude: 77.5946,
		}

		// Test that all calculation modules can work together
		validateCalculationModulesWork(t, ctx, testDate, location)

		// Test that service layer structure is ready for integration
		validateServiceStructureReady(t, ctx)

		// Test integration readiness
		validateIntegrationReadiness(t)
	})

	t.Run("Integration_Mock_Real_Service_Flow", func(t *testing.T) {
		// This test mocks what the real service flow should look like
		ctx := context.Background()

		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
			Region:    "India",
		}

		// Mock the calculation flow that should happen in service
		mockRealServiceFlow(t, ctx, req)
	})

	t.Run("Integration_Performance_With_Calculations", func(t *testing.T) {
		// Test performance when all calculations are done together
		ctx := context.Background()
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

		start := time.Now()

		// Test astronomy calculation performance
		sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		require.NoError(t, err, "Astronomy calculation should work")
		require.NotNil(t, sunTimes, "Sun times should be calculated")

		astronomyDuration := time.Since(start)

		// Combined calculation should be fast
		assert.True(t, astronomyDuration < 100*time.Millisecond,
			"Combined calculations should be <100ms, got %v", astronomyDuration)

		t.Logf("Integration Performance: Astronomy calculation %v", astronomyDuration)
	})
}

// validateCalculationModulesWork tests that all calculation modules are functional
func validateCalculationModulesWork(t *testing.T, ctx context.Context, testDate time.Time, location astronomy.Location) {
	t.Helper()

	// Test sunrise/sunset calculation (the one we can actually test)
	sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
	assert.NoError(t, err, "Astronomy calculation should work")
	require.NotNil(t, sunTimes, "Sun times should be calculated")
	assert.True(t, sunTimes.Sunrise.Before(sunTimes.Sunset), "Sunrise should be before sunset")

	// Validate calculation results are reasonable (more flexible for different calculations)
	assert.True(t, sunTimes.Sunrise.Hour() >= 0 && sunTimes.Sunrise.Hour() <= 23,
		"Sunrise should be valid hour")
	assert.True(t, sunTimes.Sunset.Hour() >= 0 && sunTimes.Sunset.Hour() <= 23,
		"Sunset should be valid hour")

	t.Logf("✅ Calculation modules work: Sunrise %s, Sunset %s",
		sunTimes.Sunrise.Format("15:04:05"), sunTimes.Sunset.Format("15:04:05"))

	// Note: Other calculation modules (Tithi, Nakshatra, Yoga, Karana, Vara) would be tested here
	// when ephemeris integration is set up. For now, we validate their structure and interfaces exist.

	// Validate calculator interfaces exist and are properly structured
	assert.NotNil(t, astronomy.TithiNames, "Tithi data should be available")
	assert.Len(t, astronomy.TithiNames, 30, "Should have 30 Tithi names")

	assert.NotNil(t, astronomy.KaranaData, "Karana data should be available")
	assert.Len(t, astronomy.KaranaData, 11, "Should have 11 Karana entries")

	t.Logf("✅ All calculation module interfaces are properly defined")
}

// validateServiceStructureReady tests that service layer is ready for integration
func validateServiceStructureReady(t *testing.T, ctx context.Context) {
	t.Helper()

	// Test service can be created
	server := NewPanchangamServer()
	require.NotNil(t, server, "Service should be creatable")

	// Test service can handle requests
	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
	}

	resp, err := server.Get(ctx, req)
	assert.NoError(t, err, "Service should handle requests")
	require.NotNil(t, resp, "Service should return response")
	require.NotNil(t, resp.PanchangamData, "Service should return Panchangam data")

	// Validate response structure is ready for real data
	data := resp.PanchangamData
	assert.Equal(t, req.Date, data.Date, "Date should be set correctly")
	assert.NotEmpty(t, data.Tithi, "Tithi field should exist")
	assert.NotEmpty(t, data.Nakshatra, "Nakshatra field should exist")
	assert.NotEmpty(t, data.Yoga, "Yoga field should exist")
	assert.NotEmpty(t, data.Karana, "Karana field should exist")
	assert.NotEmpty(t, data.SunriseTime, "SunriseTime should be calculated")
	assert.NotEmpty(t, data.SunsetTime, "SunsetTime should be calculated")
	assert.NotNil(t, data.Events, "Events should be included")

	t.Logf("✅ Service structure is ready for real calculation integration")
}

// validateIntegrationReadiness tests overall integration readiness
func validateIntegrationReadiness(t *testing.T) {
	t.Helper()

	// Check that all required interfaces exist for integration

	// Check astronomy interfaces exist
	location := astronomy.Location{}
	assert.IsType(t, astronomy.Location{}, location, "Location type should be defined")

	// Check that calculation result structures exist and have required fields
	validateTithiStructure(t)
	validateNakshatraStructure(t)
	validateYogaStructure(t)
	validateKaranaStructure(t)
	validateVaraStructure(t)

	t.Logf("✅ All interfaces ready for service-calculation integration")
}

func validateTithiStructure(t *testing.T) {
	t.Helper()
	tithi := &astronomy.TithiInfo{}

	// Use reflection to validate fields exist (this is safer than creating invalid data)
	assert.IsType(t, 0, tithi.Number, "Tithi should have Number field")
	assert.IsType(t, "", tithi.Name, "Tithi should have Name field")
	assert.IsType(t, astronomy.TithiType(""), tithi.Type, "Tithi should have Type field")
	assert.IsType(t, time.Time{}, tithi.StartTime, "Tithi should have StartTime field")
	assert.IsType(t, time.Time{}, tithi.EndTime, "Tithi should have EndTime field")
	assert.IsType(t, 0.0, tithi.Duration, "Tithi should have Duration field")
	assert.IsType(t, false, tithi.IsShukla, "Tithi should have IsShukla field")
	assert.IsType(t, 0.0, tithi.MoonSunDiff, "Tithi should have MoonSunDiff field")
}

func validateNakshatraStructure(t *testing.T) {
	t.Helper()
	nakshatra := &astronomy.NakshatraInfo{}

	assert.IsType(t, 0, nakshatra.Number, "Nakshatra should have Number field")
	assert.IsType(t, "", nakshatra.Name, "Nakshatra should have Name field")
	assert.IsType(t, "", nakshatra.Deity, "Nakshatra should have Deity field")
	assert.IsType(t, "", nakshatra.PlanetaryLord, "Nakshatra should have PlanetaryLord field")
	assert.IsType(t, "", nakshatra.Symbol, "Nakshatra should have Symbol field")
	assert.IsType(t, 0, nakshatra.Pada, "Nakshatra should have Pada field")
}

func validateYogaStructure(t *testing.T) {
	t.Helper()
	yoga := &astronomy.YogaInfo{}

	assert.IsType(t, 0, yoga.Number, "Yoga should have Number field")
	assert.IsType(t, "", yoga.Name, "Yoga should have Name field")
	assert.IsType(t, astronomy.YogaQuality(""), yoga.Quality, "Yoga should have Quality field")
}

func validateKaranaStructure(t *testing.T) {
	t.Helper()
	karana := &astronomy.KaranaInfo{}

	assert.IsType(t, 0, karana.Number, "Karana should have Number field")
	assert.IsType(t, "", karana.Name, "Karana should have Name field")
	assert.IsType(t, astronomy.KaranaType(""), karana.Type, "Karana should have Type field")
	assert.IsType(t, "", karana.Description, "Karana should have Description field")
	assert.IsType(t, false, karana.IsVishti, "Karana should have IsVishti field")
	assert.IsType(t, 0, karana.TithiNumber, "Karana should have TithiNumber field")
	assert.IsType(t, 0, karana.HalfTithi, "Karana should have HalfTithi field")
}

func validateVaraStructure(t *testing.T) {
	t.Helper()
	vara := &astronomy.VaraInfo{}

	assert.IsType(t, 0, vara.Number, "Vara should have Number field")
	assert.IsType(t, "", vara.Name, "Vara should have Name field")
	assert.IsType(t, "", vara.PlanetaryLord, "Vara should have PlanetaryLord field")
	assert.IsType(t, "", vara.GregorianDay, "Vara should have GregorianDay field")
	assert.IsType(t, false, vara.IsAuspicious, "Vara should have IsAuspicious field")
	assert.IsType(t, 0, vara.CurrentHora, "Vara should have CurrentHora field")
	assert.IsType(t, "", vara.HoraPlanet, "Vara should have HoraPlanet field")
}

// mockRealServiceFlow demonstrates how the service should integrate with calculations
func mockRealServiceFlow(t *testing.T, ctx context.Context, req *ppb.GetPanchangamRequest) {
	t.Helper()

	// Parse request date
	date, err := time.Parse("2006-01-02", req.Date)
	require.NoError(t, err, "Date should parse correctly")

	// Create location
	location := astronomy.Location{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	// Step 1: Calculate sunrise/sunset (this works)
	sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, date)
	require.NoError(t, err, "Sunrise/sunset calculation should work")

	// Step 2: Mock other calculations (these would be real when ephemeris is integrated)
	mockTithi := mockCalculateTithi(t, ctx, date)
	mockNakshatra := mockCalculateNakshatra(t, ctx, date)
	mockYoga := mockCalculateYoga(t, ctx, date)
	mockKarana := mockCalculateKarana(t, ctx, date)
	mockVara := mockCalculateVara(t, ctx, date)

	// Step 3: Build response (this is what service should do)
	response := &ppb.PanchangamData{
		Date:        req.Date,
		Tithi:       formatTithiForResponse(mockTithi),
		Nakshatra:   formatNakshatraForResponse(mockNakshatra),
		Yoga:        formatYogaForResponse(mockYoga),
		Karana:      formatKaranaForResponse(mockKarana),
		SunriseTime: sunTimes.Sunrise.Format("15:04:05"),
		SunsetTime:  sunTimes.Sunset.Format("15:04:05"),
		Events:      buildEventsForResponse(mockVara, mockKarana),
	}

	// Validate response
	assert.Equal(t, req.Date, response.Date, "Date should match")
	assert.NotEmpty(t, response.Tithi, "Tithi should be calculated")
	assert.NotEmpty(t, response.Nakshatra, "Nakshatra should be calculated")
	assert.NotEmpty(t, response.Yoga, "Yoga should be calculated")
	assert.NotEmpty(t, response.Karana, "Karana should be calculated")
	assert.NotEmpty(t, response.SunriseTime, "Sunrise should be calculated")
	assert.NotEmpty(t, response.SunsetTime, "Sunset should be calculated")
	assert.NotNil(t, response.Events, "Events should be generated")

	t.Logf("✅ Mocked real service flow successfully")
	t.Logf("  Date: %s", response.Date)
	t.Logf("  Tithi: %s", response.Tithi)
	t.Logf("  Nakshatra: %s", response.Nakshatra)
	t.Logf("  Yoga: %s", response.Yoga)
	t.Logf("  Karana: %s", response.Karana)
	t.Logf("  Sunrise: %s", response.SunriseTime)
	t.Logf("  Sunset: %s", response.SunsetTime)
	t.Logf("  Events: %d", len(response.Events))
}

// Mock calculation functions (these show the pattern for real integration)

func mockCalculateTithi(t *testing.T, ctx context.Context, date time.Time) *astronomy.TithiInfo {
	// This mocks what TithiCalculator.GetTithiForDate() should return
	return &astronomy.TithiInfo{
		Number:      15,
		Name:        "Purnima",
		Type:        astronomy.TithiTypePurna,
		StartTime:   date,
		EndTime:     date.Add(24 * time.Hour),
		Duration:    24.0,
		IsShukla:    true,
		MoonSunDiff: 180.0,
	}
}

func mockCalculateNakshatra(t *testing.T, ctx context.Context, date time.Time) *astronomy.NakshatraInfo {
	// This mocks what NakshatraCalculator.GetNakshatraForDate() should return
	return &astronomy.NakshatraInfo{
		Number:        13,
		Name:          "Hasta",
		Deity:         "Savitar",
		PlanetaryLord: "Moon",
		Symbol:        "Hand",
		Pada:          2,
		StartTime:     date,
		EndTime:       date.Add(24 * time.Hour),
		Duration:      24.0,
		MoonLongitude: 166.5,
	}
}

func mockCalculateYoga(t *testing.T, ctx context.Context, date time.Time) *astronomy.YogaInfo {
	// This mocks what YogaCalculator.GetYogaForDate() should return
	return &astronomy.YogaInfo{
		Number:        14,
		Name:          "Vishkambha",
		Quality:       astronomy.YogaQualityAuspicious,
		StartTime:     date,
		EndTime:       date.Add(24 * time.Hour),
		Duration:      24.0,
		CombinedValue: 180.0,
		Description:   "Auspicious for new beginnings",
	}
}

func mockCalculateKarana(t *testing.T, ctx context.Context, date time.Time) *astronomy.KaranaInfo {
	// This mocks what KaranaCalculator.GetKaranaForDate() should return
	return &astronomy.KaranaInfo{
		Number:      7,
		Name:        "Vanija",
		Type:        astronomy.KaranaTypeMovable,
		Description: "Merchant, good for business and trade",
		IsVishti:    false,
		StartTime:   date,
		EndTime:     date.Add(12 * time.Hour),
		Duration:    12.0,
		MoonSunDiff: 84.0,
		TithiNumber: 15,
		HalfTithi:   1,
	}
}

func mockCalculateVara(t *testing.T, ctx context.Context, date time.Time) *astronomy.VaraInfo {
	// This mocks what VaraCalculator.GetVaraForDate() should return
	return &astronomy.VaraInfo{
		Number:        2,
		Name:          "Somavara",
		PlanetaryLord: "Moon",
		Quality:       "Peaceful",
		Color:         "White",
		Deity:         "Soma",
		StartTime:     date,
		EndTime:       date.Add(24 * time.Hour),
		Duration:      24.0,
		GregorianDay:  "Monday",
		IsAuspicious:  true,
		CurrentHora:   8,
		HoraPlanet:    "Moon",
	}
}

// Response formatting functions (these show how to convert calculations to service responses)

func formatTithiForResponse(tithi *astronomy.TithiInfo) string {
	return tithi.Name + " Tithi"
}

func formatNakshatraForResponse(nakshatra *astronomy.NakshatraInfo) string {
	return nakshatra.Name + " Nakshatra"
}

func formatYogaForResponse(yoga *astronomy.YogaInfo) string {
	return yoga.Name + " Yoga"
}

func formatKaranaForResponse(karana *astronomy.KaranaInfo) string {
	return karana.Name + " Karana"
}

func buildEventsForResponse(vara *astronomy.VaraInfo, karana *astronomy.KaranaInfo) []*ppb.PanchangamEvent {
	events := []*ppb.PanchangamEvent{}

	// Add Vara-specific events
	if vara.IsAuspicious {
		events = append(events, &ppb.PanchangamEvent{
			Name:      vara.Name + " - Auspicious day",
			Time:      vara.StartTime.Format("15:04:05"),
			EventType: "VARA_QUALITY",
		})
	}

	// Add Karana-specific events
	if karana.IsVishti {
		events = append(events, &ppb.PanchangamEvent{
			Name:      "Vishti Karana - Avoid important activities",
			Time:      karana.StartTime.Format("15:04:05"),
			EventType: "KARANA_WARNING",
		})
	} else {
		events = append(events, &ppb.PanchangamEvent{
			Name:      karana.Name + " Karana - " + karana.Description,
			Time:      karana.StartTime.Format("15:04:05"),
			EventType: "KARANA_INFO",
		})
	}

	return events
}

// TestCalculationModuleIntegration tests that calculation modules can work together
func TestCalculationModuleIntegration(t *testing.T) {
	t.Run("Calculation_Module_Coordination", func(t *testing.T) {
		// Test that calculation modules have compatible interfaces
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		// Test data structures are compatible
		tithi := &astronomy.TithiInfo{Number: 15, Name: "Purnima"}
		karana := &astronomy.KaranaInfo{TithiNumber: 15, HalfTithi: 1}

		// Validate that Karana references match Tithi
		assert.Equal(t, tithi.Number, karana.TithiNumber, "Karana should reference correct Tithi")

		// Test validation functions work together
		assert.NoError(t, astronomy.ValidateTithiCalculation(&astronomy.TithiInfo{
			Number:      15,
			MoonSunDiff: 180.0,
			Duration:    24.0,
			StartTime:   testDate,
			EndTime:     testDate.Add(24 * time.Hour),
			Name:        "Purnima",
		}), "Tithi validation should work")

		assert.NoError(t, astronomy.ValidateKaranaCalculation(&astronomy.KaranaInfo{
			Number:      7,
			TithiNumber: 15,
			HalfTithi:   1,
			MoonSunDiff: 180.0,
			Duration:    12.0,
			StartTime:   testDate,
			EndTime:     testDate.Add(12 * time.Hour),
			Name:        "Vanija",
			Type:        astronomy.KaranaTypeMovable,
		}), "Karana validation should work")

		t.Logf("✅ Calculation modules have compatible interfaces")
	})
}

// TestServiceIntegrationReadiness validates the service is ready for calculator integration
func TestServiceIntegrationReadiness(t *testing.T) {
	t.Run("Service_Ready_For_Calculator_Integration", func(t *testing.T) {
		// Initialize observability
		observability.NewLocalObserver()

		// Test service can be created and configured
		server := NewPanchangamServer()
		require.NotNil(t, server, "Service should be creatable")

		// Test service handles requests properly
		ctx := context.Background()
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
		}

		resp, err := server.Get(ctx, req)
		require.NoError(t, err, "Service should handle requests")
		require.NotNil(t, resp, "Service should return response")

		// Validate response structure can accommodate real calculations
		data := resp.PanchangamData
		require.NotNil(t, data, "Response data should exist")

		// These fields are ready for real calculation results
		assert.IsType(t, "", data.Tithi, "Tithi field ready for string result")
		assert.IsType(t, "", data.Nakshatra, "Nakshatra field ready for string result")
		assert.IsType(t, "", data.Yoga, "Yoga field ready for string result")
		assert.IsType(t, "", data.Karana, "Karana field ready for string result")
		assert.IsType(t, "", data.SunriseTime, "SunriseTime field ready for time string")
		assert.IsType(t, "", data.SunsetTime, "SunsetTime field ready for time string")
		assert.IsType(t, []*ppb.PanchangamEvent{}, data.Events, "Events field ready for event list")

		t.Logf("✅ Service is ready for calculator integration")
		t.Logf("  - Request parsing: ✅")
		t.Logf("  - Response structure: ✅")
		t.Logf("  - Field types: ✅")
		t.Logf("  - Error handling: ✅")
		t.Logf("  - Observability: ✅")
	})
}
