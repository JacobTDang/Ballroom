package main

// TopKFrequent returns the k most frequent elements in nums, in any
// order.
func TopKFrequent(nums []int, k int) []int {
	counts := make(map[int]int)
	for _, n := range nums {
		counts[n]++
	}

	buckets := make([][]int, len(nums)+1)
	for n, c := range counts {
		buckets[c] = append(buckets[c], n)
	}

	result := make([]int, 0, k)
	for i := len(buckets) - 1; i >= 0 && len(result) < k; i-- {
		for _, n := range buckets[i] {
			result = append(result, n)
			if len(result) == k {
				break
			}
		}
	}
	return result
}
