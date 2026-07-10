#include <algorithm>
#include <queue>
#include <tuple>
#include <vector>

// SwimInWater returns the minimum time t such that you can swim from the
// top-left to the bottom-right of grid, where at time t you may move
// between adjacent cells whose elevation is <= t.
int SwimInWater(std::vector<std::vector<int>>& grid) {
    int n = static_cast<int>(grid.size());
    if (n == 0) return 0;

    std::vector<std::vector<bool>> visited(n, std::vector<bool>(n, false));

    using State = std::tuple<int, int, int>; // elevation, row, col
    std::priority_queue<State, std::vector<State>, std::greater<State>> pq;
    pq.push({grid[0][0], 0, 0});
    visited[0][0] = true;

    int dirs[4][2] = {{1, 0}, {-1, 0}, {0, 1}, {0, -1}};

    while (!pq.empty()) {
        auto [elevation, row, col] = pq.top();
        pq.pop();
        if (row == n - 1 && col == n - 1) {
            return elevation;
        }
        for (auto& d : dirs) {
            int nr = row + d[0], nc = col + d[1];
            if (nr < 0 || nr >= n || nc < 0 || nc >= n || visited[nr][nc]) continue;
            visited[nr][nc] = true;
            int maxElevation = std::max(elevation, grid[nr][nc]);
            pq.push({maxElevation, nr, nc});
        }
    }
    return -1;
}
