#include <cassert>
#include <cstdio>
#include <vector>

int CarFleet(int target, const std::vector<int>& position, const std::vector<int>& speed);

int main() {
    assert(CarFleet(12, {10, 8, 0, 5, 3}, {2, 4, 1, 1, 3}) == 3);
    assert(CarFleet(10, {3}, {3}) == 1);
    assert(CarFleet(100, {0, 2, 4}, {4, 2, 1}) == 1);
    assert(CarFleet(10, {0, 4, 8}, {1, 1, 1}) == 3);
    assert(CarFleet(10, {0, 3, 6}, {5, 5, 5}) == 3);
    assert(CarFleet(20, {1, 4}, {2, 1}) == 1);
    assert(CarFleet(5, {5}, {1}) == 1);
    printf("all assertions passed\n");
    return 0;
}
