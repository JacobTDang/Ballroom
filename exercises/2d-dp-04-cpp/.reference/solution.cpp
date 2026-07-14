#include <vector>

// Change returns the number of combinations of coins (unlimited
// supply of each denomination) that sum to amount.
int Change(int amount, std::vector<int>& coins) {
    std::vector<int> dp(amount + 1, 0);
    dp[0] = 1;
    for (int coin : coins) {
        for (int x = coin; x <= amount; x++) {
            dp[x] += dp[x - coin];
        }
    }
    return dp[amount];
}
