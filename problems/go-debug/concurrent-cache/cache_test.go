package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	c := New(time.Minute)
	c.Set("a", "1")
	got, ok := c.Get("a")
	if !ok || got != "1" {
		t.Fatalf("Get(a) = %q, %v; want 1, true", got, ok)
	}
}

func TestExpiry(t *testing.T) {
	c := New(50 * time.Millisecond)
	c.Set("a", "1")
	time.Sleep(80 * time.Millisecond)
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected a to have expired")
	}
}

func TestDelete(t *testing.T) {
	c := New(time.Minute)
	c.Set("a", "1")
	c.Delete("a")
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected a to be deleted")
	}
}

func TestLen(t *testing.T) {
	c := New(time.Minute)
	c.Set("a", "1")
	c.Set("b", "2")
	c.Set("c", "3")
	if got := c.Len(); got != 3 {
		t.Fatalf("Len() = %d; want 3", got)
	}
}

func TestConcurrentAccess(t *testing.T) {
	c := New(time.Minute)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("k%d", j%10)
				c.Set(key, fmt.Sprint(i))
				c.Get(key)
			}
		}(i)
	}
	wg.Wait()
}

func TestCleanupRemovesExpired(t *testing.T) {
	c := New(20 * time.Millisecond)
	c.Set("a", "1")
	c.Set("b", "2")
	time.Sleep(150 * time.Millisecond)
	if got := c.Len(); got != 0 {
		t.Fatalf("Len() = %d after expiry; want 0", got)
	}
}

func TestMissesAreCounted(t *testing.T) {
	c := New(time.Minute)
	c.Get("x")
	c.Get("y")
	_, misses := c.Stats()
	if misses != 2 {
		t.Fatalf("misses = %d; want 2", misses)
	}
}
