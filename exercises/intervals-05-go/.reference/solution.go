package main

import "sort"

// MinMeetingRooms returns the minimum number of conference rooms
// required so that all meetings in intervals can happen without any
// two overlapping meetings sharing a room.
func MinMeetingRooms(intervals [][]int) int {
	n := len(intervals)
	if n == 0 {
		return 0
	}

	starts := make([]int, n)
	ends := make([]int, n)
	for i, interval := range intervals {
		starts[i] = interval[0]
		ends[i] = interval[1]
	}
	sort.Ints(starts)
	sort.Ints(ends)

	rooms, maxRooms := 0, 0
	i, j := 0, 0
	for i < n {
		if starts[i] < ends[j] {
			rooms++
			i++
			if rooms > maxRooms {
				maxRooms = rooms
			}
		} else {
			rooms--
			j++
		}
	}

	return maxRooms
}
