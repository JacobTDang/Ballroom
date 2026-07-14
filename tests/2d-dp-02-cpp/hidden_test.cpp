#include <cassert>
#include <cstdio>
#include <string>

int LongestCommonSubsequence(std::string text1, std::string text2);

void testClassic() {
    assert(LongestCommonSubsequence("abcde", "ace") == 3);
}

void testIdentical() {
    assert(LongestCommonSubsequence("abc", "abc") == 3);
}

void testNoCommon() {
    assert(LongestCommonSubsequence("abc", "def") == 0);
}

void testEmptyFirst() {
    assert(LongestCommonSubsequence("", "abc") == 0);
}

void testDifferentOrder() {
    assert(LongestCommonSubsequence("abc", "acb") == 2);
}

void testInterspersedNoise() {
    assert(LongestCommonSubsequence("aggtab", "gxtxayb") == 4);
}

void testRepeatedChars() {
    assert(LongestCommonSubsequence("aaaa", "aa") == 2);
}

void testSingleCharMatch() {
    assert(LongestCommonSubsequence("a", "a") == 1);
}

void testBothEmpty() {
    assert(LongestCommonSubsequence("", "") == 0);
}

int main() {
    testClassic();
    testIdentical();
    testNoCommon();
    testEmptyFirst();
    testDifferentOrder();
    testInterspersedNoise();
    testRepeatedChars();
    testSingleCharMatch();
    testBothEmpty();
    std::printf("all tests passed\n");
    return 0;
}
