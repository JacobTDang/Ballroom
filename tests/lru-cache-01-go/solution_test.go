package main

import "testing"

func TestLRUCache(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 100)
	c.Put(2, 200)
	if v, ok := c.Get(1); !ok || v != 100 {
		t.Fatalf("Get(1) = %d, %v, want 100, true", v, ok)
	}

	c.Put(3, 300) // evicts 2
	if _, ok := c.Get(2); ok {
		t.Fatal("Get(2) should be evicted")
	}
	if v, ok := c.Get(3); !ok || v != 300 {
		t.Fatalf("Get(3) = %d, %v, want 300, true", v, ok)
	}

	c.Put(4, 400) // evicts 1
	if _, ok := c.Get(1); ok {
		t.Fatal("Get(1) should be evicted")
	}
	if v, ok := c.Get(3); !ok || v != 300 {
		t.Fatalf("Get(3) = %d, %v, want 300, true", v, ok)
	}
	if v, ok := c.Get(4); !ok || v != 400 {
		t.Fatalf("Get(4) = %d, %v, want 400, true", v, ok)
	}
}
