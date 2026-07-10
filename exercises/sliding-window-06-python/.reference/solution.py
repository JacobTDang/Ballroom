from collections import deque


def max_sliding_window(nums: list[int], k: int) -> list[int]:
    """Return the maximum of each contiguous window of size k as it
    slides from the start of nums to the end."""
    dq: deque[int] = deque()  # indices into nums, values strictly decreasing
    res = []
    for i, n in enumerate(nums):
        while dq and nums[dq[-1]] < n:
            dq.pop()
        dq.append(i)
        if dq[0] <= i - k:
            dq.popleft()
        if i >= k - 1:
            res.append(nums[dq[0]])
    return res
