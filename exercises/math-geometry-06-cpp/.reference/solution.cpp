// MyPow computes x raised to the power n in O(log n) time.
double MyPow(double x, int n) {
    // Use a 64-bit exponent so negating the most-negative n never
    // overflows.
    long long exp = n;
    if (exp < 0) {
        x = 1 / x;
        exp = -exp;
    }

    double result = 1.0;
    double base = x;
    while (exp > 0) {
        if (exp % 2 == 1) {
            result *= base;
        }
        base *= base;
        exp /= 2;
    }
    return result;
}
