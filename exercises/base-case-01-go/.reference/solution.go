package main

// CountPaths counts the number of paths from the top-left to the
// bottom-right corner of grid, moving only right or down, that avoid
// cells marked 1 (blocked).
func CountPaths(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])

	var helper func(r, c int) int
	helper = func(r, c int) int {
		if r >= rows || c >= cols || grid[r][c] == 1 {
			return 0
		}
		if r == rows-1 && c == cols-1 {
			return 1
		}
		return helper(r+1, c) + helper(r, c+1)
	}

	return helper(0, 0)
}
