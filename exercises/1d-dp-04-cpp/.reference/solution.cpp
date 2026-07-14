#include <algorithm>
#include <vector>

// robLinear is the House Robber I logic for a non-circular line of
// houses.
static int robLinear(std::vector<int>& nums) {
    int prev = 0, curr = 0;
    for (int n : nums) {
        int next = curr;
        int alt = prev + n;
        if (alt > next) next = alt;
        prev = curr;
        curr = next;
    }
    return curr;
}

// RobCircular returns the maximum amount of money that can be robbed
// from houses arranged in a circle (house 0 and house n-1 are
// adjacent), given nums[i] is the money in house i, without robbing
// two adjacent houses.
int RobCircular(std::vector<int>& nums) {
    int n = nums.size();
    if (n == 1) return nums[0];

    std::vector<int> withoutLast(nums.begin(), nums.end() - 1);
    std::vector<int> withoutFirst(nums.begin() + 1, nums.end());
    return std::max(robLinear(withoutLast), robLinear(withoutFirst));
}
