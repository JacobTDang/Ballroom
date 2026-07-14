def can_partition(nums: list[int]) -> bool:
    total = sum(nums)
    if total % 2 != 0:
        return False

    target = total // 2
    reachable = [False] * (target + 1)
    reachable[0] = True

    for n in nums:
        for i in range(target, n - 1, -1):
            if reachable[i - n]:
                reachable[i] = True

    return reachable[target]
