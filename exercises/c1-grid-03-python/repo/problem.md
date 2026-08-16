# Fraud Heatmap Smoothing

The fraud team renders alert counts per city block as a heatmap. Raw
counts are noisy, so each cell is displayed as a neighborhood average:
sum the cell and its in-bounds 8-neighbors, then floor-divide (`//`)
by how many cells went into the sum (4 in a corner of a large grid,
6 on an edge, 9 in the middle). Every average is computed from the
ORIGINAL grid — values smoothed earlier in the pass must not feed into
later cells. Return the smoothed grid as a new list of lists and leave
the input grid unchanged.

Write `smooth(grid: list[list[int]]) -> list[list[int]]`.

## Examples

    smooth([[1, 1],
            [1, 3]])      -> [[1, 1],
                              [1, 1]]

    smooth([[0, 9, 0],
            [9, 0, 9],
            [0, 9, 0]])   -> [[4, 4, 4],
                              [4, 4, 4],
                              [4, 4, 4]]

    smooth([[3, 6, 9]])   -> [[4, 6, 7]]

## Constraints

- 1 <= rows, cols <= 100
- 0 <= grid[r][c] <= 1000
