#include <cassert>
#include <cstdio>
#include <vector>

int MaxSubArray(std::vector<int>& nums);

void testClassic() {
    std::vector<int> nums = {-2, 1, -3, 4, -1, 2, 1, -5, 4};
    assert(MaxSubArray(nums) == 6);
}

void testAllNegative() {
    std::vector<int> nums = {-3, -2, -1};
    assert(MaxSubArray(nums) == -1);
}

void testAllPositive() {
    std::vector<int> nums = {1, 2, 3, 4};
    assert(MaxSubArray(nums) == 10);
}

void testSingleElement() {
    std::vector<int> nums = {5};
    assert(MaxSubArray(nums) == 5);
}

void testLargeNegativeInMiddle() {
    std::vector<int> nums = {5, 4, -20, 7, 8};
    assert(MaxSubArray(nums) == 15);
}

int main() {
    testClassic();
    testAllNegative();
    testAllPositive();
    testSingleElement();
    testLargeNegativeInMiddle();
    std::printf("all tests passed\n");
    return 0;
}
