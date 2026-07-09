#include <vector>

// TwoSum returns the 1-indexed positions of the two numbers in numbers
// (sorted ascending) that add up to target.
std::vector<int> TwoSum(const std::vector<int>& numbers, int target) {
    int lo = 0, hi = static_cast<int>(numbers.size()) - 1;
    while (lo < hi) {
        int sum = numbers[lo] + numbers[hi];
        if (sum == target) {
            return {lo + 1, hi + 1};
        } else if (sum < target) {
            lo++;
        } else {
            hi--;
        }
    }
    return {};
}
