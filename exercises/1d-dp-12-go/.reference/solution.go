package main

// CanPartition returns whether nums can be partitioned into two
// subsets with equal sum.
func CanPartition(nums []int) bool {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	if sum%2 != 0 {
		return false
	}

	target := sum / 2
	reachable := make([]bool, target+1)
	reachable[0] = true

	for _, n := range nums {
		for i := target; i >= n; i-- {
			if reachable[i-n] {
				reachable[i] = true
			}
		}
	}

	return reachable[target]
}
