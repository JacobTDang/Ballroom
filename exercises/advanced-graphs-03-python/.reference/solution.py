def network_delay_time(times: list[list[int]], n: int, k: int) -> int:
    adj: list[list[tuple[int, int]]] = [[] for _ in range(n + 1)]
    for u, v, w in times:
        adj[u].append((v, w))

    dist = [float("inf")] * (n + 1)
    dist[k] = 0

    visited = [False] * (n + 1)
    for _ in range(n):
        u = -1
        for v in range(1, n + 1):
            if not visited[v] and (u == -1 or dist[v] < dist[u]):
                u = v
        if u == -1 or dist[u] == float("inf"):
            break
        visited[u] = True
        for v, w in adj[u]:
            if dist[u] + w < dist[v]:
                dist[v] = dist[u] + w

    max_dist = 0
    for v in range(1, n + 1):
        if dist[v] == float("inf"):
            return -1
        max_dist = max(max_dist, dist[v])
    return max_dist
