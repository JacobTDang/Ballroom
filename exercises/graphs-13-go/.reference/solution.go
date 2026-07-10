package main

// LadderLength returns the number of words in the shortest
// transformation sequence from beginWord to endWord, changing one
// letter at a time through words in wordList, or 0 if impossible.
func LadderLength(beginWord string, endWord string, wordList []string) int {
	wordSet := make(map[string]bool, len(wordList))
	for _, w := range wordList {
		wordSet[w] = true
	}
	if !wordSet[endWord] {
		return 0
	}

	patterns := make(map[string][]string)
	addPatterns := func(word string) {
		b := []byte(word)
		for i := 0; i < len(b); i++ {
			orig := b[i]
			b[i] = '*'
			pattern := string(b)
			patterns[pattern] = append(patterns[pattern], word)
			b[i] = orig
		}
	}
	for w := range wordSet {
		addPatterns(w)
	}
	addPatterns(beginWord)

	visited := map[string]bool{beginWord: true}
	queue := []string{beginWord}
	steps := 1

	for len(queue) > 0 {
		next := []string{}
		for _, word := range queue {
			if word == endWord {
				return steps
			}
			b := []byte(word)
			for i := 0; i < len(b); i++ {
				orig := b[i]
				b[i] = '*'
				pattern := string(b)
				for _, neighbor := range patterns[pattern] {
					if !visited[neighbor] {
						visited[neighbor] = true
						next = append(next, neighbor)
					}
				}
				b[i] = orig
			}
		}
		queue = next
		steps++
	}
	return 0
}
