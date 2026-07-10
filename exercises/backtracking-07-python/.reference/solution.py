def partition(s: str) -> list[list[str]]:
    """Return every way to split s into substrings that are all
    palindromes."""

    def is_palindrome(sub: str) -> bool:
        return sub == sub[::-1]

    res: list[list[str]] = []
    cur: list[str] = []

    def backtrack(start: int) -> None:
        if start == len(s):
            res.append(cur[:])
            return
        for end in range(start + 1, len(s) + 1):
            sub = s[start:end]
            if is_palindrome(sub):
                cur.append(sub)
                backtrack(end)
                cur.pop()

    backtrack(0)
    return res
