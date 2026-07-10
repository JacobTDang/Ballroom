#include <vector>

// MissingNumber returns the one number in [0, n] missing from nums,
// where n is nums.size().
int MissingNumber(std::vector<int>& nums) {
    int result = static_cast<int>(nums.size());
    for (int i = 0; i < static_cast<int>(nums.size()); i++) {
        result ^= i ^ nums[i];
    }
    return result;
}
