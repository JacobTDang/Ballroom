#include <cassert>
#include <cstdio>

int UniquePaths(int m, int n);

void testClassic() {
    assert(UniquePaths(3, 7) == 28);
}

void testSmall() {
    assert(UniquePaths(3, 2) == 3);
}

void testSingleCell() {
    assert(UniquePaths(1, 1) == 1);
}

void testSingleRow() {
    assert(UniquePaths(1, 5) == 1);
}

void testSingleColumn() {
    assert(UniquePaths(5, 1) == 1);
}

void testSquare() {
    assert(UniquePaths(3, 3) == 6);
}

void testLargerSquare() {
    assert(UniquePaths(10, 10) == 48620);
}

void testBoundaryMaxWithMinOther() {
    assert(UniquePaths(2, 100) == 100);
}

int main() {
    testClassic();
    testSmall();
    testSingleCell();
    testSingleRow();
    testSingleColumn();
    testSquare();
    testLargerSquare();
    testBoundaryMaxWithMinOther();
    std::printf("all tests passed\n");
    return 0;
}
