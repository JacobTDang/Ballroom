#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> product_except_self(const std::vector<int>& nums);

int main() {
    assert((product_except_self({1, 2, 3, 4}) == std::vector<int>{24, 12, 8, 6}));
    assert((product_except_self({-1, 1, 0, -3, 3}) == std::vector<int>{0, 0, 9, 0, 0}));
    assert((product_except_self({2, 3}) == std::vector<int>{3, 2}));
    assert((product_except_self({5, 0, 0, 4}) == std::vector<int>{0, 0, 0, 0}));
    assert((product_except_self({1, 1, 1, 1}) == std::vector<int>{1, 1, 1, 1}));
    printf("all assertions passed\n");
    return 0;
}
