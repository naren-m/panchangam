package cache

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRedisCacheGetRejectsUninitializedClient(t *testing.T) {
	cache := &RedisCache{ttl: time.Minute}

	_, err := cache.Get(context.Background(), "missing-entry")

	requireUninitializedCacheError(t, err)
}

func TestRedisCacheGetBatchRejectsUninitializedClient(t *testing.T) {
	cache := &RedisCache{ttl: time.Minute}

	_, err := cache.GetBatch(context.Background(), []string{"missing-entry"})

	requireUninitializedCacheError(t, err)
}

func TestRedisCacheDeleteRejectsUninitializedClient(t *testing.T) {
	cache := &RedisCache{ttl: time.Minute}

	err := cache.Delete(context.Background(), "missing-entry")

	requireUninitializedCacheError(t, err)
}

func TestRedisCacheClearRejectsUninitializedClient(t *testing.T) {
	cache := &RedisCache{ttl: time.Minute}

	err := cache.Clear(context.Background())

	requireUninitializedCacheError(t, err)
}

func TestRedisCacheGetStatsRejectsUninitializedClient(t *testing.T) {
	cache := &RedisCache{ttl: time.Minute}

	_, err := cache.GetStats(context.Background())

	requireUninitializedCacheError(t, err)
}

func TestRedisCacheCloseRejectsUninitializedClient(t *testing.T) {
	cache := &RedisCache{ttl: time.Minute}

	err := cache.Close()

	requireUninitializedCacheError(t, err)
}

func TestRedisCacheHealthCheckRejectsUninitializedClient(t *testing.T) {
	cache := &RedisCache{ttl: time.Minute}

	err := cache.HealthCheck(context.Background())

	requireUninitializedCacheError(t, err)
}

func requireUninitializedCacheError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("expected uninitialized client to return an error")
	}
	if !errors.Is(err, errRedisCacheUninitialized) {
		t.Fatalf("expected uninitialized client error, got %v", err)
	}
}
