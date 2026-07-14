#include <cassert>
#include <cstdio>
#include <vector>

int MaxProfit(std::vector<int>& prices);

void testClassic() {
    std::vector<int> prices = {1, 2, 3, 0, 2};
    assert(MaxProfit(prices) == 3);
}

void testSingleDay() {
    std::vector<int> prices = {1};
    assert(MaxProfit(prices) == 0);
}

void testMonotonicIncreasing() {
    std::vector<int> prices = {1, 2, 4};
    assert(MaxProfit(prices) == 3);
}

void testEmpty() {
    std::vector<int> prices = {};
    assert(MaxProfit(prices) == 0);
}

void testMonotonicDecreasing() {
    std::vector<int> prices = {5, 4, 3, 2, 1};
    assert(MaxProfit(prices) == 0);
}

void testTwoDaysProfit() {
    std::vector<int> prices = {1, 2};
    assert(MaxProfit(prices) == 1);
}

void testCooldownForcesWait() {
    std::vector<int> prices = {1, 4, 2, 7};
    assert(MaxProfit(prices) == 6);
}

void testLargerMultiTrade() {
    std::vector<int> prices = {6, 1, 3, 2, 4, 7};
    assert(MaxProfit(prices) == 6);
}

void testBoundaryValues() {
    std::vector<int> prices = {10000, 1};
    assert(MaxProfit(prices) == 0);
}

int main() {
    testClassic();
    testSingleDay();
    testMonotonicIncreasing();
    testEmpty();
    testMonotonicDecreasing();
    testTwoDaysProfit();
    testCooldownForcesWait();
    testLargerMultiTrade();
    testBoundaryValues();
    std::printf("all tests passed\n");
    return 0;
}
