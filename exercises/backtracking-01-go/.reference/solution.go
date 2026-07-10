package main

// Subsets returns every subset of nums (the power set).
func Subsets(nums []int) [][]int {
	var res [][]int
	var cur []int
	var backtrack func(start int)
	backtrack = func(start int) {
		res = append(res, append([]int(nil), cur...))
		for i := start; i < len(nums); i++ {
			cur = append(cur, nums[i])
			backtrack(i + 1)
			cur = cur[:len(cur)-1]
		}
	}
	backtrack(0)
	return res
}
