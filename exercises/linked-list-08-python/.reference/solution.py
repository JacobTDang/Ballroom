def find_duplicate(nums: list[int]) -> int:
    """Return the one repeated value in nums, using Floyd's cycle
    detection over the implicit index -> nums[index] linked list."""
    slow = fast = nums[0]
    while True:
        slow = nums[slow]
        fast = nums[nums[fast]]
        if slow == fast:
            break

    slow2 = nums[0]
    while slow2 != slow:
        slow2 = nums[slow2]
        slow = nums[slow]
    return slow
