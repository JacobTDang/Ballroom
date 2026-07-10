#include <cassert>
#include <cstdio>
#include <vector>

int Jump(std::vector<int>& nums);

void testClassic() {
    std::vector<int> nums = {2, 3, 1, 1, 4};
    assert(Jump(nums) == 2);
}

void testSingleElement() {
    std::vector<int> nums = {0};
    assert(Jump(nums) == 0);
}

void testAllOnes() {
    std::vector<int> nums = {1, 1, 1, 1};
    assert(Jump(nums) == 3);
}

void testBigFirstJump() {
    std::vector<int> nums = {5, 0, 0, 0, 0};
    assert(Jump(nums) == 1);
}

void testTwoElements() {
    std::vector<int> nums = {2, 1};
    assert(Jump(nums) == 1);
}

int main() {
    testClassic();
    testSingleElement();
    testAllOnes();
    testBigFirstJump();
    testTwoElements();
    std::printf("all tests passed\n");
    return 0;
}
