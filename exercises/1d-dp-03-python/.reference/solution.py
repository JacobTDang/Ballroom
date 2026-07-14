def rob(nums: list[int]) -> int:
    prev, curr = 0, 0
    for n in nums:
        next_val = curr
        alt = prev + n
        if alt > next_val:
            next_val = alt
        prev, curr = curr, next_val
    return curr
