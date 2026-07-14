#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> TwoSum(const std::vector<int>& numbers, int target);

int main() {
    assert((TwoSum({2, 7, 11, 15}, 9) == std::vector<int>{1, 2}));
    assert((TwoSum({2, 3, 4}, 6) == std::vector<int>{1, 3}));
    assert((TwoSum({-1, 0}, -1) == std::vector<int>{1, 2}));
    assert((TwoSum({3, 3}, 6) == std::vector<int>{1, 2}));
    assert((TwoSum({1, 2, 3, 4, 4, 9, 56, 90}, 8) == std::vector<int>{4, 5}));
    assert((TwoSum({-8, -3, 0, 4, 9, 13}, 5) == std::vector<int>{1, 6}));
    assert((TwoSum({-8, -3, 0, 4, 9, 13}, 4) == std::vector<int>{3, 4}));
    assert((TwoSum({-3, -1, 0, 2, 4, 5}, 1) == std::vector<int>{1, 5}));
    assert((TwoSum({5, 25, 75}, 100) == std::vector<int>{2, 3}));
    printf("all assertions passed\n");
    return 0;
}
