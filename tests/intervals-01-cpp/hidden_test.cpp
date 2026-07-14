#include <cassert>
#include <cstdio>
#include <vector>

std::vector<std::vector<int>> Insert(std::vector<std::vector<int>>& intervals, std::vector<int>& newInterval);

void check(std::vector<std::vector<int>> intervals, std::vector<int> newInterval, std::vector<std::vector<int>> want) {
    assert(Insert(intervals, newInterval) == want);
}

int main() {
    check({{1, 3}, {6, 9}}, {2, 5}, {{1, 5}, {6, 9}});
    check({{1, 2}, {3, 5}, {6, 7}, {8, 10}, {12, 16}}, {4, 8}, {{1, 2}, {3, 10}, {12, 16}});
    check({}, {5, 7}, {{5, 7}});
    check({{1, 5}}, {6, 8}, {{1, 5}, {6, 8}});
    check({{1, 5}}, {0, 0}, {{0, 0}, {1, 5}});
    check({{2, 3}, {4, 5}, {6, 7}}, {1, 10}, {{1, 10}});
    check({{1, 2}}, {2, 3}, {{1, 3}});
    check({{1, 2}}, {3, 4}, {{1, 2}, {3, 4}});
    check({{3, 5}}, {3, 5}, {{3, 5}});
    check({{1, 10}}, {3, 5}, {{1, 10}});
    check({{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}, {11, 12}, {13, 14}, {15, 16}, {17, 18}, {19, 20}}, {6, 15}, {{1, 2}, {3, 4}, {5, 16}, {17, 18}, {19, 20}});
    check({{50000, 60000}}, {0, 100000}, {{0, 100000}});
    std::printf("all assertions passed\n");
    return 0;
}
