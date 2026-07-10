#include <cassert>
#include <cstdio>
#include <climits>

int Reverse(int x);

void testClassic() {
    assert(Reverse(123) == 321);
}

void testNegative() {
    assert(Reverse(-123) == -321);
}

void testTrailingZero() {
    assert(Reverse(120) == 21);
}

void testOverflowPositive() {
    assert(Reverse(1534236469) == 0);
}

void testOverflowNegative() {
    assert(Reverse(INT_MIN) == 0);
}

void testZero() {
    assert(Reverse(0) == 0);
}

int main() {
    testClassic();
    testNegative();
    testTrailingZero();
    testOverflowPositive();
    testOverflowNegative();
    testZero();
    std::printf("all tests passed\n");
    return 0;
}
