#include <cassert>
#include <cstdio>
#include <vector>

bool SearchMatrix(const std::vector<std::vector<int>>& matrix, int target);

int main() {
    std::vector<std::vector<int>> m = {{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}};
    assert(SearchMatrix(m, 3) == true);
    assert(SearchMatrix(m, 13) == false);
    assert(SearchMatrix({{1}}, 1) == true);
    assert(SearchMatrix({{1, 3}}, 3) == true);
    assert(SearchMatrix(m, 60) == true);
    assert(SearchMatrix(m, 0) == false);
    assert(SearchMatrix(m, 1) == true);
    assert(SearchMatrix(m, 23) == true);
    printf("all assertions passed\n");
    return 0;
}
