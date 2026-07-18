def max_below_limit(values: list[int], limit: int) -> int:
    """Return the largest value in values that is <= limit, or -1 if
    no value qualifies. Currently always returns -1 — find and fix the
    bug."""
    result = float("-inf")

    def record(x: int) -> None:
        result = x

    for x in values:
        if x <= limit and x > result:
            record(x)
    return -1 if result == float("-inf") else result
