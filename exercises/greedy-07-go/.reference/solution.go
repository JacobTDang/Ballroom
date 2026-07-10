package main

// PartitionLabels returns the sizes of the parts formed by splitting s
// so that each letter appears in at most one part, in order.
func PartitionLabels(s string) []int {
	last := make(map[byte]int)
	for i := 0; i < len(s); i++ {
		last[s[i]] = i
	}

	var result []int
	start, end := 0, 0
	for i := 0; i < len(s); i++ {
		if last[s[i]] > end {
			end = last[s[i]]
		}
		if i == end {
			result = append(result, end-start+1)
			start = i + 1
		}
	}
	return result
}
