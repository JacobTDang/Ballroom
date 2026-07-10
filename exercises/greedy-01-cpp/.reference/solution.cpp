#include <vector>

// MaxSubArray returns the largest sum of any contiguous subarray of nums.
int MaxSubArray(std::vector<int>& nums) {
    int best = nums[0];
    int cur = nums[0];
    for (size_t i = 1; i < nums.size(); i++) {
        if (cur < 0) {
            cur = nums[i];
        } else {
            cur += nums[i];
        }
        if (cur > best) best = cur;
    }
    return best;
}
