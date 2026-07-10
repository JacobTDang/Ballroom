package main

import "sort"

// ThreeSum returns every unique triplet of elements in nums that sums
// to zero.
func ThreeSum(nums []int) [][]int {
	sort.Ints(nums)
	var res [][]int
	n := len(nums)
	for i := 0; i < n-2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		lo, hi := i+1, n-1
		for lo < hi {
			sum := nums[i] + nums[lo] + nums[hi]
			switch {
			case sum < 0:
				lo++
			case sum > 0:
				hi--
			default:
				res = append(res, []int{nums[i], nums[lo], nums[hi]})
				lo++
				hi--
				for lo < hi && nums[lo] == nums[lo-1] {
					lo++
				}
				for lo < hi && nums[hi] == nums[hi+1] {
					hi--
				}
			}
		}
	}
	return res
}
