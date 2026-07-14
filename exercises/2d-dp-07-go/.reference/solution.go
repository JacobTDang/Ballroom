package main

// LongestIncreasingPath returns the length of the longest strictly
// increasing path in matrix, moving 4-directionally.
func LongestIncreasingPath(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	rows, cols := len(matrix), len(matrix[0])
	memo := make([][]int, rows)
	for i := range memo {
		memo[i] = make([]int, cols)
	}

	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	var dfs func(r, c int) int
	dfs = func(r, c int) int {
		if memo[r][c] != 0 {
			return memo[r][c]
		}
		best := 1
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols && matrix[nr][nc] > matrix[r][c] {
				length := 1 + dfs(nr, nc)
				if length > best {
					best = length
				}
			}
		}
		memo[r][c] = best
		return best
	}

	result := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if length := dfs(r, c); length > result {
				result = length
			}
		}
	}
	return result
}
