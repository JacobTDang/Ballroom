#include <cassert>
#include <cstdio>
#include <vector>

int MinMeetingRooms(std::vector<std::vector<int>>& intervals);

void check(std::vector<std::vector<int>> intervals, int want) {
    assert(MinMeetingRooms(intervals) == want);
}

int main() {
    check({{0, 30}, {5, 10}, {15, 20}}, 2);
    check({{7, 10}, {2, 4}}, 1);
    check({{5, 10}, {10, 15}}, 1);
    check({}, 0);
    check({{1, 2}, {1, 2}, {1, 2}}, 3);
    check({{1, 5}}, 1);
    check({{1, 100}, {1, 100}, {1, 100}, {1, 100}, {1, 100}}, 5);
    check({{1, 10}, {2, 7}, {3, 19}, {8, 12}, {10, 20}, {11, 30}}, 4);
    std::printf("all assertions passed\n");
    return 0;
}
