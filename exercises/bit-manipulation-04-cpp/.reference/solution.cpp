#include <cstdint>

// ReverseBits returns n with its 32 bits in reversed order.
uint32_t ReverseBits(uint32_t n) {
    uint32_t result = 0;
    for (int i = 0; i < 32; i++) {
        result = (result << 1) | (n & 1);
        n >>= 1;
    }
    return result;
}
