#include <cassert>
#include <cstdio>
#include <vector>

int LengthOfLIS(std::vector<int>& nums);

void testClassic() {
    std::vector<int> nums = {10, 9, 2, 5, 3, 7, 101, 18};
    assert(LengthOfLIS(nums) == 4);
}

void testRepeatedDip() {
    std::vector<int> nums = {0, 1, 0, 3, 2, 3};
    assert(LengthOfLIS(nums) == 4);
}

void testAllEqual() {
    std::vector<int> nums = {7, 7, 7, 7};
    assert(LengthOfLIS(nums) == 1);
}

void testSingleElement() {
    std::vector<int> nums = {5};
    assert(LengthOfLIS(nums) == 1);
}

void testStrictlyDecreasing() {
    std::vector<int> nums = {5, 4, 3, 2, 1};
    assert(LengthOfLIS(nums) == 1);
}

void testStrictlyIncreasing() {
    std::vector<int> nums = {1, 2, 3, 4, 5};
    assert(LengthOfLIS(nums) == 5);
}

void testNegativeValues() {
    std::vector<int> nums = {-1, -2, 0, 1, -3, 5};
    assert(LengthOfLIS(nums) == 4);
}

void testBoundaryValues() {
    std::vector<int> nums = {-10000, 10000};
    assert(LengthOfLIS(nums) == 2);
}

void testDuplicatesBreakStreak() {
    std::vector<int> nums = {3, 3, 3, 4, 4, 5};
    assert(LengthOfLIS(nums) == 3);
}

int main() {
    testClassic();
    testRepeatedDip();
    testAllEqual();
    testSingleElement();
    testStrictlyDecreasing();
    testStrictlyIncreasing();
    testNegativeValues();
    testBoundaryValues();
    testDuplicatesBreakStreak();
    std::printf("all tests passed\n");
    return 0;
}
