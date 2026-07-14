package main

import "sort"

// EraseOverlapIntervals returns the minimum number of intervals that
// must be removed so the rest of intervals are non-overlapping.
func EraseOverlapIntervals(intervals [][]int) int {
	if len(intervals) == 0 {
		return 0
	}

	sorted := make([][]int, len(intervals))
	copy(sorted, intervals)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i][1] < sorted[j][1]
	})

	removals := 0
	lastEnd := sorted[0][1]
	for _, interval := range sorted[1:] {
		if interval[0] < lastEnd {
			removals++
		} else {
			lastEnd = interval[1]
		}
	}

	return removals
}
