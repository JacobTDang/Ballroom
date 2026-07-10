#include <algorithm>
#include <cassert>
#include <cstdio>
#include <vector>

std::vector<std::vector<int>> KClosest(std::vector<std::vector<int>>& points, int k);

std::vector<std::vector<int>> normalize(std::vector<std::vector<int>> points) {
    std::sort(points.begin(), points.end());
    return points;
}

void check(std::vector<std::vector<int>> points, int k, std::vector<std::vector<int>> want) {
    auto got = normalize(KClosest(points, k));
    assert(got == normalize(want));
}

int main() {
    check({{1, 3}, {-2, 2}}, 1, {{-2, 2}});
    check({{3, 3}, {5, -1}, {-2, 4}}, 2, {{3, 3}, {-2, 4}});
    check({{0, 1}, {1, 0}}, 2, {{0, 1}, {1, 0}});
    printf("all assertions passed\n");
    return 0;
}
