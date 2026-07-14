def insert(intervals: list[list[int]], new_interval: list[int]) -> list[list[int]]:
    """Insert new_interval into the sorted, non-overlapping intervals
    list, merging overlaps as needed, and return the resulting sorted,
    non-overlapping list."""
    result = []
    i, n = 0, len(intervals)
    start, end = new_interval[0], new_interval[1]

    while i < n and intervals[i][1] < start:
        result.append(intervals[i])
        i += 1

    while i < n and intervals[i][0] <= end:
        start = min(start, intervals[i][0])
        end = max(end, intervals[i][1])
        i += 1
    result.append([start, end])

    while i < n:
        result.append(intervals[i])
        i += 1

    return result
