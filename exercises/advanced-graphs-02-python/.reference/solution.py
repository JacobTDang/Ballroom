def min_cost_connect_points(points: list[list[int]]) -> int:
    n = len(points)
    if n <= 1:
        return 0

    in_tree = [False] * n
    min_dist = [float("inf")] * n
    min_dist[0] = 0

    total = 0
    for _ in range(n):
        u = -1
        for v in range(n):
            if not in_tree[v] and (u == -1 or min_dist[v] < min_dist[u]):
                u = v
        in_tree[u] = True
        total += min_dist[u]

        for v in range(n):
            if not in_tree[v]:
                dist = abs(points[u][0] - points[v][0]) + abs(points[u][1] - points[v][1])
                if dist < min_dist[v]:
                    min_dist[v] = dist

    return total
