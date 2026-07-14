#include <cassert>
#include <cstdio>
#include <vector>

bool CanAttendMeetings(std::vector<std::vector<int>>& intervals);

void check(std::vector<std::vector<int>> intervals, bool want) {
    assert(CanAttendMeetings(intervals) == want);
}

int main() {
    check({{0, 30}, {5, 10}, {15, 20}}, false);
    check({{7, 10}, {2, 4}}, true);
    check({{5, 10}, {10, 15}}, true);
    check({}, true);
    check({{3, 8}}, true);
    check({{13, 15}, {1, 5}, {6, 8}, {14, 20}}, false);
    check({{0, 1000000}, {999999, 1000000}}, false);
    check({{0, 10}, {10, 20}, {20, 30}, {30, 40}, {40, 50}, {50, 60}, {60, 70}}, true);
    std::printf("all assertions passed\n");
    return 0;
}
