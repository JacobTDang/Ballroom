#include <cassert>
#include <cstdio>
#include <vector>

int LongestIncreasingPath(std::vector<std::vector<int>>& matrix);

void testClassic() {
    std::vector<std::vector<int>> matrix = {{9, 9, 4}, {6, 6, 8}, {2, 1, 1}};
    assert(LongestIncreasingPath(matrix) == 4);
}

void testSecondClassic() {
    std::vector<std::vector<int>> matrix = {{3, 4, 5}, {3, 2, 6}, {2, 2, 1}};
    assert(LongestIncreasingPath(matrix) == 4);
}

void testSingleCell() {
    std::vector<std::vector<int>> matrix = {{1}};
    assert(LongestIncreasingPath(matrix) == 1);
}

void testSingleRow() {
    std::vector<std::vector<int>> matrix = {{1, 2, 3, 4}};
    assert(LongestIncreasingPath(matrix) == 4);
}

void testAllEqual() {
    std::vector<std::vector<int>> matrix = {{7, 7}, {7, 7}};
    assert(LongestIncreasingPath(matrix) == 1);
}

void testSingleColumn() {
    std::vector<std::vector<int>> matrix = {{1}, {2}, {3}};
    assert(LongestIncreasingPath(matrix) == 3);
}

void testSnakeFullTraversal() {
    std::vector<std::vector<int>> matrix = {{1, 2, 3}, {6, 5, 4}, {7, 8, 9}};
    assert(LongestIncreasingPath(matrix) == 9);
}

void testNegativeValues() {
    std::vector<std::vector<int>> matrix = {{-1, -2}, {-3, -4}};
    assert(LongestIncreasingPath(matrix) == 3);
}

int main() {
    testClassic();
    testSecondClassic();
    testSingleCell();
    testSingleRow();
    testAllEqual();
    testSingleColumn();
    testSnakeFullTraversal();
    testNegativeValues();
    std::printf("all tests passed\n");
    return 0;
}
