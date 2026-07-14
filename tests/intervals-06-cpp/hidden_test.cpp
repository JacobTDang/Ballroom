#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> MinInterval(std::vector<std::vector<int>>& intervals, std::vector<int>& queries);

void check(std::vector<std::vector<int>> intervals, std::vector<int> queries, std::vector<int> want) {
    assert(MinInterval(intervals, queries) == want);
}

int main() {
    check({{1, 4}, {2, 4}, {3, 6}, {4, 4}}, {2, 3, 4, 5}, {3, 3, 1, 4});
    check({{2, 3}, {2, 5}, {1, 8}, {20, 25}}, {2, 19, 5, 22}, {2, -1, 4, 6});
    check({{1, 10}}, {1, 10, 11}, {10, 10, -1});
    check({{5, 5}}, {5}, {1});
    check({{1, 10000000}}, {1, 10000000}, {10000000, 10000000});
    check({{1, 100}, {10, 20}, {15, 16}, {50, 60}}, {15, 55, 99, 5}, {2, 11, 100, 100});
    std::printf("all assertions passed\n");
    return 0;
}
