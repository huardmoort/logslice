package index_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/index"
)

func TestCacheGetOrBuild(t *testing.T) {
	cache := index.NewCache(5 * time.Minute)
	r := strings.NewReader(sampleLog)

	idx, err := cache.GetOrBuild("test-key", r, "")
	if err != nil {
		t.Fatalf("GetOrBuild failed: %v", err)
	}
	if len(idx.Entries) == 0 {
		t.Error("expected non-empty index")
	}
}

func TestCacheHit(t *testing.T) {
	cache := index.NewCache(5 * time.Minute)
	r := strings.NewReader(sampleLog)

	// Build once
	first, _ := cache.GetOrBuild("key", r, "")

	// Second call should return cached (reader is exhausted, but cache hit)
	second, err := cache.GetOrBuild("key", r, "")
	if err != nil {
		t.Fatalf("second GetOrBuild failed: %v", err)
	}
	if first != second {
		t.Error("expected same index pointer on cache hit")
	}
}

func TestCacheInvalidate(t *testing.T) {
	cache := index.NewCache(5 * time.Minute)
	r := strings.NewReader(sampleLog)
	cache.GetOrBuild("key", r, "") //nolint

	cache.Invalidate("key")
	_, ok := cache.Get("key")
	if ok {
		t.Error("expected cache miss after invalidation")
	}
}

func TestCacheTTLExpiry(t *testing.T) {
	cache := index.NewCache(1 * time.Millisecond)
	r := strings.NewReader(sampleLog)
	cache.GetOrBuild("key", r, "") //nolint

	time.Sleep(5 * time.Millisecond)

	_, ok := cache.Get("key")
	if ok {
		t.Error("expected cache miss after TTL expiry")
	}
}

func TestCacheZeroTTLNeverExpires(t *testing.T) {
	cache := index.NewCache(0)
	r := strings.NewReader(sampleLog)
	cache.GetOrBuild("key", r, "") //nolint

	time.Sleep(5 * time.Millisecond)

	_, ok := cache.Get("key")
	if !ok {
		t.Error("expected cache hit with zero TTL")
	}
}
