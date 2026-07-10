def subsets_with_dup(nums: list[int]) -> list[list[int]]:
    """Return every unique subset of nums, which may contain
    duplicate values."""
    nums = sorted(nums)
    res: list[list[int]] = []
    cur: list[int] = []

    def backtrack(start: int) -> None:
        res.append(cur[:])
        for i in range(start, len(nums)):
            if i > start and nums[i] == nums[i - 1]:
                continue
            cur.append(nums[i])
            backtrack(i + 1)
            cur.pop()

    backtrack(0)
    return res
