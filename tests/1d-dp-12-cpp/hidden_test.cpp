#include <cassert>
#include <cstdio>
#include <vector>

bool CanPartition(std::vector<int>& nums);

void testClassic() {
    std::vector<int> nums = {1, 5, 11, 5};
    assert(CanPartition(nums) == true);
}

void testOddSum() {
    std::vector<int> nums = {1, 2, 3, 5};
    assert(CanPartition(nums) == false);
}

void testEvenSplit() {
    std::vector<int> nums = {1, 2, 3, 4};
    assert(CanPartition(nums) == true);
}

void testTwoEqual() {
    std::vector<int> nums = {2, 2};
    assert(CanPartition(nums) == true);
}

void testSingleElement() {
    std::vector<int> nums = {4};
    assert(CanPartition(nums) == false);
}

void testAllSame() {
    std::vector<int> nums = {3, 3, 3, 3};
    assert(CanPartition(nums) == true);
}

void testEvenSumUnreachable() {
    std::vector<int> nums = {2, 2, 3, 5};
    assert(CanPartition(nums) == false);
}

void testBoundaryValues() {
    std::vector<int> nums = {100, 100, 100, 100};
    assert(CanPartition(nums) == true);
}

void testLargerMultiCombination() {
    std::vector<int> nums = {1, 2, 3, 4, 5, 6, 7};
    assert(CanPartition(nums) == true);
}

int main() {
    testClassic();
    testOddSum();
    testEvenSplit();
    testTwoEqual();
    testSingleElement();
    testAllSame();
    testEvenSumUnreachable();
    testBoundaryValues();
    testLargerMultiCombination();
    std::printf("all tests passed\n");
    return 0;
}
