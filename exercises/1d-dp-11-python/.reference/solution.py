def length_of_lis(nums: list[int]) -> int:
    n = len(nums)
    dp = [1] * n
    best = 0

    for i in range(n):
        for j in range(i):
            if nums[j] < nums[i] and dp[j] + 1 > dp[i]:
                dp[i] = dp[j] + 1
        best = max(best, dp[i])

    return best
