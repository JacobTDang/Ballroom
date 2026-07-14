#include <cassert>
#include <cstdio>
#include <vector>

int FindTargetSumWays(std::vector<int>& nums, int target);

void testClassic() {
    std::vector<int> nums = {1, 1, 1, 1, 1};
    assert(FindTargetSumWays(nums, 3) == 5);
}

void testSingle() {
    std::vector<int> nums = {1};
    assert(FindTargetSumWays(nums, 1) == 1);
}

void testUnreachable() {
    std::vector<int> nums = {1, 2, 3};
    assert(FindTargetSumWays(nums, 100) == 0);
}

void testZeros() {
    std::vector<int> nums = {0, 0, 0, 0, 0, 0, 0, 0, 1};
    assert(FindTargetSumWays(nums, 1) == 256);
}

void testZeroTarget() {
    std::vector<int> nums = {1, 1};
    assert(FindTargetSumWays(nums, 0) == 2);
}

void testNegativeTarget() {
    std::vector<int> nums = {1, 1, 1, 1, 1};
    assert(FindTargetSumWays(nums, -3) == 5);
}

void testMixedZeroTarget() {
    std::vector<int> nums = {1, 2, 1};
    assert(FindTargetSumWays(nums, 0) == 2);
}

void testParityImpossible() {
    std::vector<int> nums = {1, 2, 3};
    assert(FindTargetSumWays(nums, 1) == 0);
}

int main() {
    testClassic();
    testSingle();
    testUnreachable();
    testZeros();
    testZeroTarget();
    testNegativeTarget();
    testMixedZeroTarget();
    testParityImpossible();
    std::printf("all tests passed\n");
    return 0;
}
