package main

// Trap returns the total units of water trapped between the bars
// described by height.
func Trap(height []int) int {
	if len(height) == 0 {
		return 0
	}
	lo, hi := 0, len(height)-1
	leftMax, rightMax := height[lo], height[hi]
	total := 0
	for lo < hi {
		if leftMax < rightMax {
			lo++
			if height[lo] > leftMax {
				leftMax = height[lo]
			} else {
				total += leftMax - height[lo]
			}
		} else {
			hi--
			if height[hi] > rightMax {
				rightMax = height[hi]
			} else {
				total += rightMax - height[hi]
			}
		}
	}
	return total
}
