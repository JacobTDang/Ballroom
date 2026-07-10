#include <queue>
#include <utility>
#include <vector>

const int kInf = 2147483647;

// WallsAndGates fills every empty room in rooms with its distance to
// the nearest gate, in place. Rooms that can't reach a gate stay
// kInf.
void WallsAndGates(std::vector<std::vector<int>>& rooms) {
    int rows = static_cast<int>(rooms.size());
    int cols = static_cast<int>(rooms[0].size());
    std::queue<std::pair<int, int>> q;
    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            if (rooms[r][c] == 0) q.push({r, c});
        }
    }

    int dr[4] = {1, -1, 0, 0};
    int dc[4] = {0, 0, 1, -1};
    while (!q.empty()) {
        auto [r, c] = q.front();
        q.pop();
        for (int i = 0; i < 4; i++) {
            int nr = r + dr[i], nc = c + dc[i];
            if (nr < 0 || nr >= rows || nc < 0 || nc >= cols || rooms[nr][nc] != kInf) continue;
            rooms[nr][nc] = rooms[r][c] + 1;
            q.push({nr, nc});
        }
    }
}
