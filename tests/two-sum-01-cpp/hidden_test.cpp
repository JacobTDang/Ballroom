#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> two_sum(const std::vector<int>& nums, int target);

int main() {
    assert((two_sum({2, 7, 11, 15}, 9) == std::vector<int>{0, 1}));
    assert((two_sum({3, 2, 4}, 6) == std::vector<int>{1, 2}));
    assert((two_sum({3, 3}, 6) == std::vector<int>{0, 1}));
    assert((two_sum({1, 2, 3, 4, 5}, 9) == std::vector<int>{3, 4}));
    assert((two_sum({-3, 4, 3, 90}, 0) == std::vector<int>{0, 2}));
    assert((two_sum({0, 4, 3, 0}, 0) == std::vector<int>{0, 3}));
    assert((two_sum({2, 7}, 9) == std::vector<int>{0, 1}));
    assert((two_sum({-5, -3, -1}, -8) == std::vector<int>{0, 1}));
    assert((two_sum({-1, 1}, 0) == std::vector<int>{0, 1}));
    assert((two_sum({-1000000000, 1000000000}, 0) == std::vector<int>{0, 1}));
    assert((two_sum({1000, 2000, 3000, 4000, 5000, 7, 6000, 7000, 20, 8000}, 27) == std::vector<int>{5, 8}));
    assert((two_sum({1000000000, -999999999}, 1) == std::vector<int>{0, 1}));
    printf("all assertions passed\n");
    return 0;
}
