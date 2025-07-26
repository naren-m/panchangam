package ephemeris

import (
	"context"
	"sync"
	"time"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

// Cache defines the interface for ephemeris data caching
type Cache interface {
	// Get retrieves a value from the cache
	Get(ctx context.Context, key string) (interface{}, bool)
	
	// Set stores a value in the cache with TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration)
	
	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) bool
	
	// Clear clears all cache entries
	Clear(ctx context.Context) error
	
	// GetStats returns cache statistics
	GetStats(ctx context.Context) *CacheStats
	
	// Close closes the cache and releases resources
	Close() error
}

// CacheStats represents cache statistics
type CacheStats struct {
	Hits            int64     `json:"hits"`
	Misses          int64     `json:"misses"`
	Evictions       int64     `json:"evictions"`
	Entries         int64     `json:"entries"`
	MemoryUsage     int64     `json:"memory_usage_bytes"`
	LastAccess      time.Time `json:"last_access"`
	HitRate         float64   `json:"hit_rate"`
	AverageLatency  time.Duration `json:"average_latency"`
}

// CacheEntry represents a cached entry
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
	CreatedAt time.Time
	AccessCount int64
	LastAccess  time.Time
}

// MemoryCache implements an in-memory cache with observability
type MemoryCache struct {
	data      map[string]*CacheEntry
	mutex     sync.RWMutex
	stats     *CacheStats
	observer  observability.ObserverInterface
	cleanupTicker *time.Ticker
	stopCleanup   chan struct{}
	maxSize   int
	defaultTTL time.Duration
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache(maxSize int, defaultTTL time.Duration) *MemoryCache {
	cache := &MemoryCache{
		data:       make(map[string]*CacheEntry),
		stats:      &CacheStats{},
		observer:   observability.Observer(),
		maxSize:    maxSize,
		defaultTTL: defaultTTL,
		stopCleanup: make(chan struct{}),
	}
	
	// Start cleanup goroutine
	cache.cleanupTicker = time.NewTicker(5 * time.Minute)
	go cache.cleanupExpired()
	
	return cache
}

// Get retrieves a value from the cache
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	_, span := c.observer.CreateSpan(ctx, "ephemeris.cache.Get")
	defer span.End()
	
	span.SetAttributes(
		attribute.String("cache_key", key),
		attribute.String("operation", "get"),
	)
	
	start := time.Now()
	
	c.mutex.RLock()
	entry, exists := c.data[key]
	c.mutex.RUnlock()
	
	latency := time.Since(start)
	
	if !exists {
		c.recordMiss(latency)
		span.SetAttributes(
			attribute.Bool("cache_hit", false),
			attribute.Int64("latency_ms", latency.Milliseconds()),
		)
		span.AddEvent("Cache miss")
		return nil, false
	}
	
	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		c.mutex.Lock()
		delete(c.data, key)
		c.mutex.Unlock()
		
		c.recordMiss(latency)
		span.SetAttributes(
			attribute.Bool("cache_hit", false),
			attribute.Bool("expired", true),
			attribute.Int64("latency_ms", latency.Milliseconds()),
		)
		span.AddEvent("Cache entry expired")
		return nil, false
	}
	
	// Update access statistics
	c.mutex.Lock()
	entry.AccessCount++
	entry.LastAccess = time.Now()
	c.mutex.Unlock()
	
	c.recordHit(latency)
	span.SetAttributes(
		attribute.Bool("cache_hit", true),
		attribute.Int64("latency_ms", latency.Milliseconds()),
		attribute.Int64("access_count", entry.AccessCount),
	)
	span.AddEvent("Cache hit")
	
	return entry.Value, true
}

// Set stores a value in the cache with TTL
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	_, span := c.observer.CreateSpan(ctx, "ephemeris.cache.Set")
	defer span.End()
	
	if ttl == 0 {
		ttl = c.defaultTTL
	}
	
	span.SetAttributes(
		attribute.String("cache_key", key),
		attribute.String("operation", "set"),
		attribute.Int64("ttl_seconds", int64(ttl.Seconds())),
	)
	
	start := time.Now()
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Check if we need to evict entries
	if len(c.data) >= c.maxSize {
		c.evictLRU()
	}
	
	now := time.Now()
	entry := &CacheEntry{
		Value:       value,
		ExpiresAt:   now.Add(ttl),
		CreatedAt:   now,
		AccessCount: 1,
		LastAccess:  now,
	}
	
	c.data[key] = entry
	c.stats.Entries = int64(len(c.data))
	
	latency := time.Since(start)
	
	span.SetAttributes(
		attribute.Bool("success", true),
		attribute.Int64("latency_ms", latency.Milliseconds()),
		attribute.Int64("cache_size", c.stats.Entries),
	)
	span.AddEvent("Cache entry stored")
}

// Delete removes a value from the cache
func (c *MemoryCache) Delete(ctx context.Context, key string) bool {
	_, span := c.observer.CreateSpan(ctx, "ephemeris.cache.Delete")
	defer span.End()
	
	span.SetAttributes(
		attribute.String("cache_key", key),
		attribute.String("operation", "delete"),
	)
	
	c.mutex.Lock()
	_, exists := c.data[key]
	if exists {
		delete(c.data, key)
		c.stats.Entries = int64(len(c.data))
	}
	c.mutex.Unlock()
	
	span.SetAttributes(
		attribute.Bool("found", exists),
		attribute.Int64("cache_size", c.stats.Entries),
	)
	
	if exists {
		span.AddEvent("Cache entry deleted")
	} else {
		span.AddEvent("Cache entry not found")
	}
	
	return exists
}

// Clear clears all cache entries
func (c *MemoryCache) Clear(ctx context.Context) error {
	_, span := c.observer.CreateSpan(ctx, "ephemeris.cache.Clear")
	defer span.End()
	
	span.SetAttributes(attribute.String("operation", "clear"))
	
	c.mutex.Lock()
	entriesCleared := len(c.data)
	c.data = make(map[string]*CacheEntry)
	c.stats.Entries = 0
	c.mutex.Unlock()
	
	span.SetAttributes(
		attribute.Int("entries_cleared", entriesCleared),
		attribute.Bool("success", true),
	)
	span.AddEvent("Cache cleared")
	
	return nil
}

// GetStats returns cache statistics
func (c *MemoryCache) GetStats(ctx context.Context) *CacheStats {
	_, span := c.observer.CreateSpan(ctx, "ephemeris.cache.GetStats")
	defer span.End()
	
	c.mutex.RLock()
	stats := &CacheStats{
		Hits:        c.stats.Hits,
		Misses:      c.stats.Misses,
		Evictions:   c.stats.Evictions,
		Entries:     c.stats.Entries,
		MemoryUsage: c.stats.MemoryUsage,
		LastAccess:  c.stats.LastAccess,
		HitRate:     c.stats.HitRate,
		AverageLatency: c.stats.AverageLatency,
	}
	c.mutex.RUnlock()
	
	// Calculate hit rate
	total := stats.Hits + stats.Misses
	if total > 0 {
		stats.HitRate = float64(stats.Hits) / float64(total)
	}
	
	span.SetAttributes(
		attribute.Int64("cache_hits", stats.Hits),
		attribute.Int64("cache_misses", stats.Misses),
		attribute.Int64("cache_entries", stats.Entries),
		attribute.Float64("hit_rate", stats.HitRate),
	)
	span.AddEvent("Cache statistics retrieved")
	
	return stats
}

// Close closes the cache and releases resources
func (c *MemoryCache) Close() error {
	close(c.stopCleanup)
	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
	}
	
	c.mutex.Lock()
	c.data = nil
	c.mutex.Unlock()
	
	return nil
}

// recordHit records a cache hit
func (c *MemoryCache) recordHit(latency time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.stats.Hits++
	c.stats.LastAccess = time.Now()
	c.updateAverageLatency(latency)
}

// recordMiss records a cache miss
func (c *MemoryCache) recordMiss(latency time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.stats.Misses++
	c.stats.LastAccess = time.Now()
	c.updateAverageLatency(latency)
}

// updateAverageLatency updates the average latency using exponential moving average
func (c *MemoryCache) updateAverageLatency(latency time.Duration) {
	alpha := 0.1 // EMA smoothing factor
	if c.stats.AverageLatency == 0 {
		c.stats.AverageLatency = latency
	} else {
		c.stats.AverageLatency = time.Duration(
			float64(c.stats.AverageLatency)*(1-alpha) + float64(latency)*alpha,
		)
	}
}

// evictLRU evicts the least recently used entry
func (c *MemoryCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time
	
	for key, entry := range c.data {
		if oldestKey == "" || entry.LastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.LastAccess
		}
	}
	
	if oldestKey != "" {
		delete(c.data, oldestKey)
		c.stats.Evictions++
		c.stats.Entries = int64(len(c.data))
	}
}

// cleanupExpired removes expired entries periodically
func (c *MemoryCache) cleanupExpired() {
	for {
		select {
		case <-c.cleanupTicker.C:
			c.mutex.Lock()
			now := time.Now()
			for key, entry := range c.data {
				if now.After(entry.ExpiresAt) {
					delete(c.data, key)
				}
			}
			c.stats.Entries = int64(len(c.data))
			c.mutex.Unlock()
		case <-c.stopCleanup:
			return
		}
	}
}

// NoOpCache is a cache that doesn't cache anything (for testing)
type NoOpCache struct{}

// NewNoOpCache creates a new no-op cache
func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

// Get always returns false (no cache)
func (c *NoOpCache) Get(ctx context.Context, key string) (interface{}, bool) {
	return nil, false
}

// Set does nothing
func (c *NoOpCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	// No-op
}

// Delete always returns false
func (c *NoOpCache) Delete(ctx context.Context, key string) bool {
	return false
}

// Clear does nothing
func (c *NoOpCache) Clear(ctx context.Context) error {
	return nil
}

// GetStats returns empty stats
func (c *NoOpCache) GetStats(ctx context.Context) *CacheStats {
	return &CacheStats{}
}

// Close does nothing
func (c *NoOpCache) Close() error {
	return nil
}