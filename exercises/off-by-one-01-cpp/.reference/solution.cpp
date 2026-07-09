#include <vector>

// Returns the largest value in v.
int max_of(const std::vector<int>& v) {
    int best = v[0];
    for (size_t i = 1; i < v.size(); i++) {
        if (v[i] > best) best = v[i];
    }
    return best;
}
