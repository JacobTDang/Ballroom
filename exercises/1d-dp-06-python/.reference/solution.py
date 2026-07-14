def _count_expansions(s: str, l: int, r: int) -> int:
    """Grow outward from the center between indices l and r, counting
    one palindrome for every successful expansion."""
    count = 0
    while l >= 0 and r < len(s) and s[l] == s[r]:
        count += 1
        l -= 1
        r += 1
    return count


def count_substrings(s: str) -> int:
    count = 0
    for i in range(len(s)):
        count += _count_expansions(s, i, i)
        count += _count_expansions(s, i, i + 1)
    return count
