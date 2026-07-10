package main

import "sort"

// KClosest returns the k points from points closest to the origin,
// in any order.
func KClosest(points [][]int, k int) [][]int {
	sorted := make([][]int, len(points))
	copy(sorted, points)
	sort.Slice(sorted, func(i, j int) bool {
		di := sorted[i][0]*sorted[i][0] + sorted[i][1]*sorted[i][1]
		dj := sorted[j][0]*sorted[j][0] + sorted[j][1]*sorted[j][1]
		return di < dj
	})
	return sorted[:k]
}
