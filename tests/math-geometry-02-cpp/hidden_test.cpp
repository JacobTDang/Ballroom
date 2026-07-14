#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> SpiralOrder(std::vector<std::vector<int>>& matrix);

void test3x3() {
    std::vector<std::vector<int>> matrix = {
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    };
    std::vector<int> want = {1, 2, 3, 6, 9, 8, 7, 4, 5};
    assert(SpiralOrder(matrix) == want);
}

void test3x4() {
    std::vector<std::vector<int>> matrix = {
        {1, 2, 3, 4},
        {5, 6, 7, 8},
        {9, 10, 11, 12},
    };
    std::vector<int> want = {1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7};
    assert(SpiralOrder(matrix) == want);
}

void testSingleRow() {
    std::vector<std::vector<int>> matrix = {{1, 2, 3, 4}};
    std::vector<int> want = {1, 2, 3, 4};
    assert(SpiralOrder(matrix) == want);
}

void testSingleColumn() {
    std::vector<std::vector<int>> matrix = {{1}, {2}, {3}};
    std::vector<int> want = {1, 2, 3};
    assert(SpiralOrder(matrix) == want);
}

void testSingleElement() {
    std::vector<std::vector<int>> matrix = {{7}};
    std::vector<int> want = {7};
    assert(SpiralOrder(matrix) == want);
}

void test4x3() {
    std::vector<std::vector<int>> matrix = {
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
        {10, 11, 12},
    };
    std::vector<int> want = {1, 2, 3, 6, 9, 12, 11, 10, 7, 4, 5, 8};
    assert(SpiralOrder(matrix) == want);
}

void test2x2() {
    std::vector<std::vector<int>> matrix = {
        {1, 2},
        {3, 4},
    };
    std::vector<int> want = {1, 2, 4, 3};
    assert(SpiralOrder(matrix) == want);
}

void testNegativeValues() {
    std::vector<std::vector<int>> matrix = {
        {-1, -2},
        {-3, -4},
    };
    std::vector<int> want = {-1, -2, -4, -3};
    assert(SpiralOrder(matrix) == want);
}

void test4x4() {
    std::vector<std::vector<int>> matrix = {
        {1, 2, 3, 4},
        {5, 6, 7, 8},
        {9, 10, 11, 12},
        {13, 14, 15, 16},
    };
    std::vector<int> want = {1, 2, 3, 4, 8, 12, 16, 15, 14, 13, 9, 5, 6, 7, 11, 10};
    assert(SpiralOrder(matrix) == want);
}

int main() {
    test3x3();
    test3x4();
    testSingleRow();
    testSingleColumn();
    testSingleElement();
    test4x3();
    test2x2();
    testNegativeValues();
    test4x4();
    std::printf("all tests passed\n");
    return 0;
}
