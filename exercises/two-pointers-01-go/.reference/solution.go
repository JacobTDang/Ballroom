package main

// TwoSum returns the 1-indexed positions of the two numbers in numbers
// (sorted ascending) that add up to target.
func TwoSum(numbers []int, target int) []int {
	lo, hi := 0, len(numbers)-1
	for lo < hi {
		sum := numbers[lo] + numbers[hi]
		switch {
		case sum == target:
			return []int{lo + 1, hi + 1}
		case sum < target:
			lo++
		default:
			hi--
		}
	}
	return nil
}
