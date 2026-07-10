package main

import "sort"

// SubsetsWithDup returns every unique subset of nums, which may
// contain duplicate values.
func SubsetsWithDup(nums []int) [][]int {
	sorted := append([]int(nil), nums...)
	sort.Ints(sorted)

	var res [][]int
	var cur []int
	var backtrack func(start int)
	backtrack = func(start int) {
		res = append(res, append([]int(nil), cur...))
		for i := start; i < len(sorted); i++ {
			if i > start && sorted[i] == sorted[i-1] {
				continue
			}
			cur = append(cur, sorted[i])
			backtrack(i + 1)
			cur = cur[:len(cur)-1]
		}
	}
	backtrack(0)
	return res
}
