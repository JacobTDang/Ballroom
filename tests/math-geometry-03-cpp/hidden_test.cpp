#include <cassert>
#include <cstdio>
#include <vector>

void SetZeroes(std::vector<std::vector<int>>& matrix);

void testClassic() {
    std::vector<std::vector<int>> matrix = {
        {1, 1, 1},
        {1, 0, 1},
        {1, 1, 1},
    };
    std::vector<std::vector<int>> want = {
        {1, 0, 1},
        {0, 0, 0},
        {1, 0, 1},
    };
    SetZeroes(matrix);
    assert(matrix == want);
}

void testTwoZeroes() {
    std::vector<std::vector<int>> matrix = {
        {0, 1, 2, 0},
        {3, 4, 5, 2},
        {1, 3, 1, 5},
    };
    std::vector<std::vector<int>> want = {
        {0, 0, 0, 0},
        {0, 4, 5, 0},
        {0, 3, 1, 0},
    };
    SetZeroes(matrix);
    assert(matrix == want);
}

void testSingleZero() {
    std::vector<std::vector<int>> matrix = {
        {1, 0},
        {1, 1},
    };
    std::vector<std::vector<int>> want = {
        {0, 0},
        {1, 0},
    };
    SetZeroes(matrix);
    assert(matrix == want);
}

void testNoZero() {
    std::vector<std::vector<int>> matrix = {
        {1, 2},
        {3, 4},
    };
    std::vector<std::vector<int>> want = {
        {1, 2},
        {3, 4},
    };
    SetZeroes(matrix);
    assert(matrix == want);
}

void testAllZeros() {
    std::vector<std::vector<int>> matrix = {
        {0, 0},
        {0, 0},
    };
    std::vector<std::vector<int>> want = {
        {0, 0},
        {0, 0},
    };
    SetZeroes(matrix);
    assert(matrix == want);
}

void testCornerZero() {
    std::vector<std::vector<int>> matrix = {
        {0, 1},
        {1, 1},
    };
    std::vector<std::vector<int>> want = {
        {0, 0},
        {0, 1},
    };
    SetZeroes(matrix);
    assert(matrix == want);
}

void testSingleRow() {
    std::vector<std::vector<int>> matrix = {{1, 0, 3}};
    std::vector<std::vector<int>> want = {{0, 0, 0}};
    SetZeroes(matrix);
    assert(matrix == want);
}

void testSingleColumn() {
    std::vector<std::vector<int>> matrix = {{1}, {0}, {3}};
    std::vector<std::vector<int>> want = {{0}, {0}, {0}};
    SetZeroes(matrix);
    assert(matrix == want);
}

void testNegativeValues() {
    std::vector<std::vector<int>> matrix = {
        {-1, 0},
        {-2, -3},
    };
    std::vector<std::vector<int>> want = {
        {0, 0},
        {-2, 0},
    };
    SetZeroes(matrix);
    assert(matrix == want);
}

int main() {
    testClassic();
    testTwoZeroes();
    testSingleZero();
    testNoZero();
    testAllZeros();
    testCornerZero();
    testSingleRow();
    testSingleColumn();
    testNegativeValues();
    std::printf("all tests passed\n");
    return 0;
}
