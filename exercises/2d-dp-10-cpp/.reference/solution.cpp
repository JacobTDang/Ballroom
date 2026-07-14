#include <algorithm>
#include <vector>

// MaxCoins returns the maximum coins obtainable by bursting all
// balloons in nums, where bursting balloon i yields
// nums[left] * nums[i] * nums[right] using the current neighbors.
int MaxCoins(std::vector<int>& nums) {
    int n = nums.size();
    std::vector<int> balloons(n + 2, 1);
    for (int i = 0; i < n; i++) balloons[i + 1] = nums[i];

    std::vector<std::vector<int>> dp(n + 2, std::vector<int>(n + 2, 0));

    for (int length = 2; length <= n + 1; length++) {
        for (int l = 0; l + length <= n + 1; l++) {
            int r = l + length;
            for (int k = l + 1; k < r; k++) {
                int coins = dp[l][k] + dp[k][r] + balloons[l] * balloons[k] * balloons[r];
                dp[l][r] = std::max(dp[l][r], coins);
            }
        }
    }
    return dp[0][n + 1];
}
