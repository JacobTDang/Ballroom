def next_higher(spend: list[int]) -> list[int]:
    """Return, per day, the days until a strictly higher spend (0 if never)."""
    result = [0] * len(spend)
    pending = []  # indices whose next higher day has not appeared yet
    for j, value in enumerate(spend):
        while pending and spend[pending[-1]] < value:
            i = pending.pop()
            result[i] = j - i
        pending.append(j)
    return result
