#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> CountBits(int n);

void testClassic() {
    std::vector<int> want = {0, 1, 1, 2, 1};
    assert(CountBits(4) == want);
}

void testSmall() {
    std::vector<int> want = {0, 1, 1};
    assert(CountBits(2) == want);
}

void testZero() {
    std::vector<int> want = {0};
    assert(CountBits(0) == want);
}

void testLarger() {
    std::vector<int> want = {0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4};
    assert(CountBits(15) == want);
}

int main() {
    testClassic();
    testSmall();
    testZero();
    testLarger();
    std::printf("all tests passed\n");
    return 0;
}
