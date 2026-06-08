package ephemeris

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthChecker(t *testing.T) {
	primary := NewJPLProvider()
	fallback := NewSwissProvider()

	checker := NewHealthChecker([]EphemerisProvider{primary, fallback})

	t.Run("initial status", func(t *testing.T) {
		assert.NotNil(t, checker)
	})

	t.Run("start and stop", func(t *testing.T) {
		checker.Start()
		assert.True(t, checker.isRunning)

		checker.Stop()
		assert.False(t, checker.isRunning)
	})

	t.Run("individual status", func(t *testing.T) {
		newChecker := NewHealthChecker([]EphemerisProvider{primary, fallback})
		newChecker.checkHealth()

		status, found := newChecker.GetStatus("JPL DE440")
		assert.True(t, found)
		assert.True(t, status.Available)
	})

	t.Run("metrics", func(t *testing.T) {
		metricsChecker := NewHealthChecker([]EphemerisProvider{primary, fallback})
		metricsChecker.checkHealth()

		metrics := metricsChecker.GetMetrics()
		assert.NotNil(t, metrics)
		assert.Equal(t, 2, metrics["total_providers"])
		assert.Equal(t, 2, metrics["healthy_providers"])
		assert.Equal(t, 0, metrics["unhealthy_providers"])
		assert.Equal(t, 100.0, metrics["health_percentage"])
	})

	t.Run("add and remove providers", func(t *testing.T) {
		addRemoveChecker := NewHealthChecker([]EphemerisProvider{primary, fallback})
		addRemoveChecker.checkHealth()

		statuses := addRemoveChecker.GetAllStatuses()
		assert.Len(t, statuses, 2)

		newProvider := NewJPLProvider()
		addRemoveChecker.AddProvider(newProvider)
		addRemoveChecker.checkHealth()

		statuses = addRemoveChecker.GetAllStatuses()
		assert.Len(t, statuses, 2)

		addRemoveChecker.RemoveProvider("JPL DE440")

		statuses = addRemoveChecker.GetAllStatuses()
		assert.Len(t, statuses, 1)
	})
}

func TestHealthCheckerMetricsWithNoProviders(t *testing.T) {
	checker := NewHealthChecker(nil)

	metrics := checker.GetMetrics()

	assert.Equal(t, 0, metrics["total_providers"])
	assert.Equal(t, 0, metrics["healthy_providers"])
	assert.Equal(t, 0, metrics["unhealthy_providers"])
	assert.Equal(t, 0.0, metrics["health_percentage"])
}

func TestHealthCheckerAggregateMethodsIgnoreNilStatuses(t *testing.T) {
	checker := NewHealthChecker(nil)
	checker.statuses["nil-provider"] = nil

	assert.NotPanics(t, func() {
		assert.False(t, checker.IsHealthy())
	})

	assert.NotPanics(t, func() {
		metrics := checker.GetMetrics()
		assert.Equal(t, 0, metrics["healthy_providers"])
		assert.Equal(t, int64(0), metrics["average_response_time_ms"])
	})

	assert.NotPanics(t, func() {
		checker.checkHealth()
	})
}

func TestNewHealthCheckerSkipsNilProviders(t *testing.T) {
	checker := NewHealthChecker([]EphemerisProvider{
		nil,
		&countingHealthProvider{name: "real-provider"},
	})

	if assert.Len(t, checker.providers, 1) {
		assert.Equal(t, "real-provider", checker.providers[0].GetProviderName())
	}
}

func TestNewHealthCheckerSkipsDuplicateProviderNames(t *testing.T) {
	checker := NewHealthChecker([]EphemerisProvider{
		&countingHealthProvider{name: "duplicate-provider"},
		&countingHealthProvider{name: "duplicate-provider"},
	})

	checker.checkHealth()

	assert.Len(t, checker.providers, 1)
	metrics := checker.GetMetrics()
	assert.Equal(t, 1, metrics["total_providers"])
	assert.Equal(t, 100.0, metrics["health_percentage"])
}

func TestNewHealthCheckerCopiesProviderList(t *testing.T) {
	providers := []EphemerisProvider{
		&countingHealthProvider{name: "original-provider"},
	}
	checker := NewHealthChecker(providers)

	providers[0] = &countingHealthProvider{name: "changed-provider"}

	checker.checkHealth()

	_, foundOriginal := checker.GetStatus("original-provider")
	_, foundChanged := checker.GetStatus("changed-provider")
	assert.True(t, foundOriginal)
	assert.False(t, foundChanged)
}

func TestHealthCheckerGetStatusReturnsCopy(t *testing.T) {
	checker := NewHealthChecker([]EphemerisProvider{
		&countingHealthProvider{name: "copy-provider"},
	})
	checker.checkHealth()

	status, found := checker.GetStatus("copy-provider")
	assert.True(t, found)

	status.Available = false
	status.ErrorMessage = "mutated by caller"

	current, found := checker.GetStatus("copy-provider")
	assert.True(t, found)
	assert.True(t, current.Available)
	assert.NotEqual(t, "mutated by caller", current.ErrorMessage)
}

func TestHealthCheckerGetAllStatusesReturnsCopies(t *testing.T) {
	checker := NewHealthChecker([]EphemerisProvider{
		&countingHealthProvider{name: "copy-provider"},
	})
	checker.checkHealth()

	statuses := checker.GetAllStatuses()
	statuses["copy-provider"].Available = false
	statuses["copy-provider"].ErrorMessage = "mutated by caller"

	current := checker.GetAllStatuses()
	assert.True(t, current["copy-provider"].Available)
	assert.NotEqual(t, "mutated by caller", current["copy-provider"].ErrorMessage)
}

func TestHealthCheckerRemoveProviderSkipsNilProviders(t *testing.T) {
	checker := NewHealthChecker([]EphemerisProvider{
		&countingHealthProvider{name: "real-provider"},
	})
	checker.providers = append([]EphemerisProvider{nil}, checker.providers...)

	assert.NotPanics(t, func() {
		checker.RemoveProvider("real-provider")
	})

	assert.Len(t, checker.providers, 1)
	assert.Nil(t, checker.providers[0])
}

func TestHealthCheckerAddProviderSkipsDuplicateNames(t *testing.T) {
	checker := NewHealthChecker([]EphemerisProvider{
		&countingHealthProvider{name: "duplicate-provider"},
	})

	checker.AddProvider(&countingHealthProvider{name: "duplicate-provider"})
	checker.checkHealth()

	assert.Len(t, checker.providers, 1)
	metrics := checker.GetMetrics()
	assert.Equal(t, 1, metrics["total_providers"])
	assert.Equal(t, 100.0, metrics["health_percentage"])
}

func TestHealthCheckerRestartKeepsPeriodicChecksRunning(t *testing.T) {
	provider := &countingHealthProvider{}
	checker := NewHealthChecker([]EphemerisProvider{provider})
	checker.SetCheckInterval(10 * time.Millisecond)
	checker.SetTimeout(100 * time.Millisecond)

	checker.Start()
	waitForHealthChecks(t, provider, 2)
	checker.Stop()

	checksAfterStop := provider.checks.Load()

	checker.Start()
	defer checker.Stop()

	waitForHealthChecks(t, provider, checksAfterStop+2)
}

func TestHealthCheckerStopClearsTicker(t *testing.T) {
	checker := NewHealthChecker(nil)
	checker.SetCheckInterval(10 * time.Millisecond)

	checker.Start()
	checker.Stop()

	assert.False(t, checker.isRunning)
	assert.Nil(t, checker.ticker)

	checker.SetCheckInterval(20 * time.Millisecond)

	assert.Equal(t, 20*time.Millisecond, checker.interval)
	assert.Nil(t, checker.ticker)
}

func TestHealthCheckerIgnoresNonPositiveIntervals(t *testing.T) {
	checker := NewHealthChecker(nil)
	originalInterval := checker.interval

	checker.SetCheckInterval(0)
	assert.Equal(t, originalInterval, checker.interval)

	checker.SetCheckInterval(10 * time.Millisecond)
	checker.Start()
	defer checker.Stop()

	assert.NotPanics(t, func() {
		checker.SetCheckInterval(0)
	})
	assert.Equal(t, 10*time.Millisecond, checker.interval)

	assert.NotPanics(t, func() {
		checker.SetCheckInterval(-time.Millisecond)
	})
	assert.Equal(t, 10*time.Millisecond, checker.interval)
}

func TestHealthCheckerIgnoresNonPositiveTimeouts(t *testing.T) {
	checker := NewHealthChecker(nil)
	originalTimeout := checker.timeout

	checker.SetTimeout(0)
	assert.Equal(t, originalTimeout, checker.timeout)

	checker.SetTimeout(10 * time.Millisecond)

	checker.SetTimeout(-time.Millisecond)
	assert.Equal(t, 10*time.Millisecond, checker.timeout)
}

func TestHealthCheckerIntervalChangeKeepsPeriodicChecksRunning(t *testing.T) {
	provider := &countingHealthProvider{}
	checker := NewHealthChecker([]EphemerisProvider{provider})
	checker.SetCheckInterval(50 * time.Millisecond)
	checker.SetTimeout(100 * time.Millisecond)

	checker.Start()
	defer checker.Stop()
	waitForHealthChecks(t, provider, 1)

	checker.SetCheckInterval(10 * time.Millisecond)

	waitForHealthChecks(t, provider, 3)
}

func TestHealthCheckerConcurrentProviderChanges(t *testing.T) {
	checker := NewHealthChecker([]EphemerisProvider{
		&countingHealthProvider{name: "base-provider"},
	})
	checker.SetTimeout(10 * time.Millisecond)

	var wg sync.WaitGroup
	start := make(chan struct{})

	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start

		for i := 0; i < 200; i++ {
			checker.checkHealth()
			checker.GetMetrics()
		}
	}()

	go func() {
		defer wg.Done()
		<-start

		for i := 0; i < 200; i++ {
			name := fmt.Sprintf("provider-%d", i)
			checker.AddProvider(&countingHealthProvider{name: name})
			checker.SetTimeout(time.Duration(10+i%5) * time.Millisecond)
			checker.RemoveProvider(name)
		}
	}()

	close(start)
	wg.Wait()
}

type countingHealthProvider struct {
	checks atomic.Int64
	name   string
}

func (p *countingHealthProvider) GetPlanetaryPositions(ctx context.Context, jd JulianDay) (*PlanetaryPositions, error) {
	return &PlanetaryPositions{JulianDay: jd}, nil
}

func (p *countingHealthProvider) GetSunPosition(ctx context.Context, jd JulianDay) (*SolarPosition, error) {
	return &SolarPosition{JulianDay: jd}, nil
}

func (p *countingHealthProvider) GetMoonPosition(ctx context.Context, jd JulianDay) (*LunarPosition, error) {
	return &LunarPosition{JulianDay: jd}, nil
}

func (p *countingHealthProvider) IsAvailable(ctx context.Context) bool {
	p.checks.Add(1)
	return true
}

func (p *countingHealthProvider) GetDataRange() (startJD, endJD JulianDay) {
	return 1, 2
}

func (p *countingHealthProvider) GetHealthStatus(ctx context.Context) (*HealthStatus, error) {
	return &HealthStatus{}, nil
}

func (p *countingHealthProvider) GetProviderName() string {
	if p.name != "" {
		return p.name
	}
	return "Counting Provider"
}

func (p *countingHealthProvider) GetVersion() string {
	return "test"
}

func (p *countingHealthProvider) Close() error {
	return nil
}

func waitForHealthChecks(t *testing.T, provider *countingHealthProvider, minimum int64) {
	t.Helper()

	deadline := time.Now().Add(250 * time.Millisecond)
	for time.Now().Before(deadline) {
		if provider.checks.Load() >= minimum {
			return
		}
		time.Sleep(time.Millisecond)
	}

	t.Fatalf("expected at least %d health checks, got %d", minimum, provider.checks.Load())
}
