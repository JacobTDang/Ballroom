#include <cassert>
#include <cstdio>
#include <vector>

int RobCircular(std::vector<int>& nums);

void testClassic() {
    std::vector<int> nums = {2, 3, 2};
    assert(RobCircular(nums) == 3);
}

void testFourHouses() {
    std::vector<int> nums = {1, 2, 3, 1};
    assert(RobCircular(nums) == 4);
}

void testThreeInARow() {
    std::vector<int> nums = {1, 2, 3};
    assert(RobCircular(nums) == 3);
}

void testSingleHouse() {
    std::vector<int> nums = {5};
    assert(RobCircular(nums) == 5);
}

void testTwoHouses() {
    std::vector<int> nums = {5, 10};
    assert(RobCircular(nums) == 10);
}

void testAllZeros() {
    std::vector<int> nums = {0, 0, 0, 0};
    assert(RobCircular(nums) == 0);
}

void testLargerAlternating() {
    std::vector<int> nums = {2, 3, 2, 3, 2, 3, 2};
    assert(RobCircular(nums) == 9);
}

void testBoundaryMaxValues() {
    std::vector<int> nums = {1000, 1000, 1000};
    assert(RobCircular(nums) == 1000);
}

int main() {
    testClassic();
    testFourHouses();
    testThreeInARow();
    testSingleHouse();
    testTwoHouses();
    testAllZeros();
    testLargerAlternating();
    testBoundaryMaxValues();
    std::printf("all tests passed\n");
    return 0;
}
