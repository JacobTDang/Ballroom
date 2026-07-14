#include <cassert>
#include <cstdio>
#include <vector>

int OrangesRotting(std::vector<std::vector<int>>& grid);

int main() {
    std::vector<std::vector<int>> a = {{2, 1, 1}, {1, 1, 0}, {0, 1, 1}};
    assert(OrangesRotting(a) == 4);
    std::vector<std::vector<int>> b = {{2, 1, 1}, {0, 1, 1}, {1, 0, 1}};
    assert(OrangesRotting(b) == -1);
    std::vector<std::vector<int>> c = {{0, 2}};
    assert(OrangesRotting(c) == 0);
    std::vector<std::vector<int>> d = {{0}};
    assert(OrangesRotting(d) == 0);
    std::vector<std::vector<int>> e = {{2, 1, 1}, {1, 1, 1}, {1, 1, 2}};
    assert(OrangesRotting(e) == 2);
    std::vector<std::vector<int>> f = {{2, 2}, {2, 2}};
    assert(OrangesRotting(f) == 0);
    printf("all assertions passed\n");
    return 0;
}
