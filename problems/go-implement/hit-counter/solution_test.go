package solution

import (
	"sync"
	"testing"
	"time"
)

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func at(d time.Duration) time.Time { return base.Add(d) }

func TestCountWindowBoundaries(t *testing.T) {
	h := NewHitCounter(10 * time.Second)
	h.Hit(at(0))
	h.Hit(at(5 * time.Second))
	h.Hit(at(10 * time.Second))
	tests := []struct {
		name string
		now  time.Duration
		want int
	}{
		{"just before first ages out (future hit not counted)", 9*time.Second + 999*time.Millisecond, 2},
		{"first ages out exactly", 10 * time.Second, 2},
		{"second ages out exactly", 15 * time.Second, 1},
		{"all aged out", 20 * time.Second, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fresh counter per case so non-decreasing time holds.
			c := NewHitCounter(10 * time.Second)
			c.Hit(at(0))
			c.Hit(at(5 * time.Second))
			c.Hit(at(10 * time.Second))
			if got := c.Count(at(tt.now)); got != tt.want {
				t.Errorf("Count(%v) = %d, want %d", tt.now, got, tt.want)
			}
		})
	}
}

func TestEmptyCounter(t *testing.T) {
	h := NewHitCounter(time.Minute)
	if got := h.Count(base); got != 0 {
		t.Errorf("Count on empty = %d, want 0", got)
	}
}

func TestExampleFromReadme(t *testing.T) {
	h := NewHitCounter(5 * time.Minute)
	h.Hit(at(0))
	h.Hit(at(time.Minute))
	h.Hit(at(2 * time.Minute))
	if got := h.Count(at(2 * time.Minute)); got != 3 {
		t.Errorf("Count(+2m) = %d, want 3", got)
	}
	if got := h.Count(at(5 * time.Minute)); got != 2 {
		t.Errorf("Count(+5m) = %d, want 2", got)
	}
	if got := h.Count(at(time.Hour)); got != 0 {
		t.Errorf("Count(+1h) = %d, want 0", got)
	}
}

func TestSameTimestampHits(t *testing.T) {
	h := NewHitCounter(time.Second)
	for i := 0; i < 50; i++ {
		h.Hit(at(0))
	}
	if got := h.Count(at(0)); got != 50 {
		t.Errorf("Count = %d, want 50", got)
	}
	if got := h.Count(at(999 * time.Millisecond)); got != 50 {
		t.Errorf("Count = %d, want 50", got)
	}
	if got := h.Count(at(time.Second)); got != 0 {
		t.Errorf("Count = %d, want 0", got)
	}
}

func TestSubSecondPrecision(t *testing.T) {
	h := NewHitCounter(100 * time.Millisecond)
	h.Hit(at(0))
	h.Hit(at(60 * time.Millisecond))
	if got := h.Count(at(99 * time.Millisecond)); got != 2 {
		t.Errorf("Count(99ms) = %d, want 2", got)
	}
	if got := h.Count(at(100 * time.Millisecond)); got != 1 {
		t.Errorf("Count(100ms) = %d, want 1", got)
	}
	if got := h.Count(at(160 * time.Millisecond)); got != 0 {
		t.Errorf("Count(160ms) = %d, want 0", got)
	}
}

func TestSlidingOverLongStream(t *testing.T) {
	h := NewHitCounter(10 * time.Second)
	for i := 0; i < 1000; i++ {
		now := at(time.Duration(i) * time.Second)
		h.Hit(now)
		want := i + 1
		if want > 10 {
			want = 10
		}
		if got := h.Count(now); got != want {
			t.Fatalf("at second %d: Count = %d, want %d", i, got, want)
		}
	}
}

func TestConcurrentHits(t *testing.T) {
	h := NewHitCounter(time.Minute)
	const goroutines, perG = 16, 200
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < perG; i++ {
				h.Hit(at(0))
			}
		}()
	}
	wg.Wait()
	if got := h.Count(at(30 * time.Second)); got != goroutines*perG {
		t.Errorf("Count = %d, want %d", got, goroutines*perG)
	}
}

func TestConcurrentHitAndCount(t *testing.T) {
	h := NewHitCounter(time.Minute)
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				h.Hit(at(0))
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				h.Count(at(0))
			}
		}()
	}
	wg.Wait()
	if got := h.Count(at(0)); got != 8*200 {
		t.Errorf("Count = %d, want %d", got, 8*200)
	}
}
