package main

// TwoSum returns the indices of the two numbers in nums that add up to
// target.
func TwoSum(nums []int, target int) []int {
	seen := make(map[int]int, len(nums))
	for i, n := range nums {
		complement := target - n
		if j, ok := seen[complement]; ok {
			return []int{j, i}
		}
		seen[n] = i
	}
	return nil
}
