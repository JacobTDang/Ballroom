from collections import deque

INF = 2147483647


def walls_and_gates(rooms: list[list[int]]) -> None:
    """Fill every empty room in rooms with its distance to the
    nearest gate, in place. Rooms that can't reach a gate stay INF."""
    rows, cols = len(rooms), len(rooms[0])
    queue = deque()
    for r in range(rows):
        for c in range(cols):
            if rooms[r][c] == 0:
                queue.append((r, c))

    while queue:
        r, c = queue.popleft()
        for dr, dc in ((1, 0), (-1, 0), (0, 1), (0, -1)):
            nr, nc = r + dr, c + dc
            if not (0 <= nr < rows and 0 <= nc < cols) or rooms[nr][nc] != INF:
                continue
            rooms[nr][nc] = rooms[r][c] + 1
            queue.append((nr, nc))
