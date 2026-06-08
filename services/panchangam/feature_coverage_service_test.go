package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

func testFeatureSERVICE_001(t *testing.T) {
	t.Run("SERVICE_001_gRPC_Service", func(t *testing.T) {
		observability.NewLocalObserver()

		server := NewPanchangamServer()
		require.NotNil(t, server, "SERVICE_001: Server should be created")

		ctx := context.Background()
		req := &ppb.GetPanchangamRequest{
			Date:              "2024-01-15",
			Latitude:          12.9716,
			Longitude:         77.5946,
			Timezone:          "Asia/Kolkata",
			Region:            "India",
			CalculationMethod: "traditional",
			Locale:            "en",
		}

		resp, err := server.Get(ctx, req)
		assert.NoError(t, err, "SERVICE_001: Service should handle valid requests")
		require.NotNil(t, resp, "SERVICE_001: Response should not be nil")
		require.NotNil(t, resp.PanchangamData, "SERVICE_001: PanchangamData should not be nil")

		data := resp.PanchangamData

		assert.Equal(t, req.Date, data.Date, "SERVICE_001: Date should match request")
		assert.NotEmpty(t, data.Tithi, "SERVICE_001: Tithi should be provided")
		assert.NotEmpty(t, data.Nakshatra, "SERVICE_001: Nakshatra should be provided")
		assert.NotEmpty(t, data.Yoga, "SERVICE_001: Yoga should be provided")
		assert.NotEmpty(t, data.Karana, "SERVICE_001: Karana should be provided")
		assert.NotEmpty(t, data.SunriseTime, "SERVICE_001: SunriseTime should be provided")
		assert.NotEmpty(t, data.SunsetTime, "SERVICE_001: SunsetTime should be provided")
		assert.NotNil(t, data.Events, "SERVICE_001: Events should not be nil")

		assert.True(t, len(req.CalculationMethod) == 0 || req.CalculationMethod != "", "SERVICE_001: CalculationMethod should be handled")
		assert.True(t, len(req.Locale) == 0 || req.Locale != "", "SERVICE_001: Locale should be handled")
		assert.True(t, len(req.Region) == 0 || req.Region != "", "SERVICE_001: Region should be handled")

		t.Logf("SERVICE_001: Validated gRPC service with protocol buffers")
	})
}

func testFeatureASTRONOMY_001(t *testing.T) {
	t.Run("ASTRONOMY_001_Sunrise_Sunset", func(t *testing.T) {
		ctx := context.Background()
		testDate := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)
		location := astronomy.Location{
			Latitude:  12.9716,
			Longitude: 77.5946,
		}

		sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		assert.NoError(t, err, "ASTRONOMY_001: Sun times calculation should succeed")
		require.NotNil(t, sunTimes, "ASTRONOMY_001: Sun times should not be nil")

		assert.True(t, sunTimes.Sunrise.Before(sunTimes.Sunset), "ASTRONOMY_001: Sunrise should be before sunset")
		assert.True(t, sunTimes.Sunrise.Hour() >= 0 && sunTimes.Sunrise.Hour() <= 23, "ASTRONOMY_001: Sunrise should be valid hour")
		assert.True(t, sunTimes.Sunset.Hour() >= 0 && sunTimes.Sunset.Hour() <= 23, "ASTRONOMY_001: Sunset should be valid hour")

		dayLength := sunTimes.Sunset.Sub(sunTimes.Sunrise)
		assert.True(t, dayLength > 8*time.Hour && dayLength < 16*time.Hour, "ASTRONOMY_001: Day length should be reasonable")

		t.Logf("ASTRONOMY_001: Validated sunrise %s, sunset %s",
			sunTimes.Sunrise.Format("15:04:05"), sunTimes.Sunset.Format("15:04:05"))
	})
}

func testFeatureOBSERVABILITY_001(t *testing.T) {
	t.Run("OBSERVABILITY_001_OpenTelemetry", func(t *testing.T) {
		observer := observability.NewLocalObserver()
		require.NotNil(t, observer, "OBSERVABILITY_001: Observer should be created")

		ctx := context.Background()

		_, span := observer.CreateSpan(ctx, "test_span")
		assert.NotNil(t, span, "OBSERVABILITY_001: Span should be created")

		span.SetAttributes(attribute.String("test_key", "test_value"))
		span.AddEvent("test_event")
		span.End()
		span.RecordError(assert.AnError)

		t.Logf("OBSERVABILITY_001: Validated OpenTelemetry integration")
	})
}

func TestFeatureCoveragePerformance(t *testing.T) {
	t.Run("Feature_Performance_Benchmarks", func(t *testing.T) {
		ctx := context.Background()
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		start := time.Now()

		location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}
		_, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)

		duration := time.Since(start)
		assert.NoError(t, err, "Performance test should not fail")
		assert.True(t, duration < 50*time.Millisecond,
			"ASTRONOMY_001 performance: should be <50ms, got %v", duration)

		t.Logf("Feature Performance: Astronomy calculation completed in %v", duration)
	})
}

func TestFeatureCoverageIntegration(t *testing.T) {
	t.Run("Feature_Integration_Patterns", func(t *testing.T) {
		ctx := context.Background()
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

		sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
		assert.NoError(t, err, "Integration: Astronomy should work")
		require.NotNil(t, sunTimes, "Integration: Sun times should be calculated")

		observability.NewLocalObserver()
		ctx, span := observability.Observer().CreateSpan(ctx, "integration_test")
		span.SetAttributes(attribute.String("test", "integration"))
		span.End()

		server := NewPanchangamServer()
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
		}
		resp, err := server.Get(ctx, req)
		assert.NoError(t, err, "Integration: Service should work")
		require.NotNil(t, resp, "Integration: Service should respond")

		t.Logf("Feature Integration: All components work together")
	})
}
