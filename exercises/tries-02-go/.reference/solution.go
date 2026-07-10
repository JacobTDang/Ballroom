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
	node := d
	for i := 0; i < len(word); i++ {
		c := word[i]
		if node.children[c] == nil {
			node.children[c] = NewWordDictionary()
		}
		node = node.children[c]
	}
	node.isEnd = true
}

func (d *WordDictionary) Search(word string) bool {
	return d.searchFrom(word, 0)
}

func (d *WordDictionary) searchFrom(word string, idx int) bool {
	node := d
	for i := idx; i < len(word); i++ {
		c := word[i]
		if c == '.' {
			for _, child := range node.children {
				if child.searchFrom(word, i+1) {
					return true
				}
			}
			return false
		}
		if node.children[c] == nil {
			return false
		}
		node = node.children[c]
	}
	return node.isEnd
}
