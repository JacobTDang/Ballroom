#include <cassert>
#include <cstdio>
#include <vector>

bool CanJump(std::vector<int>& nums);

void testClassic() {
    std::vector<int> nums = {2, 3, 1, 1, 4};
    assert(CanJump(nums) == true);
}

void testClassicFalse() {
    std::vector<int> nums = {3, 2, 1, 0, 4};
    assert(CanJump(nums) == false);
}

void testSingleElement() {
    std::vector<int> nums = {0};
    assert(CanJump(nums) == true);
}

void testZeroAtStartBlocks() {
    std::vector<int> nums = {0, 1};
    assert(CanJump(nums) == false);
}

void testBigFirstJumpCoversRest() {
    std::vector<int> nums = {5, 0, 0, 0, 0};
    assert(CanJump(nums) == true);
}

int main() {
    testClassic();
    testClassicFalse();
    testSingleElement();
    testZeroAtStartBlocks();
    testBigFirstJumpCoversRest();
    std::printf("all tests passed\n");
    return 0;
}
