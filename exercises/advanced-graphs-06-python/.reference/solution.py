def find_cheapest_price(n: int, flights: list[list[int]], src: int, dst: int, k: int) -> int:
    inf = float("inf")
    dist = [inf] * n
    dist[src] = 0

    # Bellman-Ford limited to exactly k+1 relaxation rounds. Each round
    # must relax edges using a SNAPSHOT of the previous round's
    # distances, not the array being updated in place during that same
    # round, or a single round could silently chain multiple edges
    # together and violate the stop limit.
    for _ in range(k + 1):
        prev = dist[:]

        for u, v, price in flights:
            if prev[u] == inf:
                continue
            if prev[u] + price < dist[v]:
                dist[v] = prev[u] + price

    return -1 if dist[dst] == inf else dist[dst]
