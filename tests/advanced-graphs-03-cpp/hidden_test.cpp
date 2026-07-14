#include <cassert>
#include <cstdio>
#include <vector>

int NetworkDelayTime(std::vector<std::vector<int>>& times, int n, int k);

void testClassic() {
    std::vector<std::vector<int>> times = {{2, 1, 1}, {2, 3, 1}, {3, 4, 1}};
    assert(NetworkDelayTime(times, 4, 2) == 2);
}

void testSingleEdgeReachable() {
    std::vector<std::vector<int>> times = {{1, 2, 1}};
    assert(NetworkDelayTime(times, 2, 1) == 1);
}

void testUnreachable() {
    std::vector<std::vector<int>> times = {{1, 2, 1}};
    assert(NetworkDelayTime(times, 2, 2) == -1);
}

void testShortestOfMultiplePaths() {
    std::vector<std::vector<int>> times = {{1, 2, 1}, {2, 3, 2}, {1, 3, 4}};
    assert(NetworkDelayTime(times, 3, 1) == 3);
}

void testSingleNodeNoEdges() {
    std::vector<std::vector<int>> times = {};
    assert(NetworkDelayTime(times, 1, 1) == 0);
}

int main() {
    testClassic();
    testSingleEdgeReachable();
    testUnreachable();
    testShortestOfMultiplePaths();
    testSingleNodeNoEdges();
    std::printf("all tests passed\n");
    return 0;
}
