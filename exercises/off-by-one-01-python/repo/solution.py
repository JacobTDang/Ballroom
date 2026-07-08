def max_of(values: list[int]) -> int:
    """Return the largest value in values. Currently crashes — find and
    fix the bug."""
    best = values[0]
    for i in range(len(values) + 1):
        if values[i] > best:
            best = values[i]
    return best
