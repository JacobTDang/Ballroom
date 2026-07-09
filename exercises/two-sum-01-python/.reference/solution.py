def two_sum(nums: list[int], target: int) -> list[int]:
    """Return the indices of the two numbers in nums that add up to
    target."""
    seen: dict[int, int] = {}
    for i, n in enumerate(nums):
        complement = target - n
        if complement in seen:
            return [seen[complement], i]
        seen[n] = i
    return []
