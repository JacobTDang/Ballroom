#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> PlusOne(std::vector<int>& digits);

void testSimple() {
    std::vector<int> digits = {1, 2, 3};
    std::vector<int> want = {1, 2, 4};
    assert(PlusOne(digits) == want);
}

void testAllNines() {
    std::vector<int> digits = {9, 9, 9};
    std::vector<int> want = {1, 0, 0, 0};
    assert(PlusOne(digits) == want);
}

void testSingleZero() {
    std::vector<int> digits = {0};
    std::vector<int> want = {1};
    assert(PlusOne(digits) == want);
}

void testTrailingNine() {
    std::vector<int> digits = {1, 2, 9};
    std::vector<int> want = {1, 3, 0};
    assert(PlusOne(digits) == want);
}

void testSingleNine() {
    std::vector<int> digits = {9};
    std::vector<int> want = {1, 0};
    assert(PlusOne(digits) == want);
}

void testPartialTrailingNines() {
    std::vector<int> digits = {1, 9, 9};
    std::vector<int> want = {2, 0, 0};
    assert(PlusOne(digits) == want);
}

void testMixedNoCarryPastStop() {
    std::vector<int> digits = {9, 8, 9, 9};
    std::vector<int> want = {9, 9, 0, 0};
    assert(PlusOne(digits) == want);
}

void testSingleDigitNotNine() {
    std::vector<int> digits = {5};
    std::vector<int> want = {6};
    assert(PlusOne(digits) == want);
}

void testLargerAllNines() {
    std::vector<int> digits = {9, 9, 9, 9, 9};
    std::vector<int> want = {1, 0, 0, 0, 0, 0};
    assert(PlusOne(digits) == want);
}

int main() {
    testSimple();
    testAllNines();
    testSingleZero();
    testTrailingNine();
    testSingleNine();
    testPartialTrailingNines();
    testMixedNoCarryPastStop();
    testSingleDigitNotNine();
    testLargerAllNines();
    std::printf("all tests passed\n");
    return 0;
}
