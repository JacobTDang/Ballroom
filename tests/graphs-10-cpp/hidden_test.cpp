#include <cassert>
#include <cstdio>
#include <vector>

bool ValidTree(int n, std::vector<std::vector<int>>& edges);

void testValid() {
    std::vector<std::vector<int>> edges = {{0, 1}, {0, 2}, {0, 3}, {1, 4}};
    assert(ValidTree(5, edges) == true);
}

void testHasCycle() {
    std::vector<std::vector<int>> edges = {{0, 1}, {1, 2}, {2, 3}, {1, 3}, {1, 4}};
    assert(ValidTree(5, edges) == false);
}

void testDisconnected() {
    std::vector<std::vector<int>> edges = {{0, 1}, {2, 3}};
    assert(ValidTree(4, edges) == false);
}

void testSingleNode() {
    std::vector<std::vector<int>> edges = {};
    assert(ValidTree(1, edges) == true);
}

void testNoEdgesMultipleNodes() {
    std::vector<std::vector<int>> edges = {};
    assert(ValidTree(3, edges) == false);
}

int main() {
    testValid();
    testHasCycle();
    testDisconnected();
    testSingleNode();
    testNoEdgesMultipleNodes();
    std::printf("all tests passed\n");
    return 0;
}
