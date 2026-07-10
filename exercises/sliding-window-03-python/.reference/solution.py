from collections import defaultdict


def character_replacement(s: str, k: int) -> int:
    """Return the length of the longest substring of s that can be
    made to contain only one repeating letter after at most k
    character replacements."""
    count: dict[str, int] = defaultdict(int)
    left = 0
    max_freq = 0
    best = 0
    for right, c in enumerate(s):
        count[c] += 1
        max_freq = max(max_freq, count[c])
        while right - left + 1 - max_freq > k:
            count[s[left]] -= 1
            left += 1
        best = max(best, right - left + 1)
    return best
