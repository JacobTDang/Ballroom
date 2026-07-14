#include <algorithm>
#include <vector>

// Merge merges all overlapping intervals in intervals and returns the
// resulting sorted, non-overlapping list.
std::vector<std::vector<int>> Merge(std::vector<std::vector<int>>& intervals) {
    if (intervals.empty()) return {};

    std::vector<std::vector<int>> sorted = intervals;
    std::sort(sorted.begin(), sorted.end(), [](const std::vector<int>& a, const std::vector<int>& b) {
        return a[0] < b[0];
    });

    std::vector<std::vector<int>> result;
    result.push_back(sorted[0]);
    for (size_t i = 1; i < sorted.size(); i++) {
        auto& last = result.back();
        if (sorted[i][0] <= last[1]) {
            last[1] = std::max(last[1], sorted[i][1]);
        } else {
            result.push_back(sorted[i]);
        }
    }

    return result;
}
