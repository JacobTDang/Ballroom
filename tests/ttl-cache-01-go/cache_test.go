package main

import "testing"

func TestBasicPutGet(t *testing.T) {
	c := NewTTLCache(2, 1000)
	c.PutAt("a", 1, 0)
	if v, ok := c.GetAt("a", 10); !ok || v != 1 {
		t.Fatalf("GetAt(a) = %d,%v want 1,true", v, ok)
	}
	if _, ok := c.GetAt("missing", 10); ok {
		t.Fatal("GetAt(missing) reported a hit")
	}
}

func TestLRUEvictionRespectsRecency(t *testing.T) {
	c := NewTTLCache(2, 100000)
	c.PutAt("a", 1, 0)
	c.PutAt("b", 2, 1)
	c.GetAt("a", 2)     // a is now more recent than b
	c.PutAt("c", 3, 3)  // must evict b, not a
	if _, ok := c.GetAt("b", 4); ok {
		t.Fatal("b survived eviction despite being least recently used")
	}
	if v, ok := c.GetAt("a", 4); !ok || v != 1 {
		t.Fatal("a was evicted despite being recently used")
	}
	if v, ok := c.GetAt("c", 4); !ok || v != 3 {
		t.Fatal("c missing right after insert")
	}
}

func TestTTLExpiryFromWriteTime(t *testing.T) {
	c := NewTTLCache(2, 100)
	c.PutAt("a", 1, 0)
	if _, ok := c.GetAt("a", 99); !ok {
		t.Fatal("entry expired early (99ms old, ttl 100)")
	}
	if _, ok := c.GetAt("a", 100); ok {
		t.Fatal("entry alive at exactly ttl -- expiry is >= ttl after write")
	}
}

func TestGetRefreshesRecencyNotExpiry(t *testing.T) {
	c := NewTTLCache(2, 100)
	c.PutAt("a", 1, 0)
	if _, ok := c.GetAt("a", 99); !ok {
		t.Fatal("setup: a should be alive at 99")
	}
	if _, ok := c.GetAt("a", 100); ok {
		t.Fatal("a alive at 100 -- Get must not extend the TTL")
	}
}

func TestExpiredEntriesDoNotOccupyCapacity(t *testing.T) {
	c := NewTTLCache(2, 100)
	c.PutAt("a", 1, 0)
	c.PutAt("b", 2, 0)
	// Both are corpses at t=200; adding two new keys must not touch a
	// live entry (there are none) and both newcomers must fit.
	c.PutAt("x", 10, 200)
	c.PutAt("y", 20, 201)
	if v, ok := c.GetAt("x", 202); !ok || v != 10 {
		t.Fatal("x missing -- an expired corpse was counted against capacity")
	}
	if v, ok := c.GetAt("y", 202); !ok || v != 20 {
		t.Fatal("y missing -- an expired corpse was counted against capacity")
	}
}

func TestRewriteResetsValueAndExpiry(t *testing.T) {
	c := NewTTLCache(2, 100)
	c.PutAt("a", 1, 0)
	c.PutAt("a", 2, 50) // rewrite: new value, expiry now runs from 50
	if v, ok := c.GetAt("a", 149); !ok || v != 2 {
		t.Fatalf("GetAt(a, 149) = %d,%v want 2,true (rewritten at 50, ttl 100)", v, ok)
	}
	if _, ok := c.GetAt("a", 150); ok {
		t.Fatal("a alive at 150 -- rewrite at 50 sets expiry to 150")
	}
}
