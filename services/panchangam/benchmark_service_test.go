package panchangam

import (
	"context"
	"testing"
	"time"

	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/require"
)

func benchmarkServiceLayer(b *testing.B) {
	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
		Region:    "India",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		start := time.Now()
		resp, err := server.Get(ctx, req)
		duration := time.Since(start)

		if err != nil {
			b.Fatalf("service request failed: %v", err)
		}
		if resp == nil {
			b.Fatal("service response is nil")
		}
		if i == 0 {
			b.Logf("Service response time: %v", duration)
		}
	}
}

func benchmarkServiceResponseTarget(b *testing.B) {
	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()
		resp, err := server.Get(ctx, req)
		duration := time.Since(start)

		if err != nil {
			b.Fatalf("service request failed: %v", err)
		}
		if resp == nil {
			b.Fatal("service response is nil")
		}
		if duration > 500*time.Millisecond {
			b.Errorf("Service response exceeded 500ms target: %v", duration)
		}
	}
}

func benchmarkEndToEndTarget(b *testing.B) {
	server := NewPanchangamServer()
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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()
		resp, err := server.Get(ctx, req)

		if err != nil {
			b.Fatalf("service request failed: %v", err)
		}
		require.NotNil(b, resp)
		require.NotNil(b, resp.PanchangamData)

		data := resp.PanchangamData
		require.Equal(b, req.Date, data.Date)
		require.NotEmpty(b, data.Tithi)
		require.NotEmpty(b, data.Nakshatra)
		require.NotEmpty(b, data.Yoga)
		require.NotEmpty(b, data.Karana)
		require.NotEmpty(b, data.SunriseTime)
		require.NotEmpty(b, data.SunsetTime)

		duration := time.Since(start)
		if duration > 500*time.Millisecond {
			b.Errorf("End-to-end response exceeded 500ms target: %v", duration)
		}
	}
}

func BenchmarkConcurrentFeatureAccess(b *testing.B) {
	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
	}

	b.ResetTimer()
	b.SetParallelism(10)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()
			resp, err := server.Get(ctx, req)
			duration := time.Since(start)

			if err == nil && resp != nil && duration > 1*time.Second {
				b.Errorf("Concurrent response exceeded 1s target: %v", duration)
			}
		}
	})
}

func BenchmarkMemoryUsage(b *testing.B) {
	server := NewPanchangamServer()
	ctx := context.Background()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, err := server.Get(ctx, req)
		if err == nil && resp != nil {
			_ = resp.PanchangamData
		}
	}
}
