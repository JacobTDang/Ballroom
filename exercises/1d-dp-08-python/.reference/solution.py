def coin_change(coins: list[int], amount: int) -> int:
    sentinel = amount + 1
    dp = [sentinel] * (amount + 1)
    dp[0] = 0

    for i in range(1, amount + 1):
        for c in coins:
            if c <= i and dp[i - c] + 1 < dp[i]:
                dp[i] = dp[i - c] + 1

    return -1 if dp[amount] == sentinel else dp[amount]
