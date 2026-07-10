def max_area(height: list[int]) -> int:
    """Return the largest amount of water a container formed by two
    lines in height (with the x-axis) can hold."""
    lo, hi = 0, len(height) - 1
    best = 0
    while lo < hi:
        h = min(height[lo], height[hi])
        best = max(best, h * (hi - lo))
        if height[lo] < height[hi]:
            lo += 1
        else:
            hi -= 1
    return best
