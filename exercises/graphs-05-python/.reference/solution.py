from collections import deque


def oranges_rotting(grid: list[list[int]]) -> int:
    """Return the minimum number of minutes until no cell in grid has
    a fresh orange, or -1 if some fresh orange can never rot."""
    rows, cols = len(grid), len(grid[0])
    queue = deque()
    fresh = 0
    for r in range(rows):
        for c in range(cols):
            if grid[r][c] == 2:
                queue.append((r, c))
            elif grid[r][c] == 1:
                fresh += 1
    if fresh == 0:
        return 0

    minutes = 0
    while queue and fresh > 0:
        for _ in range(len(queue)):
            r, c = queue.popleft()
            for dr, dc in ((1, 0), (-1, 0), (0, 1), (0, -1)):
                nr, nc = r + dr, c + dc
                if not (0 <= nr < rows and 0 <= nc < cols) or grid[nr][nc] != 1:
                    continue
                grid[nr][nc] = 2
                fresh -= 1
                queue.append((nr, nc))
        minutes += 1

    return -1 if fresh > 0 else minutes
