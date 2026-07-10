#include <cassert>
#include <cstdio>
#include <vector>

int MaxArea(const std::vector<int>& height);

int main() {
    assert(MaxArea({1, 8, 6, 2, 5, 4, 8, 3, 7}) == 49);
    assert(MaxArea({1, 1}) == 1);
    assert(MaxArea({4, 3, 2, 1, 4}) == 16);
    assert(MaxArea({1, 2, 1}) == 2);
    assert(MaxArea({1, 2, 4, 3}) == 4);
    assert(MaxArea({1, 3, 2, 5, 25, 24, 5}) == 24);
    printf("all assertions passed\n");
    return 0;
}
