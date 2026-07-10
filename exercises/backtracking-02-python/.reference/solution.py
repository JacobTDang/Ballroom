def combination_sum(candidates: list[int], target: int) -> list[list[int]]:
    """Return every unique combination of candidates (each usable
    unlimited times) that sums to target."""
    res: list[list[int]] = []
    cur: list[int] = []

    def backtrack(start: int, remain: int) -> None:
        if remain == 0:
            res.append(cur[:])
            return
        if remain < 0:
            return
        for i in range(start, len(candidates)):
            cur.append(candidates[i])
            backtrack(i, remain - candidates[i])
            cur.pop()

    backtrack(0, target)
    return res
