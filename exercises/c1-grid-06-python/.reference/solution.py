def rotate(grid: list[list[int]]) -> list[list[int]]:
    """Return the R x C grid rotated 90 degrees clockwise as a new C x R grid."""
    rows, cols = len(grid), len(grid[0])
    return [[grid[rows - 1 - j][i] for j in range(rows)] for i in range(cols)]
