#include <algorithm>
#include <cassert>
#include <cstdio>
#include <vector>

std::vector<std::vector<int>> PacificAtlantic(std::vector<std::vector<int>>& heights);

std::vector<std::vector<int>> normalize(std::vector<std::vector<int>> lists) {
    std::sort(lists.begin(), lists.end());
    return lists;
}

void testClassic() {
    std::vector<std::vector<int>> heights = {
        {1, 2, 2, 3, 5},
        {3, 2, 3, 4, 4},
        {2, 4, 5, 3, 1},
        {6, 7, 1, 4, 5},
        {5, 1, 1, 2, 4},
    };
    std::vector<std::vector<int>> want = {{0, 4}, {1, 3}, {1, 4}, {2, 2}, {3, 0}, {3, 1}, {4, 0}};

    auto got = normalize(PacificAtlantic(heights));
    auto wantNorm = normalize(want);
    assert(got == wantNorm);
}

void testSingleCell() {
    std::vector<std::vector<int>> heights = {{1}};
    auto got = normalize(PacificAtlantic(heights));
    std::vector<std::vector<int>> want = {{0, 0}};
    assert(got == want);
}

int main() {
    testClassic();
    testSingleCell();
    std::printf("all tests passed\n");
    return 0;
}
