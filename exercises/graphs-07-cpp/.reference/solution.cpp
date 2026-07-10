#include <functional>
#include <vector>

// Solve captures all regions of 'O' surrounded by 'X' by flipping
// them to 'X' in place. Regions connected to the border are left
// untouched.
void Solve(std::vector<std::vector<char>>& board) {
    int rows = static_cast<int>(board.size());
    if (rows == 0) return;
    int cols = static_cast<int>(board[0].size());
    if (cols == 0) return;

    std::function<void(int, int)> dfs = [&](int r, int c) {
        if (r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] != 'O') return;
        board[r][c] = '#';
        dfs(r + 1, c);
        dfs(r - 1, c);
        dfs(r, c + 1);
        dfs(r, c - 1);
    };

    for (int c = 0; c < cols; c++) {
        dfs(0, c);
        dfs(rows - 1, c);
    }
    for (int r = 0; r < rows; r++) {
        dfs(r, 0);
        dfs(r, cols - 1);
    }

    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            if (board[r][c] == 'O') {
                board[r][c] = 'X';
            } else if (board[r][c] == '#') {
                board[r][c] = 'O';
            }
        }
    }
}
