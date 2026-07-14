#include <cassert>
#include <cstdio>
#include <vector>

int MaxProfit(const std::vector<int>& prices);

int main() {
    assert(MaxProfit({7, 1, 5, 3, 6, 4}) == 5);
    assert(MaxProfit({7, 6, 4, 3, 1}) == 0);
    assert(MaxProfit({2, 4, 1}) == 2);
    assert(MaxProfit({1}) == 0);
    assert(MaxProfit({3, 3, 3, 3}) == 0);
    assert(MaxProfit({1, 2, 4, 2, 5, 7, 2, 4, 9, 0}) == 8);
    assert(MaxProfit({}) == 0);
    assert(MaxProfit({3, 1, 4, 1, 5, 9, 2, 6}) == 8);
    assert(MaxProfit({2, 1, 2, 1, 0, 1, 2}) == 2);
    printf("all assertions passed\n");
    return 0;
}
