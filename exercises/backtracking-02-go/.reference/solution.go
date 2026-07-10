package main

// CombinationSum returns every unique combination of candidates
// (each usable unlimited times) that sums to target.
func CombinationSum(candidates []int, target int) [][]int {
	var res [][]int
	var cur []int
	var backtrack func(start, remain int)
	backtrack = func(start, remain int) {
		if remain == 0 {
			res = append(res, append([]int(nil), cur...))
			return
		}
		if remain < 0 {
			return
		}
		for i := start; i < len(candidates); i++ {
			cur = append(cur, candidates[i])
			backtrack(i, remain-candidates[i])
			cur = cur[:len(cur)-1]
		}
	}
	backtrack(0, target)
	return res
}
