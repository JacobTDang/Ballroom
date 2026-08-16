from __future__ import annotations

from bisect import bisect_right


def pair_within_budget(prices: list[int], budget: int) -> tuple[int, int] | None:
    """Return the index pair with the largest total within budget, or None."""
    values = sorted(prices)
    lo, hi = 0, len(values) - 1
    best = None
    while lo < hi:
        total = values[lo] + values[hi]
        if total <= budget:
            if best is None or total > best:
                best = total
            lo += 1
        else:
            hi -= 1
    if best is None:
        return None

    positions = {}
    for index, price in enumerate(prices):
        positions.setdefault(price, []).append(index)
    for i, price in enumerate(prices):
        matches = positions.get(best - price)
        if not matches:
            continue
        at = bisect_right(matches, i)
        if at < len(matches):
            return i, matches[at]
    raise AssertionError("unreachable: best total must come from some pair")
