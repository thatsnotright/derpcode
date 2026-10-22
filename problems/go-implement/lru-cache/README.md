# LRU Cache

Implement a fixed-capacity cache that evicts the **least recently used** entry when it is full.

Both reading a key (`Get`) and writing a key (`Put`) count as "using" it. `Get` and `Put` should run in O(1) average time.

## Methods to implement

```go
// NewLRUCache returns an empty cache holding at most capacity entries.
func NewLRUCache(capacity int) *LRUCache

// Get returns the value for key and true if present, marking the key as most
// recently used. If the key is absent it returns 0, false.
func (c *LRUCache) Get(key int) (int, bool)

// Put inserts or updates key and marks it as most recently used. If the cache
// is over capacity afterwards, the least recently used entry is evicted.
func (c *LRUCache) Put(key, value int)

// Len returns the number of entries currently stored.
func (c *LRUCache) Len() int
```

## Example

```go
c := NewLRUCache(2)
c.Put(1, 10)
c.Put(2, 20)
c.Get(1)      // 10, true   (key 1 is now most recent)
c.Put(3, 30)  // evicts key 2
c.Get(2)      // 0, false
c.Get(3)      // 30, true
```

## Constraints

- `capacity >= 0`. A cache with capacity 0 never stores anything.
- Keys and values are `int`; a stored value of `0` is distinct from a missing key.
- Single-threaded use is fine.

## Follow-ups

1. Make the cache safe for concurrent use. How would you reduce lock contention?
2. Add per-entry TTL expiry. When do you evict expired entries?
3. Add an `OnEvict(key, value int)` callback. What are the pitfalls of calling user code while holding a lock?
