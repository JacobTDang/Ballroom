#include <cassert>
#include <cstdio>
#include <string>

bool IsMatch(std::string s, std::string p);

void testNoStarMismatch() {
    assert(IsMatch("aa", "a") == false);
}

void testStarRepeat() {
    assert(IsMatch("aa", "a*") == true);
}

void testClassic() {
    assert(IsMatch("aab", "c*a*b") == true);
}

void testLongNoMatch() {
    assert(IsMatch("mississippi", "mis*is*p*.") == false);
}

void testBothEmpty() {
    assert(IsMatch("", "") == true);
}

void testEmptyStringStarZero() {
    assert(IsMatch("", "a*") == true);
}

void testDotMatchesAny() {
    assert(IsMatch("a", ".") == true);
}

void testDotStarMatchesAll() {
    assert(IsMatch("ab", ".*") == true);
}

void testLongerMatch() {
    assert(IsMatch("mississippi", "mis*is*ip*.") == true);
}

void testDotStarTrailingLiteralFails() {
    assert(IsMatch("ab", ".*c") == false);
}

int main() {
    testNoStarMismatch();
    testStarRepeat();
    testClassic();
    testLongNoMatch();
    testBothEmpty();
    testEmptyStringStarZero();
    testDotMatchesAny();
    testDotStarMatchesAll();
    testLongerMatch();
    testDotStarTrailingLiteralFails();
    std::printf("all tests passed\n");
    return 0;
}
