package main

// TTLCache: each entry remembers its write time (expiry never moves)
// and its last-touched sequence (recency does). A Put at capacity
// first drops expired corpses -- they don't deserve an eviction -- and
// only then evicts the least recently used live entry. O(n) eviction
// scan: fine at exercise scale, and the honest place to say "a real
// one pairs the map with a linked list".
type entry struct {
	value     int
	writtenAt int64
	touched   int64 // monotonic sequence, not time: recency ordering
}

type TTLCache struct {
	capacity  int
	ttlMillis int64
	items     map[string]*entry
	seq       int64
}

func NewTTLCache(capacity int, ttlMillis int64) *TTLCache {
	return &TTLCache{capacity: capacity, ttlMillis: ttlMillis, items: make(map[string]*entry)}
}

func (c *TTLCache) expired(e *entry, nowMillis int64) bool {
	return nowMillis-e.writtenAt >= c.ttlMillis
}

func (c *TTLCache) PutAt(key string, value int, nowMillis int64) {
	c.seq++
	if e, ok := c.items[key]; ok {
		e.value = value
		e.writtenAt = nowMillis
		e.touched = c.seq
		return
	}

	// Purge corpses first; only evict a live entry if still needed.
	for k, e := range c.items {
		if c.expired(e, nowMillis) {
			delete(c.items, k)
		}
	}
	if len(c.items) >= c.capacity {
		var lruKey string
		var lruSeq int64 = 1 << 62
		for k, e := range c.items {
			if e.touched < lruSeq {
				lruSeq = e.touched
				lruKey = k
			}
		}
		delete(c.items, lruKey)
	}
	c.items[key] = &entry{value: value, writtenAt: nowMillis, touched: c.seq}
}

func (c *TTLCache) GetAt(key string, nowMillis int64) (int, bool) {
	e, ok := c.items[key]
	if !ok || c.expired(e, nowMillis) {
		return 0, false
	}
	c.seq++
	e.touched = c.seq // recency refreshed; writtenAt deliberately not
	return e.value, true
}
