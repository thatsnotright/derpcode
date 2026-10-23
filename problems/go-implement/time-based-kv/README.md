# Time-Based Key-Value Store

Implement a key-value store that keeps every value a key has ever been set to, tagged with a timestamp, and can look up the value a key had at any point in time.

## Methods to implement

```go
// NewTimeMap returns an empty store.
func NewTimeMap() *TimeMap

// Set records that key had value at the given timestamp. Setting the same
// (key, timestamp) again overwrites the earlier value.
func (m *TimeMap) Set(key, value string, timestamp int)

// Get returns the value whose timestamp is the largest one <= timestamp for
// key, or "" if there is none.
func (m *TimeMap) Get(key string, timestamp int) string
```

## Example

```go
m := NewTimeMap()
m.Set("foo", "bar", 1)
m.Set("foo", "baz", 4)
m.Get("foo", 1)  // "bar"
m.Get("foo", 3)  // "bar"  (latest timestamp <= 3 is 1)
m.Get("foo", 4)  // "baz"
m.Get("foo", 0)  // ""     (nothing at or before 0)
m.Get("nope", 9) // ""
```

## Constraints

- Timestamps for a key **may arrive out of order**: `Set(k, v, 10)` can be followed by `Set(k, w, 5)`.
- Timestamps are `int` and may be zero or negative.
- `Get` should be O(log n) in the number of entries for the key.
- Single-threaded use is fine.

## Follow-ups

1. If timestamps are guaranteed to arrive in increasing order, how does `Set` simplify?
2. Make the store safe for concurrent use with many readers. Which lock would you pick?
3. Old entries pile up forever. How would you add a retention policy without breaking `Get` for recent times?
