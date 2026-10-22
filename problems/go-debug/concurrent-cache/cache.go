// Package cache provides a small in-memory cache with per-entry expiry.
package cache

import (
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
}

// Cache is a string cache whose entries expire after a fixed TTL.
type Cache struct {
	mu     sync.RWMutex
	items  map[string]entry
	ttl    time.Duration
	hits   int
	misses int
}

// New returns a Cache whose entries live for ttl.
func New(ttl time.Duration) *Cache {
	c := &Cache{items: make(map[string]entry), ttl: ttl}
	go func() {
		for range time.Tick(ttl) {
			c.Cleanup()
		}
	}()
	return c
}

// Set stores value under key.
func (c *Cache) Set(key, value string) {
	c.mu.RLock()
	c.items[key] = entry{value: value, expiresAt: time.Now().Add(c.ttl)}
	c.mu.RUnlock()
}

// Get returns the value for key if it is present and has not expired.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	e, ok := c.items[key]
	if !ok {
		c.misses++
		return "", false
	}
	if c.expired(e) {
		delete(c.items, key)
		c.misses++
		c.mu.Unlock()
		return "", false
	}
	c.hits++
	c.mu.Unlock()
	return e.value, true
}

// Delete removes key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) expired(e entry) bool {
	return time.Since(e.expiresAt) > c.ttl
}

// Cleanup removes expired entries and returns how many were removed.
func (c *Cache) Cleanup() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	removed := 0
	for k, e := range c.items {
		if c.expired(e) {
			delete(c.items, k)
			removed++
		}
	}
	return removed
}

// Len returns the number of entries currently stored.
func (c *Cache) Len() int {
	return len(c.items)
}

// Stats returns the number of cache hits and misses so far.
func (c *Cache) Stats() (hits, misses int) {
	return c.hits, c.misses
}
