def longest_consecutive(nums: list[int]) -> int:
    """Return the length of the longest run of consecutive integers
    present in nums (order doesn't matter, duplicates don't count
    extra)."""
    num_set = set(nums)
    longest = 0
    for n in num_set:
        if n - 1 in num_set:
            continue  # n isn't the start of a sequence
        length = 1
        while n + length in num_set:
            length += 1
        longest = max(longest, length)
    return longest
