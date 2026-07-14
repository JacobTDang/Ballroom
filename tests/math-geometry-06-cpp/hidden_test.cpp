#include <cassert>
#include <climits>
#include <cmath>
#include <cstdio>

double MyPow(double x, int n);

static bool approxEqual(double a, double b) {
    return std::fabs(a - b) < 1e-6;
}

int main() {
    assert(approxEqual(MyPow(2.0, 10), 1024.0));
    assert(approxEqual(MyPow(2.1, 3), 9.261));
    assert(approxEqual(MyPow(2.0, -2), 0.25));
    assert(approxEqual(MyPow(0.5, 0), 1.0));
    assert(approxEqual(MyPow(-2.0, 3), -8.0));
    // x = 1 keeps the expected result exact regardless of exponent
    // magnitude, while still exercising the negate-the-most-negative-
    // exponent overflow edge case.
    assert(approxEqual(MyPow(1.0, INT_MIN), 1.0));
    assert(approxEqual(MyPow(-2.0, -2), 0.25));
    assert(approxEqual(MyPow(3.0, 5), 243.0));
    assert(approxEqual(MyPow(1.5, 2), 2.25));
    std::printf("all assertions passed\n");
    return 0;
}
