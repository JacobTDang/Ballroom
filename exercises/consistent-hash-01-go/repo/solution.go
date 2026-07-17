package main

import (
	"hash/fnv"
	"sort"
)

// Ring maps keys to nodes. Adding or removing a node should only
// remap the keys in its neighborhood.
//
// TODO: hash(key) %% len(nodes) remaps almost EVERY key whenever the
// node count changes -- the exact failure consistent hashing exists
// to fix. (vnodes is ignored here too.)
type Ring struct {
	nodes []string
}

func NewRing(vnodes int) *Ring {
	return &Ring{}
}

func (r *Ring) AddNode(name string) {
	r.nodes = append(r.nodes, name)
	sort.Strings(r.nodes)
}

func (r *Ring) RemoveNode(name string) {
	for i, n := range r.nodes {
		if n == name {
			r.nodes = append(r.nodes[:i], r.nodes[i+1:]...)
			return
		}
	}
}

func (r *Ring) Lookup(key string) string {
	if len(r.nodes) == 0 {
		return ""
	}
	h := fnv.New32a()
	h.Write([]byte(key))
	return r.nodes[int(h.Sum32())%len(r.nodes)]
}
