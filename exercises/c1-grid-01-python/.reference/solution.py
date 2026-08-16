def walk(grid: list[list[int]], moves: str) -> tuple[int, int]:
    """Return the robot's final (row, col) after applying the move string."""
    deltas = {"U": (-1, 0), "D": (1, 0), "L": (0, -1), "R": (0, 1)}
    rows, cols = len(grid), len(grid[0])
    r = c = 0
    for move in moves:
        dr, dc = deltas[move]
        nr, nc = r + dr, c + dc
        if 0 <= nr < rows and 0 <= nc < cols and grid[nr][nc] == 0:
            r, c = nr, nc
    return (r, c)
