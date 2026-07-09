#include <string>
#include <vector>

// Reports whether the filled cells of a 9x9 Sudoku board satisfy
// Sudoku's placement rules (no digit repeated within a row, column, or
// 3x3 box). Empty cells are '.'.
bool is_valid_sudoku(const std::vector<std::string>& board) {
    bool rows[9][9] = {};
    bool cols[9][9] = {};
    bool boxes[9][9] = {};

    for (int r = 0; r < 9; r++) {
        for (int c = 0; c < 9; c++) {
            char ch = board[r][c];
            if (ch == '.') continue;
            int d = ch - '1';
            int b = (r / 3) * 3 + c / 3;
            if (rows[r][d] || cols[c][d] || boxes[b][d]) return false;
            rows[r][d] = cols[c][d] = boxes[b][d] = true;
        }
    }
    return true;
}
