#include <cassert>
#include <cstdio>
#include <vector>

int MinCostClimbingStairs(std::vector<int>& cost);

void testThree() {
    std::vector<int> cost = {10, 15, 20};
    assert(MinCostClimbingStairs(cost) == 15);
}

void testTen() {
    std::vector<int> cost = {1, 100, 1, 1, 1, 100, 1, 1, 100, 1};
    assert(MinCostClimbingStairs(cost) == 6);
}

void testTwoEqual() {
    std::vector<int> cost = {0, 0};
    assert(MinCostClimbingStairs(cost) == 0);
}

void testBoundaryMaxValues() {
    std::vector<int> cost = {999, 999};
    assert(MinCostClimbingStairs(cost) == 999);
}

void testLargerAscending() {
    std::vector<int> cost = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
    assert(MinCostClimbingStairs(cost) == 25);
}

int main() {
    testThree();
    testTen();
    testTwoEqual();
    testBoundaryMaxValues();
    testLargerAscending();
    std::printf("all tests passed\n");
    return 0;
}
