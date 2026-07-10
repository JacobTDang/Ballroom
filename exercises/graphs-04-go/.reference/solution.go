package main

// Inf represents an empty room whose distance to the nearest gate is
// not yet known.
const Inf = 2147483647

// WallsAndGates fills every empty room in rooms with its distance to
// the nearest gate, in place. Rooms that can't reach a gate stay Inf.
func WallsAndGates(rooms [][]int) {
	rows, cols := len(rooms), len(rooms[0])
	var queue [][2]int
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if rooms[r][c] == 0 {
				queue = append(queue, [2]int{r, c})
			}
		}
	}

	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, d := range dirs {
			nr, nc := cur[0]+d[0], cur[1]+d[1]
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols || rooms[nr][nc] != Inf {
				continue
			}
			rooms[nr][nc] = rooms[cur[0]][cur[1]] + 1
			queue = append(queue, [2]int{nr, nc})
		}
	}
}
