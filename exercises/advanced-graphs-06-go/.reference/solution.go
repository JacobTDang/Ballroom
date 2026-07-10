package main

// FindCheapestPrice returns the cheapest price from src to dst using at
// most k stops (at most k+1 edges), or -1 if impossible.
func FindCheapestPrice(n int, flights [][]int, src int, dst int, k int) int {
	const inf = 1 << 30

	dist := make([]int, n)
	for i := range dist {
		dist[i] = inf
	}
	dist[src] = 0

	// Bellman-Ford limited to exactly k+1 relaxation rounds. Each round
	// must relax edges using a SNAPSHOT of the previous round's
	// distances, not the array being updated in place during that same
	// round, or a single round could silently chain multiple edges
	// together and violate the stop limit.
	for round := 0; round <= k; round++ {
		prev := make([]int, n)
		copy(prev, dist)

		for _, f := range flights {
			u, v, price := f[0], f[1], f[2]
			if prev[u] == inf {
				continue
			}
			if prev[u]+price < dist[v] {
				dist[v] = prev[u] + price
			}
		}
	}

	if dist[dst] == inf {
		return -1
	}
	return dist[dst]
}
