def can_attend_meetings(intervals: list[list[int]]) -> bool:
    """Return True if a person could attend every meeting in intervals
    without any two of them overlapping."""
    sorted_intervals = sorted(intervals, key=lambda interval: interval[0])

    for i in range(1, len(sorted_intervals)):
        if sorted_intervals[i][0] < sorted_intervals[i - 1][1]:
            return False
    return True
