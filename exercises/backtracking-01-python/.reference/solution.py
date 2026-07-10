def subsets(nums: list[int]) -> list[list[int]]:
    """Return every subset of nums (the power set)."""
    res: list[list[int]] = []
    cur: list[int] = []

    def backtrack(start: int) -> None:
        res.append(cur[:])
        for i in range(start, len(nums)):
            cur.append(nums[i])
            backtrack(i + 1)
            cur.pop()

    backtrack(0)
    return res
