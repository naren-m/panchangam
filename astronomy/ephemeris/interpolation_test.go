package ephemeris

import (
	"testing"
	"time"
)

func createTestManager(tb testing.TB) *Manager {
	tb.Helper()

	primary := NewJPLProvider()
	fallback := NewSwissProvider()
	cache := NewMemoryCache(100, 1*time.Hour)

	return NewManager(primary, fallback, cache)
}
