import math


def min_eating_speed(piles: list[int], h: int) -> int:
    """Return the minimum bananas-per-hour eating speed that lets
    Koko finish every pile within h hours."""
    lo, hi = 1, max(piles)
    while lo < hi:
        mid = lo + (hi - lo) // 2
        hours = sum(math.ceil(p / mid) for p in piles)
        if hours <= h:
            hi = mid
        else:
            lo = mid + 1
    return lo
