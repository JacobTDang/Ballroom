#include <cassert>
#include <cstdint>
#include <cstdio>

int HammingWeight(uint32_t n);

void testClassic() {
    assert(HammingWeight(11) == 3);
}

void testZero() {
    assert(HammingWeight(0) == 0);
}

void testAllOnes() {
    assert(HammingWeight(4294967295u) == 32);
}

void testPowerOfTwo() {
    assert(HammingWeight(1u << 31) == 1);
}

void testAlternatingBits() {
    assert(HammingWeight(0xAAAAAAAAu) == 16);
}

int main() {
    testClassic();
    testZero();
    testAllOnes();
    testPowerOfTwo();
    testAlternatingBits();
    std::printf("all tests passed\n");
    return 0;
}
