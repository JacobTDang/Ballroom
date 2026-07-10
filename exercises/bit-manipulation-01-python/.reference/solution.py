def single_number(nums: list[int]) -> int:
    result = 0
    for n in nums:
        result ^= n
    return result
