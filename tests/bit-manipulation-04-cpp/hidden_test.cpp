#include <cassert>
#include <cstdint>
#include <cstdio>

uint32_t ReverseBits(uint32_t n);

void testOne() {
    assert(ReverseBits(1u) == 2147483648u);
}

void testZero() {
    assert(ReverseBits(0u) == 0u);
}

void testAllOnes() {
    assert(ReverseBits(4294967295u) == 4294967295u);
}

void testTwo() {
    assert(ReverseBits(2u) == 1073741824u);
}

void testClassic() {
    assert(ReverseBits(43261596u) == 964176192u);
}

int main() {
    testOne();
    testZero();
    testAllOnes();
    testTwo();
    testClassic();
    std::printf("all tests passed\n");
    return 0;
}
