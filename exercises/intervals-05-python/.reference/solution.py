def min_meeting_rooms(intervals: list[list[int]]) -> int:
    """Return the minimum number of conference rooms required so that
    all meetings in intervals can happen without any two overlapping
    meetings sharing a room."""
    n = len(intervals)
    if n == 0:
        return 0

    starts = sorted(interval[0] for interval in intervals)
    ends = sorted(interval[1] for interval in intervals)

    rooms = 0
    max_rooms = 0
    i, j = 0, 0
    while i < n:
        if starts[i] < ends[j]:
            rooms += 1
            i += 1
            max_rooms = max(max_rooms, rooms)
        else:
            rooms -= 1
            j += 1

    return max_rooms
