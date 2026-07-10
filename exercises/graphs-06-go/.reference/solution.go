package main

// PacificAtlantic returns every cell [r, c] from which water can
// flow to both the Pacific (top/left edges) and Atlantic
// (bottom/right edges) oceans.
func PacificAtlantic(heights [][]int) [][]int {
	rows, cols := len(heights), len(heights[0])
	pacific := make([][]bool, rows)
	atlantic := make([][]bool, rows)
	for i := range pacific {
		pacific[i] = make([]bool, cols)
		atlantic[i] = make([]bool, cols)
	}

	var dfs func(r, c int, visited [][]bool, prevHeight int)
	dfs = func(r, c int, visited [][]bool, prevHeight int) {
		if r < 0 || r >= rows || c < 0 || c >= cols || visited[r][c] || heights[r][c] < prevHeight {
			return
		}
		visited[r][c] = true
		dfs(r+1, c, visited, heights[r][c])
		dfs(r-1, c, visited, heights[r][c])
		dfs(r, c+1, visited, heights[r][c])
		dfs(r, c-1, visited, heights[r][c])
	}

	for c := 0; c < cols; c++ {
		dfs(0, c, pacific, heights[0][c])
		dfs(rows-1, c, atlantic, heights[rows-1][c])
	}
	for r := 0; r < rows; r++ {
		dfs(r, 0, pacific, heights[r][0])
		dfs(r, cols-1, atlantic, heights[r][cols-1])
	}

	var res [][]int
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if pacific[r][c] && atlantic[r][c] {
				res = append(res, []int{r, c})
			}
		}
	}
	return res
}
