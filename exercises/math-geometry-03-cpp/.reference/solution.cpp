#include <vector>

// SetZeroes sets the entire row and column of any zero element to zero,
// in place.
void SetZeroes(std::vector<std::vector<int>>& matrix) {
    int rows = static_cast<int>(matrix.size());
    if (rows == 0) return;
    int cols = static_cast<int>(matrix[0].size());

    std::vector<bool> zeroRow(rows, false);
    std::vector<bool> zeroCol(cols, false);

    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            if (matrix[r][c] == 0) {
                zeroRow[r] = true;
                zeroCol[c] = true;
            }
        }
    }

    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            if (zeroRow[r] || zeroCol[c]) {
                matrix[r][c] = 0;
            }
        }
    }
}
