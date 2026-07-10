package main

import "sort"

// IsNStraightHand returns whether hand can be rearranged into groups of
// groupSize consecutive cards.
func IsNStraightHand(hand []int, groupSize int) bool {
	if len(hand)%groupSize != 0 {
		return false
	}

	count := make(map[int]int)
	for _, c := range hand {
		count[c]++
	}

	keys := make([]int, 0, len(count))
	for k := range count {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		need := count[k]
		if need == 0 {
			continue
		}
		for i := 0; i < groupSize; i++ {
			c := k + i
			if count[c] < need {
				return false
			}
			count[c] -= need
		}
	}
	return true
}
