#include <vector>

// RotateImage rotates the n x n matrix 90 degrees clockwise, in place.
void RotateImage(std::vector<std::vector<int>>& matrix) {
    int n = static_cast<int>(matrix.size());

    // Transpose the matrix (reflect across the main diagonal).
    for (int i = 0; i < n; i++) {
        for (int j = i + 1; j < n; j++) {
            std::swap(matrix[i][j], matrix[j][i]);
        }
    }

    // Reverse each row to complete the clockwise rotation.
    for (int i = 0; i < n; i++) {
        for (int l = 0, r = n - 1; l < r; l++, r--) {
            std::swap(matrix[i][l], matrix[i][r]);
        }
    }
}
