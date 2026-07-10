#include <vector>

// FindMin returns the minimum element of a sorted array that has
// been rotated between 1 and nums.size() times.
int FindMin(const std::vector<int>& nums) {
    int lo = 0, hi = static_cast<int>(nums.size()) - 1;
    while (lo < hi) {
        int mid = lo + (hi - lo) / 2;
        if (nums[mid] > nums[hi]) {
            lo = mid + 1;
        } else {
            hi = mid;
        }
    }
    return nums[lo];
}
