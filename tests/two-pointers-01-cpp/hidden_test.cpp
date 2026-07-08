#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> TwoSum(const std::vector<int>& numbers, int target);

int main() {
    assert((TwoSum({2, 7, 11, 15}, 9) == std::vector<int>{1, 2}));
    assert((TwoSum({2, 3, 4}, 6) == std::vector<int>{1, 3}));
    assert((TwoSum({-1, 0}, -1) == std::vector<int>{1, 2}));
    printf("all assertions passed\n");
    return 0;
}
