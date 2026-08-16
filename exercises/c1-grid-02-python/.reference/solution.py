def spiral(grid: list[list[int]]) -> list[int]:
    """Return the grid's values in clockwise spiral order from the top-left."""
    top, bottom = 0, len(grid) - 1
    left, right = 0, len(grid[0]) - 1
    out = []
    while top <= bottom and left <= right:
        for c in range(left, right + 1):
            out.append(grid[top][c])
        top += 1
        for r in range(top, bottom + 1):
            out.append(grid[r][right])
        right -= 1
        if top <= bottom:
            for c in range(right, left - 1, -1):
                out.append(grid[bottom][c])
            bottom -= 1
        if left <= right:
            for r in range(bottom, top - 1, -1):
                out.append(grid[r][left])
            left += 1
    return out
