#include <cassert>
#include <cstdio>
#include <vector>

int MaxCoins(std::vector<int>& nums);

void testClassic() {
    std::vector<int> nums = {3, 1, 5, 8};
    assert(MaxCoins(nums) == 167);
}

void testTwoBalloons() {
    std::vector<int> nums = {1, 5};
    assert(MaxCoins(nums) == 10);
}

void testSingleBalloon() {
    std::vector<int> nums = {7};
    assert(MaxCoins(nums) == 7);
}

void testOnes() {
    std::vector<int> nums = {1, 1};
    assert(MaxCoins(nums) == 2);
}

void testAllOnesLarger() {
    std::vector<int> nums = {1, 1, 1, 1};
    assert(MaxCoins(nums) == 4);
}

void testZeroValueBalloon() {
    std::vector<int> nums = {3, 0, 5};
    assert(MaxCoins(nums) == 20);
}

void testLargerAscending() {
    std::vector<int> nums = {1, 2, 3, 4, 5};
    assert(MaxCoins(nums) == 110);
}

void testDescending() {
    std::vector<int> nums = {5, 3, 1};
    assert(MaxCoins(nums) == 25);
}

int main() {
    testClassic();
    testTwoBalloons();
    testSingleBalloon();
    testOnes();
    testAllOnesLarger();
    testZeroValueBalloon();
    testLargerAscending();
    testDescending();
    std::printf("all tests passed\n");
    return 0;
}
