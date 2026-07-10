package main

// LargestRectangleArea returns the area of the largest rectangle that
// fits under the histogram described by heights.
func LargestRectangleArea(heights []int) int {
	type entry struct{ idx, height int }
	var stack []entry
	best := 0
	n := len(heights)
	for i := 0; i <= n; i++ {
		h := 0
		if i < n {
			h = heights[i]
		}
		start := i
		for len(stack) > 0 && stack[len(stack)-1].height >= h {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if area := top.height * (i - top.idx); area > best {
				best = area
			}
			start = top.idx
		}
		stack = append(stack, entry{start, h})
	}
	return best
}
