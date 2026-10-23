# Hit Counter

Implement a counter that records hits (e.g. page views) and reports how many occurred within a sliding time window. It must be **safe for concurrent use**.

Time is passed in explicitly, so the implementation must not read the wall clock or sleep.

## Methods to implement

```go
// NewHitCounter returns a counter whose window is the given duration.
func NewHitCounter(window time.Duration) *HitCounter

// Hit records one hit at time t.
func (h *HitCounter) Hit(t time.Time)

// Count returns the number of recorded hits in the half-open interval
// (t-window, t]: a hit exactly window before t is no longer counted, and a
// hit exactly at t is.
func (h *HitCounter) Count(t time.Time) int
```

## Example

```go
h := NewHitCounter(5 * time.Minute)
base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
h.Hit(base)
h.Hit(base.Add(time.Minute))
h.Hit(base.Add(2 * time.Minute))
h.Count(base.Add(2 * time.Minute)) // 3
h.Count(base.Add(5 * time.Minute)) // 2  (the hit at base has aged out)
h.Count(base.Add(time.Hour))       // 0
```

## Constraints

- Calls are made with non-decreasing times (`Hit` and `Count` times never go backwards relative to calls already made), so old hits can be discarded.
- Many hits may share the same timestamp.
- Memory should stay proportional to the hits inside the window, not to all hits ever recorded.
- Methods may be called from many goroutines; `go test -race` must be clean.

## Follow-ups

1. If you get millions of hits per second, storing each timestamp is too heavy. How would you bucket hits (e.g. per-second ring buffer) and what precision do you give up?
2. Reduce lock contention: sharded counters? atomics?
3. How would you make this work across several servers?
