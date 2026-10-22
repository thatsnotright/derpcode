# Token Bucket Rate Limiter

Implement a per-key rate limiter using the **token bucket** algorithm.

Each key has its own bucket that holds up to `capacity` tokens and starts **full**. Tokens refill continuously at `refillPerSec` tokens per second (fractional tokens accumulate), never exceeding `capacity`. A request for `n` tokens is allowed only if the bucket currently holds at least `n` tokens; if allowed, the tokens are consumed. A denied request consumes nothing.

Time comes from an injected clock function so behavior is deterministic and testable; never call `time.Now` directly.

## Methods to implement

```go
// NewRateLimiter creates a limiter where every key gets a bucket of the given
// capacity, refilled at refillPerSec tokens per second. now supplies the
// current time.
func NewRateLimiter(capacity int, refillPerSec float64, now func() time.Time) *RateLimiter

// Allow reports whether one token could be taken from key's bucket.
func (r *RateLimiter) Allow(key string) bool

// AllowN reports whether n tokens could be taken from key's bucket at once.
// n <= 0 or n > capacity is never allowed.
func (r *RateLimiter) AllowN(key string, n int) bool
```

## Example

```go
// capacity 3, refill 1 token/sec
rl := NewRateLimiter(3, 1, clock)
rl.Allow("u") // true  (3 -> 2)
rl.Allow("u") // true  (2 -> 1)
rl.Allow("u") // true  (1 -> 0)
rl.Allow("u") // false
// advance clock by 2s
rl.Allow("u") // true  (2 -> 1)
rl.Allow("v") // true  (different key, own full bucket)
```

## Constraints

- Must be safe for concurrent use by multiple goroutines.
- If the clock goes backwards, treat it as no time passing.
- Keys are arbitrary strings; buckets are created lazily.

## Follow-ups

1. Buckets for keys that are never seen again leak memory. How would you evict idle buckets?
2. Compare token bucket with fixed-window, sliding-window-log and leaky bucket. Trade-offs?
3. How would you make this work across several servers sharing a limit?
