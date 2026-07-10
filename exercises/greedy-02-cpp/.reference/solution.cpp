#include <vector>

// CanJump returns whether the last index of nums is reachable, where
// nums[i] is the maximum jump length from index i.
bool CanJump(std::vector<int>& nums) {
    int farthest = 0;
    for (size_t i = 0; i < nums.size(); i++) {
        if ((int)i > farthest) return false;
        if ((int)i + nums[i] > farthest) farthest = (int)i + nums[i];
    }
    return true;
}
