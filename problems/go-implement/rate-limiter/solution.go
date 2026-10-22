package solution

import "time"

// RateLimiter is a per-key token bucket rate limiter.
type RateLimiter struct {
	// TODO: add fields
	now func() time.Time
}

// NewRateLimiter creates a limiter where every key gets a bucket of the given
// capacity, refilled at refillPerSec tokens per second. now supplies the
// current time.
func NewRateLimiter(capacity int, refillPerSec float64, now func() time.Time) *RateLimiter {
	// TODO: implement
	return &RateLimiter{now: now}
}

// Allow reports whether one token could be taken from key's bucket.
func (r *RateLimiter) Allow(key string) bool {
	// TODO: implement
	return false
}

// AllowN reports whether n tokens could be taken from key's bucket at once.
// n <= 0 or n > capacity is never allowed.
func (r *RateLimiter) AllowN(key string, n int) bool {
	// TODO: implement
	return false
}
