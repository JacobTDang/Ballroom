#include <cstdlib>
#include <stdexcept>
#include <vector>

// Returns the largest absolute difference between two adjacent
// elements in v. Throws std::invalid_argument if v has fewer than two
// elements. Currently crashes under AddressSanitizer on some inputs —
// find and fix the bug.
int max_adjacent_diff(const std::vector<int>& v) {
    if (v.empty()) {
        throw std::invalid_argument("max_adjacent_diff: need at least two values");
    }
    int best = std::abs(v[1] - v[0]);
    for (size_t i = 1; i + 1 < v.size(); i++) {
        int d = std::abs(v[i + 1] - v[i]);
        if (d > best) best = d;
    }
    return best;
}
