def smooth(grid: list[list[int]]) -> list[list[int]]:
    """Return a new grid of floor-averaged neighborhoods read from the original."""
    rows, cols = len(grid), len(grid[0])
    out = []
    for r in range(rows):
        row = []
        for c in range(cols):
            total = count = 0
            for dr in (-1, 0, 1):
                for dc in (-1, 0, 1):
                    nr, nc = r + dr, c + dc
                    if 0 <= nr < rows and 0 <= nc < cols:
                        total += grid[nr][nc]
                        count += 1
            row.append(total // count)
        out.append(row)
    return out
