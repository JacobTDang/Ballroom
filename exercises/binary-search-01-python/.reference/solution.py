def search(nums: list[int], target: int) -> int:
    """Return the index of target in the sorted nums, or -1 if it's
    not present."""
    lo, hi = 0, len(nums) - 1
    while lo <= hi:
        mid = lo + (hi - lo) // 2
        if nums[mid] == target:
            return mid
        elif nums[mid] < target:
            lo = mid + 1
        else:
            hi = mid - 1
    return -1
