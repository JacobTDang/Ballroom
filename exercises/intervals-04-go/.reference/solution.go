package main

import "sort"

// CanAttendMeetings returns true if a person could attend every
// meeting in intervals without any two of them overlapping.
func CanAttendMeetings(intervals [][]int) bool {
	sorted := make([][]int, len(intervals))
	copy(sorted, intervals)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i][0] < sorted[j][0]
	})

	for i := 1; i < len(sorted); i++ {
		if sorted[i][0] < sorted[i-1][1] {
			return false
		}
	}
	return true
}
