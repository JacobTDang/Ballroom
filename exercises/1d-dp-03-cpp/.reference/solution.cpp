#include <vector>

// Rob returns the maximum amount of money that can be robbed from
// houses arranged in a line, given nums[i] is the money in house i,
// without robbing two adjacent houses.
int Rob(std::vector<int>& nums) {
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
