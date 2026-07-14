#include <cassert>
#include <cstdio>
#include <vector>

int Change(int amount, std::vector<int>& coins);

void testClassic() {
    std::vector<int> coins = {1, 2, 5};
    assert(Change(5, coins) == 4);
}

void testNoWay() {
    std::vector<int> coins = {2};
    assert(Change(3, coins) == 0);
}

void testZeroAmount() {
    std::vector<int> coins = {1, 2, 3};
    assert(Change(0, coins) == 1);
}

void testExactSingleCoin() {
    std::vector<int> coins = {10};
    assert(Change(10, coins) == 1);
}

void testLargerAmount() {
    std::vector<int> coins = {1, 2, 5};
    assert(Change(10, coins) == 10);
}

void testSingleCoinNoDivide() {
    std::vector<int> coins = {3};
    assert(Change(7, coins) == 0);
}

void testMoreDenominations() {
    std::vector<int> coins = {2, 5, 3, 6};
    assert(Change(10, coins) == 5);
}

void testBoundaryAmountSingleCoin() {
    std::vector<int> coins = {1};
    assert(Change(500, coins) == 1);
}

int main() {
    testClassic();
    testNoWay();
    testZeroAmount();
    testExactSingleCoin();
    testLargerAmount();
    testSingleCoinNoDivide();
    testMoreDenominations();
    testBoundaryAmountSingleCoin();
    std::printf("all tests passed\n");
    return 0;
}
