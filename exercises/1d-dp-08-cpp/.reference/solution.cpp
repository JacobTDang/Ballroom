#include <vector>

// CoinChange returns the fewest number of coins from coins (unlimited
// supply of each denomination) needed to make up amount, or -1 if
// amount cannot be made up by any combination of the coins.
int CoinChange(std::vector<int>& coins, int amount) {
    int sentinel = amount + 1;

    std::vector<int> dp(amount + 1, sentinel);
    dp[0] = 0;

    for (int i = 1; i <= amount; i++) {
        for (int c : coins) {
            if (c <= i && dp[i - c] + 1 < dp[i]) {
                dp[i] = dp[i - c] + 1;
            }
        }
    }

    if (dp[amount] == sentinel) return -1;
    return dp[amount];
}
