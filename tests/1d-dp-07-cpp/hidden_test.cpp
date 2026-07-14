#include <cassert>
#include <cstdio>
#include <string>

int NumDecodings(std::string s);

void testTwoWays() {
    assert(NumDecodings("12") == 2);
}

void testThreeWays() {
    assert(NumDecodings("226") == 3);
}

void testLeadingZero() {
    assert(NumDecodings("06") == 0);
}

void testTwoDigitOnly() {
    assert(NumDecodings("10") == 1);
}

void testSingleDigit() {
    assert(NumDecodings("5") == 1);
}

void testLoneZero() {
    assert(NumDecodings("0") == 0);
}

void testJustOverTwentySix() {
    assert(NumDecodings("27") == 1);
}

void testUnresolvableZeroPair() {
    assert(NumDecodings("100") == 0);
}

void testLongerMultipleWays() {
    assert(NumDecodings("11106") == 2);
}

int main() {
    testTwoWays();
    testThreeWays();
    testLeadingZero();
    testTwoDigitOnly();
    testSingleDigit();
    testLoneZero();
    testJustOverTwentySix();
    testUnresolvableZeroPair();
    testLongerMultipleWays();
    std::printf("all tests passed\n");
    return 0;
}
