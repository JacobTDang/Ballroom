#include <algorithm>
#include <functional>
#include <vector>

// MaxAreaOfIsland returns the number of cells in the largest
// 4-directionally connected island of 1s in grid, or 0 if there is
// no island.
int MaxAreaOfIsland(std::vector<std::vector<int>>& grid) {
    int rows = static_cast<int>(grid.size());
    int cols = static_cast<int>(grid[0].size());

    std::function<int(int, int)> dfs = [&](int r, int c) -> int {
        if (r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] != 1) return 0;
        grid[r][c] = 0;
        return 1 + dfs(r + 1, c) + dfs(r - 1, c) + dfs(r, c + 1) + dfs(r, c - 1);
    };

    int best = 0;
    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            if (grid[r][c] == 1) best = std::max(best, dfs(r, c));
        }
    }
    return best;
}
