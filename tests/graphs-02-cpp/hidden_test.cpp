#include <cassert>
#include <cstdio>
#include <vector>

int MaxAreaOfIsland(std::vector<std::vector<int>>& grid);

int main() {
    std::vector<std::vector<int>> grid = {
        {0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
        {0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
        {0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0},
        {0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
    };
    assert(MaxAreaOfIsland(grid) == 6);

    std::vector<std::vector<int>> empty = {{0, 0, 0, 0, 0, 0, 0, 0}};
    assert(MaxAreaOfIsland(empty) == 0);

    std::vector<std::vector<int>> single = {{1}};
    assert(MaxAreaOfIsland(single) == 1);
    printf("all assertions passed\n");
    return 0;
}
