def erase_overlap_intervals(intervals: list[list[int]]) -> int:
    """Return the minimum number of intervals that must be removed so
    the rest of intervals are non-overlapping."""
    if not intervals:
        return 0

    sorted_intervals = sorted(intervals, key=lambda interval: interval[1])

    removals = 0
    last_end = sorted_intervals[0][1]
    for start, end in sorted_intervals[1:]:
        if start < last_end:
            removals += 1
        else:
            last_end = end

    return removals
