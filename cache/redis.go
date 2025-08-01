package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/naren-m/panchangam/log"
)

var logger = log.Logger()

// RedisCache represents a Redis-based cache
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// PanchangamCacheData represents cached panchangam data
type PanchangamCacheData struct {
	Date         string    `json:"date"`
	Tithi        string    `json:"tithi"`
	Nakshatra    string    `json:"nakshatra"`
	Yoga         string    `json:"yoga"`
	Karana       string    `json:"karana"`
	SunriseTime  string    `json:"sunrise_time"`
	SunsetTime   string    `json:"sunset_time"`
	Events       []Event   `json:"events"`
	CachedAt     time.Time `json:"cached_at"`
}

// Event represents a panchangam event
type Event struct {
	Name      string `json:"name"`
	Time      string `json:"time"`
	EventType string `json:"event_type"`
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(addr, password string, db int, ttl time.Duration) (*RedisCache, error) {
	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis cache connected successfully", "addr", addr, "db", db, "ttl", ttl)

	return &RedisCache{
		client: rdb,
		ttl:    ttl,
	}, nil
}

// GenerateCacheKey generates a cache key for panchangam data
func (r *RedisCache) GenerateCacheKey(date, region, method string, lat, lng float64) string {
	return fmt.Sprintf("panchangam:%s:%s:%s:%.4f:%.4f", date, region, method, lat, lng)
}

// Get retrieves cached panchangam data
func (r *RedisCache) Get(ctx context.Context, key string) (*PanchangamCacheData, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get cache key %s: %w", key, err)
	}

	var data PanchangamCacheData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		logger.Error("Failed to unmarshal cached data", "key", key, "error", err)
		// Delete corrupted cache entry
		r.client.Del(ctx, key)
		return nil, nil
	}

	// Check if data is still fresh (additional staleness check)
	if time.Since(data.CachedAt) > r.ttl {
		logger.Debug("Cache entry expired", "key", key, "cached_at", data.CachedAt)
		r.client.Del(ctx, key)
		return nil, nil
	}

	logger.Debug("Cache hit", "key", key, "cached_at", data.CachedAt)
	return &data, nil
}

// Set stores panchangam data in cache
func (r *RedisCache) Set(ctx context.Context, key string, data *PanchangamCacheData) error {
	data.CachedAt = time.Now()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	if err := r.client.Set(ctx, key, jsonData, r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache key %s: %w", key, err)
	}

	logger.Debug("Cache set", "key", key, "ttl", r.ttl)
	return nil
}

// GetBatch retrieves multiple cached entries
func (r *RedisCache) GetBatch(ctx context.Context, keys []string) (map[string]*PanchangamCacheData, error) {
	if len(keys) == 0 {
		return make(map[string]*PanchangamCacheData), nil
	}

	// Use pipeline for batch operations
	pipe := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to execute batch get: %w", err)
	}

	result := make(map[string]*PanchangamCacheData)
	
	for i, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue // Cache miss for this key
			}
			logger.Error("Failed to get batch cache key", "key", keys[i], "error", err)
			continue
		}

		var data PanchangamCacheData
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			logger.Error("Failed to unmarshal batch cached data", "key", keys[i], "error", err)
			continue
		}

		// Check staleness
		if time.Since(data.CachedAt) <= r.ttl {
			result[keys[i]] = &data
		}
	}

	logger.Debug("Batch cache operation", "requested", len(keys), "hits", len(result))
	return result, nil
}

// SetBatch stores multiple entries in cache
func (r *RedisCache) SetBatch(ctx context.Context, data map[string]*PanchangamCacheData) error {
	if len(data) == 0 {
		return nil
	}

	pipe := r.client.Pipeline()
	now := time.Now()

	for key, cacheData := range data {
		cacheData.CachedAt = now

		jsonData, err := json.Marshal(cacheData)
		if err != nil {
			logger.Error("Failed to marshal batch cache data", "key", key, "error", err)
			continue
		}

		pipe.Set(ctx, key, jsonData, r.ttl)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute batch set: %w", err)
	}

	logger.Debug("Batch cache set", "count", len(data), "ttl", r.ttl)
	return nil
}

// Delete removes a cache entry
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Clear removes all cache entries with panchangam prefix
func (r *RedisCache) Clear(ctx context.Context) error {
	keys, err := r.client.Keys(ctx, "panchangam:*").Result()
	if err != nil {
		return fmt.Errorf("failed to get cache keys: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	logger.Info("Cache cleared", "keys_deleted", len(keys))
	return nil
}

// GetStats returns cache statistics
func (r *RedisCache) GetStats(ctx context.Context) (map[string]interface{}, error) {
	info, err := r.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis stats: %w", err)
	}

	// Count panchangam cache keys
	keys, err := r.client.Keys(ctx, "panchangam:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to count cache keys: %w", err)
	}

	stats := map[string]interface{}{
		"cache_keys_count": len(keys),
		"ttl_seconds":      int(r.ttl.Seconds()),
		"redis_info":       info,
	}

	return stats, nil
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// HealthCheck performs a health check on the cache
func (r *RedisCache) HealthCheck(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}