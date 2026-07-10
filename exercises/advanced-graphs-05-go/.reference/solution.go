package main

import "sort"

// AlienOrder derives a valid character ordering for the alien alphabet
// implied by words, assumed sorted according to that unknown ordering.
// Returns "" if no valid ordering exists.
func AlienOrder(words []string) string {
	adj := make(map[byte]map[byte]bool)
	inDegree := make(map[byte]int)

	for _, w := range words {
		for i := 0; i < len(w); i++ {
			c := w[i]
			if _, ok := adj[c]; !ok {
				adj[c] = make(map[byte]bool)
				inDegree[c] = 0
			}
		}
	}

	for i := 0; i < len(words)-1; i++ {
		w1, w2 := words[i], words[i+1]
		minLen := len(w1)
		if len(w2) < minLen {
			minLen = len(w2)
		}
		if len(w1) > len(w2) && w1[:minLen] == w2[:minLen] {
			return ""
		}
		for j := 0; j < minLen; j++ {
			if w1[j] != w2[j] {
				if !adj[w1[j]][w2[j]] {
					adj[w1[j]][w2[j]] = true
					inDegree[w2[j]]++
				}
				break
			}
		}
	}

	var queue []byte
	for c, d := range inDegree {
		if d == 0 {
			queue = append(queue, c)
		}
	}
	sort.Slice(queue, func(i, j int) bool { return queue[i] < queue[j] })

	var order []byte
	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]
		order = append(order, c)

		neighbors := make([]byte, 0, len(adj[c]))
		for n := range adj[c] {
			neighbors = append(neighbors, n)
		}
		sort.Slice(neighbors, func(i, j int) bool { return neighbors[i] < neighbors[j] })

		for _, n := range neighbors {
			inDegree[n]--
			if inDegree[n] == 0 {
				queue = append(queue, n)
			}
		}
	}

	if len(order) != len(inDegree) {
		return ""
	}
	return string(order)
}
