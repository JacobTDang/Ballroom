# Vault Regions

Floor sensors map the vault as a grid: `1` is an occupied cell
(safe-deposit shelving), `0` is clear floor. Two occupied cells belong
to the same region when they touch up, down, left, or right — diagonal
contact does not connect them. Count how many distinct regions the
floor holds.

Write `count_regions(grid: list[list[int]]) -> int`.

## Examples

    count_regions([[1, 1, 0],
                   [0, 1, 0],
                   [0, 0, 1]])    -> 2

    count_regions([[1, 0, 1],
                   [0, 1, 0],
                   [1, 0, 1]])    -> 5

    count_regions([[1, 1, 1, 1]]) -> 1

## Constraints

- 1 <= rows, cols <= 100
- grid[r][c] is 0 or 1
- a single region may wind through every occupied cell of a 100 x 100 grid
