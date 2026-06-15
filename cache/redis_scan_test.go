package cache

import (
	"os"
	"strings"
	"testing"
)

func TestRedisCacheUsesScanInsteadOfBlockingKeys(t *testing.T) {
	content := readRedisSource(t)
	if strings.Contains(content, ".Keys(") {
		t.Fatal("Redis cache code must use SCAN instead of the blocking KEYS command")
	}
	if !strings.Contains(content, ".Scan(ctx, cursor, \"panchangam:*\",") {
		t.Fatal("Redis cache code should scan panchangam:* keys")
	}
}

func TestRedisCacheHandlesCleanupDeleteErrors(t *testing.T) {
	if strings.Contains(readRedisSource(t), "\tr.client.Del(ctx, key)\n") {
		t.Fatal("cache cleanup deletes must handle or log Redis delete errors")
	}
}

func TestRedisCacheClosesClientWhenPingFails(t *testing.T) {
	if !strings.Contains(readRedisSource(t), "_ = rdb.Close()") {
		t.Fatal("Redis cache constructor should close the client when the startup ping fails")
	}
}

func TestRedisCacheBatchCleanupDeletesInvalidEntries(t *testing.T) {
	content := readRedisSource(t)
	if !strings.Contains(content, "Failed to delete corrupted batch cache entry") {
		t.Fatal("batch cache reads should delete corrupted entries after logging the decode error")
	}
	if !strings.Contains(content, "Failed to delete expired batch cache entry") {
		t.Fatal("batch cache reads should delete expired entries after logging the staleness")
	}
}

func TestRedisCacheGetReportsInvalidCachedJSON(t *testing.T) {
	content := readRedisSource(t)
	if !strings.Contains(content, "failed to decode cache key %s") {
		t.Fatal("single cache reads should return a clear error when cached JSON is invalid")
	}
}

func TestRedisCacheBatchGetReportsCommandErrors(t *testing.T) {
	content := readRedisSource(t)
	if !strings.Contains(content, "failed to get batch cache key %s") {
		t.Fatal("batch cache reads should return a clear keyed error for Redis command failures")
	}
}

func TestRedisCacheBatchGetReportsInvalidCachedJSON(t *testing.T) {
	content := readRedisSource(t)
	if !strings.Contains(content, "failed to decode batch cache key %s") {
		t.Fatal("batch cache reads should return a clear keyed error when cached JSON is invalid")
	}
}

func TestRedisCacheBatchSetReportsMarshalErrors(t *testing.T) {
	content := readRedisSource(t)
	if !strings.Contains(content, "failed to marshal batch cache data for key %s") {
		t.Fatal("batch cache writes should return a clear keyed error when marshaling fails")
	}
}

func readRedisSource(t *testing.T) string {
	t.Helper()

	source, err := os.ReadFile("redis.go")
	if err != nil {
		t.Fatalf("read redis.go: %v", err)
	}

	return string(source)
}
