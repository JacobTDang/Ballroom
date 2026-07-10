#include <algorithm>
#include <vector>

// KClosest returns the k points from points closest to the origin,
// in any order.
std::vector<std::vector<int>> KClosest(std::vector<std::vector<int>>& points, int k) {
    std::vector<std::vector<int>> sorted = points;
    std::sort(sorted.begin(), sorted.end(), [](const std::vector<int>& a, const std::vector<int>& b) {
        return a[0] * a[0] + a[1] * a[1] < b[0] * b[0] + b[1] * b[1];
    });
    sorted.resize(k);
    return sorted;
}
