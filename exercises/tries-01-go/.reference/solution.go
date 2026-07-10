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
	node := t
	for i := 0; i < len(word); i++ {
		c := word[i]
		if node.children[c] == nil {
			node.children[c] = NewTrie()
		}
		node = node.children[c]
	}
	node.isEnd = true
}

func (t *Trie) Search(word string) bool {
	node := t.find(word)
	return node != nil && node.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
	return t.find(prefix) != nil
}

func (t *Trie) find(s string) *Trie {
	node := t
	for i := 0; i < len(s); i++ {
		c := s[i]
		if node.children[c] == nil {
			return nil
		}
		node = node.children[c]
	}
	return node
}
