def align(t: int, k: int) -> int:
    """Rounds t down to the start of its k-wide bucket. Currently wrong
    for negative t -- find and fix the bug."""
    return int(t / k) * k
