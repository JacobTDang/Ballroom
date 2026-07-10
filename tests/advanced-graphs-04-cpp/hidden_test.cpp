#include <cassert>
#include <cstdio>
#include <vector>

int SwimInWater(std::vector<std::vector<int>>& grid);

void testSmall() {
    std::vector<std::vector<int>> grid = {{0, 2}, {1, 3}};
    assert(SwimInWater(grid) == 3);
}

void testLarger() {
    std::vector<std::vector<int>> grid = {
        {0, 1, 2, 3, 4},
        {24, 23, 22, 21, 5},
        {12, 13, 14, 15, 16},
        {11, 17, 18, 19, 20},
        {10, 9, 8, 7, 6}
    };
    assert(SwimInWater(grid) == 16);
}

void testSingleCell() {
    std::vector<std::vector<int>> grid = {{0}};
    assert(SwimInWater(grid) == 0);
}

int main() {
    testSmall();
    testLarger();
    testSingleCell();
    std::printf("all tests passed\n");
    return 0;
}
