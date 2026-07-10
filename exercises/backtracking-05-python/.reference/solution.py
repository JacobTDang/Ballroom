def combination_sum2(candidates: list[int], target: int) -> list[list[int]]:
    """Return every unique combination of candidates (each usable at
    most once) that sums to target."""
    candidates = sorted(candidates)
    res: list[list[int]] = []
    cur: list[int] = []

    def backtrack(start: int, remain: int) -> None:
        if remain == 0:
            res.append(cur[:])
            return
        for i in range(start, len(candidates)):
            if i > start and candidates[i] == candidates[i - 1]:
                continue
            if candidates[i] > remain:
                break
            cur.append(candidates[i])
            backtrack(i + 1, remain - candidates[i])
            cur.pop()

    backtrack(0, target)
    return res
