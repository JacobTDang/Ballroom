def find_target_sum_ways(nums: list[int], target: int) -> int:
    total = sum(nums)
    if target > total or target < -total:
        return 0
    if (total + target) % 2 != 0:
        return 0
    subset_sum = (total + target) // 2

    dp = [0] * (subset_sum + 1)
    dp[0] = 1
    for n in nums:
        for s in range(subset_sum, n - 1, -1):
            dp[s] += dp[s - n]
    return dp[subset_sum]
