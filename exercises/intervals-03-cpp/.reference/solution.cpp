#include <algorithm>
#include <vector>

// EraseOverlapIntervals returns the minimum number of intervals that
// must be removed so the rest of intervals are non-overlapping.
int EraseOverlapIntervals(std::vector<std::vector<int>>& intervals) {
    if (intervals.empty()) return 0;

    std::vector<std::vector<int>> sorted = intervals;
    std::sort(sorted.begin(), sorted.end(), [](const std::vector<int>& a, const std::vector<int>& b) {
        return a[1] < b[1];
    });

    int removals = 0;
    int lastEnd = sorted[0][1];
    for (size_t i = 1; i < sorted.size(); i++) {
        if (sorted[i][0] < lastEnd) {
            removals++;
        } else {
            lastEnd = sorted[i][1];
        }
    }

    return removals;
}
