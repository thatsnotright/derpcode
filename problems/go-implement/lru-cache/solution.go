package solution

// LRUCache is a fixed-capacity cache that evicts the least recently used entry.
type LRUCache struct {
	// TODO: add fields
}

// NewLRUCache returns an empty cache holding at most capacity entries.
func NewLRUCache(capacity int) *LRUCache {
	// TODO: implement
	return nil
}

// Get returns the value for key and true if present, marking the key as most
// recently used. If the key is absent it returns 0, false.
func (c *LRUCache) Get(key int) (int, bool) {
	// TODO: implement
	return 0, false
}

// Put inserts or updates key and marks it as most recently used. If the cache
// is over capacity afterwards, the least recently used entry is evicted.
func (c *LRUCache) Put(key, value int) {
	// TODO: implement
}

// Len returns the number of entries currently stored.
func (c *LRUCache) Len() int {
	// TODO: implement
	return 0
}
