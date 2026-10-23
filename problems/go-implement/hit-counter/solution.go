package solution

import "time"

// HitCounter counts hits inside a sliding time window. It is safe for
// concurrent use.
type HitCounter struct {
	// TODO: add fields
}

// NewHitCounter returns a counter whose window is the given duration.
func NewHitCounter(window time.Duration) *HitCounter {
	// TODO: implement
	return &HitCounter{}
}

// Hit records one hit at time t.
func (h *HitCounter) Hit(t time.Time) {
	// TODO: implement
}

// Count returns the number of recorded hits in the half-open interval
// (t-window, t]: a hit exactly window before t is no longer counted, and a
// hit exactly at t is.
func (h *HitCounter) Count(t time.Time) int {
	// TODO: implement
	return 0
}
