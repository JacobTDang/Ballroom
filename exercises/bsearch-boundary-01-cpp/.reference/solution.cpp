#include <vector>

// Returns the index of the first element in v that is >= target, or
// v.size() if every element is smaller.
int first_at_least(const std::vector<int>& v, int target) {
    int lo = 0, hi = static_cast<int>(v.size());
    while (lo < hi) {
        int mid = (lo + hi) / 2;
        if (v[mid] < target) {
            lo = mid + 1;
        } else {
            hi = mid;
        }
    }
    return lo;
}
