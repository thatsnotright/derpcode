package solution

import (
	"sync"
	"time"
)

type bucket struct {
	tokens float64
	last   time.Time
}

// RateLimiter is a per-key token bucket rate limiter.
type RateLimiter struct {
	capacity     float64
	refillPerSec float64
	now          func() time.Time

	mu      sync.Mutex
	buckets map[string]*bucket
}

// NewRateLimiter creates a limiter where every key gets a bucket of the given
// capacity, refilled at refillPerSec tokens per second. now supplies the
// current time.
func NewRateLimiter(capacity int, refillPerSec float64, now func() time.Time) *RateLimiter {
	return &RateLimiter{
		capacity:     float64(capacity),
		refillPerSec: refillPerSec,
		now:          now,
		buckets:      make(map[string]*bucket),
	}
}

// Allow reports whether one token could be taken from key's bucket.
func (r *RateLimiter) Allow(key string) bool {
	return r.AllowN(key, 1)
}

// AllowN reports whether n tokens could be taken from key's bucket at once.
// n <= 0 or n > capacity is never allowed.
func (r *RateLimiter) AllowN(key string, n int) bool {
	if n <= 0 || float64(n) > r.capacity {
		return false
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()
	b, ok := r.buckets[key]
	if !ok {
		b = &bucket{tokens: r.capacity, last: now}
		r.buckets[key] = b
	}

	if elapsed := now.Sub(b.last); elapsed > 0 {
		b.tokens += elapsed.Seconds() * r.refillPerSec
		if b.tokens > r.capacity {
			b.tokens = r.capacity
		}
		b.last = now
	}

	if b.tokens < float64(n) {
		return false
	}
	b.tokens -= float64(n)
	return true
}
