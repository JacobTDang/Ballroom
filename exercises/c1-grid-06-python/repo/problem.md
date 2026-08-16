# Rotate the Card Art

Custom card designs are stored as a pixel grid with R rows and C
columns. The print shop needs the artwork turned 90 degrees clockwise:
the first row of the input becomes the last column of the output, so an
R x C grid comes back as a C x R grid. Return the rotated grid as a new
list of lists; the input grid must stay unchanged.

Write `rotate(grid: list[list[int]]) -> list[list[int]]`.

## Examples

    rotate([[1, 2],
            [3, 4]])      -> [[3, 1],
                              [4, 2]]

    rotate([[1, 2, 3],
            [4, 5, 6]])   -> [[4, 1],
                              [5, 2],
                              [6, 3]]

    rotate([[1, 2, 3]])   -> [[1],
                              [2],
                              [3]]

## Constraints

- 1 <= rows, cols <= 100
- 0 <= grid[r][c] <= 255
