package main

// RotateImage rotates the n x n matrix 90 degrees clockwise, in place.
func RotateImage(matrix [][]int) {
	n := len(matrix)

	// Transpose the matrix (reflect across the main diagonal).
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	// Reverse each row to complete the clockwise rotation.
	for i := 0; i < n; i++ {
		for l, r := 0, n-1; l < r; l, r = l+1, r-1 {
			matrix[i][l], matrix[i][r] = matrix[i][r], matrix[i][l]
		}
	}
}
