package main

// SearchMatrix reports whether target is present in matrix, treating
// it as one flattened sorted sequence.
func SearchMatrix(matrix [][]int, target int) bool {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return false
	}
	rows, cols := len(matrix), len(matrix[0])
	lo, hi := 0, rows*cols-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		val := matrix[mid/cols][mid%cols]
		switch {
		case val == target:
			return true
		case val < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return false
}
