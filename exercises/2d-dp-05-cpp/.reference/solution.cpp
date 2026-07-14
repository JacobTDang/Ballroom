#include <vector>

// FindTargetSumWays counts the number of ways to assign + or - to
// each num in nums so that the resulting expression evaluates to
// target.
int FindTargetSumWays(std::vector<int>& nums, int target) {
    int total = 0;
    for (int n : nums) total += n;
    if (target > total || target < -total) return 0;
    if ((total + target) % 2 != 0) return 0;
    int subsetSum = (total + target) / 2;

    std::vector<int> dp(subsetSum + 1, 0);
    dp[0] = 1;
    for (int n : nums) {
        for (int s = subsetSum; s >= n; s--) {
            dp[s] += dp[s - n];
        }
    }
    return dp[subsetSum];
}
