#include <queue>
#include <utility>
#include <vector>

// OrangesRotting returns the minimum number of minutes until no cell
// in grid has a fresh orange, or -1 if some fresh orange can never
// rot.
int OrangesRotting(std::vector<std::vector<int>>& grid) {
    int rows = static_cast<int>(grid.size());
    int cols = static_cast<int>(grid[0].size());
    std::queue<std::pair<int, int>> q;
    int fresh = 0;
    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            if (grid[r][c] == 2) {
                q.push({r, c});
            } else if (grid[r][c] == 1) {
                fresh++;
            }
        }
    }
    if (fresh == 0) return 0;

    int dr[4] = {1, -1, 0, 0};
    int dc[4] = {0, 0, 1, -1};
    int minutes = 0;
    while (!q.empty() && fresh > 0) {
        int size = static_cast<int>(q.size());
        for (int i = 0; i < size; i++) {
            auto [r, c] = q.front();
            q.pop();
            for (int j = 0; j < 4; j++) {
                int nr = r + dr[j], nc = c + dc[j];
                if (nr < 0 || nr >= rows || nc < 0 || nc >= cols || grid[nr][nc] != 1) continue;
                grid[nr][nc] = 2;
                fresh--;
                q.push({nr, nc});
            }
        }
        minutes++;
    }

    return fresh > 0 ? -1 : minutes;
}
