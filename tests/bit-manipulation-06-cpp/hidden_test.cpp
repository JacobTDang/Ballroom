#include <cassert>
#include <cstdio>
#include <climits>

int GetSum(int a, int b);

void testClassic() {
    assert(GetSum(1, 1) == 2);
}

void testPositivePositive() {
    assert(GetSum(2, 3) == 5);
}

void testNegativePositiveCancel() {
    assert(GetSum(-1, 1) == 0);
}

void testTwoNegatives() {
    assert(GetSum(-5, -7) == -12);
}

void testWithZero() {
    assert(GetSum(0, 0) == 0);
}

void testMaxInt32Bounds() {
    assert(GetSum(INT_MAX, 0) == INT_MAX);
    assert(GetSum(INT_MIN, 0) == INT_MIN);
}

int main() {
    testClassic();
    testPositivePositive();
    testNegativePositiveCancel();
    testTwoNegatives();
    testWithZero();
    testMaxInt32Bounds();
    std::printf("all tests passed\n");
    return 0;
}
