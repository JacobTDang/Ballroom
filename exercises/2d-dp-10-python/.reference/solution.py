def max_coins(nums: list[int]) -> int:
    n = len(nums)
    balloons = [1] + nums + [1]
    dp = [[0] * (n + 2) for _ in range(n + 2)]

    for length in range(2, n + 2):
        for l in range(0, n + 2 - length):
            r = l + length
            for k in range(l + 1, r):
                coins = dp[l][k] + dp[k][r] + balloons[l] * balloons[k] * balloons[r]
                dp[l][r] = max(dp[l][r], coins)
    return dp[0][n + 1]
