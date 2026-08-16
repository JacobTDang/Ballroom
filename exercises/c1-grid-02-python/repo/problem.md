# Spiral Statement

A statement renderer prints the cells of a grid in a clockwise spiral:
start at the top-left cell, go right across the top row, then down the
right edge, left across the bottom row, up the left edge, and keep
circling inward until every cell has been printed exactly once. The
grid may be any R x C shape, including a single row or a single column.

Write `spiral(grid: list[list[int]]) -> list[int]` returning the cell
values in that order.

## Examples

    spiral([[1, 2, 3],
            [4, 5, 6],
            [7, 8, 9]])   -> [1, 2, 3, 6, 9, 8, 7, 4, 5]

    spiral([[1, 2, 3],
            [4, 5, 6]])   -> [1, 2, 3, 6, 5, 4]

    spiral([[7],
            [8],
            [9]])         -> [7, 8, 9]

## Constraints

- 1 <= rows, cols <= 100
- -1_000_000 <= grid[r][c] <= 1_000_000
