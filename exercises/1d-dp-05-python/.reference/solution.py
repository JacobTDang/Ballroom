def _expand_around_center(s: str, l: int, r: int) -> int:
    """Grow outward from the center between indices l and r (l == r for
    an odd-length center, r == l + 1 for an even-length center) and
    return the length of the palindrome found."""
    while l >= 0 and r < len(s) and s[l] == s[r]:
        l -= 1
        r += 1
    return r - l - 1


def longest_palindrome(s: str) -> str:
    if not s:
        return ""
    start, end = 0, 0
    for i in range(len(s)):
        len1 = _expand_around_center(s, i, i)
        len2 = _expand_around_center(s, i, i + 1)
        max_len = max(len1, len2)
        if max_len > end - start + 1:
            start = i - (max_len - 1) // 2
            end = i + max_len // 2
    return s[start:end + 1]
