#include <cassert>
#include <cstdio>
#include <vector>

int FindCheapestPrice(int n, std::vector<std::vector<int>>& flights, int src, int dst, int k);

void testOneStop() {
    std::vector<std::vector<int>> flights = {{0, 1, 100}, {1, 2, 100}, {2, 0, 100}, {1, 3, 600}, {2, 3, 200}};
    assert(FindCheapestPrice(4, flights, 0, 3, 1) == 700);
}

void testCheaperViaStop() {
    std::vector<std::vector<int>> flights = {{0, 1, 100}, {1, 2, 100}, {0, 2, 500}};
    assert(FindCheapestPrice(3, flights, 0, 2, 1) == 200);
}

void testNoStopsAllowed() {
    std::vector<std::vector<int>> flights = {{0, 1, 100}, {1, 2, 100}, {0, 2, 500}};
    assert(FindCheapestPrice(3, flights, 0, 2, 0) == 500);
}

int main() {
    testOneStop();
    testCheaperViaStop();
    testNoStopsAllowed();
    std::printf("all tests passed\n");
    return 0;
}
