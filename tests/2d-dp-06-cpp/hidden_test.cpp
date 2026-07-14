#include <cassert>
#include <cstdio>
#include <string>

bool IsInterleave(std::string s1, std::string s2, std::string s3);

void testClassic() {
    assert(IsInterleave("aabcc", "dbbca", "aadbbcbcac") == true);
}

void testNotInterleaved() {
    assert(IsInterleave("aabcc", "dbbca", "aadbbbaccc") == false);
}

void testAllEmpty() {
    assert(IsInterleave("", "", "") == true);
}

void testOneEmpty() {
    assert(IsInterleave("a", "", "a") == true);
}

void testLengthMismatch() {
    assert(IsInterleave("abc", "def", "abcde") == false);
}

void testFirstEmptyMatch() {
    assert(IsInterleave("", "abc", "abc") == true);
}

void testFirstEmptyMismatch() {
    assert(IsInterleave("", "abc", "abd") == false);
}

void testAmbiguousMultipleWays() {
    assert(IsInterleave("ab", "ab", "abab") == true);
}

void testRequiresBacktrackChoice() {
    assert(IsInterleave("ab", "ab", "aabb") == true);
}

int main() {
    testClassic();
    testNotInterleaved();
    testAllEmpty();
    testOneEmpty();
    testLengthMismatch();
    testFirstEmptyMatch();
    testFirstEmptyMismatch();
    testAmbiguousMultipleWays();
    testRequiresBacktrackChoice();
    std::printf("all tests passed\n");
    return 0;
}
