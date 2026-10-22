package solution

import "container/list"

type entry struct {
	key, value int
}

// LRUCache is a fixed-capacity cache that evicts the least recently used entry.
type LRUCache struct {
	capacity int
	order    *list.List // front = most recently used
	items    map[int]*list.Element
}

// NewLRUCache returns an empty cache holding at most capacity entries.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		order:    list.New(),
		items:    make(map[int]*list.Element),
	}
}

// Get returns the value for key and true if present, marking the key as most
// recently used. If the key is absent it returns 0, false.
func (c *LRUCache) Get(key int) (int, bool) {
	el, ok := c.items[key]
	if !ok {
		return 0, false
	}
	c.order.MoveToFront(el)
	return el.Value.(*entry).value, true
}

// Put inserts or updates key and marks it as most recently used. If the cache
// is over capacity afterwards, the least recently used entry is evicted.
func (c *LRUCache) Put(key, value int) {
	if c.capacity <= 0 {
		return
	}
	if el, ok := c.items[key]; ok {
		el.Value.(*entry).value = value
		c.order.MoveToFront(el)
		return
	}
	c.items[key] = c.order.PushFront(&entry{key, value})
	if c.order.Len() > c.capacity {
		oldest := c.order.Back()
		c.order.Remove(oldest)
		delete(c.items, oldest.Value.(*entry).key)
	}
}

// Len returns the number of entries currently stored.
func (c *LRUCache) Len() int {
	return c.order.Len()
}
