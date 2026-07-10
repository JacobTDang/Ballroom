#include <algorithm>
#include <vector>

// ThreeSum returns every unique triplet of elements in nums that sums
// to zero.
std::vector<std::vector<int>> ThreeSum(std::vector<int>& nums) {
    std::sort(nums.begin(), nums.end());
    std::vector<std::vector<int>> res;
    int n = static_cast<int>(nums.size());
    for (int i = 0; i < n - 2; i++) {
        if (i > 0 && nums[i] == nums[i - 1]) continue;
        int lo = i + 1, hi = n - 1;
        while (lo < hi) {
            int sum = nums[i] + nums[lo] + nums[hi];
            if (sum < 0) {
                lo++;
            } else if (sum > 0) {
                hi--;
            } else {
                res.push_back({nums[i], nums[lo], nums[hi]});
                lo++;
                hi--;
                while (lo < hi && nums[lo] == nums[lo - 1]) lo++;
                while (lo < hi && nums[hi] == nums[hi + 1]) hi--;
            }
        }
    }
    return res;
}
