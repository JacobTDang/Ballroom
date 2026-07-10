#include <cassert>
#include <cstdio>
#include <vector>

int MinCostConnectPoints(std::vector<std::vector<int>>& points);

void testClassic() {
    std::vector<std::vector<int>> points = {{0, 0}, {2, 2}, {3, 10}, {5, 2}, {7, 0}};
    assert(MinCostConnectPoints(points) == 20);
}

void testThreePoints() {
    std::vector<std::vector<int>> points = {{3, 12}, {-2, 5}, {-4, 1}};
    assert(MinCostConnectPoints(points) == 18);
}

void testSinglePoint() {
    std::vector<std::vector<int>> points = {{0, 0}};
    assert(MinCostConnectPoints(points) == 0);
}

int main() {
    testClassic();
    testThreePoints();
    testSinglePoint();
    std::printf("all tests passed\n");
    return 0;
}
