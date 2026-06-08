package ephemeris

import "github.com/naren-m/panchangam/observability"

// Manager manages multiple ephemeris providers with fallback and caching
type Manager struct {
	primary       EphemerisProvider
	fallback      EphemerisProvider
	cache         Cache
	observer      observability.ObserverInterface
	healthChecker *HealthChecker
}

// NewManager creates a new ephemeris manager
func NewManager(primary, fallback EphemerisProvider, cache Cache) *Manager {
	manager := &Manager{
		primary:  primary,
		fallback: fallback,
		cache:    cache,
		observer: observability.Observer(),
	}

	// Initialize health checker
	manager.healthChecker = NewHealthChecker([]EphemerisProvider{primary, fallback})

	return manager
}
