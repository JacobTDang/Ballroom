def min_cost_climbing_stairs(cost: list[int]) -> int:
    n = len(cost)
    prev, curr = 0, 0
    for i in range(2, n + 1):
        next_cost = curr + cost[i - 1]
        alt = prev + cost[i - 2]
        if alt < next_cost:
            next_cost = alt
        prev, curr = curr, next_cost
    return curr
