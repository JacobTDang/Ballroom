def merge(intervals: list[list[int]]) -> list[list[int]]:
    """Merge all overlapping intervals in intervals and return the
    resulting sorted, non-overlapping list."""
    if not intervals:
        return []

    sorted_intervals = sorted(intervals, key=lambda interval: interval[0])

    result = [list(sorted_intervals[0])]
    for start, end in sorted_intervals[1:]:
        last = result[-1]
        if start <= last[1]:
            last[1] = max(last[1], end)
        else:
            result.append([start, end])

    return result
