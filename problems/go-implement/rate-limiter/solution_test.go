package solution

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type fakeClock struct {
	t time.Time
}

func newFakeClock() *fakeClock {
	return &fakeClock{t: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
}

func (c *fakeClock) Now() time.Time          { return c.t }
func (c *fakeClock) Advance(d time.Duration) { c.t = c.t.Add(d) }

func TestStartsFullThenDenies(t *testing.T) {
	clk := newFakeClock()
	rl := NewRateLimiter(3, 1, clk.Now)
	for i := 0; i < 3; i++ {
		if !rl.Allow("u") {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	if rl.Allow("u") {
		t.Error("4th request should be denied")
	}
}

func TestRefillOverTime(t *testing.T) {
	clk := newFakeClock()
	rl := NewRateLimiter(3, 1, clk.Now)
	for i := 0; i < 3; i++ {
		rl.Allow("u")
	}
	clk.Advance(2 * time.Second)
	if !rl.Allow("u") || !rl.Allow("u") {
		t.Error("two tokens should have refilled")
	}
	if rl.Allow("u") {
		t.Error("only two tokens should have refilled")
	}
}

func TestFractionalRefillAccumulates(t *testing.T) {
	clk := newFakeClock()
	rl := NewRateLimiter(2, 2, clk.Now) // 2 tokens per second
	rl.Allow("u")
	rl.Allow("u")
	clk.Advance(250 * time.Millisecond) // 0.5 token
	if rl.Allow("u") {
		t.Error("half a token is not enough")
	}
	clk.Advance(250 * time.Millisecond) // 1.0 token total
	if !rl.Allow("u") {
		t.Error("accumulated fractions should yield one token")
	}
}

func TestRefillCappedAtCapacity(t *testing.T) {
	clk := newFakeClock()
	rl := NewRateLimiter(2, 10, clk.Now)
	clk.Advance(time.Hour)
	if !rl.Allow("u") || !rl.Allow("u") {
		t.Fatal("first two requests should be allowed")
	}
	if rl.Allow("u") {
		t.Error("tokens must not exceed capacity after long idle")
	}
}

func TestKeysAreIndependent(t *testing.T) {
	clk := newFakeClock()
	rl := NewRateLimiter(1, 1, clk.Now)
	if !rl.Allow("a") {
		t.Error("a first request should be allowed")
	}
	if rl.Allow("a") {
		t.Error("a second request should be denied")
	}
	if !rl.Allow("b") {
		t.Error("b should have its own bucket")
	}
}

func TestAllowN(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want bool
	}{
		{"within capacity", 3, true},
		{"exactly capacity", 5, true},
		{"above capacity", 6, false},
		{"zero", 0, false},
		{"negative", -1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := NewRateLimiter(5, 1, newFakeClock().Now)
			if got := rl.AllowN("u", tt.n); got != tt.want {
				t.Errorf("AllowN(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestDeniedRequestConsumesNothing(t *testing.T) {
	clk := newFakeClock()
	rl := NewRateLimiter(5, 1, clk.Now)
	if !rl.AllowN("u", 3) {
		t.Fatal("AllowN(3) should succeed")
	}
	if rl.AllowN("u", 3) { // only 2 left
		t.Fatal("AllowN(3) with 2 tokens should be denied")
	}
	if !rl.AllowN("u", 2) {
		t.Error("denied request must not have consumed tokens")
	}
}

func TestClockGoingBackwardsDoesNotRefill(t *testing.T) {
	clk := newFakeClock()
	rl := NewRateLimiter(1, 1, clk.Now)
	rl.Allow("u")
	clk.Advance(-10 * time.Second)
	if rl.Allow("u") {
		t.Error("backwards clock must not grant tokens")
	}
	clk.Advance(10 * time.Second) // back to the original time
	if rl.Allow("u") {
		t.Error("no real time has passed since the last allowed request")
	}
}

func TestConcurrentUse(t *testing.T) {
	clk := newFakeClock() // frozen: only reads happen concurrently
	rl := NewRateLimiter(50, 1, clk.Now)
	var allowed int64
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if rl.Allow("shared") {
				atomic.AddInt64(&allowed, 1)
			}
		}()
	}
	wg.Wait()
	if allowed != 50 {
		t.Errorf("allowed = %d, want exactly 50", allowed)
	}
}
