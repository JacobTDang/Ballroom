package main

// BloomFilter: a bit array + k hashes. Add sets k bits; MightContain
// checks them -- "definitely not" or "probably yes".
//
// TODO: one weak hash into a fixed 64-slot table, ignoring both
// parameters -- everything added is found (no false negatives), but
// the table saturates instantly and almost every absent key collides.
type BloomFilter struct {
	table [64]bool
}

func NewBloomFilter(bits, hashes int) *BloomFilter {
	return &BloomFilter{}
}

func (b *BloomFilter) hash(key string) int {
	h := 0
	for _, c := range key {
		h += int(c)
	}
	return h % 64
}

func (b *BloomFilter) Add(key string) {
	b.table[b.hash(key)] = true
}

func (b *BloomFilter) MightContain(key string) bool {
	return b.table[b.hash(key)]
}
