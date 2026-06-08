//go:build integration
// +build integration

package gateway

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestConcurrentRequestPerformance(t *testing.T) {
	handler := newTestGatewayHandler(t)
	concurrentRequests := 50
	timeout := 5 * time.Second

	start := time.Now()
	var wg sync.WaitGroup
	errors := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			query := buildQueryString("2024-01-15", 12.9716, 77.5946, "Asia/Kolkata")
			req := httptest.NewRequest("GET", "/api/v1/panchangam?"+query, nil)
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != http.StatusOK {
				errors <- &testError{msg: "Request failed"}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)

	errorCount := 0
	for range errors {
		errorCount++
	}

	if errorCount > 0 {
		t.Errorf("%d out of %d concurrent requests failed", errorCount, concurrentRequests)
	}

	if duration > timeout {
		t.Errorf("Concurrent requests took %v, expected < %v", duration, timeout)
	}

	t.Logf("Concurrent performance: %d requests in %v (%.2f req/sec)",
		concurrentRequests, duration, float64(concurrentRequests)/duration.Seconds())
}

func TestResponseTimeTarget(t *testing.T) {
	handler := newTestGatewayHandler(t)
	iterations := 100
	var totalDuration time.Duration

	for i := 0; i < iterations; i++ {
		query := buildQueryString("2024-01-15", 12.9716, 77.5946, "Asia/Kolkata")
		req := httptest.NewRequest("GET", "/api/v1/panchangam?"+query, nil)
		w := httptest.NewRecorder()

		start := time.Now()
		handler(w, req)
		duration := time.Since(start)

		totalDuration += duration

		if w.Code != http.StatusOK {
			t.Errorf("Request %d failed with status %d", i, w.Code)
		}
	}

	avgDuration := totalDuration / time.Duration(iterations)
	t.Logf("Average response time: %v over %d requests", avgDuration, iterations)

	targetDuration := 500 * time.Millisecond
	if avgDuration > targetDuration {
		t.Errorf("Average response time %v exceeds target %v", avgDuration, targetDuration)
	}
}

func TestErrorRecoveryTime(t *testing.T) {
	handler := newTestGatewayHandler(t)
	maxRetries := 3
	retryDelay := 500 * time.Millisecond

	start := time.Now()

	for attempt := 0; attempt < maxRetries; attempt++ {
		query := buildQueryString("2024-01-15", 12.9716, 77.5946, "Asia/Kolkata")
		req := httptest.NewRequest("GET", "/api/v1/panchangam?"+query, nil)
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code == http.StatusOK {
			break
		}

		if attempt < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	duration := time.Since(start)

	if duration > 3*time.Second {
		t.Errorf("Error recovery took %v, expected <3s", duration)
	}

	t.Logf("Error recovery completed in %v", duration)
}
