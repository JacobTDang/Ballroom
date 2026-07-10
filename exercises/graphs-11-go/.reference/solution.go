package main

// CountComponents returns the number of connected components in the
// undirected graph of n nodes described by edges.
func CountComponents(n int, edges [][]int) int {
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

	components := n
	for _, e := range edges {
		rootA, rootB := find(e[0]), find(e[1])
		if rootA != rootB {
			parent[rootA] = rootB
			components--
		}
	}
	return components
}
