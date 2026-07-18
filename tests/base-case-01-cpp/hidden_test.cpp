#include <cassert>
#include <cstdio>
#include <vector>

int count_paths(const std::vector<std::vector<int>>& grid);

int main() {
    assert(count_paths({{0, 0}, {0, 0}}) == 2);
    assert(count_paths({{0}}) == 1);
    assert(count_paths({{0, 0}, {0, 1}}) == 0);
    // Anti-overfit: the destination cell itself is open, but every
    // route to it is blocked. A fix that just hardcodes the
    // destination base case to 1 without preserving the blocked-cell
    // check must still get 0 here.
    assert(count_paths({{0, 1}, {1, 0}}) == 0);
    assert(count_paths({{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}) == 2);
    assert(count_paths({{0, 0, 0}}) == 1);
    printf("all assertions passed\n");
    return 0;
}
