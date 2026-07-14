#include <algorithm>
#include <functional>
#include <vector>

// LongestIncreasingPath returns the length of the longest strictly
// increasing path in matrix, moving 4-directionally.
int LongestIncreasingPath(std::vector<std::vector<int>>& matrix) {
    if (matrix.empty() || matrix[0].empty()) return 0;
    int rows = matrix.size(), cols = matrix[0].size();
    std::vector<std::vector<int>> memo(rows, std::vector<int>(cols, 0));
    int dr[4] = {1, -1, 0, 0};
    int dc[4] = {0, 0, 1, -1};

    std::function<int(int, int)> dfs = [&](int r, int c) -> int {
        if (memo[r][c] != 0) return memo[r][c];
        int best = 1;
        for (int k = 0; k < 4; k++) {
            int nr = r + dr[k], nc = c + dc[k];
            if (nr >= 0 && nr < rows && nc >= 0 && nc < cols && matrix[nr][nc] > matrix[r][c]) {
                int length = 1 + dfs(nr, nc);
                if (length > best) best = length;
            }
        }
        memo[r][c] = best;
        return best;
    };

    int result = 0;
    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            result = std::max(result, dfs(r, c));
        }
    }
    return result;
}
