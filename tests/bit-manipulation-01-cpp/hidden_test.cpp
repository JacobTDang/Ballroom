#include <cassert>
#include <cstdio>
#include <vector>

int SingleNumber(std::vector<int>& nums);

void testClassic() {
    std::vector<int> nums = {2, 2, 1};
    assert(SingleNumber(nums) == 1);
}

void testLongerMix() {
    std::vector<int> nums = {4, 1, 2, 1, 2};
    assert(SingleNumber(nums) == 4);
}

void testSingleElement() {
    std::vector<int> nums = {7};
    assert(SingleNumber(nums) == 7);
}

void testNegativeNumbers() {
    std::vector<int> nums = {-1, -1, -2};
    assert(SingleNumber(nums) == -2);
}

int main() {
    testClassic();
    testLongerMix();
    testSingleElement();
    testNegativeNumbers();
    std::printf("all tests passed\n");
    return 0;
}
