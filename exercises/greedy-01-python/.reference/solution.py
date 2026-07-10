def max_sub_array(nums: list[int]) -> int:
    best = nums[0]
    cur = nums[0]
    for n in nums[1:]:
        if cur < 0:
            cur = n
        else:
            cur += n
        best = max(best, cur)
    return best
