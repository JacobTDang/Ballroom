#include <functional>
#include <vector>

// PacificAtlantic returns every cell [r, c] from which water can
// flow to both the Pacific (top/left edges) and Atlantic
// (bottom/right edges) oceans.
std::vector<std::vector<int>> PacificAtlantic(std::vector<std::vector<int>>& heights) {
    int rows = static_cast<int>(heights.size());
    int cols = static_cast<int>(heights[0].size());
    std::vector<std::vector<bool>> pacific(rows, std::vector<bool>(cols, false));
    std::vector<std::vector<bool>> atlantic(rows, std::vector<bool>(cols, false));

    std::function<void(int, int, std::vector<std::vector<bool>>&, int)> dfs =
        [&](int r, int c, std::vector<std::vector<bool>>& visited, int prevHeight) {
            if (r < 0 || r >= rows || c < 0 || c >= cols || visited[r][c] ||
                heights[r][c] < prevHeight) {
                return;
            }
            visited[r][c] = true;
            dfs(r + 1, c, visited, heights[r][c]);
            dfs(r - 1, c, visited, heights[r][c]);
            dfs(r, c + 1, visited, heights[r][c]);
            dfs(r, c - 1, visited, heights[r][c]);
        };

    for (int c = 0; c < cols; c++) {
        dfs(0, c, pacific, heights[0][c]);
        dfs(rows - 1, c, atlantic, heights[rows - 1][c]);
    }
    for (int r = 0; r < rows; r++) {
        dfs(r, 0, pacific, heights[r][0]);
        dfs(r, cols - 1, atlantic, heights[r][cols - 1]);
    }

    std::vector<std::vector<int>> res;
    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            if (pacific[r][c] && atlantic[r][c]) res.push_back({r, c});
        }
    }
    return res;
}
