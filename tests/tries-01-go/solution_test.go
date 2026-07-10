package main

import "testing"

func TestTrie(t *testing.T) {
	trie := NewTrie()
	trie.Insert("apple")
	if !trie.Search("apple") {
		t.Error("Search(apple) = false, want true")
	}
	if trie.Search("app") {
		t.Error("Search(app) = true, want false")
	}
	if !trie.StartsWith("app") {
		t.Error("StartsWith(app) = false, want true")
	}
	trie.Insert("app")
	if !trie.Search("app") {
		t.Error("Search(app) = false, want true (inserted since)")
	}
}

func TestTrie_StartsWithFalseForUnrelatedPrefix(t *testing.T) {
	trie := NewTrie()
	trie.Insert("banana")
	if trie.StartsWith("ban ") {
		t.Error("StartsWith with trailing space should not match")
	}
	if !trie.StartsWith("ban") {
		t.Error("StartsWith(ban) = false, want true")
	}
	if trie.Search("ban") {
		t.Error("Search(ban) = true, want false (only a prefix, not inserted)")
	}
}

func TestTrie_EmptyTrieHasNoMatches(t *testing.T) {
	trie := NewTrie()
	if trie.Search("anything") {
		t.Error("Search on empty trie should be false")
	}
	if trie.StartsWith("a") {
		t.Error("StartsWith on empty trie should be false")
	}
}
