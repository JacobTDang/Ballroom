def _rob_linear(nums: list[int]) -> int:
    """House Robber I logic for a non-circular line of houses."""
    prev, curr = 0, 0
    for n in nums:
        next_val = curr
        alt = prev + n
        if alt > next_val:
            next_val = alt
        prev, curr = curr, next_val
    return curr


def rob_circular(nums: list[int]) -> int:
    n = len(nums)
    if n == 1:
        return nums[0]
    return max(_rob_linear(nums[:-1]), _rob_linear(nums[1:]))
