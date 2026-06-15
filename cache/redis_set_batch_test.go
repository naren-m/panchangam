package cache

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRedisCacheSetBatchRejectsNilEntry(t *testing.T) {
	cache := &RedisCache{ttl: time.Minute}

	err := cache.SetBatch(context.Background(), map[string]*PanchangamCacheData{
		"nil-entry": nil,
	})

	if err == nil {
		t.Fatal("expected nil batch entry to return an error")
	}
	if !strings.Contains(err.Error(), "nil cache data for key nil-entry") {
		t.Fatalf("expected keyed nil-entry error, got %v", err)
	}
}

func TestRedisCacheSetBatchRejectsUninitializedClient(t *testing.T) {
	cache := &RedisCache{ttl: time.Minute}

	err := cache.SetBatch(context.Background(), map[string]*PanchangamCacheData{
		"valid-entry": {Date: "2026-06-05"},
	})

	requireUninitializedCacheError(t, err)
}
