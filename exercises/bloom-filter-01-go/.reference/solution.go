package main

import "hash/fnv"

// BloomFilter with double hashing: two independent FNV-1a hashes
// generate all k probe positions as h1 + i*h2 -- the standard trick
// (Kirsch-Mitzenmacher) that makes k hashes as good as k independent
// ones without computing k real hashes.
type BloomFilter struct {
	bits   []bool
	hashes int
}

func NewBloomFilter(bits, hashes int) *BloomFilter {
	return &BloomFilter{bits: make([]bool, bits), hashes: hashes}
}

// mix is a splitmix64-style finalizer: FNV alone clusters badly on
// similar keys over power-of-two table sizes -- the avalanche step is
// what makes the k derived probes behave independently.
func mix(h uint64) uint64 {
	h ^= h >> 30
	h *= 0xBF58476D1CE4E5B9
	h ^= h >> 27
	h *= 0x94D049BB133111EB
	h ^= h >> 31
	return h
}

func (b *BloomFilter) positions(key string) []int {
	f := fnv.New64a()
	f.Write([]byte(key))
	h1 := mix(f.Sum64())
	h2 := mix(h1^0x9E3779B97F4A7C15) | 1 // odd, so it cycles the whole table

	out := make([]int, b.hashes)
	for i := 0; i < b.hashes; i++ {
		out[i] = int((h1 + uint64(i)*h2) % uint64(len(b.bits)))
	}
	return out
}

func (b *BloomFilter) Add(key string) {
	for _, p := range b.positions(key) {
		b.bits[p] = true
	}
}

func (b *BloomFilter) MightContain(key string) bool {
	for _, p := range b.positions(key) {
		if !b.bits[p] {
			return false
		}
	}
	return true
}
