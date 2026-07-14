#include <vector>

// LengthOfLIS returns the length of the longest strictly increasing
// subsequence of nums.
int LengthOfLIS(std::vector<int>& nums) {
    int n = nums.size();
    std::vector<int> dp(n, 1);
    int best = 0;

    for (int i = 0; i < n; i++) {
        for (int j = 0; j < i; j++) {
            if (nums[j] < nums[i] && dp[j] + 1 > dp[i]) {
                dp[i] = dp[j] + 1;
            }
        }
        if (dp[i] > best) best = dp[i];
    }

    return best;
}
