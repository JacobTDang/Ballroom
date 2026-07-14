#include <cassert>
#include <cstdio>
#include <vector>

int Trap(const std::vector<int>& height);

int main() {
    assert(Trap({0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}) == 6);
    assert(Trap({4, 2, 0, 3, 2, 5}) == 9);
    assert(Trap({}) == 0);
    assert(Trap({1, 2, 3, 4, 5}) == 0);
    assert(Trap({5, 4, 3, 2, 1}) == 0);
    assert(Trap({3, 0, 3}) == 3);
    assert(Trap({2, 0, 2}) == 2);
    assert(Trap({5}) == 0);
    assert(Trap({1, 0, 1}) == 1);
    assert(Trap({4, 4, 4, 4}) == 0);
    assert(Trap({5, 2, 1, 2, 1, 5}) == 14);
    printf("all assertions passed\n");
    return 0;
}
