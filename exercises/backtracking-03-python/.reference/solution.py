def permute(nums: list[int]) -> list[list[int]]:
    """Return every permutation of nums."""
    res: list[list[int]] = []
    cur: list[int] = []
    used = [False] * len(nums)

    def backtrack() -> None:
        if len(cur) == len(nums):
            res.append(cur[:])
            return
        for i, n in enumerate(nums):
            if used[i]:
                continue
            used[i] = True
            cur.append(n)
            backtrack()
            cur.pop()
            used[i] = False

    backtrack()
    return res
