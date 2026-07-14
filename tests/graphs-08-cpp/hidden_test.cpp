#include <cassert>
#include <cstdio>
#include <vector>

bool CanFinish(int numCourses, std::vector<std::vector<int>>& prerequisites);

void testNoCycle() {
    std::vector<std::vector<int>> prereqs = {{1, 0}};
    assert(CanFinish(2, prereqs) == true);
}

void testCycle() {
    std::vector<std::vector<int>> prereqs = {{1, 0}, {0, 1}};
    assert(CanFinish(2, prereqs) == false);
}

void testNoPrerequisites() {
    std::vector<std::vector<int>> prereqs = {};
    assert(CanFinish(5, prereqs) == true);
}

void testLongerCycle() {
    std::vector<std::vector<int>> prereqs = {{1, 0}, {2, 1}, {3, 2}, {0, 3}};
    assert(CanFinish(4, prereqs) == false);
}

void testDiamondDAGNoCycle() {
    std::vector<std::vector<int>> prereqs = {{1, 0}, {2, 0}, {3, 1}, {3, 2}};
    assert(CanFinish(4, prereqs) == true);
}

int main() {
    testNoCycle();
    testCycle();
    testNoPrerequisites();
    testLongerCycle();
    testDiamondDAGNoCycle();
    std::printf("all tests passed\n");
    return 0;
}
