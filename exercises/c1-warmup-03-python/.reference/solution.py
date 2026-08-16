def longest_streak(amounts: list[int]) -> int:
    """Return the length of the longest strictly increasing consecutive run."""
    if not amounts:
        return 0
    best = current = 1
    for prev, nxt in zip(amounts, amounts[1:]):
        current = current + 1 if nxt > prev else 1
        best = max(best, current)
    return best
