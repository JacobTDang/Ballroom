def count_paths(grid: list[list[int]]) -> int:
    """Count the number of paths from the top-left to the bottom-right
    of grid, moving only right or down, that avoid cells marked 1
    (blocked)."""
    rows, cols = len(grid), len(grid[0])

    def helper(r: int, c: int) -> int:
        if r >= rows or c >= cols or grid[r][c] == 1:
            return 0
        if r == rows - 1 and c == cols - 1:
            return 1
        return helper(r + 1, c) + helper(r, c + 1)

    return helper(0, 0)
