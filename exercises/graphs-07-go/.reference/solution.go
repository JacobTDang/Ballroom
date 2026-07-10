package main

// Solve captures all regions of 'O' surrounded by 'X' by flipping
// them to 'X' in place. Regions connected to the border are left
// untouched.
func Solve(board [][]byte) {
	rows, cols := len(board), len(board[0])
	if rows == 0 || cols == 0 {
		return
	}

	var dfs func(r, c int)
	dfs = func(r, c int) {
		if r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] != 'O' {
			return
		}
		board[r][c] = '#'
		dfs(r+1, c)
		dfs(r-1, c)
		dfs(r, c+1)
		dfs(r, c-1)
	}

	for c := 0; c < cols; c++ {
		dfs(0, c)
		dfs(rows-1, c)
	}
	for r := 0; r < rows; r++ {
		dfs(r, 0)
		dfs(r, cols-1)
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			switch board[r][c] {
			case 'O':
				board[r][c] = 'X'
			case '#':
				board[r][c] = 'O'
			}
		}
	}
}
