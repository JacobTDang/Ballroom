package main

// Insert inserts newInterval into the sorted, non-overlapping intervals
// list, merging overlaps as needed, and returns the resulting sorted,
// non-overlapping list.
func Insert(intervals [][]int, newInterval []int) [][]int {
	result := [][]int{}
	i, n := 0, len(intervals)
	start, end := newInterval[0], newInterval[1]

	for i < n && intervals[i][1] < start {
		result = append(result, intervals[i])
		i++
	}

	for i < n && intervals[i][0] <= end {
		if intervals[i][0] < start {
			start = intervals[i][0]
		}
		if intervals[i][1] > end {
			end = intervals[i][1]
		}
		i++
	}
	result = append(result, []int{start, end})

	for i < n {
		result = append(result, intervals[i])
		i++
	}

	return result
}
