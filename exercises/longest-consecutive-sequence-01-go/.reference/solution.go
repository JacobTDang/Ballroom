package main

// LongestConsecutive returns the length of the longest run of
// consecutive integers present in nums (order doesn't matter, duplicates
// don't count extra).
func LongestConsecutive(nums []int) int {
	set := make(map[int]bool, len(nums))
	for _, n := range nums {
		set[n] = true
	}

	longest := 0
	for n := range set {
		if set[n-1] {
			continue // n isn't the start of a sequence
		}
		length := 1
		for set[n+length] {
			length++
		}
		if length > longest {
			longest = length
		}
	}
	return longest
}
