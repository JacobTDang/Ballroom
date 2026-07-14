#include <unordered_set>

static int SumOfSquaredDigits(int n) {
    int sum = 0;
    while (n > 0) {
        int digit = n % 10;
        sum += digit * digit;
        n /= 10;
    }
    return sum;
}

// IsHappy reports whether n is a happy number.
bool IsHappy(int n) {
    std::unordered_set<int> seen;
    while (n != 1 && seen.find(n) == seen.end()) {
        seen.insert(n);
        n = SumOfSquaredDigits(n);
    }
    return n == 1;
}
