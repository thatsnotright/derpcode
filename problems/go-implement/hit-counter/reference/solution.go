package solution

import (
	"sync"
	"time"
)

// HitCounter counts hits inside a sliding time window. It is safe for
// concurrent use.
type HitCounter struct {
	window time.Duration

	mu   sync.Mutex
	hits []time.Time // non-decreasing; hits older than the window are dropped
}

// NewHitCounter returns a counter whose window is the given duration.
func NewHitCounter(window time.Duration) *HitCounter {
	return &HitCounter{window: window}
}

// Hit records one hit at time t.
func (h *HitCounter) Hit(t time.Time) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.hits = append(h.hits, t)
}

// Count returns the number of recorded hits in the half-open interval
// (t-window, t]: a hit exactly window before t is no longer counted, and a
// hit exactly at t is.
func (h *HitCounter) Count(t time.Time) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	cutoff := t.Add(-h.window)
	drop := 0
	for drop < len(h.hits) && !h.hits[drop].After(cutoff) {
		drop++
	}
	if drop > 0 {
		// Copy down so the backing array does not retain expired hits forever.
		n := copy(h.hits, h.hits[drop:])
		h.hits = h.hits[:n]
	}

	count := 0
	for _, hit := range h.hits {
		if !hit.After(t) {
			count++
		}
	}
	return count
}
