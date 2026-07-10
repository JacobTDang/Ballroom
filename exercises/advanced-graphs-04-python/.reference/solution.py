import heapq


def swim_in_water(grid: list[list[int]]) -> int:
    n = len(grid)
    if n == 0:
        return 0

    visited = [[False] * n for _ in range(n)]
    heap = [(grid[0][0], 0, 0)]
    visited[0][0] = True

    dirs = [(1, 0), (-1, 0), (0, 1), (0, -1)]

    while heap:
        elevation, row, col = heapq.heappop(heap)
        if row == n - 1 and col == n - 1:
            return elevation
        for dr, dc in dirs:
            nr, nc = row + dr, col + dc
            if 0 <= nr < n and 0 <= nc < n and not visited[nr][nc]:
                visited[nr][nc] = True
                heapq.heappush(heap, (max(elevation, grid[nr][nc]), nr, nc))

    return -1
