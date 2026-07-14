#include <vector>

// UniquePaths returns the number of unique paths from the top-left to
// the bottom-right of an m x n grid, moving only right or down.
int UniquePaths(int m, int n) {
    std::vector<int> dp(n, 1);
    for (int r = 1; r < m; r++) {
        for (int c = 1; c < n; c++) {
            dp[c] += dp[c - 1];
        }
    }
    return dp[n - 1];
}
