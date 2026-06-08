package ephemeris

import "github.com/naren-m/panchangam/observability"

func init() {
	// Initialize observability for testing.
	observability.NewLocalObserver()
}
