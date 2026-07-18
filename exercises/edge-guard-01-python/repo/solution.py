def max_adjacent_diff(values: list[int]) -> int:
    """Return the largest absolute difference between two adjacent
    elements in values. Raises ValueError if values has fewer than two
    elements. Currently crashes on some inputs — find and fix the
    bug."""
    if len(values) == 0:
        raise ValueError("max_adjacent_diff: need at least two values")
    best = abs(values[1] - values[0])
    for i in range(1, len(values) - 1):
        diff = abs(values[i + 1] - values[i])
        if diff > best:
            best = diff
    return best
