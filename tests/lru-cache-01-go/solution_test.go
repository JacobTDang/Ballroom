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

func TestLRUCache_UpdatingExistingKeyDoesNotEvict(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(1, 10) // update, not a new insert -- must not evict 2
	if v, ok := c.Get(2); !ok || v != 2 {
		t.Fatalf("Get(2) = %d, %v, want 2, true -- should not have been evicted by an update", v, ok)
	}
	if v, ok := c.Get(1); !ok || v != 10 {
		t.Fatalf("Get(1) = %d, %v, want 10, true", v, ok)
	}
}

func TestLRUCache_GetRefreshesRecency(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Get(1)    // 1 is now most recently used
	c.Put(3, 3) // should evict 2, not 1
	if _, ok := c.Get(2); ok {
		t.Fatal("Get(2) should be evicted (1 was refreshed by Get)")
	}
	if v, ok := c.Get(1); !ok || v != 1 {
		t.Fatalf("Get(1) = %d, %v, want 1, true -- should have survived", v, ok)
	}
}

func TestLRUCache_CapacityOneEvictsImmediately(t *testing.T) {
	c := NewLRUCache(1)
	c.Put(1, 1)
	c.Put(2, 2)
	if _, ok := c.Get(1); ok {
		t.Fatal("Get(1) should be evicted, capacity is 1")
	}
	if v, ok := c.Get(2); !ok || v != 2 {
		t.Fatalf("Get(2) = %d, %v, want 2, true", v, ok)
	}
}

func TestLRUCache_MissingKeyReturnsFalse(t *testing.T) {
	c := NewLRUCache(2)
	if _, ok := c.Get(999); ok {
		t.Fatal("expected Get on a never-inserted key to return ok=false")
	}
}
