package main

// LRUCache is a fixed-capacity cache that evicts the least recently used
// entry when it's full.
type LRUCache struct {
	capacity int
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{capacity: capacity}
}

// Get returns the value for key, or (0, false) if not present. Accessing
// a key marks it as most recently used.
func (c *LRUCache) Get(key int) (int, bool) {
	// TODO: implement
	return 0, false
}

// Put inserts or updates key with value, marking it most recently used.
// If inserting a new key would exceed capacity, evict the least recently
// used entry first.
func (c *LRUCache) Put(key, value int) {
	// TODO: implement
}
