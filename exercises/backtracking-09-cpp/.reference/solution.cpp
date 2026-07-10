#include <functional>
#include <string>
#include <vector>

// SolveNQueens returns every distinct board configuration that
// places n queens on an n x n board with no two attacking each other.
std::vector<std::vector<std::string>> SolveNQueens(int n) {
    std::vector<std::vector<std::string>> res;
    std::vector<bool> cols(n, false), diag1(2 * n, false), diag2(2 * n, false);
    std::vector<std::string> board(n, std::string(n, '.'));

    std::function<void(int)> backtrack = [&](int r) {
        if (r == n) {
            res.push_back(board);
            return;
        }
        for (int c = 0; c < n; c++) {
            if (cols[c] || diag1[r + c] || diag2[r - c + n]) continue;
            cols[c] = diag1[r + c] = diag2[r - c + n] = true;
            board[r][c] = 'Q';
            backtrack(r + 1);
            board[r][c] = '.';
            cols[c] = diag1[r + c] = diag2[r - c + n] = false;
        }
    };
    backtrack(0);
    return res;
}
