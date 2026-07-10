#include <vector>

// SingleNumber returns the element of nums that appears exactly once,
// given every other element appears exactly twice.
int SingleNumber(std::vector<int>& nums) {
    int result = 0;
    for (int n : nums) {
        result ^= n;
    }
    return result;
}
