def max_of(values: list[int]) -> int:
    """Return the largest value in values."""
    best = values[0]
    for i in range(1, len(values)):
        if values[i] > best:
            best = values[i]
    return best
