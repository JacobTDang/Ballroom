#include <cassert>
#include <cstdio>
#include <vector>

void RotateImage(std::vector<std::vector<int>>& matrix);

void test3x3() {
    std::vector<std::vector<int>> matrix = {
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    };
    std::vector<std::vector<int>> want = {
        {7, 4, 1},
        {8, 5, 2},
        {9, 6, 3},
    };
    RotateImage(matrix);
    assert(matrix == want);
}

void test2x2() {
    std::vector<std::vector<int>> matrix = {
        {1, 2},
        {3, 4},
    };
    std::vector<std::vector<int>> want = {
        {3, 1},
        {4, 2},
    };
    RotateImage(matrix);
    assert(matrix == want);
}

void test1x1() {
    std::vector<std::vector<int>> matrix = {{5}};
    std::vector<std::vector<int>> want = {{5}};
    RotateImage(matrix);
    assert(matrix == want);
}

void test4x4() {
    std::vector<std::vector<int>> matrix = {
        {1, 2, 3, 4},
        {5, 6, 7, 8},
        {9, 10, 11, 12},
        {13, 14, 15, 16},
    };
    std::vector<std::vector<int>> want = {
        {13, 9, 5, 1},
        {14, 10, 6, 2},
        {15, 11, 7, 3},
        {16, 12, 8, 4},
    };
    RotateImage(matrix);
    assert(matrix == want);
}

void testNegativeValues() {
    std::vector<std::vector<int>> matrix = {
        {-1, -2},
        {-3, -4},
    };
    std::vector<std::vector<int>> want = {
        {-3, -1},
        {-4, -2},
    };
    RotateImage(matrix);
    assert(matrix == want);
}

void testAllSameValues() {
    std::vector<std::vector<int>> matrix = {
        {7, 7, 7},
        {7, 7, 7},
        {7, 7, 7},
    };
    std::vector<std::vector<int>> want = {
        {7, 7, 7},
        {7, 7, 7},
        {7, 7, 7},
    };
    RotateImage(matrix);
    assert(matrix == want);
}

void testWithZero() {
    std::vector<std::vector<int>> matrix = {
        {0, 1},
        {2, 3},
    };
    std::vector<std::vector<int>> want = {
        {2, 0},
        {3, 1},
    };
    RotateImage(matrix);
    assert(matrix == want);
}

void test5x5() {
    std::vector<std::vector<int>> matrix = {
        {1, 2, 3, 4, 5},
        {6, 7, 8, 9, 10},
        {11, 12, 13, 14, 15},
        {16, 17, 18, 19, 20},
        {21, 22, 23, 24, 25},
    };
    std::vector<std::vector<int>> want = {
        {21, 16, 11, 6, 1},
        {22, 17, 12, 7, 2},
        {23, 18, 13, 8, 3},
        {24, 19, 14, 9, 4},
        {25, 20, 15, 10, 5},
    };
    RotateImage(matrix);
    assert(matrix == want);
}

int main() {
    test3x3();
    test2x2();
    test1x1();
    test4x4();
    testNegativeValues();
    testAllSameValues();
    testWithZero();
    test5x5();
    std::printf("all tests passed\n");
    return 0;
}
