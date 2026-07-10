package main

// ValidTree reports whether the n nodes and given undirected edges
// form a valid tree (connected, no cycles).
func ValidTree(n int, edges [][]int) bool {
	if len(edges) != n-1 {
		return false
	}

	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}

	var find func(x int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}

	for _, e := range edges {
		rootA, rootB := find(e[0]), find(e[1])
		if rootA == rootB {
			return false
		}
		parent[rootA] = rootB
	}
	return true
}
