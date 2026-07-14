#include <cassert>
#include <cstdio>
#include <vector>

std::vector<std::vector<int>> Merge(std::vector<std::vector<int>>& intervals);

void check(std::vector<std::vector<int>> intervals, std::vector<std::vector<int>> want) {
    assert(Merge(intervals) == want);
}

int main() {
    check({{1, 3}, {2, 6}, {8, 10}, {15, 18}}, {{1, 6}, {8, 10}, {15, 18}});
    check({{1, 4}, {4, 5}}, {{1, 5}});
    check({{15, 18}, {2, 6}, {1, 3}, {8, 10}}, {{1, 6}, {8, 10}, {15, 18}});
    check({{1, 4}}, {{1, 4}});
    check({{1, 10}, {2, 3}, {4, 5}}, {{1, 10}});
    check({{1, 2}, {3, 4}, {5, 6}}, {{1, 2}, {3, 4}, {5, 6}});
    check({{0, 1}, {9999, 10000}}, {{0, 1}, {9999, 10000}});
    check({{1, 3}, {2, 4}, {5, 7}, {6, 8}, {10, 12}, {15, 20}, {18, 25}, {30, 31}}, {{1, 4}, {5, 8}, {10, 12}, {15, 25}, {30, 31}});
    std::printf("all assertions passed\n");
    return 0;
}
