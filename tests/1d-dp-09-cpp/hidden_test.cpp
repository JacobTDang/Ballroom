#include <cassert>
#include <cstdio>
#include <vector>

int MaxProduct(std::vector<int>& nums);

void testClassic() {
    std::vector<int> nums = {2, 3, -2, 4};
    assert(MaxProduct(nums) == 6);
}

void testZeroSplits() {
    std::vector<int> nums = {-2, 0, -1};
    assert(MaxProduct(nums) == 0);
}

void testTwoNegativesFlip() {
    std::vector<int> nums = {-2, 3, -4};
    assert(MaxProduct(nums) == 24);
}

void testSingleNegative() {
    std::vector<int> nums = {-5};
    assert(MaxProduct(nums) == -5);
}

void testAllPositive() {
    std::vector<int> nums = {1, 2, 3, 4};
    assert(MaxProduct(nums) == 24);
}

void testSinglePositive() {
    std::vector<int> nums = {7};
    assert(MaxProduct(nums) == 7);
}

void testMultipleZeroSplitIslands() {
    std::vector<int> nums = {0, 2, 0, 3, 0};
    assert(MaxProduct(nums) == 3);
}

void testWholeArrayEvenNegatives() {
    std::vector<int> nums = {2, -3, 4, -5};
    assert(MaxProduct(nums) == 120);
}

int main() {
    testClassic();
    testZeroSplits();
    testTwoNegativesFlip();
    testSingleNegative();
    testAllPositive();
    testSinglePositive();
    testMultipleZeroSplitIslands();
    testWholeArrayEvenNegatives();
    std::printf("all tests passed\n");
    return 0;
}
