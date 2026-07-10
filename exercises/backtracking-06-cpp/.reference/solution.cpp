#include <functional>
#include <string>
#include <vector>

// Exist reports whether word can be traced out on board via
// sequentially adjacent cells, each cell used at most once.
bool Exist(std::vector<std::vector<char>>& board, std::string word) {
    int rows = static_cast<int>(board.size());
    int cols = static_cast<int>(board[0].size());

    std::function<bool(int, int, int)> dfs = [&](int r, int c, int idx) -> bool {
        if (idx == static_cast<int>(word.size())) return true;
        if (r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] != word[idx]) return false;
        char tmp = board[r][c];
        board[r][c] = '#';
        bool found = dfs(r + 1, c, idx + 1) || dfs(r - 1, c, idx + 1) ||
                     dfs(r, c + 1, idx + 1) || dfs(r, c - 1, idx + 1);
        board[r][c] = tmp;
        return found;
    };

    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            if (dfs(r, c, 0)) return true;
        }
    }
    return false;
}
