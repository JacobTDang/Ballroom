package main

// trieNode is a prefix-tree node used to search for every word in
// words simultaneously during the board DFS.
type trieNode struct {
	children map[byte]*trieNode
	word     string // non-empty exactly at nodes completing a full word
}

func newTrieNode() *trieNode {
	return &trieNode{children: make(map[byte]*trieNode)}
}

// FindWords returns every word from words that can be traced out on
// board via sequentially adjacent cells, each cell used at most once
// per word.
func FindWords(board [][]byte, words []string) []string {
	root := newTrieNode()
	for _, w := range words {
		node := root
		for i := 0; i < len(w); i++ {
			c := w[i]
			if node.children[c] == nil {
				node.children[c] = newTrieNode()
			}
			node = node.children[c]
		}
		node.word = w
	}

	rows, cols := len(board), len(board[0])
	var res []string

	var dfs func(r, c int, node *trieNode)
	dfs = func(r, c int, node *trieNode) {
		if r < 0 || r >= rows || c < 0 || c >= cols {
			return
		}
		ch := board[r][c]
		if ch == '#' {
			return
		}
		next, ok := node.children[ch]
		if !ok {
			return
		}
		if next.word != "" {
			res = append(res, next.word)
			next.word = "" // don't report the same word twice
		}
		board[r][c] = '#'
		dfs(r+1, c, next)
		dfs(r-1, c, next)
		dfs(r, c+1, next)
		dfs(r, c-1, next)
		board[r][c] = ch
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			dfs(r, c, root)
		}
	}
	return res
}
