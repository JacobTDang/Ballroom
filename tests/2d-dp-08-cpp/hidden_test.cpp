#include <cassert>
#include <cstdio>
#include <string>

int NumDistinct(std::string s, std::string t);

void testClassic() {
    assert(NumDistinct("rabbbit", "rabbit") == 3);
}

void testSecondClassic() {
    assert(NumDistinct("babgbag", "bag") == 5);
}

void testExactMatch() {
    assert(NumDistinct("abc", "abc") == 1);
}

void testTargetLonger() {
    assert(NumDistinct("abc", "abcd") == 0);
}

void testEmptyTarget() {
    assert(NumDistinct("abc", "") == 1);
}

void testEmptySource() {
    assert(NumDistinct("", "abc") == 0);
}

void testBothEmpty() {
    assert(NumDistinct("", "") == 1);
}

void testRepeatedCharsCombinatoric() {
    assert(NumDistinct("aaaa", "aa") == 6);
}

void testSingleCharManyOccurrences() {
    assert(NumDistinct("aaa", "a") == 3);
}

int main() {
    testClassic();
    testSecondClassic();
    testExactMatch();
    testTargetLonger();
    testEmptyTarget();
    testEmptySource();
    testBothEmpty();
    testRepeatedCharsCombinatoric();
    testSingleCharManyOccurrences();
    std::printf("all tests passed\n");
    return 0;
}
