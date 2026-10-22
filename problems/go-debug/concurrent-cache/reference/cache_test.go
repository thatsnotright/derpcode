package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func newTestCache(t *testing.T, ttl time.Duration) *Cache {
	t.Helper()
	c := New(ttl)
	t.Cleanup(c.Close)
	return c
}

func TestSetAndGet(t *testing.T) {
	c := newTestCache(t, time.Minute)
	c.Set("a", "1")
	got, ok := c.Get("a")
	if !ok || got != "1" {
		t.Fatalf("Get(a) = %q, %v; want 1, true", got, ok)
	}
}

func TestGetMissing(t *testing.T) {
	c := newTestCache(t, time.Minute)
	if got, ok := c.Get("nope"); ok || got != "" {
		t.Fatalf("Get(nope) = %q, %v; want empty, false", got, ok)
	}
	// A miss must not leave the lock held.
	c.Set("a", "1")
	if _, ok := c.Get("a"); !ok {
		t.Fatal("Get(a) after a miss should succeed")
	}
}

func TestExpiry(t *testing.T) {
	tests := []struct {
		name  string
		sleep time.Duration
		want  bool
	}{
		{"before ttl", 0, true},
		{"after ttl", 120 * time.Millisecond, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestCache(t, 100*time.Millisecond)
			c.Set("a", "1")
			time.Sleep(tt.sleep)
			if _, ok := c.Get("a"); ok != tt.want {
				t.Fatalf("Get(a) ok = %v; want %v", ok, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	c := newTestCache(t, time.Minute)
	c.Set("a", "1")
	c.Delete("a")
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected a to be deleted")
	}
}

func TestLen(t *testing.T) {
	c := newTestCache(t, time.Minute)
	c.Set("a", "1")
	c.Set("b", "2")
	c.Set("c", "3")
	if got := c.Len(); got != 3 {
		t.Fatalf("Len() = %d; want 3", got)
	}
}

func TestConcurrentAccess(t *testing.T) {
	c := newTestCache(t, time.Minute)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("k%d", j%10)
				c.Set(key, fmt.Sprint(i))
				c.Get(key)
				c.Len()
				c.Stats()
				c.Cleanup()
			}
		}(i)
	}
	wg.Wait()
}

func TestCleanupRemovesExpired(t *testing.T) {
	c := newTestCache(t, 20*time.Millisecond)
	c.Set("a", "1")
	c.Set("b", "2")
	deadline := time.Now().Add(2 * time.Second)
	for c.Len() != 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if got := c.Len(); got != 0 {
		t.Fatalf("Len() = %d after expiry; want 0", got)
	}
}

func TestStats(t *testing.T) {
	c := newTestCache(t, time.Minute)
	c.Get("x")
	c.Get("y")
	c.Set("z", "1")
	c.Get("z")
	hits, misses := c.Stats()
	if hits != 1 || misses != 2 {
		t.Fatalf("Stats() = %d, %d; want 1, 2", hits, misses)
	}
}
