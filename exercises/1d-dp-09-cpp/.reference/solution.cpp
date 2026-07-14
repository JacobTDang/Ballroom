#include <algorithm>
#include <vector>

// MaxProduct returns the largest product of any contiguous non-empty
// subarray of nums.
int MaxProduct(std::vector<int>& nums) {
    int result = nums[0];
    int curMax = nums[0], curMin = nums[0];

    for (size_t i = 1; i < nums.size(); i++) {
        int n = nums[i];
        if (n < 0) std::swap(curMax, curMin);
        curMax = std::max(n, curMax * n);
        curMin = std::min(n, curMin * n);
        result = std::max(result, curMax);
    }

    return result;
}
