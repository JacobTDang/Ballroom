package main

// Grid is a 2-D grid of ints, rows x cols, 0-indexed.
type Grid struct {
	rows [][]int
}

// NewGrid returns a rows x cols grid with every cell set to 0.
func NewGrid(rows, cols int) *Grid {
	g := &Grid{rows: make([][]int, rows)}
	for i := range g.rows {
		g.rows[i] = make([]int, cols)
	}
	return g
}

func (g *Grid) Get(r, c int) int { return g.rows[r][c] }
func (g *Grid) Set(r, c, v int)  { g.rows[r][c] = v }

// Snapshot returns an independent copy of the grid's current cell
// values, for a caller that wants to edit the live grid and later
// compare against this saved state.
func (g *Grid) Snapshot() [][]int {
	dst := make([][]int, len(g.rows))
	for i, row := range g.rows {
		dst[i] = make([]int, len(row))
		copy(dst[i], row)
	}
	return dst
}
