#include <cassert>
#include <cstdio>
#include <vector>

int MissingNumber(std::vector<int>& nums);

void testClassic() {
    std::vector<int> nums = {3, 0, 1};
    assert(MissingNumber(nums) == 2);
}

void testMissingAtEnd() {
    std::vector<int> nums = {0, 1};
    assert(MissingNumber(nums) == 2);
}

void testMissingAtStart() {
    std::vector<int> nums = {1, 2, 3};
    assert(MissingNumber(nums) == 0);
}

void testSingleElement() {
    std::vector<int> nums = {0};
    assert(MissingNumber(nums) == 1);
}

void testMissingInMiddle() {
    std::vector<int> nums = {0, 1, 3, 4, 5};
    assert(MissingNumber(nums) == 2);
}

void testLargerShuffled() {
    std::vector<int> nums = {9, 6, 4, 2, 3, 5, 7, 0, 1};
    assert(MissingNumber(nums) == 8);
}

int main() {
    testClassic();
    testMissingAtEnd();
    testMissingAtStart();
    testSingleElement();
    testMissingInMiddle();
    testLargerShuffled();
    std::printf("all tests passed\n");
    return 0;
}
