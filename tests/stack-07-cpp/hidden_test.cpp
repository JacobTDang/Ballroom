#include <cassert>
#include <cstdio>
#include <vector>

int LargestRectangleArea(const std::vector<int>& heights);

int main() {
    assert(LargestRectangleArea({2, 1, 5, 6, 2, 3}) == 10);
    assert(LargestRectangleArea({2, 4}) == 4);
    assert(LargestRectangleArea({1}) == 1);
    assert(LargestRectangleArea({0, 0}) == 0);
    assert(LargestRectangleArea({5, 5, 5, 5}) == 20);
    assert(LargestRectangleArea({5, 4, 3, 2, 1}) == 9);
    assert(LargestRectangleArea({1, 2, 3, 4, 5}) == 9);
    printf("all assertions passed\n");
    return 0;
}
