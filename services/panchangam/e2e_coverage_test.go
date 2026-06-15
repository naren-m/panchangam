package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

// TestEndToEndFeatureCoverage provides end-to-end validation of all documented features
// This validates complete user workflows from request to response
func TestEndToEndFeatureCoverage(t *testing.T) {
	// Initialize observability for testing
	observability.NewLocalObserver()

	t.Run("E2E_Complete_Panchangam_Request", func(t *testing.T) {
		// Test complete end-to-end workflow for a Panchangam request
		testCompleteWorkflow(t)
	})

	t.Run("E2E_Multiple_Locations", func(t *testing.T) {
		// Test E2E workflow across different geographical locations
		testMultipleLocations(t)
	})

	t.Run("E2E_Multiple_Dates", func(t *testing.T) {
		// Test E2E workflow across different dates and seasons
		testMultipleDates(t)
	})

	t.Run("E2E_Feature_Consistency", func(t *testing.T) {
		// Test that all features work consistently together
		testFeatureConsistency(t)
	})

	t.Run("E2E_User_Scenarios", func(t *testing.T) {
		// Test real-world user scenarios
		testUserScenarios(t)
	})
}

// testCompleteWorkflow validates the complete end-to-end workflow
func testCompleteWorkflow(t *testing.T) {
	server := NewPanchangamServer()
	ctx := context.Background()

	// Step 1: User makes request for Panchangam data
	req := &ppb.GetPanchangamRequest{
		Date:              "2024-01-15",
		Latitude:          12.9716, // Bangalore coordinates
		Longitude:         77.5946,
		Timezone:          "Asia/Kolkata",
		Region:            "India",
		CalculationMethod: "traditional",
		Locale:            "en",
	}

	// Step 2: Service processes request
	start := time.Now()
	resp, err := server.Get(ctx, req)
	duration := time.Since(start)

	// Step 3: Validate response
	require.NoError(t, err, "E2E: Request should succeed")
	require.NotNil(t, resp, "E2E: Response should not be nil")
	require.NotNil(t, resp.PanchangamData, "E2E: Panchangam data should be present")

	data := resp.PanchangamData

	// Step 4: Validate all core features are present
	validateCoreFeatures(t, data, req)

	// Step 5: Validate performance
	assert.True(t, duration < 1*time.Second, "E2E: Response should be fast (<1s), got %v", duration)

	// Step 6: Validate observability
	validateObservability(t)

	t.Logf("E2E: Complete workflow validated in %v", duration)
	t.Logf("Request: %s at (%.4f, %.4f)", req.Date, req.Latitude, req.Longitude)
	t.Logf("Response: All 5 Panchangam elements + astronomy data")
}

// validateCoreFeatures validates all core Panchangam features are present
func validateCoreFeatures(t *testing.T, data *ppb.PanchangamData, req *ppb.GetPanchangamRequest) {
	t.Helper()

	// Validate TITHI_001: Lunar day calculation
	assert.NotEmpty(t, data.Tithi, "E2E: TITHI_001 should be present")
	assert.True(t, len(data.Tithi) > 0, "E2E: Tithi should have content")

	// Validate NAKSHATRA_001: Lunar mansion calculation
	assert.NotEmpty(t, data.Nakshatra, "E2E: NAKSHATRA_001 should be present")
	assert.True(t, len(data.Nakshatra) > 0, "E2E: Nakshatra should have content")

	// Validate YOGA_001: Auspicious combinations
	assert.NotEmpty(t, data.Yoga, "E2E: YOGA_001 should be present")
	assert.True(t, len(data.Yoga) > 0, "E2E: Yoga should have content")

	// Validate KARANA_001: Half-Tithi divisions
	assert.NotEmpty(t, data.Karana, "E2E: KARANA_001 should be present")
	assert.True(t, len(data.Karana) > 0, "E2E: Karana should have content")

	// Note: VARA_001 is not directly exposed in current service, but would be included in events

	// Validate ASTRONOMY_001: Sunrise/sunset calculations
	assert.NotEmpty(t, data.SunriseTime, "E2E: ASTRONOMY_001 sunrise should be present")
	assert.NotEmpty(t, data.SunsetTime, "E2E: ASTRONOMY_001 sunset should be present")

	// Validate time format
	_, err := time.Parse("15:04:05", data.SunriseTime)
	assert.NoError(t, err, "E2E: Sunrise time should be valid format")
	_, err = time.Parse("15:04:05", data.SunsetTime)
	assert.NoError(t, err, "E2E: Sunset time should be valid format")

	// Validate SERVICE_001: Service layer
	assert.Equal(t, req.Date, data.Date, "E2E: SERVICE_001 date handling")
	assert.NotNil(t, data.Events, "E2E: SERVICE_001 events should be included")

	t.Logf("E2E: All core features validated")
	t.Logf("Tithi: %s", data.Tithi)
	t.Logf("Nakshatra: %s", data.Nakshatra)
	t.Logf("Yoga: %s", data.Yoga)
	t.Logf("Karana: %s", data.Karana)
	t.Logf("Sunrise: %s", data.SunriseTime)
	t.Logf("Sunset: %s", data.SunsetTime)
	t.Logf("Events: %d", len(data.Events))
}

// validateObservability validates OBSERVABILITY_001: OpenTelemetry integration
func validateObservability(t *testing.T) {
	t.Helper()

	// Test that observability is working
	ctx := context.Background()
	observer := observability.Observer()
	require.NotNil(t, observer, "E2E: Observer should be available")

	// Test span creation
	_, span := observer.CreateSpan(ctx, "e2e_test_span")
	assert.NotNil(t, span, "E2E: Span should be created")

	// Test span attributes
	span.SetAttributes(attribute.String("e2e_test", "validation"))

	// Test span events
	span.AddEvent("e2e_validation")

	// Test span completion
	span.End()

	t.Logf("E2E: Observability validated")
}
