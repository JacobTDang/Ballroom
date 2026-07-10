package main

// WordDictionary supports adding words and searching, where a search
// query may use '.' to match any single character.
type WordDictionary struct {
	children map[byte]*WordDictionary
	isEnd    bool
}

func NewWordDictionary() *WordDictionary {
	return &WordDictionary{children: make(map[byte]*WordDictionary)}
}

func (d *WordDictionary) AddWord(word string) {
	// TODO: implement
}

func (d *WordDictionary) Search(word string) bool {
	// TODO: implement
	return false
}
