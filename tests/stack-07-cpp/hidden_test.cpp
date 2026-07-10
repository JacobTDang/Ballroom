#include <cassert>
#include <cstdio>
#include <vector>

int LargestRectangleArea(const std::vector<int>& heights);

int main() {
    assert(LargestRectangleArea({2, 1, 5, 6, 2, 3}) == 10);
    assert(LargestRectangleArea({2, 4}) == 4);
    assert(LargestRectangleArea({1}) == 1);
    assert(LargestRectangleArea({0, 0}) == 0);
    printf("all assertions passed\n");
    return 0;
}
