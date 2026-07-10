def find_min(nums: list[int]) -> int:
    """Return the minimum element of a sorted array that has been
    rotated between 1 and len(nums) times."""
    lo, hi = 0, len(nums) - 1
    while lo < hi:
        mid = lo + (hi - lo) // 2
        if nums[mid] > nums[hi]:
            lo = mid + 1
        else:
            hi = mid
    return nums[lo]
