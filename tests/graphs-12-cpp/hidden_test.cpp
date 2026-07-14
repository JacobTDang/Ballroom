#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> FindRedundantConnection(std::vector<std::vector<int>>& edges);

void testTriangle() {
    std::vector<std::vector<int>> edges = {{1, 2}, {1, 3}, {2, 3}};
    auto got = FindRedundantConnection(edges);
    std::vector<int> want = {2, 3};
    assert(got == want);
}

void testLaterCycle() {
    std::vector<std::vector<int>> edges = {{1, 2}, {2, 3}, {3, 4}, {1, 4}, {1, 5}};
    auto got = FindRedundantConnection(edges);
    std::vector<int> want = {1, 4};
    assert(got == want);
}

void testMergingComponents() {
    std::vector<std::vector<int>> edges = {{1, 4}, {3, 4}, {1, 3}, {1, 2}, {4, 5}};
    auto got = FindRedundantConnection(edges);
    std::vector<int> want = {1, 3};
    assert(got == want);
}

int main() {
    testTriangle();
    testLaterCycle();
    testMergingComponents();
    std::printf("all tests passed\n");
    return 0;
}
