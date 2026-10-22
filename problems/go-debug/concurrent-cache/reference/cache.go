// Package cache provides a small in-memory cache with per-entry expiry.
package cache

import (
	"sync"
	"time"
)

const defaultCleanupInterval = time.Minute

type entry struct {
	value     string
	expiresAt time.Time
}

// Cache is a string cache whose entries expire after a fixed TTL.
// It is safe for concurrent use. Call Close to stop the background cleanup.
type Cache struct {
	mu     sync.RWMutex
	items  map[string]entry
	ttl    time.Duration
	hits   int
	misses int

	stop chan struct{}
	once sync.Once
}

// New returns a Cache whose entries live for ttl.
func New(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]entry),
		ttl:   ttl,
		stop:  make(chan struct{}),
	}
	interval := ttl
	if interval <= 0 {
		interval = defaultCleanupInterval
	}
	go c.janitor(interval)
	return c
}

// Close stops the background cleanup goroutine. It is safe to call more than once.
func (c *Cache) Close() {
	c.once.Do(func() { close(c.stop) })
}

func (c *Cache) janitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.Cleanup()
		case <-c.stop:
			return
		}
	}
}

// Set stores value under key.
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = entry{value: value, expiresAt: time.Now().Add(c.ttl)}
}

// Get returns the value for key if it is present and has not expired.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[key]
	if !ok {
		c.misses++
		return "", false
	}
	if expired(e, time.Now()) {
		delete(c.items, key)
		c.misses++
		return "", false
	}
	c.hits++
	return e.value, true
}

// Delete removes key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func expired(e entry, now time.Time) bool {
	return !now.Before(e.expiresAt)
}

// Cleanup removes expired entries and returns how many were removed.
func (c *Cache) Cleanup() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	removed := 0
	for k, e := range c.items {
		if expired(e, now) {
			delete(c.items, k)
			removed++
		}
	}
	return removed
}

// Len returns the number of entries currently stored.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Stats returns the number of cache hits and misses so far.
func (c *Cache) Stats() (hits, misses int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hits, c.misses
}
