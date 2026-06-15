package ephemeris

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache(3, 1*time.Second)
	ctx := context.Background()

	t.Run("basic operations", func(t *testing.T) {
		cache.Set(ctx, "key1", "value1", 0)
		value, found := cache.Get(ctx, "key1")
		assert.True(t, found)
		assert.Equal(t, "value1", value)

		_, found = cache.Get(ctx, "nonexistent")
		assert.False(t, found)
	})

	t.Run("ttl expiration", func(t *testing.T) {
		cache.Set(ctx, "key2", "value2", time.Second)

		value, found := cache.Get(ctx, "key2")
		assert.True(t, found)
		assert.Equal(t, "value2", value)

		cache.Set(ctx, "expired-key", "expired-value", -time.Nanosecond)
		_, found = cache.Get(ctx, "expired-key")
		assert.False(t, found)
	})

	t.Run("lru eviction", func(t *testing.T) {
		cache.Set(ctx, "key3", "value3", 0)
		cache.Set(ctx, "key4", "value4", 0)
		cache.Set(ctx, "key5", "value5", 0)

		cache.Get(ctx, "key3")

		cache.Set(ctx, "key6", "value6", 0)

		_, found := cache.Get(ctx, "key3")
		assert.True(t, found)

		_, found = cache.Get(ctx, "key4")
		assert.False(t, found)
	})

	t.Run("cache stats", func(t *testing.T) {
		stats := cache.GetStats(ctx)
		assert.NotNil(t, stats)
		assert.True(t, stats.Hits > 0)
		assert.True(t, stats.Misses > 0)
	})

	t.Run("clear cache", func(t *testing.T) {
		cache.Set(ctx, "key7", "value7", 0)
		err := cache.Clear(ctx)
		assert.NoError(t, err)

		_, found := cache.Get(ctx, "key7")
		assert.False(t, found)
	})

	t.Run("close cache", func(t *testing.T) {
		err := cache.Close()
		assert.NoError(t, err)
	})
}

func TestMemoryCacheCloseIsSafeToCallTwice(t *testing.T) {
	cache := NewMemoryCache(3, time.Second)

	err := cache.Close()
	assert.NoError(t, err)

	err = cache.Close()
	assert.NoError(t, err)
}

func TestMemoryCacheSetAfterCloseDoesNotPanic(t *testing.T) {
	cache := NewMemoryCache(3, time.Second)
	ctx := context.Background()
	cache.Set(ctx, "existing-key", "existing-value", time.Second)

	assert.NoError(t, cache.Close())
	assert.Equal(t, int64(0), cache.GetStats(ctx).Entries)
	assert.NoError(t, cache.Clear(ctx))

	assert.NotPanics(t, func() {
		cache.Set(ctx, "late-key", "late-value", time.Second)
	})

	_, found := cache.Get(ctx, "late-key")
	assert.False(t, found)
}

func TestMemoryCacheGetAfterCloseDoesNotRecordMiss(t *testing.T) {
	cache := NewMemoryCache(3, time.Second)
	ctx := context.Background()

	assert.NoError(t, cache.Close())

	_, found := cache.Get(ctx, "late-key")
	assert.False(t, found)
	assert.Equal(t, int64(0), cache.GetStats(ctx).Misses)
}

func TestMemoryCacheExpiredGetUpdatesEntryCount(t *testing.T) {
	cache := NewMemoryCache(3, time.Second)
	defer func() {
		assert.NoError(t, cache.Close())
	}()

	ctx := context.Background()
	cache.Set(ctx, "expired-key", "expired-value", -time.Nanosecond)
	assert.Equal(t, int64(1), cache.GetStats(ctx).Entries)

	_, found := cache.Get(ctx, "expired-key")
	assert.False(t, found)
	assert.Equal(t, int64(0), cache.GetStats(ctx).Entries)
}

func TestMemoryCacheWithZeroMaxSizeDoesNotStoreEntries(t *testing.T) {
	cache := NewMemoryCache(0, time.Second)
	defer func() {
		assert.NoError(t, cache.Close())
	}()

	ctx := context.Background()

	assert.NotPanics(t, func() {
		cache.Set(ctx, "key", "value", time.Second)
	})

	_, found := cache.Get(ctx, "key")
	assert.False(t, found)
	assert.Equal(t, int64(0), cache.GetStats(ctx).Entries)
}

func TestMemoryCacheUpdatingExistingKeyDoesNotEvictOtherEntries(t *testing.T) {
	cache := NewMemoryCache(2, time.Second)
	defer func() {
		assert.NoError(t, cache.Close())
	}()

	ctx := context.Background()
	cache.Set(ctx, "first", "first-value", time.Second)
	cache.Set(ctx, "second", "second-value", time.Second)

	_, found := cache.Get(ctx, "first")
	assert.True(t, found)

	cache.Set(ctx, "first", "updated-first-value", time.Second)

	value, found := cache.Get(ctx, "first")
	assert.True(t, found)
	assert.Equal(t, "updated-first-value", value)

	value, found = cache.Get(ctx, "second")
	assert.True(t, found)
	assert.Equal(t, "second-value", value)
	assert.Equal(t, int64(2), cache.GetStats(ctx).Entries)
}

func TestMemoryCacheConcurrentGetDoesNotRace(t *testing.T) {
	cache := NewMemoryCache(3, time.Second)
	defer func() {
		assert.NoError(t, cache.Close())
	}()

	ctx := context.Background()
	cache.Set(ctx, "shared-key", "shared-value", time.Second)

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			for j := 0; j < 100; j++ {
				value, found := cache.Get(ctx, "shared-key")
				assert.True(t, found)
				assert.Equal(t, "shared-value", value)
			}
		}()
	}

	close(start)
	wg.Wait()
}

func TestMemoryCacheConcurrentSetDeleteDoesNotRace(t *testing.T) {
	cache := NewMemoryCache(3, time.Second)
	defer func() {
		assert.NoError(t, cache.Close())
	}()

	ctx := context.Background()
	var wg sync.WaitGroup
	start := make(chan struct{})

	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start

		for i := 0; i < 500; i++ {
			cache.Set(ctx, "shared-key", "shared-value", time.Second)
		}
	}()

	go func() {
		defer wg.Done()
		<-start

		for i := 0; i < 500; i++ {
			cache.Delete(ctx, "shared-key")
		}
	}()

	close(start)
	wg.Wait()
}

func TestNoOpCache(t *testing.T) {
	cache := NewNoOpCache()
	ctx := context.Background()

	t.Run("no-op operations", func(t *testing.T) {
		cache.Set(ctx, "key", "value", 0)

		_, found := cache.Get(ctx, "key")
		assert.False(t, found)

		deleted := cache.Delete(ctx, "key")
		assert.False(t, deleted)

		err := cache.Clear(ctx)
		assert.NoError(t, err)

		stats := cache.GetStats(ctx)
		assert.Equal(t, int64(0), stats.Hits)
		assert.Equal(t, int64(0), stats.Misses)

		err = cache.Close()
		assert.NoError(t, err)
	})
}
