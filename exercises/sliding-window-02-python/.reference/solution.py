def length_of_longest_substring(s: str) -> int:
    """Return the length of the longest substring of s with no
    repeating characters."""
    last_seen: dict[str, int] = {}
    left = 0
    best = 0
    for right, c in enumerate(s):
        if c in last_seen and last_seen[c] >= left:
            left = last_seen[c] + 1
        last_seen[c] = right
        best = max(best, right - left + 1)
    return best
