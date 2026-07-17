package main

// ProcessAll applies fn to every job and returns the results in input
// order, using at most `workers` concurrent workers.
//
// TODO: this version is sequential -- one job at a time, no workers at
// all. Parallelize it without breaking the ordering.
func ProcessAll(jobs []int, workers int, fn func(int) int) []int {
	results := make([]int, len(jobs))
	for i, j := range jobs {
		results[i] = fn(j)
	}
	return results
}
