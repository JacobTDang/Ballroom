def missing_number(nums: list[int]) -> int:
    result = len(nums)
    for i, v in enumerate(nums):
        result ^= i ^ v
    return result
