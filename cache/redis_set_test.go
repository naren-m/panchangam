package cache

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRedisCacheSetRejectsNilData(t *testing.T) {
	cache := &RedisCache{ttl: time.Minute}

	err := cache.Set(context.Background(), "nil-entry", nil)

	if err == nil {
		t.Fatal("expected nil cache data to return an error")
	}
	if !strings.Contains(err.Error(), "nil cache data for key nil-entry") {
		t.Fatalf("expected keyed nil-entry error, got %v", err)
	}
}

func TestRedisCacheSetRejectsUninitializedClient(t *testing.T) {
	cache := &RedisCache{ttl: time.Minute}

	err := cache.Set(context.Background(), "valid-entry", &PanchangamCacheData{
		Date: "2026-06-05",
	})

	requireUninitializedCacheError(t, err)
}
