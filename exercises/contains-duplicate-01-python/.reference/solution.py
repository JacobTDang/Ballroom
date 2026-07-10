def contains_duplicate(nums: list[int]) -> bool:
    """Return True if any value appears at least twice in nums."""
    seen = set()
    for n in nums:
        if n in seen:
            return True
        seen.add(n)
    return False
