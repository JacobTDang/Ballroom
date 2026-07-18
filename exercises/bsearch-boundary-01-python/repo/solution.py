def first_at_least(values: list[int], target: int) -> int:
    """Return the index of the first element in values that is >=
    target, or len(values) if every element is smaller. Currently
    returns an index one too small for some targets — find and fix
    the bug."""
    lo, hi = 0, len(values) - 1
    while lo < hi:
        mid = (lo + hi) // 2
        if values[mid] < target:
            lo = mid + 1
        else:
            hi = mid
    return lo
