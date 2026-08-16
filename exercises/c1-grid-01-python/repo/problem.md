# Branch Floor Walk

A delivery robot rolls around a branch floor plan given as a grid of
ints: `0` is open floor, `1` is a wall. Rows are numbered top to bottom
and columns left to right, and the robot starts on cell `(0, 0)`, which
is always open. It then follows a move string one character at a time:
`U` is row - 1, `D` is row + 1, `L` is col - 1, `R` is col + 1. A move
that would land on a wall or leave the grid is skipped — the robot stays
where it is and the next move is processed as usual.

Write `walk(grid: list[list[int]], moves: str) -> tuple[int, int]`
returning the robot's final `(row, col)`.

## Examples

    walk([[0, 0, 1],
          [0, 1, 0],
          [0, 0, 0]], "DDRR")   -> (2, 2)

    walk([[0, 1],
          [0, 0]], "RDR")       -> (1, 1)   # the first R hits a wall and is skipped

    walk([[0, 0]], "ULDRR")     -> (0, 1)   # U, L, D and the last R would leave the grid

## Constraints

- 1 <= rows, cols <= 100
- grid[r][c] is 0 or 1, and grid[0][0] == 0
- 0 <= len(moves) <= 10_000, characters are only U, D, L, R
