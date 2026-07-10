def product_except_self(nums: list[int]) -> list[int]:
    """Return answer where answer[i] is the product of every element in
    nums except nums[i], without using division."""
    n = len(nums)
    result = [1] * n

    prefix = 1
    for i in range(n):
        result[i] = prefix
        prefix *= nums[i]

    suffix = 1
    for i in range(n - 1, -1, -1):
        result[i] *= suffix
        suffix *= nums[i]

    return result
