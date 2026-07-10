#include <cassert>
#include <cstdio>
#include <string>

bool CheckValidString(std::string s);

void testSimple() {
    assert(CheckValidString("()") == true);
}

void testStarBalances() {
    assert(CheckValidString("(*))") == true);
}

void testUnbalanced() {
    assert(CheckValidString("(()") == false);
}

void testAllStars() {
    assert(CheckValidString("***") == true);
}

void testSingleClose() {
    assert(CheckValidString(")") == false);
}

int main() {
    testSimple();
    testStarBalances();
    testUnbalanced();
    testAllStars();
    testSingleClose();
    std::printf("all tests passed\n");
    return 0;
}
