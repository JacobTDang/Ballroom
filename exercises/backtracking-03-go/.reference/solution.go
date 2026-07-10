package main

// Permute returns every permutation of nums.
func Permute(nums []int) [][]int {
	var res [][]int
	var cur []int
	used := make([]bool, len(nums))
	var backtrack func()
	backtrack = func() {
		if len(cur) == len(nums) {
			res = append(res, append([]int(nil), cur...))
			return
		}
		for i, n := range nums {
			if used[i] {
				continue
			}
			used[i] = true
			cur = append(cur, n)
			backtrack()
			cur = cur[:len(cur)-1]
			used[i] = false
		}
	}
	backtrack()
	return res
}
