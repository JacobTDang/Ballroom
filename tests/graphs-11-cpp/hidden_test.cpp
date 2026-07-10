#include <cassert>
#include <cstdio>
#include <vector>

int CountComponents(int n, std::vector<std::vector<int>>& edges);

void testClassic() {
    std::vector<std::vector<int>> edges = {{0, 1}, {1, 2}, {3, 4}};
    assert(CountComponents(5, edges) == 2);
}

void testAllConnected() {
    std::vector<std::vector<int>> edges = {{0, 1}, {1, 2}, {2, 3}};
    assert(CountComponents(4, edges) == 1);
}

void testNoEdges() {
    std::vector<std::vector<int>> edges = {};
    assert(CountComponents(4, edges) == 4);
}

void testSingleNode() {
    std::vector<std::vector<int>> edges = {};
    assert(CountComponents(1, edges) == 1);
}

int main() {
    testClassic();
    testAllConnected();
    testNoEdges();
    testSingleNode();
    std::printf("all tests passed\n");
    return 0;
}
