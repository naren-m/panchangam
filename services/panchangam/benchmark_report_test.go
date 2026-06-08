package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/require"
)

func BenchmarkFeatureCoverageReport(b *testing.B) {
	b.Run("Performance_Report_Generation", func(b *testing.B) {
		ctx := context.Background()
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		location := astronomy.Location{Latitude: 12.9716, Longitude: 77.5946}

		server := NewPanchangamServer()
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			astronomyStart := time.Now()
			sunTimes, err := astronomy.CalculateSunTimesWithContext(ctx, location, testDate)
			astronomyDuration := time.Since(astronomyStart)

			require.NoError(b, err)
			require.NotNil(b, sunTimes)

			serviceStart := time.Now()
			resp, serviceErr := server.Get(ctx, req)
			serviceDuration := time.Since(serviceStart)

			if serviceErr == nil && resp != nil && i == 0 {
				b.Logf("Performance Report:")
				b.Logf("  Astronomy calculation: %v (target: <100ms)", astronomyDuration)
				b.Logf("  Service response: %v (target: <500ms)", serviceDuration)
				b.Logf("  Combined overhead: %v", serviceDuration-astronomyDuration)
			}
		}
	})
}
