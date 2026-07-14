#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> MaxSlidingWindow(const std::vector<int>& nums, int k);

int main() {
    assert((MaxSlidingWindow({1, 3, -1, -3, 5, 3, 6, 7}, 3) ==
            std::vector<int>{3, 3, 5, 5, 6, 7}));
    assert((MaxSlidingWindow({1}, 1) == std::vector<int>{1}));
    assert((MaxSlidingWindow({1, -1}, 1) == std::vector<int>{1, -1}));
    assert((MaxSlidingWindow({9, 11}, 2) == std::vector<int>{11}));
    assert((MaxSlidingWindow({4, -2}, 2) == std::vector<int>{4}));
    assert((MaxSlidingWindow({1, 3, 1, 2, 0, 5}, 3) == std::vector<int>{3, 3, 2, 5}));
    assert((MaxSlidingWindow({7, 2, 4}, 2) == std::vector<int>{7, 4}));
    assert((MaxSlidingWindow({1, 2, 3, 4, 5}, 5) == std::vector<int>{5}));
    assert((MaxSlidingWindow({-7, -8, 7, 5, 7, 1, 6, 0}, 4) ==
            std::vector<int>{7, 7, 7, 7, 7}));
    printf("all assertions passed\n");
    return 0;
}
