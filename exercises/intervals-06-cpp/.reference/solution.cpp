#include <algorithm>
#include <functional>
#include <queue>
#include <vector>

// MinInterval returns, for each query, the size of the smallest
// interval that contains it (left <= query <= right), or -1 if no
// interval contains it. The result is in the same order as queries.
std::vector<int> MinInterval(std::vector<std::vector<int>>& intervals, std::vector<int>& queries) {
    std::vector<std::vector<int>> sorted = intervals;
    std::sort(sorted.begin(), sorted.end(), [](const std::vector<int>& a, const std::vector<int>& b) {
        return a[0] < b[0];
    });

    int n = static_cast<int>(queries.size());
    std::vector<std::pair<int, int>> order(n);  // {value, original index}
    for (int i = 0; i < n; i++) order[i] = {queries[i], i};
    std::sort(order.begin(), order.end());

    std::vector<int> result(n);
    // Min-heap of {size, end}, smallest interval size on top.
    std::priority_queue<std::pair<int, int>, std::vector<std::pair<int, int>>, std::greater<>> heap;

    size_t i = 0;
    for (auto& [val, idx] : order) {
        while (i < sorted.size() && sorted[i][0] <= val) {
            int left = sorted[i][0], right = sorted[i][1];
            heap.push({right - left + 1, right});
            i++;
        }

        while (!heap.empty() && heap.top().second < val) {
            heap.pop();
        }

        result[idx] = heap.empty() ? -1 : heap.top().first;
    }

    return result;
}
