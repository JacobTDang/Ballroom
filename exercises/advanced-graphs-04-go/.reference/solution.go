package main

import "container/heap"

type waterCell struct {
	elevation int
	row, col  int
}

type cellHeap []waterCell

func (h cellHeap) Len() int            { return len(h) }
func (h cellHeap) Less(i, j int) bool  { return h[i].elevation < h[j].elevation }
func (h cellHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *cellHeap) Push(x interface{}) { *h = append(*h, x.(waterCell)) }
func (h *cellHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// SwimInWater returns the minimum time t such that you can swim from the
// top-left to the bottom-right of grid, where at time t you may move
// between adjacent cells whose elevation is <= t.
func SwimInWater(grid [][]int) int {
	n := len(grid)
	if n == 0 {
		return 0
	}

	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, n)
	}

	h := &cellHeap{{elevation: grid[0][0], row: 0, col: 0}}
	heap.Init(h)
	visited[0][0] = true

	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for h.Len() > 0 {
		cur := heap.Pop(h).(waterCell)
		if cur.row == n-1 && cur.col == n-1 {
			return cur.elevation
		}
		for _, d := range dirs {
			nr, nc := cur.row+d[0], cur.col+d[1]
			if nr < 0 || nr >= n || nc < 0 || nc >= n || visited[nr][nc] {
				continue
			}
			visited[nr][nc] = true
			maxElevation := cur.elevation
			if grid[nr][nc] > maxElevation {
				maxElevation = grid[nr][nc]
			}
			heap.Push(h, waterCell{elevation: maxElevation, row: nr, col: nc})
		}
	}
	return -1
}
