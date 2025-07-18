package ephemeris

import (
	"context"
	"sync"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

// HealthChecker monitors the health of ephemeris providers
type HealthChecker struct {
	providers    []EphemerisProvider
	statuses     map[string]*HealthStatus
	mutex        sync.RWMutex
	observer     observability.ObserverInterface
	ticker       *time.Ticker
	stopChan     chan struct{}
	interval     time.Duration
	timeout      time.Duration
	isRunning    bool
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(providers []EphemerisProvider) *HealthChecker {
	return &HealthChecker{
		providers: providers,
		statuses:  make(map[string]*HealthStatus),
		observer:  observability.Observer(),
		interval:  30 * time.Second,
		timeout:   5 * time.Second,
		stopChan:  make(chan struct{}),
	}
}

// Start starts the health checking routine
func (h *HealthChecker) Start() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	if h.isRunning {
		return
	}
	
	h.isRunning = true
	h.ticker = time.NewTicker(h.interval)
	
	// Initial health check
	go h.checkHealth()
	
	// Start periodic health checks
	go h.run()
}

// Stop stops the health checking routine
func (h *HealthChecker) Stop() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	if !h.isRunning {
		return
	}
	
	h.isRunning = false
	
	// Close the channel if it's not already closed
	select {
	case <-h.stopChan:
		// Already closed
	default:
		close(h.stopChan)
	}
	
	if h.ticker != nil {
		h.ticker.Stop()
	}
}

// GetStatus returns the health status of a provider
func (h *HealthChecker) GetStatus(providerName string) (*HealthStatus, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	status, exists := h.statuses[providerName]
	return status, exists
}

// GetAllStatuses returns all health statuses
func (h *HealthChecker) GetAllStatuses() map[string]*HealthStatus {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	statuses := make(map[string]*HealthStatus)
	for name, status := range h.statuses {
		statuses[name] = status
	}
	
	return statuses
}

// IsHealthy returns true if all providers are healthy
func (h *HealthChecker) IsHealthy() bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	for _, status := range h.statuses {
		if !status.Available {
			return false
		}
	}
	
	return len(h.statuses) > 0
}

// run runs the health checking loop
func (h *HealthChecker) run() {
	for {
		select {
		case <-h.ticker.C:
			h.checkHealth()
		case <-h.stopChan:
			return
		}
	}
}

// checkHealth checks the health of all providers
func (h *HealthChecker) checkHealth() {
	ctx, span := h.observer.CreateSpan(context.Background(), "ephemeris.health.CheckHealth")
	defer span.End()
	
	span.SetAttributes(
		attribute.Int("provider_count", len(h.providers)),
		attribute.String("operation", "health_check"),
	)
	
	var wg sync.WaitGroup
	
	for _, provider := range h.providers {
		if provider == nil {
			continue
		}
		
		wg.Add(1)
		go func(p EphemerisProvider) {
			defer wg.Done()
			h.checkProviderHealth(ctx, p)
		}(provider)
	}
	
	wg.Wait()
	
	// Update overall health status
	healthyCount := 0
	for _, status := range h.statuses {
		if status.Available {
			healthyCount++
		}
	}
	
	span.SetAttributes(
		attribute.Int("healthy_providers", healthyCount),
		attribute.Int("total_providers", len(h.statuses)),
		attribute.Bool("overall_healthy", healthyCount > 0),
	)
	
	span.AddEvent("Health check completed")
}

// checkProviderHealth checks the health of a single provider
func (h *HealthChecker) checkProviderHealth(ctx context.Context, provider EphemerisProvider) {
	ctx, span := h.observer.CreateSpan(ctx, "ephemeris.health.CheckProvider")
	defer span.End()
	
	providerName := provider.GetProviderName()
	
	span.SetAttributes(
		attribute.String("provider_name", providerName),
		attribute.String("provider_version", provider.GetVersion()),
	)
	
	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()
	
	start := time.Now()
	
	// Check if provider is available
	available := provider.IsAvailable(timeoutCtx)
	responseTime := time.Since(start)
	
	status := &HealthStatus{
		Available:    available,
		LastCheck:    time.Now(),
		ResponseTime: responseTime,
		Version:      provider.GetVersion(),
		Source:       providerName,
	}
	
	// Get data range if available
	if available {
		startJD, endJD := provider.GetDataRange()
		status.DataStartJD = float64(startJD)
		status.DataEndJD = float64(endJD)
		
		// Try to get detailed health status
		if detailedStatus, err := provider.GetHealthStatus(timeoutCtx); err == nil {
			status.ErrorMessage = detailedStatus.ErrorMessage
		}
	} else {
		status.ErrorMessage = "Provider not available"
	}
	
	// Update status
	h.mutex.Lock()
	h.statuses[providerName] = status
	h.mutex.Unlock()
	
	span.SetAttributes(
		attribute.Bool("available", available),
		attribute.Int64("response_time_ms", responseTime.Milliseconds()),
		attribute.Float64("data_start_jd", status.DataStartJD),
		attribute.Float64("data_end_jd", status.DataEndJD),
	)
	
	if available {
		span.AddEvent("Provider health check passed")
	} else {
		span.AddEvent("Provider health check failed")
	}
}

// SetCheckInterval sets the health check interval
func (h *HealthChecker) SetCheckInterval(interval time.Duration) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	h.interval = interval
	
	if h.ticker != nil {
		h.ticker.Stop()
		h.ticker = time.NewTicker(interval)
	}
}

// SetTimeout sets the health check timeout
func (h *HealthChecker) SetTimeout(timeout time.Duration) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	h.timeout = timeout
}

// AddProvider adds a new provider to health checking
func (h *HealthChecker) AddProvider(provider EphemerisProvider) {
	if provider == nil {
		return
	}
	
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	h.providers = append(h.providers, provider)
	
	// Immediately check the new provider
	if h.isRunning {
		go h.checkProviderHealth(context.Background(), provider)
	}
}

// RemoveProvider removes a provider from health checking
func (h *HealthChecker) RemoveProvider(providerName string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	// Remove from providers list
	for i, provider := range h.providers {
		if provider.GetProviderName() == providerName {
			h.providers = append(h.providers[:i], h.providers[i+1:]...)
			break
		}
	}
	
	// Remove from statuses
	delete(h.statuses, providerName)
}

// GetMetrics returns health check metrics
func (h *HealthChecker) GetMetrics() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	metrics := make(map[string]interface{})
	
	totalProviders := len(h.providers)
	healthyProviders := 0
	totalResponseTime := time.Duration(0)
	
	for _, status := range h.statuses {
		if status.Available {
			healthyProviders++
		}
		totalResponseTime += status.ResponseTime
	}
	
	var avgResponseTime time.Duration
	if len(h.statuses) > 0 {
		avgResponseTime = totalResponseTime / time.Duration(len(h.statuses))
	}
	
	metrics["total_providers"] = totalProviders
	metrics["healthy_providers"] = healthyProviders
	metrics["unhealthy_providers"] = totalProviders - healthyProviders
	metrics["health_percentage"] = float64(healthyProviders) / float64(totalProviders) * 100
	metrics["average_response_time_ms"] = avgResponseTime.Milliseconds()
	metrics["check_interval_seconds"] = h.interval.Seconds()
	metrics["timeout_seconds"] = h.timeout.Seconds()
	metrics["last_check"] = time.Now()
	
	return metrics
}