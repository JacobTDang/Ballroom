def min_removals(s: str) -> int:
    """Return the minimum removals needed to balance the bracket string."""
    unmatched_open = unmatched_close = 0
    for ch in s:
        if ch == "(":
            unmatched_open += 1
        elif unmatched_open:
            unmatched_open -= 1
        else:
            unmatched_close += 1
    return unmatched_open + unmatched_close
