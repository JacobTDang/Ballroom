#include <cassert>
#include <cstdio>
#include <string>

std::string LongestPalindrome(std::string s);

void testOddTie() {
    assert(LongestPalindrome("babad") == "bab");
}

void testEven() {
    assert(LongestPalindrome("cbbd") == "bb");
}

void testSingleChar() {
    assert(LongestPalindrome("a") == "a");
}

void testWholeString() {
    assert(LongestPalindrome("abba") == "abba");
}

void testAllSameLonger() {
    assert(LongestPalindrome("aaaaa") == "aaaaa");
}

void testNoRepeat() {
    assert(LongestPalindrome("abcde") == "a");
}

void testBuriedInLargerString() {
    assert(LongestPalindrome("zzabcbayy") == "abcba");
}

int main() {
    testOddTie();
    testEven();
    testSingleChar();
    testWholeString();
    testAllSameLonger();
    testNoRepeat();
    testBuriedInLargerString();
    std::printf("all tests passed\n");
    return 0;
}
