package main

import "sort"

// CombinationSum2 returns every unique combination of candidates
// (each usable at most once) that sums to target.
func CombinationSum2(candidates []int, target int) [][]int {
	sorted := append([]int(nil), candidates...)
	sort.Ints(sorted)

	var res [][]int
	var cur []int
	var backtrack func(start, remain int)
	backtrack = func(start, remain int) {
		if remain == 0 {
			res = append(res, append([]int(nil), cur...))
			return
		}
		for i := start; i < len(sorted); i++ {
			if i > start && sorted[i] == sorted[i-1] {
				continue
			}
			if sorted[i] > remain {
				break
			}
			cur = append(cur, sorted[i])
			backtrack(i+1, remain-sorted[i])
			cur = cur[:len(cur)-1]
		}
	}
	backtrack(0, target)
	return res
}
