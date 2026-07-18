#include <vector>

// Recursive helper for count_paths: counts paths from (r, c) to the
// bottom-right corner, moving only right or down, avoiding cells
// marked 1 (blocked).
int count_paths_helper(const std::vector<std::vector<int>>& grid, int r, int c, int rows, int cols) {
    if (r >= rows || c >= cols || grid[r][c] == 1) {
        return 0;
    }
    if (r == rows - 1 && c == cols - 1) {
        return 1;
    }
    return count_paths_helper(grid, r + 1, c, rows, cols) + count_paths_helper(grid, r, c + 1, rows, cols);
}

// Counts the number of paths from the top-left to the bottom-right
// corner of grid, moving only right or down, that avoid cells marked
// 1 (blocked).
int count_paths(const std::vector<std::vector<int>>& grid) {
    int rows = static_cast<int>(grid.size());
    int cols = static_cast<int>(grid[0].size());
    return count_paths_helper(grid, 0, 0, rows, cols);
}
