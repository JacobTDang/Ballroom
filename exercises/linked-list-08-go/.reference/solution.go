package main

// FindDuplicate returns the one repeated value in nums, using Floyd's
// cycle detection over the implicit index -> nums[index] linked list.
func FindDuplicate(nums []int) int {
	slow, fast := nums[0], nums[0]
	for {
		slow = nums[slow]
		fast = nums[nums[fast]]
		if slow == fast {
			break
		}
	}
	slow2 := nums[0]
	for slow2 != slow {
		slow2 = nums[slow2]
		slow = nums[slow]
	}
	return slow
}
