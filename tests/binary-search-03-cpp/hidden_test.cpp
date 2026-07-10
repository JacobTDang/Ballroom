#include <cassert>
#include <cstdio>
#include <vector>

int MinEatingSpeed(const std::vector<int>& piles, int h);

int main() {
    assert(MinEatingSpeed({3, 6, 7, 11}, 8) == 4);
    assert(MinEatingSpeed({30, 11, 23, 4, 20}, 5) == 30);
    assert(MinEatingSpeed({30, 11, 23, 4, 20}, 6) == 23);
    assert(MinEatingSpeed({1000000000}, 2) == 500000000);
    assert(MinEatingSpeed({1}, 1) == 1);
    printf("all assertions passed\n");
    return 0;
}
