#include <cassert>
#include <cstdio>
#include <string>

int CountSubstrings(std::string s);

void testThreeDistinct() {
    assert(CountSubstrings("abc") == 3);
}

void testAllSame() {
    assert(CountSubstrings("aaa") == 6);
}

void testOddPalindrome() {
    assert(CountSubstrings("aba") == 4);
}

void testSingleChar() {
    assert(CountSubstrings("z") == 1);
}

void testTwoSame() {
    assert(CountSubstrings("aa") == 3);
}

void testTwoDifferent() {
    assert(CountSubstrings("ab") == 2);
}

void testLargerAllSame() {
    assert(CountSubstrings("aaaaa") == 15);
}

void testNestedPalindromes() {
    assert(CountSubstrings("aabaa") == 9);
}

int main() {
    testThreeDistinct();
    testAllSame();
    testOddPalindrome();
    testSingleChar();
    testTwoSame();
    testTwoDifferent();
    testLargerAllSame();
    testNestedPalindromes();
    std::printf("all tests passed\n");
    return 0;
}
