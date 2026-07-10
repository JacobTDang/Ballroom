package main

// NetworkDelayTime returns the minimum time for a signal starting at node
// k to reach every node in a directed weighted graph of n nodes described
// by times[i] = [u, v, w], or -1 if some node is unreachable.
func NetworkDelayTime(times [][]int, n int, k int) int {
	const inf = 1 << 30

	adj := make([][][2]int, n+1) // adj[u] = list of {v, w}
	for _, t := range times {
		u, v, w := t[0], t[1], t[2]
		adj[u] = append(adj[u], [2]int{v, w})
	}

	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = inf
	}
	dist[k] = 0

	visited := make([]bool, n+1)
	for count := 0; count < n; count++ {
		u := -1
		for v := 1; v <= n; v++ {
			if !visited[v] && (u == -1 || dist[v] < dist[u]) {
				u = v
			}
		}
		if u == -1 || dist[u] == inf {
			break
		}
		visited[u] = true
		for _, edge := range adj[u] {
			v, w := edge[0], edge[1]
			if dist[u]+w < dist[v] {
				dist[v] = dist[u] + w
			}
		}
	}

	maxDist := 0
	for v := 1; v <= n; v++ {
		if dist[v] == inf {
			return -1
		}
		if dist[v] > maxDist {
			maxDist = dist[v]
		}
	}
	return maxDist
}
