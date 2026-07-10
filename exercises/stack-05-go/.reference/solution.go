package main

// DailyTemperatures returns, for each day, how many days until a
// warmer temperature, or 0 if there isn't one.
func DailyTemperatures(temperatures []int) []int {
	res := make([]int, len(temperatures))
	var stack []int // indices, decreasing temperature
	for i, temp := range temperatures {
		for len(stack) > 0 && temperatures[stack[len(stack)-1]] < temp {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			res[top] = i - top
		}
		stack = append(stack, i)
	}
	return res
}
