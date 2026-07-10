package main

// SolveNQueens returns every distinct board configuration that
// places n queens on an n x n board with no two attacking each other.
func SolveNQueens(n int) [][]string {
	var res [][]string
	cols := make([]bool, n)
	diag1 := make([]bool, 2*n) // r+c
	diag2 := make([]bool, 2*n) // r-c+n

	board := make([][]byte, n)
	for i := range board {
		board[i] = make([]byte, n)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}

	var backtrack func(r int)
	backtrack = func(r int) {
		if r == n {
			rows := make([]string, n)
			for i, row := range board {
				rows[i] = string(row)
			}
			res = append(res, rows)
			return
		}
		for c := 0; c < n; c++ {
			if cols[c] || diag1[r+c] || diag2[r-c+n] {
				continue
			}
			cols[c], diag1[r+c], diag2[r-c+n] = true, true, true
			board[r][c] = 'Q'
			backtrack(r + 1)
			board[r][c] = '.'
			cols[c], diag1[r+c], diag2[r-c+n] = false, false, false
		}
	}
	backtrack(0)
	return res
}
