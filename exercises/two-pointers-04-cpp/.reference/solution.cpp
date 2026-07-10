#include <algorithm>
#include <vector>

// MaxArea returns the largest amount of water a container formed by
// two lines in height (with the x-axis) can hold.
int MaxArea(const std::vector<int>& height) {
    int lo = 0, hi = static_cast<int>(height.size()) - 1;
    int best = 0;
    while (lo < hi) {
        int h = std::min(height[lo], height[hi]);
        best = std::max(best, h * (hi - lo));
        if (height[lo] < height[hi]) {
            lo++;
        } else {
            hi--;
        }
    }
    return best;
}
