#include <vector>

// Jump returns the minimum number of jumps needed to reach the last
// index of nums, where nums[i] is the maximum jump length from index i.
int Jump(std::vector<int>& nums) {
    int jumps = 0;
    int curEnd = 0;
    int farthest = 0;
    for (int i = 0; i + 1 < (int)nums.size(); i++) {
        if (i + nums[i] > farthest) farthest = i + nums[i];
        if (i == curEnd) {
            jumps++;
            curEnd = farthest;
        }
    }
    return jumps;
}
