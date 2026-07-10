#include <vector>

// Trap returns the total units of water trapped between the bars
// described by height.
int Trap(const std::vector<int>& height) {
    if (height.empty()) return 0;
    int lo = 0, hi = static_cast<int>(height.size()) - 1;
    int leftMax = height[lo], rightMax = height[hi];
    int total = 0;
    while (lo < hi) {
        if (leftMax < rightMax) {
            lo++;
            if (height[lo] > leftMax) {
                leftMax = height[lo];
            } else {
                total += leftMax - height[lo];
            }
        } else {
            hi--;
            if (height[hi] > rightMax) {
                rightMax = height[hi];
            } else {
                total += rightMax - height[hi];
            }
        }
    }
    return total;
}
