#include <cassert>
#include <cstdio>
#include <vector>

int Rob(std::vector<int>& nums);

void testClassic() {
    std::vector<int> nums = {1, 2, 3, 1};
    assert(Rob(nums) == 4);
}

void testLarger() {
    std::vector<int> nums = {2, 7, 9, 3, 1};
    assert(Rob(nums) == 12);
}

void testSingleHouse() {
    std::vector<int> nums = {5};
    assert(Rob(nums) == 5);
}

void testTwoHouses() {
    std::vector<int> nums = {2, 1};
    assert(Rob(nums) == 2);
}

void testAllZeros() {
    std::vector<int> nums = {0, 0, 0, 0};
    assert(Rob(nums) == 0);
}

void testLargerMixedValues() {
    std::vector<int> nums = {5, 5, 10, 100, 10, 5};
    assert(Rob(nums) == 110);
}

void testBoundaryMaxValues() {
    std::vector<int> nums = {1000, 1000};
    assert(Rob(nums) == 1000);
}

int main() {
    testClassic();
    testLarger();
    testSingleHouse();
    testTwoHouses();
    testAllZeros();
    testLargerMixedValues();
    testBoundaryMaxValues();
    std::printf("all tests passed\n");
    return 0;
}
