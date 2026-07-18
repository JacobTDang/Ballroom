#include <limits>
#include <vector>

// Returns the largest value in v that is <= limit, or -1 if no value
// qualifies.
int max_below_limit(const std::vector<int>& v, int limit) {
    int result = std::numeric_limits<int>::min();
    for (int x : v) {
        if (x <= limit && x > result) {
            result = x;
            if (result == limit) {
                break;
            }
        }
    }
    return result == std::numeric_limits<int>::min() ? -1 : result;
}
