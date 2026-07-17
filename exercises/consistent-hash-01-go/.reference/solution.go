package main

import (
	"fmt"
	"hash/fnv"
	"sort"
)

// Ring: each node occupies `vnodes` positions on a 32-bit hash ring
// (virtual nodes smooth the balance); a key belongs to the first
// position clockwise from its hash. Sorted positions + binary search
// answer that in O(log n). Removing a node deletes only its own
// positions, so every other key keeps its owner -- that's the whole
// property.
type Ring struct {
	vnodes    int
	positions []uint32
	owner     map[uint32]string
}

func NewRing(vnodes int) *Ring {
	return &Ring{vnodes: vnodes, owner: make(map[uint32]string)}
}

func hash32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (r *Ring) AddNode(name string) {
	for i := 0; i < r.vnodes; i++ {
		p := hash32(fmt.Sprintf("%s#%d", name, i))
		if _, taken := r.owner[p]; taken {
			continue // vanishingly rare 32-bit collision: first owner keeps it
		}
		r.owner[p] = name
		r.positions = append(r.positions, p)
	}
	sort.Slice(r.positions, func(i, j int) bool { return r.positions[i] < r.positions[j] })
}

func (r *Ring) RemoveNode(name string) {
	keep := r.positions[:0]
	for _, p := range r.positions {
		if r.owner[p] == name {
			delete(r.owner, p)
			continue
		}
		keep = append(keep, p)
	}
	r.positions = keep
}

func (r *Ring) Lookup(key string) string {
	if len(r.positions) == 0 {
		return ""
	}
	h := hash32(key)
	i := sort.Search(len(r.positions), func(i int) bool { return r.positions[i] >= h })
	if i == len(r.positions) {
		i = 0 // wrap around the ring
	}
	return r.owner[r.positions[i]]
}
