package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServicePerformance(t *testing.T) {
	observability.NewLocalObserver()

	server := NewPanchangamServer()

	t.Run("Functional_Service_Performance", func(t *testing.T) {
		ctx := context.Background()

		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
		}

		_, _ = server.Get(ctx, req)

		start := time.Now()
		resp, err := server.Get(ctx, req)
		duration := time.Since(start)

		assert.NoError(t, err, "Performance test should not fail")
		require.NotNil(t, resp, "Response should not be nil")
		assert.True(t, duration < 500*time.Millisecond,
			"Service response should be under 500ms, got %v", duration)

		t.Logf("Service response time: %v", duration)
	})

	t.Run("Functional_Service_Concurrent_Performance", func(t *testing.T) {
		ctx := context.Background()

		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
		}

		concurrency := 10
		results := make(chan time.Duration, concurrency)

		start := time.Now()
		for i := 0; i < concurrency; i++ {
			go func() {
				reqStart := time.Now()
				resp, err := server.Get(ctx, req)
				reqDuration := time.Since(reqStart)

				assert.NoError(t, err, "Concurrent request should not fail")
				require.NotNil(t, resp, "Concurrent response should not be nil")

				results <- reqDuration
			}()
		}

		var totalDuration time.Duration
		var maxDuration time.Duration
		for i := 0; i < concurrency; i++ {
			duration := <-results
			totalDuration += duration
			if duration > maxDuration {
				maxDuration = duration
			}
		}

		totalTime := time.Since(start)
		avgDuration := totalDuration / time.Duration(concurrency)

		assert.True(t, maxDuration < 1*time.Second,
			"Max concurrent response should be under 1s, got %v", maxDuration)
		assert.True(t, avgDuration < 600*time.Millisecond,
			"Average concurrent response should be under 600ms, got %v", avgDuration)

		t.Logf("Concurrent performance (%d requests):", concurrency)
		t.Logf("  Total time: %v", totalTime)
		t.Logf("  Average duration: %v", avgDuration)
		t.Logf("  Max duration: %v", maxDuration)
	})
}
