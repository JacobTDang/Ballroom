package main

// SetZeroes sets the entire row and column of any zero element to zero,
// in place.
func SetZeroes(matrix [][]int) {
	rows := len(matrix)
	if rows == 0 {
		return
	}
	cols := len(matrix[0])

	zeroRow := make([]bool, rows)
	zeroCol := make([]bool, cols)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if matrix[r][c] == 0 {
				zeroRow[r] = true
				zeroCol[c] = true
			}
		}
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if zeroRow[r] || zeroCol[c] {
				matrix[r][c] = 0
			}
		}
	}
}
