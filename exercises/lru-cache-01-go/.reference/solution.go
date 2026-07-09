package main

// lruNode is one entry in the intrusive doubly linked list used to track
// recency order.
type lruNode struct {
	key, value int
	prev, next *lruNode
}

// LRUCache is a fixed-capacity cache that evicts the least recently used
// entry when it's full.
type LRUCache struct {
	capacity int
	items    map[int]*lruNode

	// head.next is the most recently used node; tail.prev is the least
	// recently used. head/tail are sentinels, never removed.
	head, tail *lruNode
}

func NewLRUCache(capacity int) *LRUCache {
	head := &lruNode{}
	tail := &lruNode{}
	head.next = tail
	tail.prev = head
	return &LRUCache{capacity: capacity, items: make(map[int]*lruNode), head: head, tail: tail}
}

func (c *LRUCache) unlink(n *lruNode) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (c *LRUCache) pushFront(n *lruNode) {
	n.next = c.head.next
	n.prev = c.head
	c.head.next.prev = n
	c.head.next = n
}

// Get returns the value for key, or (0, false) if not present. Accessing
// a key marks it as most recently used.
func (c *LRUCache) Get(key int) (int, bool) {
	n, ok := c.items[key]
	if !ok {
		return 0, false
	}
	c.unlink(n)
	c.pushFront(n)
	return n.value, true
}

// Put inserts or updates key with value, marking it most recently used.
// If inserting a new key would exceed capacity, evict the least recently
// used entry first.
func (c *LRUCache) Put(key, value int) {
	if n, ok := c.items[key]; ok {
		n.value = value
		c.unlink(n)
		c.pushFront(n)
		return
	}
	if len(c.items) >= c.capacity {
		lru := c.tail.prev
		c.unlink(lru)
		delete(c.items, lru.key)
	}
	n := &lruNode{key: key, value: value}
	c.items[key] = n
	c.pushFront(n)
}
