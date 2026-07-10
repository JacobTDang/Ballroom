#include <cstdint>

// GetSum returns a + b without using the '+' or '-' operators.
int GetSum(int aIn, int bIn) {
    uint32_t a = static_cast<uint32_t>(aIn);
    uint32_t b = static_cast<uint32_t>(bIn);
    while (b != 0) {
        uint32_t carry = (a & b) << 1;
        a = a ^ b;
        b = carry;
    }
    return static_cast<int>(a);
}
