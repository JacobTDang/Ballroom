from collections import Counter


def min_window(s: str, t: str) -> str:
    """Return the shortest substring of s containing every character
    of t (with duplicates), or "" if no such substring exists."""
    if not t or len(s) < len(t):
        return ""
    need = Counter(t)
    required = len(need)
    have = 0
    window: dict[str, int] = {}

    best_len = -1
    best_start = 0
    left = 0
    for right, c in enumerate(s):
        window[c] = window.get(c, 0) + 1
        if c in need and window[c] == need[c]:
            have += 1

        while have == required:
            if best_len == -1 or right - left + 1 < best_len:
                best_len = right - left + 1
                best_start = left
            left_char = s[left]
            window[left_char] -= 1
            if left_char in need and window[left_char] < need[left_char]:
                have -= 1
            left += 1

    if best_len == -1:
        return ""
    return s[best_start : best_start + best_len]
