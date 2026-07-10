#include <climits>
#include <cstdint>

// Reverse returns x with its digits reversed, or 0 if the reversed
// value falls outside the signed 32-bit integer range.
int Reverse(int x) {
    int64_t result = 0;
    while (x != 0) {
        int digit = x % 10;
        x /= 10;
        result = result * 10 + digit;
        if (result < INT_MIN || result > INT_MAX) {
            return 0;
        }
    }
    return static_cast<int>(result);
}
