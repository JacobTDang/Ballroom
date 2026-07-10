def trap(height: list[int]) -> int:
    """Return the total units of water trapped between the bars
    described by height."""
    if not height:
        return 0
    lo, hi = 0, len(height) - 1
    left_max, right_max = height[lo], height[hi]
    total = 0
    while lo < hi:
        if left_max < right_max:
            lo += 1
            if height[lo] > left_max:
                left_max = height[lo]
            else:
                total += left_max - height[lo]
        else:
            hi -= 1
            if height[hi] > right_max:
                right_max = height[hi]
            else:
                total += right_max - height[hi]
    return total
