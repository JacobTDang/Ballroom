#include <utility>
#include <vector>

// LargestRectangleArea returns the area of the largest rectangle that
// fits under the histogram described by heights.
int LargestRectangleArea(const std::vector<int>& heights) {
    std::vector<std::pair<int, int>> stack;  // (start index, height)
    int best = 0;
    int n = static_cast<int>(heights.size());
    for (int i = 0; i <= n; i++) {
        int h = (i < n) ? heights[i] : 0;
        int start = i;
        while (!stack.empty() && stack.back().second >= h) {
            auto [idx, height] = stack.back();
            stack.pop_back();
            int area = height * (i - idx);
            if (area > best) best = area;
            start = idx;
        }
        stack.emplace_back(start, h);
    }
    return best;
}
