package main

import (
	"fmt"
	"testing"
)

func buildRing() *Ring {
	r := NewRing(100)
	r.AddNode("node-a")
	r.AddNode("node-b")
	r.AddNode("node-c")
	return r
}

func TestDeterministicLookup(t *testing.T) {
	r := buildRing()
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		first := r.Lookup(key)
		if first == "" {
			t.Fatalf("Lookup(%q) returned no node on a populated ring", key)
		}
		if again := r.Lookup(key); again != first {
			t.Fatalf("Lookup(%q) unstable: %q then %q", key, first, again)
		}
	}
}

func TestReasonableBalance(t *testing.T) {
	r := buildRing()
	counts := map[string]int{}
	const keys = 10000
	for i := 0; i < keys; i++ {
		counts[r.Lookup(fmt.Sprintf("key-%d", i))]++
	}
	for _, node := range []string{"node-a", "node-b", "node-c"} {
		share := float64(counts[node]) / keys
		if share < 0.10 || share > 0.60 {
			t.Fatalf("%s owns %.0f%% of keys, want a sane share (10-60%%) with 100 vnodes", node, share*100)
		}
	}
}

func TestAddRemapsOnlyANeighborhoodAndRemoveRestores(t *testing.T) {
	r := buildRing()
	const keys = 10000
	before := make(map[string]string, keys)
	for i := 0; i < keys; i++ {
		k := fmt.Sprintf("key-%d", i)
		before[k] = r.Lookup(k)
	}

	r.AddNode("node-d")
	moved := 0
	for k, owner := range before {
		if r.Lookup(k) != owner {
			moved++
		}
	}
	if moved*2 >= keys {
		t.Fatalf("adding one node moved %d of %d keys -- %%N-style rehashing moves nearly everything; consistent hashing must not", moved, keys)
	}
	if moved*20 < keys {
		t.Fatalf("adding a node with equal vnodes moved only %d of %d keys -- the new node isn't taking its share", moved, keys)
	}

	r.RemoveNode("node-d")
	for k, owner := range before {
		if got := r.Lookup(k); got != owner {
			t.Fatalf("after removing node-d, Lookup(%q) = %q, want the original %q restored exactly", k, got, owner)
		}
	}
}

func TestEmptyRing(t *testing.T) {
	r := NewRing(100)
	if got := r.Lookup("anything"); got != "" {
		t.Fatalf("Lookup on an empty ring = %q, want \"\"", got)
	}
}
