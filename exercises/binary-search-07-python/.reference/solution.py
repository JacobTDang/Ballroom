def find_median_sorted_arrays(nums1: list[int], nums2: list[int]) -> float:
    """Return the median of the two sorted arrays nums1 and nums2
    combined."""
    if len(nums1) > len(nums2):
        nums1, nums2 = nums2, nums1
    m, n = len(nums1), len(nums2)
    lo, hi = 0, m
    half = (m + n + 1) // 2
    inf = float("inf")

    while lo <= hi:
        i = lo + (hi - lo) // 2
        j = half - i

        max_left_a = nums1[i - 1] if i > 0 else -inf
        min_right_a = nums1[i] if i < m else inf
        max_left_b = nums2[j - 1] if j > 0 else -inf
        min_right_b = nums2[j] if j < n else inf

        if max_left_a <= min_right_b and max_left_b <= min_right_a:
            if (m + n) % 2 == 1:
                return float(max(max_left_a, max_left_b))
            return (max(max_left_a, max_left_b) + min(min_right_a, min_right_b)) / 2.0
        elif max_left_a > min_right_b:
            hi = i - 1
        else:
            lo = i + 1
    return 0.0  # unreachable for valid, well-formed input
