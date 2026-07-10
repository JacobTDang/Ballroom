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

int main() {
    testClassic();
    testSingleEdgeReachable();
    testUnreachable();
    std::printf("all tests passed\n");
    return 0;
}
