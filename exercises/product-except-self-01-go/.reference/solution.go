package main

// ProductExceptSelf returns answer where answer[i] is the product of
// every element in nums except nums[i], without using division.
func ProductExceptSelf(nums []int) []int {
	n := len(nums)
	result := make([]int, n)

	prefix := 1
	for i := 0; i < n; i++ {
		result[i] = prefix
		prefix *= nums[i]
	}
	suffix := 1
	for i := n - 1; i >= 0; i-- {
		result[i] *= suffix
		suffix *= nums[i]
	}
	return result
}
