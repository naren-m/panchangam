package panchangam

import (
	"context"
	"testing"
	"time"

	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndToEndErrorHandling(t *testing.T) {
	server := NewPanchangamServer()
	ctx := context.Background()

	t.Run("E2E_Invalid_Requests", func(t *testing.T) {
		invalidRequests := []struct {
			name    string
			request *ppb.GetPanchangamRequest
		}{
			{
				name: "Invalid_Date",
				request: &ppb.GetPanchangamRequest{
					Date:      "invalid-date",
					Latitude:  12.9716,
					Longitude: 77.5946,
				},
			},
			{
				name: "Invalid_Latitude",
				request: &ppb.GetPanchangamRequest{
					Date:      "2024-01-15",
					Latitude:  91.0,
					Longitude: 77.5946,
				},
			},
			{
				name: "Invalid_Longitude",
				request: &ppb.GetPanchangamRequest{
					Date:      "2024-01-15",
					Latitude:  12.9716,
					Longitude: 181.0,
				},
			},
		}

		for _, test := range invalidRequests {
			t.Run("E2E_Error_"+test.name, func(t *testing.T) {
				resp, err := server.Get(ctx, test.request)

				assert.Error(t, err, "E2E: Invalid request should return error")
				assert.Nil(t, resp, "E2E: Invalid request should not return response")

				t.Logf("E2E: Error scenario %s handled correctly", test.name)
			})
		}
	})
}

func TestEndToEndPerformance(t *testing.T) {
	server := NewPanchangamServer()
	ctx := context.Background()

	t.Run("E2E_Performance_Targets", func(t *testing.T) {
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
		}

		_, _ = server.Get(ctx, req)

		measurements := []time.Duration{}
		successCount := 0

		for i := 0; i < 10; i++ {
			start := time.Now()
			resp, err := server.Get(ctx, req)
			duration := time.Since(start)

			if err == nil && resp != nil {
				measurements = append(measurements, duration)
				successCount++
			}
		}

		assert.True(t, len(measurements) >= 1, "E2E: Should have performance measurements")

		if len(measurements) == 0 {
			return
		}

		var total time.Duration
		for _, d := range measurements {
			total += d
		}
		average := total / time.Duration(len(measurements))

		assert.True(t, average < 1*time.Second, "E2E: Average response time should be <1s, got %v", average)

		t.Logf("E2E: Performance validated")
		t.Logf("Successful requests: %d/10", successCount)
		t.Logf("Average response time: %v", average)
		t.Logf("Performance target: <1s")
	})
}

func TestEndToEndDataQuality(t *testing.T) {
	server := NewPanchangamServer()
	ctx := context.Background()

	t.Run("E2E_Data_Quality", func(t *testing.T) {
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
		}

		resp, err := server.Get(ctx, req)
		require.NoError(t, err, "E2E: Data quality request should succeed")
		require.NotNil(t, resp, "E2E: Response should not be nil")
		require.NotNil(t, resp.PanchangamData, "E2E: Data should not be nil")

		data := resp.PanchangamData

		assert.Equal(t, req.Date, data.Date, "E2E: Date should match request")
		assert.True(t, len(data.Tithi) > 0, "E2E: Tithi should have content")
		assert.True(t, len(data.Nakshatra) > 0, "E2E: Nakshatra should have content")
		assert.True(t, len(data.Yoga) > 0, "E2E: Yoga should have content")
		assert.True(t, len(data.Karana) > 0, "E2E: Karana should have content")

		_, err = time.Parse("15:04:05", data.SunriseTime)
		assert.NoError(t, err, "E2E: Sunrise time should be valid format")
		_, err = time.Parse("15:04:05", data.SunsetTime)
		assert.NoError(t, err, "E2E: Sunset time should be valid format")

		assert.NotNil(t, data.Events, "E2E: Events should not be nil")
		for i, event := range data.Events {
			assert.NotEmpty(t, event.Name, "E2E: Event %d should have name", i)
			assert.NotEmpty(t, event.Time, "E2E: Event %d should have time", i)
			assert.NotEmpty(t, event.EventType, "E2E: Event %d should have type", i)
		}

		t.Logf("E2E: Data quality validated")
		t.Logf("All fields present and properly formatted")
		t.Logf("Events: %d items with proper structure", len(data.Events))
	})
}
