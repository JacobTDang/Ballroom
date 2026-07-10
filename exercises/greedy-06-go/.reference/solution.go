package main

// MergeTriplets returns whether target can be formed by taking the
// elementwise max of some subset of triplets.
func MergeTriplets(triplets [][]int, target []int) bool {
	matched := make(map[int]bool)
	for _, t := range triplets {
		if t[0] > target[0] || t[1] > target[1] || t[2] > target[2] {
			continue
		}
		for i := 0; i < 3; i++ {
			if t[i] == target[i] {
				matched[i] = true
			}
		}
	}
	return len(matched) == 3
}
