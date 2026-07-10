#include <vector>

// SearchMatrix reports whether target is present in matrix, treating
// it as one flattened sorted sequence.
bool SearchMatrix(const std::vector<std::vector<int>>& matrix, int target) {
    if (matrix.empty() || matrix[0].empty()) return false;
    int rows = static_cast<int>(matrix.size());
    int cols = static_cast<int>(matrix[0].size());
    int lo = 0, hi = rows * cols - 1;
    while (lo <= hi) {
        int mid = lo + (hi - lo) / 2;
        int val = matrix[mid / cols][mid % cols];
        if (val == target) return true;
        if (val < target) {
            lo = mid + 1;
        } else {
            hi = mid - 1;
        }
    }
    return false;
}
