def step(grid: list[list[str]], n: int) -> list[list[str]]:
    """Return a new grid after n simultaneous right-shift ticks."""
    rows, cols = len(grid), len(grid[0])
    current = [row[:] for row in grid]
    for _ in range(n):
        nxt = [["."] * cols for _ in range(rows)]
        for r in range(rows):
            for c in range(cols):
                if current[r][c] != ">":
                    continue
                if c == cols - 1:
                    continue  # exits the grid
                if current[r][c + 1] == ".":
                    nxt[r][c + 1] = ">"
                else:
                    nxt[r][c] = ">"
        current = nxt
    return current
