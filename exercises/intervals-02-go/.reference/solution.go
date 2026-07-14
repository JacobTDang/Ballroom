package main

import "sort"

// Merge merges all overlapping intervals in intervals and returns the
// resulting sorted, non-overlapping list.
func Merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}

	sorted := make([][]int, len(intervals))
	copy(sorted, intervals)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i][0] < sorted[j][0]
	})

	result := [][]int{{sorted[0][0], sorted[0][1]}}
	for _, interval := range sorted[1:] {
		last := result[len(result)-1]
		if interval[0] <= last[1] {
			if interval[1] > last[1] {
				last[1] = interval[1]
			}
		} else {
			result = append(result, []int{interval[0], interval[1]})
		}
	}

	return result
}
