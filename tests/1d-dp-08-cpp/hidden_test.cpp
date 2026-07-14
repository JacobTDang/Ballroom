#include <cassert>
#include <cstdio>
#include <vector>

int CoinChange(std::vector<int>& coins, int amount);

void testClassic() {
    std::vector<int> coins = {1, 2, 5};
    assert(CoinChange(coins, 11) == 3);
}

void testImpossible() {
    std::vector<int> coins = {2};
    assert(CoinChange(coins, 3) == -1);
}

void testZeroAmount() {
    std::vector<int> coins = {1};
    assert(CoinChange(coins, 0) == 0);
}

void testSingleCoinExact() {
    std::vector<int> coins = {3, 7};
    assert(CoinChange(coins, 6) == 2);
}

void testLargeAmountOnlyOnes() {
    std::vector<int> coins = {1};
    assert(CoinChange(coins, 10000) == 10000);
}

void testUnreachableAmount() {
    std::vector<int> coins = {3, 5};
    assert(CoinChange(coins, 7) == -1);
}

void testMixedDenominations() {
    std::vector<int> coins = {1, 5, 10, 25};
    assert(CoinChange(coins, 63) == 6);
}

int main() {
    testClassic();
    testImpossible();
    testZeroAmount();
    testSingleCoinExact();
    testLargeAmountOnlyOnes();
    testUnreachableAmount();
    testMixedDenominations();
    std::printf("all tests passed\n");
    return 0;
}
