package main

// Trie is a prefix tree over lowercase English letters.
type Trie struct {
	children map[byte]*Trie
	isEnd    bool
}

func NewTrie() *Trie {
	return &Trie{children: make(map[byte]*Trie)}
}

func (t *Trie) Insert(word string) {
	// TODO: implement
}

func (t *Trie) Search(word string) bool {
	// TODO: implement
	return false
}

func (t *Trie) StartsWith(prefix string) bool {
	// TODO: implement
	return false
}
