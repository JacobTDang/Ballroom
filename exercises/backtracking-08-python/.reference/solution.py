PHONE_LETTERS = {
    "2": "abc",
    "3": "def",
    "4": "ghi",
    "5": "jkl",
    "6": "mno",
    "7": "pqrs",
    "8": "tuv",
    "9": "wxyz",
}


def letter_combinations(digits: str) -> list[str]:
    """Return every letter combination that digits could represent
    on a phone keypad."""
    if not digits:
        return []
    res: list[str] = []
    cur: list[str] = []

    def backtrack(idx: int) -> None:
        if idx == len(digits):
            res.append("".join(cur))
            return
        for c in PHONE_LETTERS[digits[idx]]:
            cur.append(c)
            backtrack(idx + 1)
            cur.pop()

    backtrack(0)
    return res
