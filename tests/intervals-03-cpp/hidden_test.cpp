#include <cassert>
#include <cstdio>
#include <vector>

int EraseOverlapIntervals(std::vector<std::vector<int>>& intervals);

void check(std::vector<std::vector<int>> intervals, int want) {
    assert(EraseOverlapIntervals(intervals) == want);
}

int main() {
    check({{1, 2}, {2, 3}, {3, 4}, {1, 3}}, 1);
    check({{1, 2}, {1, 2}, {1, 2}}, 2);
    check({{1, 2}, {2, 3}}, 0);
    check({{1, 2}}, 0);
    check({{1, 2}, {3, 4}, {5, 6}}, 0);
    check({{1, 100}, {11, 22}, {1, 11}, {2, 12}}, 2);
    check({{-50000, -49999}, {-49999, 50000}}, 0);
    check({{1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6}}, 4);
    std::printf("all assertions passed\n");
    return 0;
}
