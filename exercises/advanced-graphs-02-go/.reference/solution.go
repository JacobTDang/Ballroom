package main

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// MinCostConnectPoints returns the minimum total cost to connect all
// points, where the cost between two points is their Manhattan distance
// (the minimum spanning tree over the complete graph of points).
func MinCostConnectPoints(points [][]int) int {
	n := len(points)
	if n <= 1 {
		return 0
	}

	inTree := make([]bool, n)
	minDist := make([]int, n)
	for i := range minDist {
		minDist[i] = 1 << 30
	}
	minDist[0] = 0

	total := 0
	for count := 0; count < n; count++ {
		u := -1
		for v := 0; v < n; v++ {
			if !inTree[v] && (u == -1 || minDist[v] < minDist[u]) {
				u = v
			}
		}
		inTree[u] = true
		total += minDist[u]

		for v := 0; v < n; v++ {
			if !inTree[v] {
				dist := abs(points[u][0]-points[v][0]) + abs(points[u][1]-points[v][1])
				if dist < minDist[v] {
					minDist[v] = dist
				}
			}
		}
	}
	return total
}
