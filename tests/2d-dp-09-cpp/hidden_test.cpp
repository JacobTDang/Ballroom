#include <cassert>
#include <cstdio>
#include <string>

int MinDistance(std::string word1, std::string word2);

void testClassic() {
    assert(MinDistance("horse", "ros") == 3);
}

void testSecondClassic() {
    assert(MinDistance("intention", "execution") == 5);
}

void testEmptyFirst() {
    assert(MinDistance("", "abc") == 3);
}

void testIdentical() {
    assert(MinDistance("abc", "abc") == 0);
}

void testBothEmpty() {
    assert(MinDistance("", "") == 0);
}

void testEmptySecond() {
    assert(MinDistance("abc", "") == 3);
}

void testSingleCharReplace() {
    assert(MinDistance("a", "b") == 1);
}

void testPureInsertion() {
    assert(MinDistance("cat", "cats") == 1);
}

void testMultiOpMix() {
    assert(MinDistance("sunday", "saturday") == 3);
}

int main() {
    testClassic();
    testSecondClassic();
    testEmptyFirst();
    testIdentical();
    testBothEmpty();
    testEmptySecond();
    testSingleCharReplace();
    testPureInsertion();
    testMultiOpMix();
    std::printf("all tests passed\n");
    return 0;
}
