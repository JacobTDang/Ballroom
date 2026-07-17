package main

// TTLCache: LRU eviction at capacity, plus per-entry expiry ttl after
// the write.
//
// TODO: a plain map -- no eviction, no expiry, no recency. Every rule
// in the problem statement is still yours to build.
type TTLCache struct {
	capacity  int
	ttlMillis int64
	items     map[string]int
}

func NewTTLCache(capacity int, ttlMillis int64) *TTLCache {
	return &TTLCache{capacity: capacity, ttlMillis: ttlMillis, items: make(map[string]int)}
}

func (c *TTLCache) PutAt(key string, value int, nowMillis int64) {
	c.items[key] = value
}

func (c *TTLCache) GetAt(key string, nowMillis int64) (int, bool) {
	v, ok := c.items[key]
	return v, ok
}
