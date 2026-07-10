package main

// FindRedundantConnection returns the edge that can be removed to
// turn the graph back into a tree, using union-find to detect the
// first edge that connects two already-connected nodes.
func FindRedundantConnection(edges [][]int) []int {
	n := len(edges)
	parent := make([]int, n+1)
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
			return e
		}
		parent[rootA] = rootB
	}
	return nil
}
