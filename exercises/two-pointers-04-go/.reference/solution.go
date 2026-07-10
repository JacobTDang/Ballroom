package main

// MaxArea returns the largest amount of water a container formed by
// two lines in height (with the x-axis) can hold.
func MaxArea(height []int) int {
	lo, hi := 0, len(height)-1
	best := 0
	for lo < hi {
		h := height[lo]
		if height[hi] < h {
			h = height[hi]
		}
		area := h * (hi - lo)
		if area > best {
			best = area
		}
		if height[lo] < height[hi] {
			lo++
		} else {
			hi--
		}
	}
	return best
}
