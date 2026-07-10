package main

// IsValidSudoku reports whether the filled cells of a 9x9 Sudoku board
// satisfy Sudoku's placement rules (no digit repeated within a row,
// column, or 3x3 box). Empty cells are '.'.
func IsValidSudoku(board []string) bool {
	var rows, cols, boxes [9][9]bool

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			ch := board[r][c]
			if ch == '.' {
				continue
			}
			d := ch - '1'
			b := (r/3)*3 + c/3
			if rows[r][d] || cols[c][d] || boxes[b][d] {
				return false
			}
			rows[r][d] = true
			cols[c][d] = true
			boxes[b][d] = true
		}
	}
	return true
}
