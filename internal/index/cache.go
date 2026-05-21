package index

import (
	"io"
	"sync"
	"time"
)

// CacheEntry holds a built index along with its build time.
type CacheEntry struct {
	idx       *Index
	builtAt   time.Time
	format    string
}

// Cache stores pre-built indexes keyed by a string identifier (e.g. file path).
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewCache creates a Cache with the given TTL for entries.
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
}

// Get returns a cached index if present and not expired.
func (c *Cache) Get(key string) (*Index, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	if c.ttl > 0 && time.Since(e.builtAt) > c.ttl {
		return nil, false
	}
	return e.idx, true
}

// Set stores a built index under the given key.
func (c *Cache) Set(key string, idx *Index, format string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = &CacheEntry{
		idx:     idx,
		builtAt: time.Now(),
		format:  format,
	}
}

// GetOrBuild returns a cached index or builds one from r, caching the result.
func (c *Cache) GetOrBuild(key string, r io.ReadSeeker, format string) (*Index, error) {
	if idx, ok := c.Get(key); ok {
		return idx, nil
	}
	idx, err := Build(r, format)
	if err != nil {
		return nil, err
	}
	c.Set(key, idx, format)
	return idx, nil
}

// Invalidate removes a cached entry.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}
