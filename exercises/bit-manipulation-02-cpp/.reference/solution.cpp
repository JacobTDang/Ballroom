#include <cstdint>

// HammingWeight returns the number of set bits ('1's) in the binary
// representation of n.
int HammingWeight(uint32_t n) {
    int count = 0;
    while (n != 0) {
        n &= n - 1;
        count++;
    }
    return count;
}
