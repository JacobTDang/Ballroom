def largest_rectangle_area(heights: list[int]) -> int:
    """Return the area of the largest rectangle that fits under the
    histogram described by heights."""
    stack: list[tuple[int, int]] = []  # (start index, height)
    best = 0
    n = len(heights)
    for i in range(n + 1):
        h = heights[i] if i < n else 0
        start = i
        while stack and stack[-1][1] >= h:
            idx, height = stack.pop()
            best = max(best, height * (i - idx))
            start = idx
        stack.append((start, h))
    return best
